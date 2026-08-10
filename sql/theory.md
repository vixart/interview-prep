# SQL: теория

Справочник по ключевым словам и конструкциям SQL — что делает, как пишется, где подвох.
Собрано по слайдам из карточек в этой папке и разложено по темам.
Практика — в [`questions.md`](questions.md).

Диалект по умолчанию **PostgreSQL**; заметные отличия MySQL 8 и SQL Server отмечены.

**Содержание**

1. [DDL: объекты базы](#1-ddl-объекты-базы) · 2. [DML: изменение данных](#2-dml-изменение-данных) ·
3. [Выборка и фильтрация](#3-выборка-и-фильтрация) · 4. [NULL и условная логика](#4-null-и-условная-логика) ·
5. [Агрегация и группировка](#5-агрегация-и-группировка) · 6. [Соединения](#6-соединения) ·
7. [Операции над множествами](#7-операции-над-множествами) · 8. [Сортировка и ограничение выборки](#8-сортировка-и-ограничение-выборки) ·
9. [Оконные функции](#9-оконные-функции) · 10. [PIVOT и UNPIVOT](#10-pivot-и-unpivot) ·
11. [Строки и числа](#11-строки-и-числа) · 12. [Дата и время](#12-дата-и-время) ·
13. [Ограничения целостности](#13-ограничения-целостности) · 14. [CTE](#14-cte-common-table-expressions) ·
15. [Порядок выполнения](#15-порядок-выполнения-запроса)

---

## 1. DDL: объекты базы

DDL (Data Definition Language) описывает **структуру**, а не содержимое.

| Команда | Что делает |
| --- | --- |
| `CREATE DATABASE` | создает базу |
| `CREATE TABLE` | создает таблицу |
| `ALTER TABLE` | меняет структуру: колонки, типы, ограничения |
| `DROP TABLE` | удаляет таблицу вместе с данными и структурой |
| `TRUNCATE TABLE` | удаляет все строки, таблица остается |
| `CREATE INDEX` / `DROP INDEX` | создает и удаляет индекс |

```sql
CREATE TABLE employees (
  id      serial PRIMARY KEY,
  name    text NOT NULL,
  salary  numeric(10,2) CHECK (salary > 0),
  dept_id int REFERENCES departments(id)
);

ALTER TABLE employees ADD COLUMN hired_at date;
ALTER TABLE employees ALTER COLUMN name TYPE varchar(200);
CREATE INDEX idx_employees_dept ON employees (dept_id);
```

**`DELETE` vs `TRUNCATE` vs `DROP`** — классический вопрос:

| | `DELETE FROM t` | `TRUNCATE t` | `DROP TABLE t` |
| --- | --- | --- | --- |
| Что удаляет | строки (можно с `WHERE`) | все строки | строки + саму таблицу |
| Тип команды | DML | DDL | DDL |
| Скорость | медленно, построчно | быстро, освобождает страницы целиком |  — |
| Триггеры | срабатывают | нет | нет |
| Откат | да, в транзакции | в PostgreSQL да, в MySQL/Oracle нет | как у диалекта |
| Счетчик `AUTO_INCREMENT` | сохраняется | сбрасывается | — |

**Про индексы**: ускоряют чтение, но замедляют `INSERT`/`UPDATE`/`DELETE` и занимают место.
Индекс не будет использован, если колонка обернута функцией (`WHERE UPPER(name) = ...`) —
для таких случаев создают функциональный индекс `CREATE INDEX ON t (lower(name))`.

---

## 2. DML: изменение данных

```sql
SELECT id, name FROM employees WHERE dept_id = 10;      -- прочитать
INSERT INTO employees (name, salary) VALUES ('Ann', 300);
UPDATE employees SET salary = salary * 1.1 WHERE dept_id = 10;
DELETE FROM employees WHERE id = 42;
```

- **`INSERT`** умеет вставлять результат запроса: `INSERT INTO target (a, b) SELECT a, b FROM source;`
- **`SELECT ... INTO`** — создать таблицу из результата запроса. В PostgreSQL идиоматичнее
  `CREATE TABLE new AS SELECT ...`, в SQL Server — `SELECT ... INTO new FROM ...`.
- **`UPDATE` и `DELETE` без `WHERE`** затрагивают всю таблицу. Привычка: сначала выполнить
  `SELECT` с тем же условием, потом менять, и все — в транзакции.
- **`MERGE`** (`UPSERT`) — за один проход вставляет, обновляет или удаляет строки приемника
  по результату соединения с источником:

```sql
MERGE INTO target t USING source s ON t.id = s.id
WHEN MATCHED THEN UPDATE SET name = s.name
WHEN NOT MATCHED THEN INSERT (id, name) VALUES (s.id, s.name);
```

В PostgreSQL то же обычно пишут короче и безопаснее относительно гонок:
`INSERT ... ON CONFLICT (id) DO UPDATE SET ...`; в MySQL — `ON DUPLICATE KEY UPDATE`.

---

## 3. Выборка и фильтрация

**`SELECT`** — что вернуть, **`FROM`** — откуда, **`WHERE`** — какие строки оставить.

| Конструкция | Смысл |
| --- | --- |
| `DISTINCT` | убрать повторяющиеся строки результата |
| `WHERE` | фильтр строк **до** группировки |
| `AND` / `OR` / `NOT` | склейка и отрицание условий |
| `BETWEEN a AND b` | диапазон, **включая обе границы** |
| `IN (...)` | значение входит в список или результат подзапроса |
| `EXISTS (...)` | подзапрос вернул хотя бы одну строку |
| `LIKE` / `ILIKE` | сопоставление с шаблоном (`%` — любые символы, `_` — один) |

```sql
SELECT DISTINCT dept_id FROM employees
WHERE  salary BETWEEN 10000 AND 20000
  AND  dept_id IN (10, 20)
  AND  NOT (name LIKE 'test%');
```

Подводные камни, которые спрашивают:

- **`BETWEEN` с датами** отрезает время: `ordered_at BETWEEN '2026-01-01' AND '2026-01-31'`
  не захватит 31 января после полуночи. Для времени — полуинтервал
  `>= '2026-01-01' AND < '2026-02-01'`;
- **`NOT IN` + NULL**: если подзапрос вернет хоть один NULL, результат будет пустым всегда.
  `NOT EXISTS` от этого свободен;
- **`IN` vs `EXISTS`**: `EXISTS` останавливается на первом совпадении и обычно выигрывает
  на больших подзапросах; `IN` удобнее на коротком списке констант;
- приоритет `AND` выше, чем у `OR` — `a OR b AND c` это `a OR (b AND c)`. Ставьте скобки.

---

## 4. NULL и условная логика

**NULL — это не значение, а его отсутствие.** Любое сравнение с ним дает не «истину»
и не «ложь», а **UNKNOWN**, поэтому строка не попадает в результат.

```sql
WHERE dept_id = NULL      -- никогда не сработает
WHERE dept_id IS NULL     -- правильно
WHERE dept_id IS NOT NULL
```

| Функция | Что делает |
| --- | --- |
| `COALESCE(a, b, c)` | первый не-NULL из списка |
| `NULLIF(a, b)` | NULL, если `a = b`, иначе `a` |
| `IIF(cond, x, y)` | сокращенный `CASE` (SQL Server; в MySQL — `IF`) |
| `CASE WHEN ... THEN ... ELSE ... END` | полноценное ветвление, стандарт |

```sql
SELECT name,
       COALESCE(salary, 0)                       AS salary,          -- NULL → 0
       amount / NULLIF(total, 0)                 AS share,           -- защита от деления на 0
       CASE WHEN salary >= 300 THEN 'senior'
            WHEN salary >= 150 THEN 'middle'
            ELSE 'junior' END                    AS grade
FROM   employees;
```

Что еще важно про NULL:

- агрегаты (`SUM`, `AVG`, `COUNT(колонка)`) **пропускают** NULL, а `COUNT(*)` считает строки;
- `NULL` не равен `NULL`, но `UNION`, `INTERSECT`, `GROUP BY` и `DISTINCT`
  считают их **одинаковыми** и схлопывают;
- сравнение, где NULL должен вести себя как значение, — `IS DISTINCT FROM`;
- `CASE` возвращает **значение**, а не выполняет действие, поэтому его можно ставить
  и в `SELECT`, и в `WHERE`, и в `ORDER BY`, и внутрь агрегата (см. пивот).

---

## 5. Агрегация и группировка

| Агрегат | Что считает |
| --- | --- |
| `COUNT(*)` | число строк |
| `COUNT(col)` | число строк, где `col` не NULL |
| `COUNT(DISTINCT col)` | число различных значений |
| `SUM` / `AVG` | сумма и среднее |
| `MIN` / `MAX` | минимум и максимум |

**`GROUP BY`** схлопывает строки с одинаковыми значениями в одну и считает по ним агрегаты.
**`HAVING`** фильтрует уже полученные группы — это `WHERE` для агрегатов.

```sql
SELECT dept_id, COUNT(*) AS people, AVG(salary) AS avg_salary
FROM   employees
WHERE  hired_at >= '2020-01-01'      -- фильтр СТРОК, до группировки
GROUP  BY dept_id
HAVING COUNT(*) > 5                  -- фильтр ГРУПП, после группировки
ORDER  BY avg_salary DESC;
```

Правило: в `SELECT` при `GROUP BY` допустимы только колонки из `GROUP BY` и агрегаты.
`SELECT *` вместе с `GROUP BY` — ошибка (MySQL исторически это позволял и выдавал
случайную строку из группы, что хуже ошибки).

**Многоуровневые итоги** — расширения `GROUP BY`:

| Конструкция | Что дает |
| --- | --- |
| `ROLLUP(a, b)` | итоги по иерархии: `(a,b)`, `(a)`, `()` — подытоги и общий итог |
| `CUBE(a, b)` | все комбинации: `(a,b)`, `(a)`, `(b)`, `()` |
| `GROUPING SETS ((a),(b))` | ровно те группировки, что перечислены |

```sql
SELECT dept_id, EXTRACT(YEAR FROM hired_at) AS y, COUNT(*)
FROM   employees
GROUP  BY ROLLUP (dept_id, y);       -- по отделу и году, по отделу, и общий итог
```

Строки итогов отличают по `GROUPING(dept_id) = 1` — иначе их не отделить от настоящего NULL.

---

## 6. Соединения

`JOIN` склеивает строки двух таблиц по условию.

```
   employees        departments
   ┌────┬───────┐   ┌────┬──────┐
   │ id │dept_id│   │ id │ name │
   ├────┼───────┤   ├────┼──────┤
   │ 1  │  10   │   │ 10 │ IT   │      INNER  → только совпавшие пары
   │ 2  │  20   │   │ 20 │ HR   │      LEFT   → + сотрудники без отдела
   │ 3  │ NULL  │   │ 30 │ Fin  │      RIGHT  → + отделы без сотрудников
   └────┴───────┘   └────┴──────┘      FULL   → и то, и другое
```

| Вид | Что возвращает |
| --- | --- |
| `INNER JOIN` | только строки, для которых нашлось совпадение в обеих таблицах |
| `LEFT JOIN` | все строки левой таблицы; где пары нет — NULL в колонках правой |
| `RIGHT JOIN` | зеркально: все строки правой |
| `FULL JOIN` | **все строки обеих таблиц**, с NULL там, где пары нет |
| `CROSS JOIN` | декартово произведение: каждая с каждой |
| `SELF JOIN` | таблица соединяется сама с собой (иерархии, пары) |

Часто путают формулировку `FULL JOIN`: это не «строки, где есть совпадение хотя бы в одной
таблице», а **объединение** — совпавшие пары плюс несовпавшие строки слева и справа.

Типовые ловушки:

- условие для правой таблицы в `LEFT JOIN` нужно писать **в `ON`, а не в `WHERE`**:
  в `WHERE` оно отбросит строки с NULL и превратит соединение в `INNER`;
- `COUNT(*)` после `LEFT JOIN` посчитает и «пустые» строки — считайте `COUNT(правая.колонка)`;
- дублирование строк при соединении «один ко многим» — проверяйте уникальность ключа справа;
- **антиджойн** «есть слева, нет справа»: `LEFT JOIN ... WHERE right.id IS NULL` или `NOT EXISTS`.

**`CROSS APPLY` / `OUTER APPLY`** (SQL Server) — соединение с подзапросом, который **зависит
от текущей строки** левой таблицы: например, «по три последних заказа на каждого клиента».
В PostgreSQL это `CROSS JOIN LATERAL` и `LEFT JOIN LATERAL ... ON true`.

```sql
SELECT c.id, o.*
FROM   customers c
CROSS  JOIN LATERAL (
  SELECT * FROM orders o WHERE o.customer_id = c.id ORDER BY ordered_at DESC LIMIT 3
) o;
```

---

## 7. Операции над множествами

Работают со **строками целиком**: число и типы колонок должны совпадать.

| Оператор | Что делает |
| --- | --- |
| `UNION` | объединение, дубликаты убираются |
| `UNION ALL` | объединение как есть, дубликаты остаются |
| `INTERSECT` | только строки, встречающиеся в обоих запросах |
| `EXCEPT` (Oracle: `MINUS`) | строки первого запроса, которых нет во втором |

`UNION ALL` **быстрее**: `UNION` вынужден сортировать или хешировать результат, чтобы найти
дубликаты. Если дублей заведомо нет — всегда `UNION ALL`.

В отличие от `=`, эти операторы считают `NULL` равными друг другу, поэтому
`INTERSECT` и `EXCEPT` удобны для сравнения таблиц «на глаз».

`ORDER BY` пишется один раз, в самом конце — он относится ко всему объединенному результату.

---

## 8. Сортировка и ограничение выборки

```sql
SELECT * FROM employees
ORDER  BY salary DESC, name ASC
LIMIT  10 OFFSET 20;
```

| Конструкция | Где |
| --- | --- |
| `LIMIT n OFFSET m` | PostgreSQL, MySQL, SQLite |
| `OFFSET m ROWS FETCH NEXT n ROWS ONLY` | стандарт SQL, SQL Server, Oracle 12c+ |
| `TOP n` | SQL Server (`SELECT TOP 10 ...`) |
| `FETCH FIRST n ROWS WITH TIES` | вернуть еще и строки, равные последней по сортировке |

- **без `ORDER BY` порядок строк не определен** — «первые 10» без сортировки бессмысленны;
- сортировка идет **после** `SELECT`, поэтому в `ORDER BY` можно ссылаться на алиасы колонок
  (в `WHERE` — нельзя);
- NULL в сортировке: в PostgreSQL по умолчанию они «больше всех» (`NULLS LAST` при `ASC`),
  порядок настраивается через `ORDER BY col DESC NULLS LAST`;
- **`OFFSET` на больших числах медленный**: база все равно читает и отбрасывает пропущенные
  строки. Для листания используют keyset-пагинацию: `WHERE (salary, id) < (:last_salary, :last_id)`.

---

## 9. Оконные функции

Оконная функция **не схлопывает строки**, как `GROUP BY`, а дописывает к каждой строке
вычисленное по «окну» значение.

```sql
функция() OVER (
    PARTITION BY отдел      -- разбить на группы, считать в каждой отдельно (необязательно)
    ORDER BY   зарплата     -- порядок внутри группы (для рангов и LAG/LEAD обязателен)
    ROWS BETWEEN ... )      -- кадр: какие строки участвуют (необязательно)
```

**Ранжирование**

| Функция | Что делает | Ничьи |
| --- | --- | --- |
| `ROW_NUMBER()` | сквозной уникальный номер | нумерует произвольно |
| `RANK()` | ранг с **пропусками** после ничьей | 1, 2, 2, **4** |
| `DENSE_RANK()` | ранг **без пропусков** | 1, 2, 2, **3** |
| `NTILE(n)` | делит строки на `n` примерно равных групп | для перцентилей и квартилей |

**Смещение по строкам**

| Функция | Что делает |
| --- | --- |
| `LAG(col, n)` | значение из строки на `n` позиций **назад** |
| `LEAD(col, n)` | значение из строки на `n` позиций **вперед** |
| `FIRST_VALUE` / `LAST_VALUE` / `NTH_VALUE(col, n)` | значение из первой / последней / N-й строки кадра |

```sql
SELECT month, revenue,
       LAG(revenue)  OVER (ORDER BY month) AS prev_month,
       revenue - LAG(revenue) OVER (ORDER BY month) AS diff
FROM   monthly_revenue;
```

**Кадр окна: `ROWS` против `RANGE`**

- `ROWS` считает **физические строки**: `ROWS BETWEEN 2 PRECEDING AND CURRENT ROW` — ровно три строки;
- `RANGE` считает **по значению** сортировки: все строки с равным значением попадают в кадр разом.

Отсюда классическая ошибка в нарастающем итоге: при нескольких строках с одной датой
`RANGE` (кадр по умолчанию!) сложит всю дату сразу, а `ROWS` — построчно.

```sql
SUM(amount) OVER (ORDER BY sale_date ROWS BETWEEN UNBOUNDED PRECEDING AND CURRENT ROW)
```

Что еще спрашивают:

- окно считается **после `WHERE` и `GROUP BY`, но до `SELECT`**, поэтому по результату окна
  нельзя фильтровать в том же запросе — нужен CTE или подзапрос;
- агрегаты тоже работают как оконные: `AVG(salary) OVER (PARTITION BY dept_id)` даст среднее
  по отделу рядом с каждой строкой, не убирая строки;
- `FETCH` в стандарте — это про курсоры («выдать очередную порцию строк»), но в повседневном
  SQL это `OFFSET ... FETCH NEXT ... ROWS ONLY` из раздела 8.

---

## 10. PIVOT и UNPIVOT

**`PIVOT`** — значения одной колонки становятся **колонками** результата.
**`UNPIVOT`** — обратная операция: несколько колонок сворачиваются в строки.

```
       PIVOT →                             ← UNPIVOT
 dept │ cnt        ┌──────┬────┬────┐    id │ q1 │ q2      id │ quarter │ amount
 ─────┼────        │  IT  │ HR │Fin │    ───┼────┼───      ───┼─────────┼───────
 IT   │  3    ⇒    │   3  │  2 │ 1  │     1 │ 10 │ 20   ⇒   1 │   q1    │  10
 HR   │  2         └──────┴────┴────┘                       1 │   q2    │  20
```

Синтаксис `PIVOT`/`UNPIVOT` есть в SQL Server и Oracle. В PostgreSQL и MySQL то же самое
пишут вручную — и это переносимо:

```sql
-- пивот: агрегат с условием
SELECT COUNT(*) FILTER (WHERE d.name = 'IT') AS it,     -- PostgreSQL
       SUM(CASE WHEN d.name = 'HR' THEN 1 ELSE 0 END) AS hr   -- работает везде
FROM   employees e JOIN departments d ON d.id = e.dept_id;

-- анпивот: UNION ALL или VALUES
SELECT id, 'q1' AS quarter, q1 AS amount FROM revenue
UNION ALL
SELECT id, 'q2', q2 FROM revenue;
```

Ограничение любого пивота: **список колонок задается в запросе заранее**. Динамический
пивот в чистом SQL невозможен — генерируют текст запроса или разворачивают в приложении.

---

## 11. Строки и числа

| Функция | Что делает |
| --- | --- |
| `CONCAT(a, b)` / `a \|\| b` | склеить строки |
| `SUBSTRING(s FROM 2 FOR 5)` | подстрока с позиции и длины |
| `REPLACE(s, from, to)` | заменить все вхождения подстроки |
| `LENGTH(s)` | длина строки (SQL Server: `LEN`) |
| `UPPER(s)` / `LOWER(s)` | регистр; `INITCAP` — с заглавных |
| `TRIM(s)` / `LTRIM` / `RTRIM` | убрать пробелы по краям |
| `POSITION(sub IN s)` | позиция подстроки |
| `ROUND(x, n)` | округление числа до `n` знаков |

```sql
SELECT CONCAT_WS(' ', first_name, last_name) AS full_name,   -- пропускает NULL
       UPPER(TRIM(email))                    AS email,
       ROUND(salary / 12.0, 2)               AS monthly
FROM   employees;
```

Нюансы:

- `||` и `CONCAT` в PostgreSQL дают **NULL, если хоть один аргумент NULL**; `CONCAT_WS` — нет.
  В SQL Server `CONCAT` NULL игнорирует, а `+` — нет;
- `LENGTH` в PostgreSQL считает **символы**, `octet_length` — байты; в MySQL наоборот:
  `CHAR_LENGTH` — символы, `LENGTH` — байты;
- позиции в SQL нумеруются **с единицы**, а не с нуля;
- `TRIM` убирает пробелы только по краям; внутри строки — `REPLACE` или `regexp_replace`;
- `ROUND` на `float` может удивить из-за двоичного представления — для денег используйте
  `numeric`/`decimal`, а не `float`.

---

## 12. Дата и время

| Функция | Что делает |
| --- | --- |
| `CURRENT_DATE` | текущая дата |
| `CURRENT_TIME` | текущее время |
| `CURRENT_TIMESTAMP` | дата и время (SQL Server: `GETDATE()`, MySQL: `NOW()`) |
| `DATEADD(unit, n, d)` | прибавить интервал (SQL Server) |
| `DATEDIFF(unit, a, b)` | разница между датами (SQL Server; в MySQL — только в днях) |
| `DATEPART(unit, d)` / `EXTRACT(unit FROM d)` | вытащить год, месяц, день, час |
| `DATE_TRUNC('month', d)` | «обрезать» дату до начала периода |

```sql
-- PostgreSQL
SELECT CURRENT_DATE - INTERVAL '6 months'      AS half_year_ago,
       EXTRACT(YEAR FROM hired_at)             AS year,
       date_trunc('month', hired_at)           AS month,
       (end_date - start_date)                 AS days;      -- date - date = число дней

-- SQL Server
SELECT DATEADD(month, -6, GETDATE()), DATEPART(year, hired_at), DATEDIFF(day, a, b);

-- MySQL
SELECT DATE_SUB(NOW(), INTERVAL 6 MONTH), YEAR(hired_at), DATEDIFF(b, a);
```

Практика:

- для группировки по периодам берут `date_trunc`, а не `EXTRACT`: сортировка остается
  хронологической и не путаются разные годы;
- условие пишется на «голую» колонку (`WHERE hired_at >= ...`), иначе индекс не работает;
- храните время в UTC с типом `timestamptz`, а часовой пояс применяйте на выводе;
- разница `timestamp - timestamp` дает интервал, а не число: приводите к дате или к дням явно.

---

## 13. Ограничения целостности

Ограничения — это правила, которые проверяет **сама база**, а не приложение.

| Ограничение | Что гарантирует |
| --- | --- |
| `PRIMARY KEY` | уникальность строки; NOT NULL + UNIQUE, ровно один на таблицу |
| `FOREIGN KEY` | значение существует в другой таблице (ссылочная целостность) |
| `UNIQUE` | все значения колонки различны; NULL допускается |
| `NOT NULL` | значение обязано быть |
| `CHECK` | произвольное условие: `CHECK (salary > 0)` |
| `DEFAULT` | значение по умолчанию при вставке |

`CONSTRAINT имя ...` дает ограничению имя — тогда его видно в ошибках и можно удалить:

```sql
ALTER TABLE employees
  ADD CONSTRAINT chk_salary_positive CHECK (salary > 0);
```

**Что делать со ссылающимися строками** при удалении или изменении родителя:

| Действие | Поведение при `DELETE` родителя |
| --- | --- |
| `NO ACTION` | ошибка (проверка может быть отложена до конца транзакции) |
| `RESTRICT` | ошибка сразу, отложить проверку нельзя |
| `CASCADE` | удалить и дочерние строки тоже |
| `SET NULL` | обнулить внешний ключ у дочерних |
| `SET DEFAULT` | поставить значение по умолчанию |

```sql
CREATE TABLE orders (
  id          serial PRIMARY KEY,
  customer_id int REFERENCES customers(id) ON DELETE CASCADE ON UPDATE CASCADE
);
```

`ON DELETE CASCADE` удобен, но опасен: одна строка может утащить за собой пол-базы.
В финансовых и аудируемых данных чаще выбирают `RESTRICT` и мягкое удаление флагом.

---

## 14. CTE (Common Table Expressions)

`WITH` дает подзапросу имя и делает многоэтажный запрос читаемым.

```sql
WITH monthly AS (
  SELECT date_trunc('month', ordered_at) AS month, SUM(amount) AS revenue
  FROM   orders GROUP BY 1
), ranked AS (
  SELECT month, revenue, RANK() OVER (ORDER BY revenue DESC) AS rnk
  FROM   monthly
)
SELECT * FROM ranked WHERE rnk <= 3;
```

- CTE читается сверху вниз, поэтому им заменяют вложенные подзапросы и разбивают логику на шаги;
- через CTE фильтруют результат оконной функции — иначе никак;
- **рекурсивный CTE** обходит иерархии и генерирует ряды:

```sql
WITH RECURSIVE tree AS (
  SELECT id, name, manager_id, 1 AS level FROM employees WHERE manager_id IS NULL
  UNION ALL
  SELECT e.id, e.name, e.manager_id, t.level + 1
  FROM   employees e JOIN tree t ON e.manager_id = t.id
)
SELECT * FROM tree;
```

Про производительность: в PostgreSQL до версии 12 CTE всегда материализовался (был барьером
для оптимизатора), с 12-й он может встраиваться — управляется `MATERIALIZED` /
`NOT MATERIALIZED`. В MySQL 8 и SQL Server CTE обычно встраивается всегда.

---

## 15. Порядок выполнения запроса

Написан запрос в одном порядке, а выполняется в другом:

```
FROM / JOIN → WHERE → GROUP BY → HAVING → оконные функции → SELECT → DISTINCT → ORDER BY → LIMIT
```

Отсюда следует почти вся практика этого файла:

- в `WHERE` **не видно алиасов** из `SELECT` и не работают оконные функции — они еще не вычислены;
- в `ORDER BY` алиасы видны, потому что он идет после `SELECT`;
- фильтр по агрегату — только в `HAVING`;
- фильтр по результату окна — только снаружи, через CTE или подзапрос;
- `WHERE` дешевле `HAVING`: он отсекает строки до группировки, поэтому все, что можно
  отфильтровать раньше, фильтруйте раньше.
