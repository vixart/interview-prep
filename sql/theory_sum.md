# SQL — шпора

Сжатая версия [`theory.md`](theory.md). Практика — [`questions.md`](questions.md).
PostgreSQL по умолчанию, отличия в скобках.

## Порядок выполнения

```
FROM/JOIN → WHERE → GROUP BY → HAVING → окна → SELECT → DISTINCT → ORDER BY → LIMIT
```

Отсюда: в `WHERE` нет алиасов и окон · фильтр агрегата → `HAVING` · фильтр окна → CTE/подзапрос ·
в `ORDER BY` алиасы уже есть.

## DDL

| | |
| --- | --- |
| `CREATE/ALTER/DROP TABLE` | создать / изменить структуру / удалить вместе с таблицей |
| `TRUNCATE` | удалить все строки, таблица остается |
| `CREATE INDEX` | быстрее чтение, медленнее запись |

`DELETE` — DML, построчно, триггеры, откат. `TRUNCATE` — DDL, быстро, без триггеров, сбрасывает счетчик.
`DROP` — сносит и структуру.

## DML

```sql
INSERT INTO t (a,b) VALUES (1,2);      INSERT INTO t SELECT ... ;
UPDATE t SET x = x*1.1 WHERE ... ;     DELETE FROM t WHERE ... ;
```

Upsert: `INSERT ... ON CONFLICT (id) DO UPDATE/NOTHING` (MySQL `ON DUPLICATE KEY`, стандарт `MERGE`).
**Без `WHERE` меняется вся таблица** — сначала тот же `WHERE` в `SELECT`, потом транзакция.

## Фильтрация

| | |
| --- | --- |
| `DISTINCT` | убрать дубли строк |
| `WHERE` | фильтр строк до группировки |
| `AND/OR/NOT` | `AND` приоритетнее `OR` → ставь скобки |
| `BETWEEN a AND b` | **включая обе границы** |
| `IN (...)` | значение в списке / подзапросе |
| `EXISTS (...)` | подзапрос вернул хоть строку |
| `LIKE` | `%` — любые символы, `_` — один |

Ловушки: `NOT IN` + NULL → пусто всегда (бери `NOT EXISTS`) · `BETWEEN` по датам режет время
(бери `>= начало AND < следующее`) · функция на колонке убивает индекс.

## NULL

`= NULL` не работает никогда → `IS NULL` / `IS NOT NULL`. Сравнение с NULL = UNKNOWN.

| | |
| --- | --- |
| `COALESCE(a,b,c)` | первый не-NULL (`IFNULL`, `ISNULL`, `NVL`) |
| `NULLIF(a,b)` | NULL если равны — защита от деления на 0 |
| `CASE WHEN ... THEN ... ELSE ... END` | ветвление, возвращает значение |
| `IIF(c,x,y)` | короткий CASE (SQL Server) |

Агрегаты NULL пропускают, `COUNT(*)` — считает. `GROUP BY`/`DISTINCT`/`UNION` считают NULL одинаковыми.

## Агрегация

| | |
| --- | --- |
| `COUNT(*)` / `COUNT(col)` | строки / непустые значения |
| `SUM AVG MIN MAX` | — |
| `GROUP BY` | схлопывает строки в группы |
| `HAVING` | фильтр **групп** (агрегатов) |
| `ROLLUP(a,b)` | итоги по иерархии + общий |
| `CUBE(a,b)` | все комбинации группировок |
| `GROUPING SETS` | только перечисленные |

В `SELECT` при `GROUP BY` — только колонки из `GROUP BY` и агрегаты.

## JOIN

| | |
| --- | --- |
| `INNER` | только совпавшие пары |
| `LEFT` | все слева + NULL справа |
| `RIGHT` | зеркально |
| `FULL` | **все строки обеих таблиц** |
| `CROSS` | каждая с каждой |
| `SELF` | таблица сама с собой (иерархии) |
| `LATERAL` / `CROSS APPLY` | подзапрос видит текущую строку (топ-N на группу) |

Ловушки: условие правой таблицы в `LEFT JOIN` пиши в `ON`, в `WHERE` оно сделает `INNER` ·
`COUNT(*)` после `LEFT JOIN` считает пустышки → `COUNT(правая.col)` · антиджойн =
`LEFT JOIN ... WHERE right.id IS NULL` или `NOT EXISTS`.

## Множества

| | |
| --- | --- |
| `UNION` | объединение без дублей (медленнее) |
| `UNION ALL` | объединение как есть — бери по умолчанию |
| `INTERSECT` | что есть в обоих |
| `EXCEPT` (`MINUS`) | что есть в первом и нет во втором |

Колонок поровну и совместимых типов. NULL считаются равными. `ORDER BY` — один, в конце.

## Сортировка и лимиты

`ORDER BY col DESC NULLS LAST` · `LIMIT n OFFSET m` (SQL Server `TOP n`, стандарт `OFFSET ... FETCH NEXT n ROWS ONLY`) ·
`FETCH FIRST n ROWS WITH TIES` — добрать равные.

Без `ORDER BY` порядок не определен. Большой `OFFSET` медленный → keyset: `WHERE (x,id) < (:x,:id)`.

## Окна

```sql
f() OVER (PARTITION BY отдел ORDER BY зп ROWS BETWEEN ... )
```

Не схлопывает строки — дописывает колонку.

| | |
| --- | --- |
| `ROW_NUMBER()` | 1,2,3,4 — всегда уникален |
| `RANK()` | 1,2,2,**4** — с пропуском |
| `DENSE_RANK()` | 1,2,2,**3** — без пропуска, «N-е значение» |
| `NTILE(n)` | делит на n групп |
| `LAG/LEAD(col,n)` | значение назад / вперед |
| `FIRST_VALUE/LAST_VALUE/NTH_VALUE` | из кадра |
| агрегаты `SUM/AVG OVER` | нарастающий итог, среднее по группе |

`ROWS` — физические строки, `RANGE` — по значению (кадр по умолчанию: сложит всю дату разом).
Для нарастающего итога: `ROWS BETWEEN UNBOUNDED PRECEDING AND CURRENT ROW`.
Фильтровать по окну в том же запросе нельзя → CTE.

## PIVOT / UNPIVOT

Строки → колонки и обратно. Синтаксис `PIVOT` есть только в SQL Server/Oracle; переносимо:

```sql
SUM(CASE WHEN d='IT' THEN 1 ELSE 0 END)      -- пивот (PG: COUNT(*) FILTER (WHERE ...))
SELECT id,'q1',q1 FROM t UNION ALL SELECT id,'q2',q2 FROM t   -- анпивот
```

Список колонок всегда задается вручную — динамического пивота в SQL нет.

## Строки и числа

`CONCAT` / `||` (NULL заражает → `CONCAT_WS`) · `SUBSTRING(s FROM 2 FOR 5)` · `REPLACE` ·
`LENGTH` (SQL Server `LEN`) · `UPPER/LOWER/INITCAP` · `TRIM/LTRIM/RTRIM` · `POSITION` · `ROUND(x,n)`.

Позиции с **1**. `TRIM` — только края. Для денег `numeric`, не `float`.

## Даты

`CURRENT_DATE/TIME/TIMESTAMP` (SQL Server `GETDATE()`, MySQL `NOW()`) ·
`EXTRACT(YEAR FROM d)` / `DATEPART` · `date_trunc('month', d)` — для группировки ·
`d + INTERVAL '6 months'` (SQL Server `DATEADD`, MySQL `DATE_SUB`) · `date - date` = дни (`DATEDIFF`).

Храни UTC (`timestamptz`), условие — на голую колонку.

## Ограничения

| | |
| --- | --- |
| `PRIMARY KEY` | уникален + NOT NULL, один на таблицу |
| `FOREIGN KEY` | ссылка существует |
| `UNIQUE` | без повторов (NULL можно) |
| `NOT NULL` / `CHECK` / `DEFAULT` | обязательность / условие / значение по умолчанию |

`ON DELETE`: `NO ACTION`/`RESTRICT` — ошибка · `CASCADE` — удалить детей · `SET NULL` · `SET DEFAULT`.
`CASCADE` удобен и опасен.

## CTE

```sql
WITH x AS (...), y AS (SELECT ... FROM x) SELECT * FROM y WHERE ...;
```

Читается сверху вниз, заменяет вложенные подзапросы, единственный способ отфильтровать окно.
`WITH RECURSIVE` — иерархии и генерация рядов (`generate_series` в PG).

## Топ ловушек

1. `NOT IN` с NULL → `NOT EXISTS`
2. `= NULL` → `IS NULL`
3. `COUNT(*)` после `LEFT JOIN` → `COUNT(правая.col)`
4. `BETWEEN` по датам → полуинтервал
5. функция на колонке в `WHERE` → индекс не работает
6. `RANK` вместо `DENSE_RANK` → дыры в номерах
7. `RANGE` вместо `ROWS` → сложит все строки с равной датой
8. `MAX(дата)` для «последней записи» → дубли, бери `ROW_NUMBER`/`DISTINCT ON`
9. агрегат в `WHERE` → только `HAVING`
10. `UPDATE`/`DELETE` без `WHERE` → вся таблица
