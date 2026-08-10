# SQL: 65 вопросов с ответами

Ответы к [`questions.txt`](questions.txt), нумерация совпадает.

Диалект по умолчанию — **PostgreSQL**; где отличия принципиальны, они помечены
(MySQL 8, SQL Server). Оконные функции работают везде, кроме MySQL < 8.

**Схема, на которой написаны все запросы**

```
employees (id, name, salary, dept_id, hired_at, manager_id)
departments(id, name)
customers  (id, name)
orders     (id, customer_id, product_id, amount, ordered_at)
products   (id, name)
logins     (user_id, login_date)
sales      (sale_date, amount)
```

Пример данных `employees`, к которому удобно возвращаться:

| id | name  | salary | dept_id | manager_id |
| -- | ----- | ------ | ------- | ---------- |
| 1  | Ann   | 300    | 10      | NULL       |
| 2  | Bob   | 200    | 10      | 1          |
| 3  | Cate  | 200    | 20      | 1          |
| 4  | Dan   | 100    | 20      | 3          |
| 5  | Erik  | 100    | NULL    | 3          |

---

### 1. Вторая по величине зарплата

```sql
SELECT DISTINCT salary
FROM   employees
ORDER  BY salary DESC
OFFSET 1 LIMIT 1;                 -- SQL Server: OFFSET 1 ROWS FETCH NEXT 1 ROWS ONLY
```

Классический вариант без окон и OFFSET:

```sql
SELECT MAX(salary) FROM employees
WHERE  salary < (SELECT MAX(salary) FROM employees);
```

Разница в поведении, о ней и спрашивают: если второй зарплаты нет, первый запрос вернет
**ноль строк**, второй — **одну строку с NULL**. И обязательно `DISTINCT`: без него при двух
сотрудниках с максимальной зарплатой «вторая» окажется той же суммой.

Вариант через оконную функцию — здесь параметр вынесен явно и равен `2`, поэтому запрос
без единой правки превращается в «N-ю зарплату» (вопрос 5).

**Что такое `DENSE_RANK()`.** Это оконная функция: она не схлопывает строки, как `GROUP BY`,
а **дописывает к каждой строке еще одну колонку** — ее порядковый номер. Номер выдается
по порядку из `ORDER BY` внутри `OVER (...)`; одинаковым значениям достается **один и тот же
номер**, а следующий номер идет **без пропуска** — отсюда и «dense», «плотный».

На данных из шапки файла:

| name | salary | `DENSE_RANK() OVER (ORDER BY salary DESC)` |
| ---- | ------ | ------------------------------------------ |
| Ann  | 300    | 1                                          |
| Bob  | 200    | 2   ← обе двухсотки получают ранг 2         |
| Cate | 200    | 2                                          |
| Dan  | 100    | 3   ← не 4: пропуска нет                    |
| Erik | 100    | 3                                          |

То есть `DENSE_RANK` нумерует не строки, а **различные значения**: ранг 2 — это буквально
«вторая по величине зарплата», ранг N — N-я. Именно поэтому фильтр `WHERE rnk = 2` и решает
задачу. Если добавить `PARTITION BY dept_id`, нумерация начнется заново в каждом отделе
(вопросы 22 и 23), а сравнение с `RANK` и `ROW_NUMBER` — в вопросе 22.

```sql
WITH ranked AS (
  SELECT e.*, DENSE_RANK() OVER (ORDER BY salary DESC) AS rnk
  FROM   employees e
)
SELECT * FROM ranked WHERE rnk = 2;      -- :n = 2
```

Ранг считается **после** `WHERE`, но **до** `SELECT`, поэтому по нему нельзя фильтровать
в том же запросе — окно и заворачивают в CTE или подзапрос.

Почему именно `DENSE_RANK`: на данных из примера (300, 200, 200, 100, 100) он даст
«вторую зарплату» = 200 и вернет обоих сотрудников с ней. `ROW_NUMBER` выдал бы одну
случайную строку из двух с одинаковой суммой, а `RANK` после пары двухсоток перескочил бы
с ранга 2 сразу на 4 — и запрос с `rnk = 3` не нашел бы ничего.

Есть и оконная функция, которая принимает номер прямо аргументом:

```sql
SELECT DISTINCT NTH_VALUE(salary, 2) OVER (
         ORDER BY salary DESC
         ROWS BETWEEN UNBOUNDED PRECEDING AND UNBOUNDED FOLLOWING
       ) AS second_salary
FROM   employees;
```

Две тонкости, из-за которых ее редко берут: без явного `ROWS BETWEEN ... UNBOUNDED FOLLOWING`
кадр по умолчанию заканчивается на текущей строке, и функция вернет NULL для первой же строки;
и считает она **строки, а не разные значения** — при двух сотрудниках с максимальной зарплатой
«второй» окажется та же сумма. Поэтому на собеседовании безопаснее `DENSE_RANK`.

### 2. Найти дубликаты

```sql
SELECT name, salary, COUNT(*) AS cnt
FROM   employees
GROUP  BY name, salary
HAVING COUNT(*) > 1;
```

`GROUP BY` — по тем колонкам, которые и определяют дубликат. Если нужны сами строки,
а не сводка, — оконная функция:

```sql
SELECT * FROM (
  SELECT e.*, COUNT(*) OVER (PARTITION BY name, salary) AS cnt FROM employees e
) t WHERE cnt > 1;
```

### 3. Удалить дубликаты, оставив по одной строке

```sql
WITH ranked AS (
  SELECT id, ROW_NUMBER() OVER (PARTITION BY name, salary ORDER BY id) AS rn
  FROM   employees
)
DELETE FROM employees
WHERE  id IN (SELECT id FROM ranked WHERE rn > 1);
```

`ORDER BY id` решает, какая строка выживет (обычно самая старая). Если первичного ключа нет,
в PostgreSQL спасает системный `ctid`:

```sql
DELETE FROM employees a USING employees b
WHERE  a.ctid > b.ctid AND a.name = b.name AND a.salary = b.salary;
```

### 4. Топ-5 зарплат

```sql
SELECT * FROM employees ORDER BY salary DESC LIMIT 5;
```

Если пятое и шестое место с одинаковой суммой, а обрезать «по живому» нельзя:

```sql
SELECT * FROM employees ORDER BY salary DESC FETCH FIRST 5 ROWS WITH TIES;
```

MySQL: `LIMIT 5`, SQL Server: `SELECT TOP 5 ... WITH TIES`.

### 5. N-я по величине зарплата

```sql
SELECT DISTINCT salary FROM employees
ORDER BY salary DESC OFFSET :n - 1 LIMIT 1;
```

Через окно — сразу со строками сотрудников:

```sql
SELECT * FROM (
  SELECT e.*, DENSE_RANK() OVER (ORDER BY salary DESC) AS rnk FROM employees e
) t WHERE rnk = :n;
```

Именно `DENSE_RANK`: `RANK` пропускает номера после дубликатов, и «третьей зарплаты»
может просто не оказаться.

### 6. Отделы и суммарный фонд зарплат

```sql
SELECT d.name, COALESCE(SUM(e.salary), 0) AS total
FROM   departments d
LEFT   JOIN employees e ON e.dept_id = d.id
GROUP  BY d.name
ORDER  BY total DESC;
```

`LEFT JOIN` + `COALESCE` — чтобы пустые отделы попали в отчет с нулем, а не исчезли.
С обычным `JOIN` их не будет.

### 7. Кто пришел за последние 6 месяцев

```sql
SELECT * FROM employees
WHERE  hired_at >= CURRENT_DATE - INTERVAL '6 months';
```

Условие пишется на **колонку без обертки**: `WHERE hired_at >= ...`, а не
`WHERE EXTRACT(MONTH FROM hired_at) ...` — иначе индекс по `hired_at` не будет использован.
MySQL: `DATE_SUB(CURDATE(), INTERVAL 6 MONTH)`, SQL Server: `DATEADD(month, -6, GETDATE())`.

### 8. Сотрудники без отдела

```sql
SELECT * FROM employees WHERE dept_id IS NULL;
```

Если «без отдела» = ссылка на несуществующий отдел:

```sql
SELECT e.* FROM employees e
WHERE  NOT EXISTS (SELECT 1 FROM departments d WHERE d.id = e.dept_id);
```

**Ловушка `NOT IN`**: `WHERE dept_id NOT IN (SELECT id FROM departments)` вернет пусто,
если в подзапросе встретится хоть один NULL. `NOT EXISTS` от этого свободен.

### 9. Количество сотрудников по отделам

```sql
SELECT d.name, COUNT(e.id) AS employees
FROM   departments d
LEFT   JOIN employees e ON e.dept_id = d.id
GROUP  BY d.name;
```

`COUNT(e.id)`, а не `COUNT(*)`: при `LEFT JOIN` пустой отдел дает одну строку с NULL,
и `COUNT(*)` посчитал бы ее за единицу.

### 10. Пивот: сотрудники по отделам в столбцах

Из «длинной» таблицы делаем «широкую»:

| было           |     | стало |     |     |
| -------------- | --- | ----- | --- | --- |
| dept \| cnt    |     | IT    | HR  | Fin |
| IT   \| 3      | →   | 3     | 2   | 1   |
| HR   \| 2      |     |       |     |     |
| Fin  \| 1      |     |       |     |     |

```sql
SELECT COUNT(*) FILTER (WHERE d.name = 'IT')  AS it,
       COUNT(*) FILTER (WHERE d.name = 'HR')  AS hr,
       COUNT(*) FILTER (WHERE d.name = 'Fin') AS fin
FROM   employees e JOIN departments d ON d.id = e.dept_id;
```

Переносимый вариант (работает везде) — агрегат по `CASE`:

```sql
SELECT SUM(CASE WHEN d.name = 'IT' THEN 1 ELSE 0 END) AS it, ...
```

Главное ограничение: **список колонок задается в запросе руками**. Динамический пивот
в чистом SQL невозможен — либо генерируют текст запроса, либо разворачивают в приложении.

### 11. Анпивот: столбцы в строки

| было                       |     | стало                    |
| -------------------------- | --- | ------------------------ |
| id \| q1 \| q2 \| q3       |     | id \| quarter \| amount  |
| 1  \| 10 \| 20 \| 30       | →   | 1  \| q1 \| 10           |
|                            |     | 1  \| q2 \| 20           |
|                            |     | 1  \| q3 \| 30           |

```sql
SELECT t.id, v.quarter, v.amount
FROM   revenue t
CROSS  JOIN LATERAL (VALUES ('q1', t.q1), ('q2', t.q2), ('q3', t.q3)) AS v(quarter, amount);
```

Универсально — через `UNION ALL`:

```sql
SELECT id, 'q1' AS quarter, q1 AS amount FROM revenue
UNION ALL SELECT id, 'q2', q2 FROM revenue
UNION ALL SELECT id, 'q3', q3 FROM revenue;
```

### 12. Рост выручки месяц к месяцу

```sql
WITH monthly AS (
  SELECT date_trunc('month', ordered_at) AS month, SUM(amount) AS revenue
  FROM   orders GROUP BY 1
)
SELECT month, revenue,
       revenue - LAG(revenue) OVER (ORDER BY month) AS diff,
       ROUND(100.0 * (revenue - LAG(revenue) OVER (ORDER BY month))
             / NULLIF(LAG(revenue) OVER (ORDER BY month), 0), 2) AS growth_pct
FROM   monthly ORDER BY month;
```

Сначала агрегируем по месяцам, потом `LAG` берет предыдущую строку. `NULLIF(..., 0)`
защищает от деления на ноль. У первого месяца рост будет NULL — так и должно быть.
Если в каких-то месяцах заказов не было, их строк не появится: нужен календарь
(см. вопрос 50) и `LEFT JOIN` к нему.

### 13. Нарастающий итог продаж по датам

```sql
SELECT sale_date, amount,
       SUM(amount) OVER (ORDER BY sale_date
                         ROWS BETWEEN UNBOUNDED PRECEDING AND CURRENT ROW) AS running_total
FROM   sales;
```

`ROWS` против `RANGE` — любимый вопрос: при нескольких строках с одной датой `RANGE`
(поведение по умолчанию) сложит **все строки этой даты сразу**, а `ROWS` — строго
по одной. Для честного нарастающего итога по дням лучше сначала свернуть по дате.

### 14. Пропуски в последовательности чисел

```sql
SELECT id + 1 AS gap_start, next_id - 1 AS gap_end
FROM (SELECT id, LEAD(id) OVER (ORDER BY id) AS next_id FROM t) s
WHERE next_id - id > 1;
```

Смотрим на разрыв между строкой и следующей. Альтернатива — сравнить с эталонным рядом:

```sql
SELECT g AS missing_id
FROM   generate_series((SELECT MIN(id) FROM t), (SELECT MAX(id) FROM t)) g
WHERE  NOT EXISTS (SELECT 1 FROM t WHERE t.id = g);
```

Первый вариант отдает диапазоны и дешев, второй — каждое отсутствующее значение.

### 15. Самая длинная серия подряд идущих дней логина

Классические **gaps and islands**. Трюк: у подряд идущих дат разность
«дата минус её порядковый номер» постоянна, и по ней можно группировать.

| login_date | row_number | date - rn  | остров |
| ---------- | ---------- | ---------- | ------ |
| 2026-01-01 | 1          | 2025-12-31 | A      |
| 2026-01-02 | 2          | 2025-12-31 | A      |
| 2026-01-03 | 3          | 2025-12-31 | A      |
| 2026-01-07 | 4          | 2026-01-03 | B      |
| 2026-01-08 | 5          | 2026-01-03 | B      |

```sql
WITH d AS (
  SELECT DISTINCT user_id, login_date FROM logins          -- дубли за день убираем
), grouped AS (
  SELECT user_id, login_date,
         login_date - (ROW_NUMBER() OVER (PARTITION BY user_id ORDER BY login_date))::int AS grp
  FROM   d
)
SELECT user_id, COUNT(*) AS streak, MIN(login_date) AS from_date, MAX(login_date) AS to_date
FROM   grouped
GROUP  BY user_id, grp
ORDER  BY streak DESC;
```

Для «самой длинной серии по каждому пользователю» оборачиваем в `DISTINCT ON (user_id)`
с сортировкой по `streak DESC`.

### 16. Дата первого логина каждого пользователя

```sql
SELECT user_id, MIN(login_date) AS first_login FROM logins GROUP BY user_id;
```

### 17. Последняя покупка каждого клиента

```sql
SELECT DISTINCT ON (customer_id) *
FROM   orders
ORDER  BY customer_id, ordered_at DESC;      -- PostgreSQL
```

Переносимо и с гарантией одной строки при одинаковых датах:

```sql
SELECT * FROM (
  SELECT o.*, ROW_NUMBER() OVER (PARTITION BY customer_id ORDER BY ordered_at DESC, id DESC) AS rn
  FROM   orders o
) t WHERE rn = 1;
```

Вариант «через `MAX` в подзапросе» вернет **две строки**, если у клиента две покупки
в одну и ту же секунду — на собеседовании это и проверяют.

### 18. Клиенты, которые покупали во все месяцы

Реляционное деление: у клиента столько же различных месяцев, сколько их всего.

```sql
SELECT customer_id
FROM   orders
GROUP  BY customer_id
HAVING COUNT(DISTINCT date_trunc('month', ordered_at))
     = (SELECT COUNT(DISTINCT date_trunc('month', ordered_at)) FROM orders);
```

Если «все месяцы» — это календарь, а не только месяцы с заказами, эталон берут
из `generate_series` (вопрос 50), иначе месяц без единого заказа выпадет из знаменателя.

### 19. Товары, которые ни разу не заказывали

```sql
SELECT p.* FROM products p
WHERE  NOT EXISTS (SELECT 1 FROM orders o WHERE o.product_id = p.id);
```

Эквивалент через соединение — `LEFT JOIN ... WHERE o.id IS NULL` (антиджойн).
`NOT IN` тут снова опасен из-за NULL в `orders.product_id`.

### 20. Кто получает больше средней зарплаты

```sql
SELECT * FROM employees
WHERE  salary > (SELECT AVG(salary) FROM employees);
```

Подзапрос без корреляции выполняется один раз. `AVG` игнорирует NULL — если это не то,
что нужно, используйте `AVG(COALESCE(salary, 0))`.

### 21. Кто получает больше среднего по своему отделу

```sql
SELECT * FROM (
  SELECT e.*, AVG(salary) OVER (PARTITION BY dept_id) AS dept_avg
  FROM   employees e
) t WHERE salary > dept_avg;
```

Окно считает среднее по отделу и **оставляет все строки** — один проход вместо
коррелированного подзапроса на каждую строку.

### 22. Ранжирование сотрудников внутри отдела по зарплате

```sql
SELECT name, dept_id, salary,
       ROW_NUMBER() OVER (PARTITION BY dept_id ORDER BY salary DESC) AS row_num,
       RANK()       OVER (PARTITION BY dept_id ORDER BY salary DESC) AS rnk,
       DENSE_RANK() OVER (PARTITION BY dept_id ORDER BY salary DESC) AS dense_rnk
FROM   employees;
```

Разница между тремя функциями — обязательный вопрос:

| salary | ROW_NUMBER | RANK | DENSE_RANK |
| ------ | ---------- | ---- | ---------- |
| 300    | 1          | 1    | 1          |
| 200    | 2          | 2    | 2          |
| 200    | 3          | 2    | 2          |
| 100    | 4          | **4**| **3**      |

`ROW_NUMBER` всегда уникален, `RANK` пропускает номера после ничьей, `DENSE_RANK` — нет.

### 23. Топ-3 зарплаты в каждом отделе

```sql
SELECT * FROM (
  SELECT e.*, DENSE_RANK() OVER (PARTITION BY dept_id ORDER BY salary DESC) AS rnk
  FROM   employees e
) t WHERE rnk <= 3;
```

`DENSE_RANK` даст «три уровня зарплат» (строк может быть больше трех при ничьих),
`ROW_NUMBER` — ровно три строки. Уточните у интервьюера, что именно нужно.

### 24. Общие записи двух таблиц

```sql
SELECT id, name FROM a
INTERSECT
SELECT id, name FROM b;
```

`INTERSECT` сравнивает строки целиком, убирает дубликаты и **считает NULL равными**
(в отличие от `=`). Если нужны все колонки одной из таблиц — обычный `JOIN` по ключу.

### 25. Записи, которые есть в одной таблице и нет в другой

```sql
SELECT id, name FROM a
EXCEPT                       -- MySQL 8: EXCEPT; Oracle: MINUS
SELECT id, name FROM b;
```

По ключу и с колонками первой таблицы:

```sql
SELECT a.* FROM a WHERE NOT EXISTS (SELECT 1 FROM b WHERE b.id = a.id);
```

### 26. Разделить полное имя на имя и фамилию

```sql
SELECT split_part(full_name, ' ', 1) AS first_name,
       split_part(full_name, ' ', 2) AS last_name
FROM   people;
```

Переносимый вариант через позицию пробела:

```sql
SELECT SUBSTRING(full_name FROM 1 FOR POSITION(' ' IN full_name) - 1) AS first_name,
       SUBSTRING(full_name FROM POSITION(' ' IN full_name) + 1)       AS last_name
FROM   people;
```

Стоит проговорить: на именах без пробела `POSITION` вернет 0 и выражение сломается,
а «Анна Мария Петрова» разложится неверно. В проде такие данные хранят разными колонками.

### 27. Склеить имя и фамилию

```sql
SELECT first_name || ' ' || last_name AS full_name FROM people;   -- стандарт, PostgreSQL
SELECT CONCAT_WS(' ', first_name, last_name) AS full_name FROM people;
```

Разница важна: `||` и `CONCAT` в PostgreSQL дают **NULL, если хоть один аргумент NULL**,
а `CONCAT_WS` NULL-и просто пропускает. В MySQL так же ведет себя `CONCAT`, в SQL Server —
`CONCAT` игнорирует NULL, а `+` нет.

### 28. Год, месяц и день из даты

```sql
SELECT EXTRACT(YEAR FROM hired_at)  AS y,
       EXTRACT(MONTH FROM hired_at) AS m,
       EXTRACT(DAY FROM hired_at)   AS d
FROM   employees;
```

MySQL: `YEAR()`, `MONTH()`, `DAY()`. SQL Server: `DATEPART(year, ...)`.
Для группировки по месяцу удобнее не разбирать на части, а свернуть:
`date_trunc('month', hired_at)` — сортировка останется хронологической.

### 29. Возраст по дате рождения

```sql
SELECT name, EXTRACT(YEAR FROM AGE(birth_date)) AS age FROM people;   -- PostgreSQL
```

Переносимо и корректно относительно дня рождения:

```sql
SELECT FLOOR(EXTRACT(DAY FROM (CURRENT_DATE - birth_date)) / 365.25) AS age FROM people;
```

Наивное `YEAR(now) - YEAR(birth)` ошибается на год у всех, кто еще не отпраздновал
день рождения в этом году, — на это и ловят.

### 30. Число рабочих дней в месяце

```sql
SELECT COUNT(*) AS working_days
FROM   generate_series(date_trunc('month', :d),
                       date_trunc('month', :d) + INTERVAL '1 month - 1 day',
                       INTERVAL '1 day') AS g(day)
WHERE  EXTRACT(ISODOW FROM g.day) < 6;        -- 6,7 = суббота и воскресенье
```

`ISODOW` удобнее `DOW`: неделя начинается с понедельника (1) и не надо помнить,
что воскресенье — это 0. Праздники так не учесть: под них заводят таблицу `holidays`
и добавляют `AND g.day NOT IN (SELECT day FROM holidays)`.

### 31. Сотрудники с одинаковой зарплатой

```sql
SELECT salary, COUNT(*) AS people, string_agg(name, ', ') AS who
FROM   employees
GROUP  BY salary
HAVING COUNT(*) > 1;
```

Если нужны пары сотрудников — самосоединение с защитой от зеркальных дублей:

```sql
SELECT a.name, b.name, a.salary
FROM   employees a JOIN employees b ON a.salary = b.salary AND a.id < b.id;
```

Условие `a.id < b.id` (а не `<>`) убирает и пару «сам с собой», и дубль «Bob–Cate / Cate–Bob».

### 32. Самый высокооплачиваемый сотрудник в каждом отделе

```sql
SELECT DISTINCT ON (dept_id) *
FROM   employees
ORDER  BY dept_id, salary DESC;               -- PostgreSQL, одна строка на отдел
```

Если при равных зарплатах нужны все:

```sql
SELECT * FROM (
  SELECT e.*, RANK() OVER (PARTITION BY dept_id ORDER BY salary DESC) AS rnk
  FROM   employees e
) t WHERE rnk = 1;
```

### 33. Руководители и их подчиненные

```sql
SELECT m.name AS manager, e.name AS employee
FROM   employees e
JOIN   employees m ON m.id = e.manager_id
ORDER  BY manager, employee;
```

Самосоединение по `manager_id`. `LEFT JOIN` вместо `JOIN` покажет и тех, у кого
руководителя нет. Для всей иерархии вглубь — рекурсивный CTE:

```sql
WITH RECURSIVE tree AS (
  SELECT id, name, manager_id, 1 AS level FROM employees WHERE manager_id IS NULL
  UNION ALL
  SELECT e.id, e.name, e.manager_id, t.level + 1
  FROM   employees e JOIN tree t ON e.manager_id = t.id
)
SELECT * FROM tree ORDER BY level;
```

### 34. Развернуть строку

```sql
SELECT REVERSE(name) FROM employees;
```

Есть в PostgreSQL, MySQL и SQL Server. Осторожно с многобайтовыми символами:
разворачивается по символам, но эмодзи из нескольких кодовых точек может распасться.

### 35. Частота слов в текстовой колонке

```sql
SELECT word, COUNT(*) AS cnt
FROM   articles,
       LATERAL regexp_split_to_table(lower(body), '\s+') AS word
WHERE  word <> ''
GROUP  BY word
ORDER  BY cnt DESC;
```

`lower` — чтобы «Go» и «go» не считались разными; для чистоты стоит убрать пунктуацию
`regexp_replace(body, '[^\w\s]', '', 'g')`. MySQL 8: рекурсивный CTE или `JSON_TABLE`.

### 36. Строки, где колонка равна NULL

```sql
SELECT * FROM employees WHERE dept_id IS NULL;
```

`= NULL` не работает никогда: сравнение с NULL дает не «истину» или «ложь», а **UNKNOWN**,
и строка не попадает в результат. Отсюда же `IS DISTINCT FROM` для сравнений,
где NULL должен считаться значением.

### 37. Заменить NULL значением по умолчанию

```sql
SELECT name, COALESCE(salary, 0) AS salary FROM employees;
```

`COALESCE` — стандарт и берет первый не-NULL из списка. Диалектные аналоги:
`IFNULL` (MySQL), `ISNULL` (SQL Server), `NVL` (Oracle). Для постоянной замены —
`ALTER TABLE ... SET DEFAULT` плюс `UPDATE` уже записанных строк.

### 38. Убрать пробелы по краям строки

```sql
SELECT TRIM(name) FROM employees;             -- по обоим краям
SELECT LTRIM(name), RTRIM(name) FROM employees;
```

`TRIM(BOTH ' ' FROM name)` — полная форма, можно обрезать любой символ:
`TRIM(BOTH '0' FROM code)`. Пробелы **внутри** строки `TRIM` не трогает,
для них нужен `REPLACE` или `regexp_replace`.

### 39. Самый часто заказываемый товар

```sql
SELECT p.name, COUNT(*) AS orders_cnt
FROM   orders o JOIN products p ON p.id = o.product_id
GROUP  BY p.name
ORDER  BY orders_cnt DESC
LIMIT  1;
```

Если товаров-лидеров может быть несколько — `FETCH FIRST 1 ROW WITH TIES`,
иначе один из них молча потеряется. «Часто» может означать и число заказов,
и суммарное количество штук (`SUM(quantity)`) — уточняющий вопрос.

### 40. Рост год к году

```sql
WITH yearly AS (
  SELECT date_trunc('year', ordered_at) AS year, SUM(amount) AS revenue
  FROM   orders GROUP BY 1
)
SELECT year, revenue,
       LAG(revenue) OVER (ORDER BY year) AS prev_year,
       ROUND(100.0 * (revenue - LAG(revenue) OVER (ORDER BY year))
             / NULLIF(LAG(revenue) OVER (ORDER BY year), 0), 2) AS yoy_pct
FROM   yearly ORDER BY year;
```

То же, что вопрос 12, но с шагом в год. Если нужно сравнивать месяц с тем же месяцем
прошлого года, берут `LAG(revenue, 12) OVER (ORDER BY month)` — и только при условии,
что в ряду нет пропущенных месяцев, иначе смещение уедет.

### 41. Скользящее среднее продаж (окно 3 дня)

```sql
SELECT sale_date, amount,
       AVG(amount) OVER (ORDER BY sale_date
                         ROWS BETWEEN 2 PRECEDING AND CURRENT ROW) AS moving_avg_3
FROM   sales;
```

`ROWS` считает **три строки**, а не три календарных дня: если в данных есть пропуски,
среднее захватит более старые даты. Честное «за трое суток» — оконный кадр по значению:

```sql
       AVG(amount) OVER (ORDER BY sale_date
                         RANGE BETWEEN INTERVAL '2 days' PRECEDING AND CURRENT ROW)
```

### 42. Разница между двумя датами в днях

```sql
SELECT (end_date - start_date) AS days FROM periods;              -- PostgreSQL, date - date = int
SELECT DATEDIFF(end_date, start_date) FROM periods;               -- MySQL
SELECT DATEDIFF(day, start_date, end_date) FROM periods;          -- SQL Server
```

Если колонки `timestamp`, вычитание даст `interval`; для целых дней:
`EXTRACT(DAY FROM (end_ts - start_ts))` либо приведение `end_ts::date - start_ts::date`.

### 43. Пересекающиеся диапазоны дат

Два отрезка пересекаются, если каждый начинается не позже, чем заканчивается другой:

```
   A: ├───────────┤
   B:        ├─────────┤        пересечение есть: A.start <= B.end AND B.start <= A.end
   C:                    ├───┤  с A не пересекается
```

```sql
SELECT a.id, b.id
FROM   bookings a JOIN bookings b
  ON   a.id < b.id                      -- без зеркальных пар
 AND   a.start_date <= b.end_date
 AND   b.start_date <= a.end_date;
```

В PostgreSQL то же самое короче и с поддержкой индекса GiST:

```sql
WHERE daterange(a.start_date, a.end_date, '[]') && daterange(b.start_date, b.end_date, '[]')
```

Стандартный оператор `(a.start, a.end) OVERLAPS (b.start, b.end)` тоже есть, но у него
границы полуоткрытые — смежные интервалы пересечением не считаются.

### 44. Вторая по частоте зарплата

```sql
SELECT salary, COUNT(*) AS cnt
FROM   employees
GROUP  BY salary
ORDER  BY cnt DESC
OFFSET 1 LIMIT 1;
```

Если несколько зарплат встречаются одинаково часто и нужна «вторая по частоте» как
уровень, а не строка, — ранжируем частоты:

```sql
SELECT salary, cnt FROM (
  SELECT salary, COUNT(*) AS cnt, DENSE_RANK() OVER (ORDER BY COUNT(*) DESC) AS rnk
  FROM   employees GROUP BY salary
) t WHERE rnk = 2;
```

### 45. Зарплата от 10k до 20k

```sql
SELECT * FROM employees WHERE salary BETWEEN 10000 AND 20000;
```

`BETWEEN` включает **обе** границы. С датами это регулярно приводит к ошибке:
`ordered_at BETWEEN '2026-01-01' AND '2026-01-31'` отрежет все, что произошло
31 января после полуночи, потому что дата приводится к `00:00:00`.
Для временных диапазонов правильнее полуинтервал:
`ordered_at >= '2026-01-01' AND ordered_at < '2026-02-01'`.

### 46. Удалить записи с дублирующимся значением одной колонки

```sql
WITH ranked AS (
  SELECT id, ROW_NUMBER() OVER (PARTITION BY email ORDER BY id) AS rn FROM users
)
DELETE FROM users WHERE id IN (SELECT id FROM ranked WHERE rn > 1);
```

Оставляем строку с наименьшим `id`. Чтобы дубли не появились снова, после чистки
вешают ограничение: `CREATE UNIQUE INDEX ON users (email);` — иначе задача вернется.

### 47. Вставить только уникальные записи из другой таблицы

```sql
INSERT INTO target (id, name)
SELECT DISTINCT s.id, s.name
FROM   source s
WHERE  NOT EXISTS (SELECT 1 FROM target t WHERE t.id = s.id);
```

`DISTINCT` убирает дубли внутри источника, `NOT EXISTS` — уже существующие в приемнике.
При наличии уникального индекса надежнее переложить работу на СУБД:

```sql
INSERT INTO target (id, name) SELECT id, name FROM source
ON CONFLICT (id) DO NOTHING;              -- MySQL: INSERT IGNORE
```

Второй вариант не ловит гонку между `SELECT` и `INSERT` в конкурентной среде — а первый ловит.

### 48. Поднять зарплату на 10% в конкретном отделе

```sql
UPDATE employees
SET    salary = salary * 1.10
WHERE  dept_id = (SELECT id FROM departments WHERE name = 'IT');
```

Три вещи, которые проверяют: есть ли `WHERE` вообще (без него обновится вся таблица),
округление (`ROUND(salary * 1.10, 2)` для денег) и то, что операция **не идемпотентна** —
повторный запуск даст +21%. В проде такие `UPDATE` выполняют в транзакции,
предварительно посмотрев тем же `WHERE`, сколько строк затронется.

### 49. Имя в верхнем и нижнем регистре

```sql
SELECT UPPER(name) AS upper_name, LOWER(name) AS lower_name FROM employees;
```

`INITCAP(name)` (PostgreSQL, Oracle) даст «Иван Петров» с заглавных.
Для регистронезависимого поиска не пишите `WHERE UPPER(name) = 'ANN'` — это убивает индекс;
лучше `ILIKE`, либо функциональный индекс `CREATE INDEX ON employees (lower(name))`.

### 50. Сгенерировать ряд дат между двумя датами

```sql
SELECT g::date AS day
FROM   generate_series('2026-01-01'::date, '2026-01-31'::date, INTERVAL '1 day') g;
```

Там, где `generate_series` нет (MySQL 8, SQL Server), — рекурсивный CTE:

```sql
WITH RECURSIVE dates AS (
  SELECT DATE '2026-01-01' AS day
  UNION ALL
  SELECT day + INTERVAL 1 DAY FROM dates WHERE day < DATE '2026-01-31'
)
SELECT * FROM dates;
```

Зачем это нужно: календарь — основа отчетов «без дыр». Агрегат по `orders` покажет
только дни с заказами; чтобы в графике были нули, делают `LEFT JOIN` календаря к данным.

### 51. Все города, в которых есть клиенты

```sql
SELECT DISTINCT city FROM customers ORDER BY city;
```

`DISTINCT` и `GROUP BY city` дадут один и тот же результат, но `GROUP BY` нужен, если рядом
нужен агрегат (`COUNT(*)` по городу). NULL попадет в результат одной строкой: `DISTINCT`
считает все NULL одинаковыми.

### 52. Доля каждого клиента в общей выручке

```sql
SELECT customer_id,
       SUM(amount) AS revenue,
       ROUND(100.0 * SUM(amount) / SUM(SUM(amount)) OVER (), 2) AS pct_of_total
FROM   orders
GROUP  BY customer_id
ORDER  BY pct_of_total DESC;
```

Ключевой трюк — `SUM(SUM(amount)) OVER ()`: внешний `SUM` здесь **оконный** и считается уже
по сгруппированным строкам, потому что окна выполняются **после** `GROUP BY`.
Пустое `OVER ()` означает «по всему результату», поэтому получается общий итог рядом с каждой строкой.
Нужна доля не в общем итоге, а **внутри группы** (товар в своей категории) — добавляется секция:
`SUM(amount) OVER (PARTITION BY category_id)`.

Без окна пришлось бы делать второй проход по таблице подзапросом
`(SELECT SUM(amount) FROM orders)` — работает, но читает данные дважды.
`100.0`, а не `100` — иначе в целочисленной арифметике доля превратится в ноль.

### 53. Разбить студентов на квартили по баллам

```sql
SELECT name, score,
       NTILE(4) OVER (ORDER BY score DESC) AS quartile
FROM   students;
```

`NTILE(n)` делит **строки** на `n` максимально равных групп: при 10 строках и `NTILE(4)`
получится 3, 3, 2, 2. Отсюда главная особенность: одинаковые баллы могут попасть
в разные квартили — граница проходит по счету строк, а не по значению.

Если нужна честная граница по значению, берут `PERCENT_RANK()` или `CUME_DIST()`,
а для «сколько процентов набрали меньше» — `PERCENTILE_CONT`.

### 54. Счета, у которых баланс уходил в минус

```sql
WITH running AS (
  SELECT account_id, transaction_date,
         SUM(amount) OVER (PARTITION BY account_id ORDER BY transaction_date
                           ROWS BETWEEN UNBOUNDED PRECEDING AND CURRENT ROW) AS balance
  FROM   transactions
)
SELECT DISTINCT account_id FROM running WHERE balance < 0;
```

Нарастающий итог **внутри каждого счета** (`PARTITION BY account_id`), затем фильтр по нему —
и обязательно через CTE, потому что по результату окна нельзя фильтровать в том же запросе.

`ROWS` вместо кадра по умолчанию здесь принципиален: при нескольких операциях в один день
`RANGE` сложил бы весь день разом и «мгновенный» уход в минус внутри дня потерялся бы.

Чтобы увидеть не только факт, но и когда именно, замените `DISTINCT` на
`SELECT account_id, MIN(transaction_date) FROM running WHERE balance < 0 GROUP BY account_id`.

### 55. Вся цепочка подчинения от директора

```sql
WITH RECURSIVE org AS (
  SELECT employee_id, name, manager_id, 1 AS level, name::text AS chain
  FROM   employees
  WHERE  manager_id IS NULL                       -- якорь: с кого начинаем
  UNION ALL
  SELECT e.employee_id, e.name, e.manager_id, o.level + 1, o.chain || ' → ' || e.name
  FROM   employees e
  JOIN   org o ON e.manager_id = o.employee_id    -- рекурсивная часть
)
SELECT employee_id, name, level, chain FROM org ORDER BY level, name;
```

Рекурсивный CTE состоит из двух частей, соединенных `UNION ALL`: **якорь** (верхушка иерархии)
и **шаг**, который присоединяет подчиненных к уже найденным. Выполнение останавливается,
когда очередной шаг не дал новых строк.

Что стоит проговорить: `UNION ALL` быстрее, но `UNION` защитит от зацикливания, если в данных
окажется цикл (сотрудник — сам себе начальник через цепочку); в PostgreSQL для этого есть
`CYCLE employee_id SET is_cycle USING path`. Поле `level` дает глубину, `chain` — читаемый путь.
Простое соединение «начальник — подчиненный» на один уровень — это вопрос 33.

### 56. Сколько дней между первым и вторым событием пользователя

```sql
WITH ranked AS (
  SELECT user_id, event_date,
         ROW_NUMBER() OVER (PARTITION BY user_id ORDER BY event_date) AS rn
  FROM   events
)
SELECT r1.user_id,
       r2.event_date - r1.event_date AS days_between   -- NULL, если второго события нет
FROM   ranked r1
LEFT   JOIN ranked r2 ON r2.user_id = r1.user_id AND r2.rn = 2
WHERE  r1.rn = 1;
```

Именно `LEFT JOIN` дает требуемый NULL для пользователей с одним событием — с обычным `JOIN`
такие строки просто исчезли бы.

Короче то же самое одним окном:

```sql
SELECT DISTINCT ON (user_id) user_id,
       LEAD(event_date) OVER (PARTITION BY user_id ORDER BY event_date) - event_date AS days_between
FROM   events
ORDER  BY user_id, event_date;
```

`LEAD` сам вернет NULL, если следующей строки нет.

### 57. Stickiness: отношение DAU к MAU по месяцам

```sql
WITH dau AS (                                    -- уникальные пользователи по дням
  SELECT date_trunc('month', activity_date) AS mo,
         activity_date,
         COUNT(DISTINCT user_id) AS d
  FROM   activity GROUP BY 1, 2
), mau AS (                                      -- уникальные пользователи за месяц
  SELECT date_trunc('month', activity_date) AS mo,
         COUNT(DISTINCT user_id) AS m
  FROM   activity GROUP BY 1
)
SELECT d.mo,
       ROUND(AVG(d.d) / MAX(m.m)::numeric, 3) AS stickiness
FROM   dau d JOIN mau m ON m.mo = d.mo
GROUP  BY d.mo
ORDER  BY d.mo;
```

Смысл метрики: какая доля месячной аудитории заходит в среднестатистический день.
0.2 значит «пользователь заходит примерно 6 дней из 30».

Две ошибки, которые тут делают: **суммировать DAU за месяц** вместо усреднения (один и тот же
человек посчитается многократно) и брать `COUNT(user_id)` вместо `COUNT(DISTINCT user_id)`,
если в таблице больше одной строки на пользователя за день.

### 58. «Друзья друзей» с ранжированием по числу общих друзей

```sql
SELECT f2.friend_id AS suggested,
       COUNT(*)     AS mutual_friends
FROM   friends f1
JOIN   friends f2 ON f2.user_id = f1.friend_id     -- шаг на второй круг знакомств
WHERE  f1.user_id = 123
  AND  f2.friend_id <> 123                          -- не предлагаем самого себя
  AND  NOT EXISTS (                                 -- и уже существующих друзей
         SELECT 1 FROM friends f
         WHERE  f.user_id = 123 AND f.friend_id = f2.friend_id)
GROUP  BY f2.friend_id
ORDER  BY mutual_friends DESC;
```

Два соединения одной таблицы = два шага по графу: `f1` дает друзей пользователя, `f2` —
друзей его друзей. `COUNT(*)` считает, через скольких общих знакомых мы пришли к кандидату, —
это и есть сила рекомендации.

Схема предполагает, что дружба хранится **в обе стороны** (`(1,2)` и `(2,1)`). Если она
записана один раз, перед соединением делают симметричное представление через `UNION ALL`
с перевернутыми колонками. `NOT EXISTS`, а не `NOT IN` — на случай NULL в данных.

### 59. Классифицировать клиентов: New / Returning / Churned

```sql
WITH stats AS (
  SELECT customer_id,
         MIN(order_date) AS first_order,
         MAX(order_date) AS last_order
  FROM   orders GROUP BY customer_id
)
SELECT customer_id, first_order, last_order,
       CASE
         WHEN first_order >= CURRENT_DATE - INTERVAL '30 days' THEN 'New'
         WHEN last_order  <  CURRENT_DATE - INTERVAL '90 days' THEN 'Churned'
         ELSE 'Returning'
       END AS status
FROM   stats;
```

Вся задача сводится к двум агрегатам (`MIN`/`MAX` даты заказа) и `CASE` поверх них.

Главное, что проверяют: **порядок веток `CASE`**. Условия пересекаются — клиент, сделавший
первый заказ вчера, формально попадает и под «New», и под «Returning». Срабатывает **первая
подошедшая** ветка, поэтому определения нужно расставить по приоритету и проговорить это вслух.
`ELSE` обязателен, иначе непопавшие строки молча получат NULL.

### 60. Кто пришел в тот же месяц и год, что и его руководитель

```sql
SELECT e.name AS employee, m.name AS manager, e.hired_at, m.hired_at AS manager_hired_at
FROM   employees e
JOIN   employees m ON m.id = e.manager_id
WHERE  date_trunc('month', e.hired_at) = date_trunc('month', m.hired_at);
```

Самосоединение по `manager_id` (как в вопросе 33) плюс сравнение дат **с точностью до месяца**.

Главная ловушка формулировки: сравнивать только месяц нельзя —
`EXTRACT(MONTH FROM e.hired_at) = EXTRACT(MONTH FROM m.hired_at)` соберет пары из разных лет.
Нужен либо `date_trunc('month', ...)`, либо явная проверка года и месяца:

```sql
WHERE EXTRACT(YEAR  FROM e.hired_at) = EXTRACT(YEAR  FROM m.hired_at)
  AND EXTRACT(MONTH FROM e.hired_at) = EXTRACT(MONTH FROM m.hired_at)
```

MySQL: `DATE_FORMAT(e.hired_at, '%Y-%m') = DATE_FORMAT(m.hired_at, '%Y-%m')`.
Обе формы применяют функцию к колонке, поэтому индекс по `hired_at` не сработает —
на больших таблицах это осознанный размен.

### 61. Сколько сотрудников, чье имя начинается и заканчивается одной буквой

```sql
SELECT COUNT(*) AS total
FROM   employees
WHERE  lower(LEFT(name, 1)) = lower(RIGHT(name, 1));
```

`LEFT`/`RIGHT` есть в PostgreSQL, MySQL и SQL Server. Стандартный вариант —
`SUBSTRING(name FROM 1 FOR 1)` и `SUBSTRING(name FROM LENGTH(name) FOR 1)`.

Что стоит проговорить: `lower()` нужен, иначе «Anna» не совпадет сама с собой;
хвостовые пробелы ломают сравнение — в реальных данных сначала `TRIM(name)`;
условие с функциями по обеим сторонам индекс не использует, это полный скан.

### 62. Отдел с самой высокой средней зарплатой

```sql
SELECT dept_id, ROUND(AVG(salary), 2) AS avg_salary
FROM   employees
GROUP  BY dept_id
ORDER  BY avg_salary DESC
LIMIT  1;
```

Сортировать по агрегату можно: `ORDER BY` выполняется после `GROUP BY`, поэтому
и алиас `avg_salary` тут виден (в `WHERE` он был бы недоступен).

Если отделов с одинаковым средним может быть несколько, `LIMIT 1` молча оставит один
случайный. Честные варианты:

```sql
... ORDER BY avg_salary DESC FETCH FIRST 1 ROW WITH TIES;   -- стандарт
```

```sql
SELECT * FROM (
  SELECT dept_id, AVG(salary) AS avg_salary,
         RANK() OVER (ORDER BY AVG(salary) DESC) AS rnk     -- окно поверх агрегата
  FROM   employees GROUP BY dept_id
) t WHERE rnk = 1;
```

Тот же шаблон отвечает на «категория с самой высокой средней ценой» и «топ-2 клиента
по сумме покупок» — меняются только таблица и агрегат.

### 63. Клиенты, которые заказывали, но ничего не возвращали

```sql
SELECT c.id, c.name
FROM   customers c
WHERE  EXISTS     (SELECT 1 FROM orders  o WHERE o.customer_id = c.id)
  AND  NOT EXISTS (SELECT 1 FROM returns r WHERE r.customer_id = c.id);
```

Два условия сразу: полусоединение (`EXISTS` — «хотя бы один заказ есть») и
антисоединение (`NOT EXISTS` — «ни одного возврата»). Ни то, ни другое не размножает строки,
в отличие от `JOIN`, поэтому `DISTINCT` не нужен.

Через соединения то же самое выглядит так — и как раз здесь легко ошибиться:

```sql
SELECT DISTINCT c.id, c.name
FROM   customers c
JOIN   orders  o ON o.customer_id = c.id      -- INNER: остались только заказавшие
LEFT   JOIN returns r ON r.customer_id = c.id
WHERE  r.customer_id IS NULL;                  -- и без возвратов
```

`DISTINCT` тут обязателен из-за размножения строк по заказам, а условие на `returns`
должно стоять в `WHERE` (`IS NULL` после `LEFT JOIN`), но никак не в `ON` — иначе фильтр
превратится в часть соединения и вернет всех заказавших.

`NOT IN (SELECT customer_id FROM returns)` в этой задаче опасен: один NULL в `returns.customer_id`
обнулит весь результат.

### 64. Каждый клиент и весь его первый заказ, включая клиентов без заказов

Отличие от вопроса 17 (последняя покупка) — в двух деталях: нужен **весь** заказ рядом
с **всеми** полями клиента, и клиенты без заказов из результата **не должны исчезнуть**.

```sql
SELECT c.*, o.*
FROM   clients c
LEFT   JOIN LATERAL (
  SELECT * FROM orders o
  WHERE  o.client_id = c.id
  ORDER  BY o.id
  LIMIT  1
) o ON true;                       -- LEFT, а не CROSS: иначе клиенты без заказов пропадут
```

Переносимый вариант — нумерация заказов внутри клиента:

```sql
WITH first_orders AS (
  SELECT o.*, ROW_NUMBER() OVER (PARTITION BY client_id ORDER BY id) AS rn
  FROM   orders o
)
SELECT c.*, o.*
FROM   clients c
LEFT   JOIN first_orders o ON o.client_id = c.id AND o.rn = 1;
```

Тут важно, что `o.rn = 1` стоит **в `ON`, а не в `WHERE`**: в `WHERE` это условие отбросило бы
NULL-строки, и `LEFT JOIN` молча превратился бы в `INNER`.

В PostgreSQL работает и короткая форма:

```sql
SELECT DISTINCT ON (c.id) c.*, o.*
FROM   clients c
LEFT   JOIN orders o ON o.client_id = c.id
ORDER  BY c.id, o.id;
```

Типичная ошибка на собеседовании — просто соединить таблицы и отсортировать по `o.id`:
запрос вернет **все** заказы каждого клиента, а не первый. Сортировка не ограничивает набор строк.

Если «первый» определяется датой, а не идентификатором, сортируйте по `created`
и добавляйте `id` вторым ключом — иначе при двух заказах с одинаковой датой выбор будет случайным.

### 65. Доля клиентов, у которых средний чек за 30 дней вырос

```sql
WITH last_30 AS (
  SELECT client_id, AVG(amount) AS avg_amount
  FROM   orders
  WHERE  created >= CURRENT_DATE - INTERVAL '30 days'
  GROUP  BY client_id
), prev_30 AS (
  SELECT client_id, AVG(amount) AS avg_amount
  FROM   orders
  WHERE  created >= CURRENT_DATE - INTERVAL '60 days'
    AND  created <  CURRENT_DATE - INTERVAL '30 days'
  GROUP  BY client_id
)
SELECT ROUND(100.0 * COUNT(*) FILTER (WHERE l.avg_amount > p.avg_amount)
             / NULLIF(COUNT(*), 0), 2) AS pct_grown
FROM   prev_30 p
JOIN   last_30 l USING (client_id);
```

Разбор по частям:

- **два CTE в одном `WITH`**, через запятую: второй `WITH` подряд — синтаксическая ошибка;
- **полуоткрытые интервалы** (`>= начало AND < конец`): при `BETWEEN` заказ на границе
  попал бы в оба периода и посчитался дважды;
- **`COUNT(*) FILTER (WHERE ...)`** считает только выросших, обычный `COUNT(*)` — знаменатель.
  Переносимо: `SUM(CASE WHEN l.avg_amount > p.avg_amount THEN 1 ELSE 0 END)`;
- **`NULLIF(..., 0)`** спасает от деления на ноль, если пересечения периодов нет;
- `100.0`, а не `100`, иначе целочисленное деление даст ноль.

Главный уточняющий вопрос — **что в знаменателе**. Здесь это клиенты, активные в **оба**
периода (`INNER JOIN`). Если доля считается от всех клиентов вообще:

```sql
SELECT ROUND(100.0 * COUNT(*) FILTER (WHERE l.avg_amount > p.avg_amount)
             / (SELECT COUNT(*) FROM clients), 2) AS pct_grown
FROM   prev_30 p JOIN last_30 l USING (client_id);
```

А если клиент заказывал только в последние 30 дней и его нужно считать «выросшим»,
`INNER JOIN` его потеряет — тогда `FULL JOIN` с `COALESCE(p.avg_amount, 0)`.
На собеседовании это стоит проговорить вслух: три разных знаменателя дают три разных ответа.

---

## Что чаще всего проваливают

| Ошибка | Как правильно |
| --- | --- |
| `NOT IN` с подзапросом, где есть NULL | `NOT EXISTS` |
| `= NULL` | `IS NULL` |
| `COUNT(*)` после `LEFT JOIN` | `COUNT(колонка_правой_таблицы)` |
| `BETWEEN` по датам с временем | полуинтервал `>= начало AND < конец_следующего` |
| функция на колонке в `WHERE` | условие на «голую» колонку, иначе индекс не работает |
| `RANK` там, где нужен `DENSE_RANK` | номера после ничьей пропускаются |
| `ROWS` и `RANGE` в окне | `RANGE` объединяет строки с равным значением сортировки |
| `MAX(дата)` для «последней записи» | `ROW_NUMBER` или `DISTINCT ON` — иначе дубли при равных датах |
| агрегат в `WHERE` | фильтр по агрегату идет в `HAVING` |
| `SELECT *` с `GROUP BY` | в списке только колонки из `GROUP BY` и агрегаты |
| `UPDATE`/`DELETE` без `WHERE` | сначала тот же `WHERE` в `SELECT`, потом транзакция |

Порядок выполнения, из которого следует половина этих правил:
`FROM → WHERE → GROUP BY → HAVING → оконные функции → SELECT → DISTINCT → ORDER BY → LIMIT`.
Поэтому в `WHERE` не видно алиасов из `SELECT` и не работают оконные функции,
а в `ORDER BY` — видно и работают.
