# Transactional Outbox

Паттерн надёжной публикации событий: как гарантировать, что событие уйдёт
в брокер ровно тогда, когда изменение действительно попало в базу.

Примеры на PostgreSQL + Go. Про `FOR UPDATE SKIP LOCKED` и индексы —
[`../sql/indexes.md`](../sql/indexes.md), про сами брокеры —
[`../messaging/`](../messaging/README.md).

---

## Проблема: dual write

Классический код «сохранили и отправили»:

```go
db.Exec("INSERT INTO orders ...")   // 1. запись в БД
broker.Publish("order.created", e)  // 2. отправка в Kafka/RabbitMQ
```

Это две разные системы и **две несвязанные операции**. Атомарности нет:

| Что упало | Результат |
|---|---|
| Упали между 1 и 2 | Заказ есть в БД, события нет. Склад не узнает, письмо не уйдёт |
| Брокер недоступен | То же самое |
| Публикация прошла, транзакция откатилась | Событие о заказе есть, заказа нет. Хуже всего |

Поменять порядок не помогает — просто меняется вид расхождения.
Распределённые транзакции (2PC) решают, но брокеры их обычно не
поддерживают, а где поддерживают — это медленно и хрупко.

---

## Решение

Ключевая идея: **не выходить за пределы одной БД-транзакции**.

Событие пишется не в брокер, а в обычную таблицу `outbox` — в **той же
транзакции**, что и бизнес-данные. Дальше отдельный процесс (relay)
читает эту таблицу и публикует.

```
┌─────────────────────────────────────┐
│          одна транзакция            │
│  INSERT INTO orders  ─┐             │
│  INSERT INTO outbox  ─┴─ COMMIT     │  ← атомарно: либо оба, либо ни одного
└─────────────────────────────────────┘
                 │
                 ▼
          ┌─────────────┐
          │   outbox    │
          └──────┬──────┘
                 │  relay читает неопубликованные
                 ▼
          ┌─────────────┐      ┌──────────┐
          │   relay     │─────►│  Kafka   │
          └──────┬──────┘      └──────────┘
                 │  помечает published_at
                 ▼
```

Если relay упал до публикации — строка осталась, событие уйдёт позже.
Если упал после публикации, но до отметки — событие уйдёт **повторно**.
Отсюда главное следствие: **гарантия at-least-once**, потребитель обязан
быть идемпотентным.

---

## Таблица

```sql
CREATE TABLE outbox (
    id             BIGSERIAL PRIMARY KEY,
    aggregate_type TEXT        NOT NULL,   -- 'order'
    aggregate_id   TEXT        NOT NULL,   -- '12345' — ключ партиционирования
    event_type     TEXT        NOT NULL,   -- 'order.created'
    payload        JSONB       NOT NULL,
    created_at     TIMESTAMPTZ NOT NULL DEFAULT now(),
    published_at   TIMESTAMPTZ                          -- NULL = не отправлено
);

-- частичный индекс: в нём лежат только неотправленные строки
CREATE INDEX idx_outbox_unpublished ON outbox (id)
WHERE published_at IS NULL;
```

Частичный индекс здесь принципиален. Таблица растёт, а relay всегда
спрашивает только `WHERE published_at IS NULL` — индекс остаётся
маленьким независимо от объёма истории.

---

## Запись события

```go
func CreateOrder(ctx context.Context, db *sql.DB, o Order) error {
    tx, err := db.BeginTx(ctx, nil)
    if err != nil {
        return err
    }
    defer tx.Rollback() //nolint:errcheck // no-op после успешного Commit

    if _, err := tx.ExecContext(ctx,
        `INSERT INTO orders (id, customer_id, amount) VALUES ($1, $2, $3)`,
        o.ID, o.CustomerID, o.Amount,
    ); err != nil {
        return err
    }

    payload, err := json.Marshal(o)
    if err != nil {
        return err
    }

    if _, err := tx.ExecContext(ctx,
        `INSERT INTO outbox (aggregate_type, aggregate_id, event_type, payload)
         VALUES ('order', $1, 'order.created', $2)`,
        o.ID, payload,
    ); err != nil {
        return err
    }

    return tx.Commit() // ← событие и заказ фиксируются вместе
}
```

Обратите внимание: в бизнес-коде **нет обращения к брокеру**. Он не знает
про Kafka вообще. Это и есть весь паттерн со стороны записи.

---

## Relay: публикация

```go
func (r *Relay) publishBatch(ctx context.Context) (int, error) {
    tx, err := r.db.BeginTx(ctx, nil)
    if err != nil {
        return 0, err
    }
    defer tx.Rollback() //nolint:errcheck

    rows, err := tx.QueryContext(ctx, `
        SELECT id, aggregate_id, event_type, payload
        FROM outbox
        WHERE published_at IS NULL
        ORDER BY id
        LIMIT 100
        FOR UPDATE SKIP LOCKED`)   // ← позволяет запустить несколько relay
    if err != nil {
        return 0, err
    }

    var events []Event
    for rows.Next() {
        var e Event
        if err := rows.Scan(&e.ID, &e.Key, &e.Type, &e.Payload); err != nil {
            rows.Close()
            return 0, err
        }
        events = append(events, e)
    }
    rows.Close()
    if err := rows.Err(); err != nil {
        return 0, err
    }
    if len(events) == 0 {
        return 0, nil
    }

    ids := make([]int64, 0, len(events))
    for _, e := range events {
        // ключ = aggregate_id: события одного заказа попадут в одну партицию
        if err := r.broker.Publish(ctx, e.Type, e.Key, e.Payload); err != nil {
            return 0, err // не коммитим — строки останутся неотправленными
        }
        ids = append(ids, e.ID)
    }

    if _, err := tx.ExecContext(ctx,
        `UPDATE outbox SET published_at = now() WHERE id = ANY($1)`,
        pq.Array(ids),
    ); err != nil {
        return 0, err
    }

    return len(events), tx.Commit()
}
```

Три детали, которые обычно и спрашивают:

- **`FOR UPDATE SKIP LOCKED`** — строки, взятые одним воркером,
  пропускаются другими. Можно держать несколько экземпляров relay без
  дублей и без взаимных блокировок.
- **`ORDER BY id`** — порядок публикации совпадает с порядком записи.
- **Ключ сообщения = `aggregate_id`** — Kafka гарантирует порядок внутри
  партиции, значит события одного заказа не перемешаются.

Цикл опроса — обычный тикер. Чтобы не терять секунды на задержке
поллинга, в PostgreSQL можно добавить `LISTEN/NOTIFY`: писатель шлёт
`NOTIFY outbox_new`, relay просыпается сразу.

---

## Два способа читать outbox

| | Polling publisher | CDC (Change Data Capture) |
|---|---|---|
| Как | `SELECT` по таблице раз в N мс | чтение WAL через Debezium / логическую репликацию |
| Нагрузка на БД | постоянные запросы | почти нулевая |
| Задержка | равна интервалу опроса | миллисекунды |
| Сложность | ~50 строк кода | отдельный сервис, Kafka Connect, эксплуатация |
| Когда брать | по умолчанию, до десятков тысяч событий/сек | большой поток, нужна минимальная задержка |

При CDC поле `published_at` не нужно: Debezium ловит сам факт `INSERT`,
и строки обычно удаляют сразу (`DELETE` после `INSERT` в той же
транзакции — в WAL запись всё равно останется).

---

## Дубли и идемпотентность

At-least-once означает, что событие может прийти дважды. Это не дефект
реализации, а свойство паттерна — распределённой системы без дублей не
бывает.

Обязанность потребителя — обработать повтор безопасно. Стандартный
приём — **Inbox** (зеркальный паттерн):

```sql
CREATE TABLE inbox (
    event_id     UUID PRIMARY KEY,
    processed_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
```

```go
// в одной транзакции: отметка о событии + бизнес-логика
res, err := tx.ExecContext(ctx,
    `INSERT INTO inbox (event_id) VALUES ($1) ON CONFLICT DO NOTHING`, e.ID)
n, _ := res.RowsAffected()
if n == 0 {
    return tx.Commit() // уже обрабатывали — молча пропускаем
}
// ... обработка ...
return tx.Commit()
```

Альтернатива, когда операция естественно идемпотентна: `UPSERT`,
установка состояния вместо инкремента, внешний ключ идемпотентности
в вызываемом API.

---

## Примеры использования

### 1. Заказ → уведомление пользователю

Самый частый случай. Заказ должен сохраниться, письмо — уйти, и одно без
другого недопустимо.

```
POST /orders
  └─ tx: INSERT orders + INSERT outbox('order.created')
       └─ relay → Kafka topic "orders"
            └─ notification-service → письмо
            └─ warehouse-service    → резерв товара
```

Плюс к надёжности: добавить нового потребителя (аналитику, антифрод)
можно без изменения кода сервиса заказов.

### 2. Синхронизация с внешним API

Списание денег в платёжном шлюзе нельзя откатить вместе с транзакцией БД,
а HTTP-вызов внутри транзакции держит её открытой на время сети — так
делать нельзя.

```go
// tx: обновляем заказ + пишем задание
INSERT INTO outbox (event_type, payload)
VALUES ('payment.capture', '{"order_id": 42, "amount": 100}');
```

Relay вызывает платёжный API уже вне транзакции, с ретраями. Транзакция
БД короткая, а задание не потеряется. Здесь обязателен ключ
идемпотентности в запросе к шлюзу — иначе повтор спишет деньги дважды.

### 3. Смена состояния в Saga

В распределённой саге каждый шаг публикует событие для следующего.
Outbox делает эти переходы надёжными: сервис атомарно меняет своё
состояние и объявляет о переходе.

```
order-service:   tx: status='pending'  + outbox('order.created')
payment-service: tx: charge saved      + outbox('payment.succeeded')
order-service:   tx: status='paid'     + outbox('order.paid')
```

Без outbox сага рвётся на любом падении между записью и публикацией и
залипает в промежуточном состоянии.

### 4. Аудит и проекции (CQRS)

Outbox-таблица — это готовый упорядоченный журнал изменений. Его можно
читать не только брокером: наполнять read-модель, поисковый индекс
Elasticsearch, витрину в аналитическом хранилище. Источник истины при
этом остаётся один — основная БД.

---

## Подводные камни

| Проблема | Что делать |
|---|---|
| Таблица бесконечно растёт | Удалять отправленное: `DELETE FROM outbox WHERE published_at < now() - interval '7 days'`. Лучше — партиционирование по дате и `DROP PARTITION` |
| `UPDATE published_at` создаёт bloat | Либо удалять строки сразу после публикации, либо следить за `autovacuum` |
| Один медленный потребитель тормозит всё | Relay не должен ждать обработки, только приёма брокером |
| «Ядовитое» сообщение блокирует очередь | Счётчик попыток + перенос в DLQ после N неудач |
| Порядок между разными агрегатами | Не гарантируется и не должен требоваться. Гарантия только внутри одного `aggregate_id` |
| Событие пишут в outbox вне транзакции | Самая частая ошибка реализации: смысл паттерна теряется полностью |

---

## Когда не нужен

- Событие некритично, потеря допустима (метрики, необязательная
  телеметрия).
- Монолит с одной БД, где «потребитель» — та же транзакция.
- Уже есть брокер с транзакционной семантикой поверх той же БД
  (например, очередь на таблицах).

---

## Для собеседования

- **Какую проблему решает:** dual write — невозможность атомарно
  записать в БД и отправить в брокер.
- **Как:** событие пишется в таблицу той же транзакцией, отдельный
  процесс публикует его асинхронно.
- **Какая гарантия:** at-least-once. Exactly-once не даёт — потребитель
  должен быть идемпотентным (Inbox / дедупликация по `event_id`).
- **Как читают таблицу:** поллинг с `FOR UPDATE SKIP LOCKED` либо CDC
  по WAL (Debezium).
- **Порядок:** сохраняется внутри одного `aggregate_id` через ключ
  партиционирования, между агрегатами — нет.
- **Цена:** задержка на интервал поллинга, рост таблицы, лишняя запись
  в каждой транзакции.
