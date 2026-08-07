# Go — конспект для собеседования

Краткая выжимка по книге «Go. Идиомы и паттерны проектирования» (Джон Боднер).
Формат: только суть + ссылки на запускаемые примеры в [`examples/`](examples/README.md).

```bash
cd examples && make check   # собрать, проверить и прогнать все примеры
```

**Содержание**

1. [Команды](#1-команды) · 2. [Типы и переменные](#2-типы-и-переменные) · 3. [Составные типы](#3-составные-типы) ·
4. [Управляющие конструкции](#4-управляющие-конструкции) · 5. [Функции и defer](#5-функции-и-defer) ·
6. [Указатели, стек, куча](#6-указатели-стек-куча) · 7. [Типы, методы, интерфейсы](#7-типы-методы-интерфейсы) ·
8. [Дженерики](#8-дженерики) · 9. [Ошибки](#9-ошибки) · 10. [Модули и пакеты](#10-модули-и-пакеты) ·
11. [Инструменты](#11-инструменты) · 12. [Конкурентность](#12-конкурентность) ·
13. [Стандартная библиотека](#13-стандартная-библиотека) · 14. [Контекст](#14-контекст) ·
15. [Тестирование](#15-тестирование) · 16. [reflect, unsafe, cgo](#16-reflect-unsafe-cgo) ·
17. [Идиомы](#17-идиомы)

---

## 1. Команды

| Команда | Что делает |
| --- | --- |
| `go mod init <path>` | создает `go.mod`, объявляет каталог модулем |
| `go build` / `go run x.go` | собрать / собрать-запустить-удалить бинарник |
| `go fmt ./...` / `go vet ./...` | форматирование / поиск подозрительных конструкций |
| `go test ./...` | тесты; `-cover`, `-race`, `-bench=.`, `-short`, `-tags` |
| `go get ./...` / `go mod tidy` | добавить зависимости из импортов / привести `go.mod` и `go.sum` в соответствие с кодом |
| `go list -m -versions <mod>` | список версий модуля |
| `go install <pkg>@latest` | поставить сторонний инструмент в `~/go/bin` |
| `go clean -cache -testcache -modcache` | очистка кэшей |
| `go doc <pkg>[.<ident>]` | документация из исходников |
| `GOOS=linux GOARCH=amd64 go build` | кросс-компиляция |

---

## 2. Типы и переменные

### Базовые типы

- **Целые**: `int8/16/32/64`, `uint8/16/32/64`, `int`/`uint` (32 или 64 бита — зависит от платформы), `uintptr`.
  `byte` = псевдоним `uint8`, `rune` = псевдоним `int32`.
- **Разные целочисленные типы нельзя смешивать** без явного преобразования — даже `int` и `int64`. Ошибка компиляции.
- **Плавающая точка**: `float32`, `float64`. `x/0` → `±Inf`, `0/0` → `NaN`.
- **Комплексные**: `complex64` (на базе float32), `complex128` (на базе float64) — [пример](examples/01_basics/complex_numbers/main.go).
- **`bool`**, **`string`**. Строки неизменяемы; сравнение `== != < >`, конкатенация `+`.
- **Нулевые значения**: числа — `0`, `bool` — `false`, `string` — `""`, указатели/срезы/map/каналы/функции/интерфейсы — `nil`.

Литералы **нетипизированы**: подходят любому совместимому типу. Приведения к `bool` нет; между числовыми типами
преобразование только явное — [type_conversion](examples/01_basics/type_conversion/main.go).
Восьмеричные литералы удобны для прав `rwxrwxrwx`, hex/бинарные — для битовых масок.
Строковые литералы: интерпретируемые `"..."` и необработанные `` `...` `` (можно переносы, escape не работает).

### var или :=

```go
var x int = 10          // многословно, нужно только если тип не выводится
var x = 10              // тип по умолчанию int
var x int               // нулевое значение
var x, y = 10, "hello"  // == x, y := 10, "hello"
var ( a int; b = 20 )   // блок объявлений
```

- `:=` работает **только внутри функции**; на уровне пакета — только `var`.
- `:=` может переприсваивать, если хотя бы одна переменная слева новая.
- `var x byte = 20` идиоматичнее, чем `x := byte(20)`, когда нужен нестандартный тип.

### const и iota

- Константы вычисляются на этапе компиляции; это лишь имена для литералов. **Сделать переменную неизменяемой в Go нельзя.**
- Нетипизированная константа ведет себя как литерал: `const x = 10` подойдет и `int`, и `float64`, и `byte`.
  Типизированная `const x int = 10` — только своему типу. [пример](examples/01_basics/const_declaration/main.go)
- `iota` — номер строки внутри блока `const` (с нуля). Строка без выражения **повторяет предыдущее выражение**,
  но счетчик все равно растет. Классика: `1 << iota` для битовых флагов. [пример](examples/01_basics/iota_enum/main.go)
- Значения `iota` не хранят во внешних системах: вставка константы в середину сдвинет все последующие.

### Затенение (shadowing)

- Переменная затеняет одноименную из внешнего блока; внешняя остается нетронутой.
- Затенить можно **имя пакета** (`fmt := "oops"`) и предопределенные идентификаторы (`true := 10`) — в языке всего 25 ключевых слов,
  а `int`, `string`, `nil`, `make`, `true` живут в universe block и перекрываются.
- Главная ловушка — `:=` в блоке `if`/`for`: правая часть берет внешнюю переменную, левая создает новую.
  [shadow_variables](examples/01_basics/shadow_variables/main.go) · [shadow_multiple_assignment](examples/01_basics/shadow_multiple_assignment/main.go) · [universe block](examples/01_basics/shadow_universe_block/main.go)

---

## 3. Составные типы

### Массивы

```go
var x [3]int                 // [0 0 0]
x := [...]int{10, 20, 30}    // размер выводится
x := [12]int{1, 5: 4, 10: 100}  // разреженный литерал
var x [2][3]int              // «многомерность» = массив массивов
```

- **Размер — часть типа**: `[3]int` и `[4]int` разные типы, нельзя приводить, нельзя написать общую функцию.
- Размер обязан быть константой. Сравнимы через `==`, если сравним элемент.
- На практике массив нужен как хранилище под срез.

### Срезы

Срез — это дескриптор: `{указатель на массив, len, cap}`. Длина не входит в тип.

- `append` при `len == cap` выделяет новый массив (до 256 элементов — рост ×2, дальше по формуле ~×1.25)
  и **копирует данные** — старый массив уходит сборщику. Поэтому емкость лучше задавать заранее.
- Сравнивать `==` можно только с `nil`. Для сравнения — `slices.Equal` / `slices.EqualFunc`.
- `len(nil-срез) == 0`, `append` в nil-срез работает.
- `var x []int` — nil-срез; `x := []int{}` — не nil, длина 0 (удобно для JSON).

```go
x := make([]int, 5)      // len 5, cap 5
x := make([]int, 0, 10)  // len 0, cap 10  ← обычно то, что нужно
clear(x)                 // обнулить элементы
```

[len/cap при росте](examples/02_composite/slice_len_cap/main.go)

### Срезание: общая память

- `y := x[1:3]` **не копирует**: подсрез смотрит в ту же память, cap считается до конца исходного массива.
- `append` в подсрез **затирает элементы родителя**, пока хватает cap. Как только cap превышен — происходит
  перевыделение, и срезы становятся независимыми (самое неприятное — поведение зависит от данных).
- Лечение — full slice expression `x[low:high:max]`: обрезает cap, следующий `append` гарантированно скопирует.

[общая память](examples/02_composite/slice_share_storage/main.go) · [append затирает родителя](examples/02_composite/slice_append_storage/main.go) · [непредсказуемость без :max](examples/02_composite/confusing_slices/main.go) · [full slice expression](examples/02_composite/full_slice_expression/main.go) · [срезание](examples/02_composite/slicing_slices/main.go)

### copy и преобразование массив ↔ срез

- `copy(dst, src)` копирует `min(len(dst), len(src))` элементов и возвращает их количество; работает и с перекрытием. [пример](examples/02_composite/copy_slice/main.go)
- Массив → срез (`arr[:]`) — **общая память**.
- Срез → массив (`[4]int(s)`) — **копия**; если массив длиннее среза, паника в рантайме.
- Срез → указатель на массив (`(*[4]int)(s)`) — снова **общая память**.

[array_conversion](examples/02_composite/array_conversion/main.go) · [slice_array_memory](examples/02_composite/slice_array_memory/main.go)

### Строки, руны, байты

- Индекс и срез строки работают в **байтах**: `s[6]` — байт, `s[4:7]` — подстрока по байтам (можно разрезать руну пополам).
- `[]byte(s)` — байты, `[]rune(s)` — кодовые точки. Срезы рун нужны редко.
- Для реальной работы с текстом — `strings` и `unicode/utf8`.
- `string(int)` для не-`rune`/`byte` блокируется `go vet` начиная с 1.15.

[string_slicing](examples/02_composite/string_slicing/main.go) · [string_to_slice](examples/02_composite/string_to_slice/main.go)

### Отображения (map)

```go
var nilMap map[string]int      // nil: читать можно, писать — паника
m := map[string]int{}          // литерал
m := make(map[int][]string, 10) // len 0, но емкость под 10
v, ok := m["key"]              // идиома «запятая-ok»
delete(m, "key")               // идемпотентно
clear(m)
```

- Чтение отсутствующего ключа возвращает нулевое значение — отсюда нужда в «запятая-ok».
- `==` не работает; сравнение через `maps.Equal` / `maps.EqualFunc`.
- Порядок обхода **случайный** (защита от Hash DoS: хеш-сид генерируется при создании map).
  Исключение: `fmt.Println` печатает map отсортированной по ключам.
- Множества делают на `map[T]bool` или на `map[T]struct{}` (экономнее по памяти, но проверка только через «запятая-ok»).

[map_read_write](examples/02_composite/map_read_write/main.go) · [map_set](examples/02_composite/map_set/main.go)

### Структуры

- Классов нет. Литерал заполняется по порядку полей или по именам.
- Структура **сравнима**, если сравнимы все поля (срез, map, функция внутри делают ее несравнимой).
- Преобразование `T1(t2)` возможно, если имена, типы и **порядок** полей совпадают; теги полей игнорируются при преобразовании.
- Анонимные структуры нужны в двух местах: маршалинг/демаршалинг внешних данных и табличные тесты.
  Анонимную структуру можно присваивать и сравнивать с именованной при совпадении полей.

---

## 4. Управляющие конструкции

### if

Можно объявить переменные, видимые только внутри `if`/`else if`/`else`:

```go
if n := rand.Intn(10); n == 0 { ... } else if n > 5 { ... }
// здесь n уже не существует
```

### for — единственный цикл

Четыре формы: полная, только с условием (`for x < 10`), бесконечная (`for`), `for-range`.
Инициализация только через `:=`. `break` выходит, `continue` завершает итерацию.

**for-range:**

- По срезу/массиву — `i, v`; по map — `k, v`; по строке — **байтовый индекс и руна** (индекс прыгает на 1–4). [for_range](examples/03_control_flow/for_range/main.go)
- Значение **копируется** — менять `v` бессмысленно.
- Порядок обхода map случайный. [iterate_map](examples/03_control_flow/iterate_map/main.go) · [iterate_string](examples/03_control_flow/iterate_string/main.go)
- С Go 1.22 переменные цикла создаются **заново на каждой итерации** (раньше была одна на весь цикл — классический баг с горутинами). [loop_var_capture](examples/10_concurrency/loop_var_capture/main.go)
- Метки: `continue outer` / `break outer` для выхода из вложенных циклов. [for_label](examples/03_control_flow/for_label/main.go)

`for-range` — для полного перебора, обычный `for` — когда надо пропускать элементы. Бесконечный `for` заменяет `do…while`.

### switch

- Ветви **не проваливаются**; провалиться явно — `fallthrough`. `break` в `case` выходит из switch
  (чтобы выйти из внешнего `for` — метка над циклом и `break label`).
- В `case` можно несколько значений через запятую; пустой `case` = ничего не делать.
- **Пустой (blank) switch** — `switch { case x < 5: ... }` — заменяет цепочку `if/else if`, когда сравнения разнородные.
  [switch](examples/03_control_flow/switch/main.go) · [blank_switch](examples/03_control_flow/blank_switch/main.go)

### goto

Есть, но почти не нужен. Нельзя перепрыгнуть через объявление переменной и нельзя прыгнуть внутрь другого блока.
Оправдан только для выхода из глубоко вложенной логики. [goto_valid](examples/03_control_flow/goto_valid/main.go)

---

## 5. Функции и defer

### Сигнатуры

- Именованных и опциональных параметров нет — передавай **структуру опций**. [пример](examples/04_functions/named_optional_parameters/main.go)
- Вариативный параметр — один и последний, внутри это срез; развернуть срез — `f(vals...)`. [пример](examples/04_functions/variadic/main.go)
- Множественный возврат: возвращать надо **все** значения; лишнее гасится `_`.
- Именованные возвращаемые значения — по сути предобъявленные переменные. Нужны в основном ради `defer`,
  который может их прочитать и изменить. **Пустой `return` при них — источник багов, не используй.** [пример](examples/04_functions/named_return_values/main.go)
- Функция — значение: `var f func(string) int`, нулевое значение `nil`, можно объявить тип `type opFunc func(int, int) int`.
  [func_value](examples/04_functions/func_value/main.go) · [anon_func](examples/04_functions/anon_func/main.go)

### Замыкания

Функция внутри функции видит и изменяет ее переменные. Отсюда — передача функций в `sort.Slice`
и фабрики функций (`makeMult`). Ловушка: `:=` внутри замыкания создает новую переменную вместо изменения внешней.

[closure](examples/04_functions/closure/main.go) · [ловушка `:=`](examples/04_functions/closure_shadow/main.go) · [closure_factory](examples/04_functions/closure_factory/main.go) · [sort_sample](examples/04_functions/sort_sample/main.go)

### defer

- Выполняется при выходе из функции, **после** вычисления `return`, порядок **LIFO**.
- Аргументы отложенного вызова вычисляются **сразу**, тело — потом. [defer_order](examples/04_functions/defer_order/main.go)
- Возвращаемые значения отложенной функции прочитать нельзя.
- Через именованную ошибку `defer` управляет commit/rollback транзакции. [defer_db_tx](examples/04_functions/defer_db_tx/defer_db.go)
- Тем же приемом обертывают все ошибки функции одним сообщением. [defer_wrap_error](examples/08_errors/defer_wrap_error/main.go)
- Идиома: функция, выделившая ресурс, возвращает closure, который его освобождает. [defer_closer](examples/04_functions/defer_closer/simple_cat_cancel.go)

### Передача по значению

В Go **всё** передается копией:

| Что передаем | Что происходит |
| --- | --- |
| простые типы, массивы, структуры | копируется целиком, изменения не видны снаружи |
| map | копируется дескриптор → изменения видны снаружи |
| срез | копируется дескриптор `{ptr, len, cap}`: **изменение элементов видно**, `append` — нет (у копии своя len; при превышении cap еще и своя память) |

[pass_value_type](examples/04_functions/pass_value_type/main.go) · [pass_map_slice](examples/04_functions/pass_map_slice/main.go)

---

## 6. Указатели, стек, куча

- `&x` — адрес, `*p` — разыменование, нулевое значение указателя `nil`, разыменование nil — паника.
- `new(T)` возвращает `*T` с нулевым значением (редко используется; чаще `&T{}`).
- Срезы, map, функции, каналы, интерфейсы внутри реализованы через указатели.
- [pointer_primer](examples/05_pointers/pointer_primer/main.go)

**Когда указатель уместен:**

- Указатель в сигнатуре — сигнал «функция изменит значение». [update_via_pointer](examples/05_pointers/update_via_pointer/main.go)
- Отличить «нет значения» от нулевого (`func findUser(id int) *User` возвращает `nil`).
- Когда API требует: `json.Unmarshal(data, &v)` — иначе изменения потеряются.
- Для больших данных: передача по указателю выигрывает примерно от **1 МБ**; на меньших объемах копия быстрее.
  [бенчмарк](examples/05_pointers/pointer_perf_bench/perf_test.go)
- Не заполняй структуру «через указатель в параметре» — возвращай значение.
- Вместо аллокаций в цикле — один переиспользуемый буфер (`data := make([]byte, 100)` + `file.Read(data)`).
  [reusable_buffer](examples/05_pointers/reusable_buffer/main.go)

### Стек и куча

- **Стек** — непрерывный блок памяти на поток; выделение = сдвиг указателя стека, освобождение при выходе из функции.
  Требует знать размер на этапе компиляции (поэтому размер массива — часть типа). Стек горутины растет динамически,
  при росте копируется целиком.
- В **кучу** попадает то, что не помещается в эту модель: escape analysis отправляет туда данные, указатель на которые
  возвращается из функции или живет дольше кадра стека.
- Мусор — данные, на которые не осталось ни одного указателя. Сборщик мусора в Go оптимизирован **на низкую задержку**
  (цикл ≤ ~500 мкс), а не на максимальную пропускную способность.
- Вторая проблема кучи — разброс данных по памяти: последовательное чтение (стек, срез структур) быстрее произвольного
  (срез указателей) — до двух порядков. Это и есть «механическая симпатия».

### Настройка сборщика

- `GOGC` (по умолчанию 100): следующая сборка при `HEAP + HEAP*GOGC/100`. Больше — реже собираем, больше памяти. `GOGC=off` выключает.
- `GOMEMLIMIT` — мягкий лимит общей памяти (`B`, `KiB`, `MiB`, `GiB`, `TiB`); по умолчанию фактически выключен.
  Связка «`GOGC=off` + `GOMEMLIMIT`» — типовая настройка для сервисов с предсказуемой нагрузкой.

---

## 7. Типы, методы, интерфейсы

### Типы и методы

```go
type Score int
type Converter func(string) Score
type Person struct{ FirstName, LastName string; Age int }

func (p Person) String() string { ... }   // p — приемник, имя в одну букву
```

- Методы объявляются **только на уровне пакета** и только для типов этого пакета. Перегрузки нет.
- Приемник по указателю — если метод меняет приемник или должен работать на nil-приемнике
  ([методы на nil-указателе](examples/06_interfaces/nil_receiver_tree/main.go)).
  Приемник по значению — если не меняет. **Если хотя бы один метод с указателем — делай все методы с указателем.**
- Go сам берет адрес/разыменовывает при вызове метода на переменной, поэтому разница видна не сразу.
  [value_vs_pointer_receiver](examples/06_interfaces/value_vs_pointer_receiver/main.go) · [auto_address_of](examples/06_interfaces/auto_address_of/main.go)
- Метод можно использовать как значение: `f1 := myAdder.AddTo` (метод-значение), `f2 := Adder.AddTo` (выражение метода, первый аргумент — приемник).
- `type HighScore Score` — новый тип, **не наследование**: методы не наследуются, присваивание требует явного преобразования.

### Встраивание (композиция)

```go
type Manager struct {
    Employee          // встроенное поле без имени
    Reports []Employee
}
// поля и методы Employee поднимаются наверх: m.ID, m.Description()
```

- Одноименное поле/метод внешней структуры **затеняет** встроенное (`o.X` vs `o.Inner.X`).
- Встраивание ≠ наследование: `var e Employee = m` не скомпилируется, нужно `m.Employee`.
- **Динамической диспетчеризации нет**: метод встроенного типа вызывает метод встроенного типа, а не переопределенный внешний.
  [embedding](examples/06_interfaces/embedding/main.go) · [no_dynamic_dispatch](examples/06_interfaces/no_dynamic_dispatch/main.go)
- Встраивание помогает реализовать интерфейс «за счет» встроенного типа.

### Интерфейсы

- Единственный абстрактный тип, реализуются **неявно** — тип не знает, какие интерфейсы он удовлетворяет.
  Отсюда: интерфейс объявляют **на стороне потребителя**, а не рядом с реализацией. Имена обычно на `-er`.
- **Набор методов**: у `*T` — методы с приемником и по значению, и по указателю; у `T` — только по значению.
  Поэтому `var i Incrementer = valueCounter` не компилируется, а `&valueCounter` — да. [method_set](examples/06_interfaces/method_set/main.go)
- Интерфейсы встраиваются друг в друга (`ReadCloser = Reader + Closer`).
- Декоратор: обернуть один `io.Reader` другим (`gzip.NewReader(r)`) — базовый паттерн стандартной библиотеки.
- **Принимай интерфейсы, возвращай структуры.** Добавление метода в интерфейс ломает всех, кто его реализует;
  добавление поля/метода в структуру — обратно совместимо. Возврат структуры экономит аллокации:
  значение, положенное в интерфейс, обычно уезжает в кучу.
- Пустой интерфейс `interface{}` (= `any`) — «любое значение». Нужен на границе (JSON, БД), внутри логики — признак проблемы.

### Интерфейс и nil

Интерфейс внутри — пара указателей `{тип, значение}`. **Интерфейс равен `nil`, только если оба поля nil.**

```go
var p *Counter          // nil
var i Incrementer       // nil
i = p                   // тип = *Counter, значение = nil
i == nil                // false ← классический баг «типизированный nil»
```

Лечится тем, что функции возвращают конкретный тип, а не заранее подготовленную переменную-интерфейс.
[interface_nil](examples/06_interfaces/interface_nil/main.go) · проверка через рефлексию: [reflect_nil_check](examples/14_reflect_unsafe/reflect_nil_check/main.go)

### Сравнение интерфейсов

`==` сравнивает пару (тип, значение). Если базовый тип несравним (срез, map, функция) — **паника в рантайме**.
То же и для ключей map интерфейсного типа. [interface_comparable](examples/06_interfaces/interface_comparable/main.go)

### Утверждение и переключатель типа

```go
i2, ok := i.(MyInt)      // всегда «запятая-ok», иначе паника
switch j := i.(type) { case nil: ...; case int: ...; case bool, rune: ... }
```

- Преобразование типа проверяется компилятором, **утверждение — в рантайме** и только для интерфейсов.
- Использовать редко: обычно чтобы проверить, реализует ли значение доп. интерфейс (`io.WriterTo`, `errors.Is/As`).
- Ловушка декоратора: обертка, реализующая только интерфейс A, «прячет» интерфейс B исходного типа.
  [type_assertions](examples/06_interfaces/type_assertions/main.go) · [type_switch](examples/06_interfaces/type_switch/main.go)

### Функциональные типы как интерфейсы

```go
type HandlerFunc func(http.ResponseWriter, *http.Request)
func (f HandlerFunc) ServeHTTP(w http.ResponseWriter, r *http.Request) { f(w, r) }
```

Так обычная функция становится `http.Handler`. Правило: одна простая функция — параметр функционального типа (`sort.Slice`);
зависимость от многих функций/состояния — интерфейс.
Точно так же сделан `LoggerAdapter` в [dependency_injection](examples/06_interfaces/dependency_injection/main.go).

### Внедрение зависимостей

Неявные интерфейсы делают DI бесплатным: конструктор принимает интерфейсы, реализации подставляются в `main`,
фреймворк не нужен (при желании — Wire). [dependency_injection](examples/06_interfaces/dependency_injection/main.go)

---

## 8. Дженерики

```go
type Stack[T any] struct{ vals []T }        // T — параметр типа, any — ограничение
func (s *Stack[T]) Push(val T)              { s.vals = append(s.vals, val) }
func (s *Stack[T]) Pop() (T, bool) {
    if len(s.vals) == 0 { var zero T; return zero, false }  // идиома нулевого значения
    ...
}
```

- Ограничения — это интерфейсы. `any` — любой тип, `comparable` — поддерживает `==`.
  [stack](examples/07_generics/stack/main.go) · [comparable_stack](examples/07_generics/comparable_stack/main.go)
- Обобщенные функции: `Map`, `Filter`, `Reduce` пишутся один раз на все типы. [map_filter_reduce](examples/07_generics/map_filter_reduce/main.go)
- Ограничением может быть любой интерфейс, в том числе с параметром типа (`Differ[T]`) — так задают
  «тип, умеющий сравнить себя с собой». [generic_interface](examples/07_generics/generic_interface/main.go)
- **Списки типов** задают допустимые базовые типы и разрешают операторы:
  ```go
  type Integer interface{ ~int | ~int8 | ~int64 | ~uint }   // ~ = «и любой тип на основе этого»
  ```
  Без `~` пользовательский `type MyInt int` не подойдет. Такой интерфейс годится только как ограничение.
  [type_terms](examples/07_generics/type_terms/main.go)
- Ограничение вида «`int` + метод `String()`» нереализуемо: базовый `int` не может иметь методов. [impossible_constraint](examples/07_generics/impossible_constraint/main.go)
- Тип выводится из аргументов; если параметр типа есть только в результате — указывай явно: `Convert[int, int64](a)`.
  [type_inference](examples/07_generics/type_inference/main.go)
- Константы должны влезать во все типы списка: `in + 1_000` не скомпилируется для `int8`.
- `comparable` **пропускает интерфейсы**, а сравнение интерфейса с несравнимым содержимым паникует в рантайме —
  типобезопасность здесь дырявая. [comparable_pitfall](examples/07_generics/comparable_pitfall/main.go)

**Чего в Go нет:** методов с собственными параметрами типа (только параметры получателя), вариативных параметров
разных типов, специализации, каррирования, метапрограммирования.

**Цена:** дженерики компилируются по «GC-shape stenciling» — для указательных типов один общий код + таблица методов,
поэтому обобщенный код иногда медленнее мономорфного. [бенчмарк](examples/07_generics/generics_perf_bench/generics_performance_test.go)

Структуры данных: [generic_tree](examples/07_generics/generic_tree/main.go) · [generic_linked_list](examples/07_generics/generic_linked_list/main.go) ·
[как это писали до дженериков](examples/07_generics/non_generic_tree/main.go) (ошибка типа всплывает только в рантайме)

---

## 9. Ошибки

- Ошибка — **последнее возвращаемое значение**, тип `error` (интерфейс с одним методом `Error() string`).
  Нет ошибки — `nil`; есть ошибка — остальным значениям нулевые значения. [error_basics](examples/08_errors/error_basics/main.go)
- Создание: `errors.New("...")`, `fmt.Errorf("%d is not even", i)`.

### Виды ошибок

| Прием | Когда | Пример |
| --- | --- | --- |
| Сигнальная ошибка `var ErrX = errors.New(...)` | обработка невозможна, вызывающий сверяет по значению | [sentinel_error](examples/08_errors/sentinel_error/main.go) |
| Свой тип ошибки | нужны поля (код, статус) | [custom_error](examples/08_errors/custom_error/main.go) |
| Обертывание `%w` | добавить контекст, сохранив исходную | [wrap_error](examples/08_errors/wrap_error/main.go) |
| `errors.Join` | несколько ошибок сразу (валидация) | [join_error](examples/08_errors/join_error/main.go) |

**Ловушка:** объявляй возвращаемый тип как `error`, а не как свой конкретный тип — иначе «пустая» структура
в интерфейсе даст `err != nil` при отсутствии ошибки. [custom_error_nil_trap](examples/08_errors/custom_error_nil_trap/main.go)

### Обертывание

- `fmt.Errorf("in fileChecker: %w", err)` — оборачивает; `%v` — только текст, цепочка теряется.
- Свой тип оборачивает через метод `Unwrap() error` (для нескольких — `Unwrap() []error`).
  [custom_wrapped_error](examples/08_errors/custom_wrapped_error/main.go) · [multi_error](examples/08_errors/multi_error/main.go)
- `errors.Unwrap` напрямую почти не вызывают — есть `Is`/`As`.
- Один `defer` вместо повторяющихся `fmt.Errorf` в каждой ветке. [defer_wrap_error](examples/08_errors/defer_wrap_error/main.go)

### errors.Is и errors.As

- `errors.Is(err, ErrNotExist)` — ищет **конкретный экземпляр/значение** по всей цепочке (сравнение через `==`).
  Для несравнимых типов или сравнения «по маске» определи свой метод `Is(target error) bool`.
  [errors_is](examples/08_errors/errors_is/main.go) · [errors_is_custom](examples/08_errors/errors_is_custom/main.go) · [errors_is_pattern_match](examples/08_errors/errors_is_pattern_match/main.go)
- `errors.As(err, &myErr)` — ищет **тип** и записывает найденное во второй параметр (обязательно указатель, иначе паника).
  Работает и с указателем на интерфейс. [errors_as](examples/08_errors/errors_as/main.go)
- Правило: экземпляр → `Is`, тип → `As`.

### panic / recover

- Паника = среда выполнения не знает, что делать дальше. Раскручивает стек, выполняя `defer`; неперехваченная паника
  в любой горутине убивает программу.
- `recover()` работает **только внутри `defer`**. [panic_recover](examples/08_errors/panic_recover/main.go)
- Использовать для фатальных ситуаций и для того, чтобы не уронить процесс из-за одного запроса:
  паника → `recover` → лог/мониторинг → `os.Exit(1)` или ответ 500. Библиотека не должна выпускать панику наружу — конвертируй в `error`.

---

## 10. Модули и пакеты

**Репозиторий → модуль → пакеты.** Модуль версионируется целиком, желательно один модуль на репозиторий.
Путь модуля = путь репозитория (`github.com/user/project`), путь импорта = путь модуля + путь каталога. Относительных импортов нет.

### go.mod

```
module github.com/author/project
go 1.25            // минимальная версия языка
toolchain go1.25.6 // какой тулчейн скачать (auto/local/goX.Y.Z, переопределяется GOTOOLCHAIN)
require ( ... )    // прямые зависимости; вторая секция — косвенные, с // indirect
```

- Go ≥ 1.21 при более новой директиве `go` **сам скачает** нужный тулчейн; более старые версии просто игнорировали ее.
- `go.sum` хранит хеши модулей и их `go.mod`; несовпадение хеша прерывает сборку.
- `replace example.com/old => ../local/old` — подмена модуля (локальные пути лучше не коммитить),
  `exclude` — запретить чужую версию, `retract` — отозвать **свою** версию.
- Мажорные версии ≥ v2 живут в отдельном каталоге `/v2` со своим `go.mod` и считаются **разными модулями**,
  поэтому v1 и v2 могут использоваться одновременно.

### Пакеты

- Экспортируется то, что начинается с заглавной буквы. `internal/` виден только внутри поддерева родителя `internal`.
- Циклических зависимостей нет. Имя пакета не обязано совпадать с каталогом; `main` импортировать нельзя.
- Имена: `names.Extract` / `names.Format` — пакет-существительное, функция-глагол; `util.ExtractNames` — антипаттерн.
  Точечный импорт `import . "fmt"` — не используй, а вот псевдоним при конфликте имен нормален:
  [import_alias](examples/09_tools/import_alias/main.go).
- Раскладка: маленький модуль — один пакет; приложение — вся логика в `internal`, в `main` минимум;
  библиотека — пакет в корне с именем репозитория (без дефисов), утилиты — в `cmd/<binary>`.
- `init()` избегать: скрытая регистрация через пустой импорт `_ "app/plugin"` делает поток данных нечитаемым.
- Go Doc: комментарий вплотную к объявлению, начинается с имени; `[fmt]` — ссылка на пакет, `[Add]` — на символ.
  [godoc_example](examples/09_tools/godoc_example/mini.go) · локальный сайт: `go install golang.org/x/pkgsite/cmd/pkgsite@latest && pkgsite`.

### Рабочие пространства (workspaces)

```bash
go work init ./workspace_app && go work use ./workspace_lib   # создает go.work
GOWORK=off go build                                            # собрать без workspace
```

Позволяет править несколько модулей локально без `replace` и без публикации версий.
`go.work` — локальный файл, в git его не коммитят.

### Прокси

`go get` ходит через `proxy.golang.org` (кэш версий, защита от удаления) и сверяет хеши с `sum.golang.org`.
`GOPROXY=direct` — качать напрямую из репозиториев, `GOPRIVATE` — для приватных.

---

## 11. Инструменты

- **Линтеры**: `go vet` (базовый) → `staticcheck` → `revive` → `golangci-lint` (агрегатор, 50+ анализаторов).
  `govulncheck` — уязвимости в зависимостях. `goimports` — сортировка импортов.
- **`//go:embed`** — встраивание файлов в бинарник (нужен импорт `embed`):
  ```go
  //go:embed passwords.txt
  var passwords string
  //go:embed help
  var helpInfo embed.FS      // all:help — включая скрытые файлы
  ```
  [embed_file](examples/09_tools/embed_file/main.go) · [embed_fs](examples/09_tools/embed_fs/main.go) · [скрытые файлы и `all:`](examples/09_tools/embed_hidden/main.go)
- **`go:generate`** — комментарий, запускающий кодогенерацию (`stringer`, `protoc-gen-go`); вызывать из Makefile,
  результат коммитить. [go_generate_stringer](examples/09_tools/go_generate_stringer/main.go)
- `go version -m <binary>` — из каких модулей и коммита собран бинарник.
- Проверить код на другой версии языка: `go install golang.org/dl/go1.19.2@latest && go1.19.2 download && go1.19.2 build`.

---

## 12. Конкурентность

Модель — CSP: «не делитесь памятью — общайтесь через каналы». Конкурентность нужна, только если она реально ускоряет.

### Горутины

- Легковесный поток, управляемый планировщиком Go (модель M:N — миллионы горутин на десятке потоков ОС).
  Стек начинается с килобайтов и растет. [goroutine_basics](examples/10_concurrency/goroutine_basics/main.go)
- Запуская горутину, сразу решай, **как она завершится** — иначе утечка: заблокированная навсегда горутина
  держит память и остается в планировщике. [goroutine_leak](examples/10_concurrency/goroutine_leak/main.go)
- Значения, меняющиеся снаружи, передавай параметром (до Go 1.22 это было обязательным для переменной цикла).
  [loop_var_capture](examples/10_concurrency/loop_var_capture/main.go)

### Каналы

```go
ch := make(chan int)      // небуферизованный
ch := make(chan int, 10)  // буферизованный: cap(ch)=10, len(ch)=сколько лежит сейчас
a := <-ch; ch <- b        // чтение / запись
func read(ch <-chan int)  // параметр только для чтения
func write(ch chan<- int) // параметр только для записи
```

| Операция | nil-канал | открытый пустой | заполненный | закрытый |
| --- | --- | --- | --- | --- |
| чтение | блокировка навсегда | блокировка | значение | нулевое значение, `ok == false` |
| запись | блокировка навсегда | буферизованный — ок, небуферизованный — блокировка до читателя | блокировка | **паника** |
| `close` | паника | ок | ок | **паника** |

- Значение читает **ровно одна** горутина. `for v := range ch` читает до закрытия.
- Закрывает канал **пишущая** сторона. Отличить нулевое значение от закрытия — `v, ok := <-ch`.
- По умолчанию используй небуферизованные каналы. Буфер оправдан, когда заранее известно число результатов
  или нужно ограничить конкурентность. [buffered_channel](examples/10_concurrency/buffered_channel/main.go)
- Каналы и мьютексы не должны торчать в публичном API (исключение — библиотеки-хелперы для конкурентности).

### select

- Ждет несколько операций сразу; из готовых выбирает **случайную** — поэтому голодания не бывает.
- Без `default` блокируется до готовности ветви; `default` делает вызов неблокирующим (в цикле сжигает CPU).
- Взаимоблокировка (`fatal error: all goroutines are asleep`) лечится именно select'ом.
  [deadlock](examples/10_concurrency/deadlock/main.go) → [select_avoids_deadlock](examples/10_concurrency/select_avoids_deadlock/main.go)
- Прочитанный до конца канал в `select` всегда «готов» и крутит цикл — обнули переменную канала (`in = nil`),
  чтобы отключить ветвь. [select_disable_case](examples/10_concurrency/select_disable_case/main.go)
- Связка `for { select { ... } }` — стандартный цикл обработки.

### Типовые паттерны

| Задача | Решение | Пример |
| --- | --- | --- |
| дождаться N горутин | `sync.WaitGroup`: `Add` **до** запуска, `defer wg.Done()` | [waitgroup](examples/10_concurrency/waitgroup/main.go) |
| собрать результаты | воркеры пишут в `out`, отдельная горутина `wg.Wait(); close(out)` | [waitgroup_gather](examples/10_concurrency/waitgroup_gather/main.go) |
| остановить горутину | `select { case <-ctx.Done(): return; case ch <- v: }` | [context_cancel](examples/10_concurrency/context_cancel/main.go) |
| таймаут на операцию | `context.WithTimeout` + канал результата **с буфером 1** (иначе утечка) | [timeout](examples/10_concurrency/timeout/main.go) |
| ограничить нагрузку (backpressure) | буферизованный канал-семафор + `default` → «нет мест» | [backpressure](examples/10_concurrency/backpressure/main.go) |
| выполнить один раз | `sync.Once` (поле, не локальная переменная), `sync.OnceFunc/OnceValue/OnceValues` | [sync_once](examples/10_concurrency/sync_once/main.go) · [sync_oncevalue](examples/10_concurrency/sync_oncevalue/main.go) |
| конвейер этапов | канал на каждый этап + контекст | [pipeline](examples/10_concurrency/pipeline/main.go) |

### Мьютексы

- `sync.Mutex`: `Lock`/`Unlock`; `sync.RWMutex`: `RLock`/`RUnlock` для читателей (много читателей ИЛИ один писатель).
  `RWMutex` выигрывает при «часто читаем, редко пишем».
- Всегда `defer mu.Unlock()` сразу после `Lock`.
- **Не реентерабельны**: повторный `Lock` в той же горутине = вечная блокировка. Не вызывай под блокировкой чужой код.
- Мьютекс, `WaitGroup`, `Once` **нельзя копировать** — передавай по указателю (`go vet` это ловит).
- Выбор: координация горутин и передача данных — каналы; защита поля структуры — мьютекс;
  доказанная проблема производительности с каналами — мьютекс. Одна задача двумя способами:
  [channel_vs_mutex/channel](examples/10_concurrency/channel_vs_mutex/channel/main.go) · [channel_vs_mutex/mutex](examples/10_concurrency/channel_vs_mutex/mutex/main.go) · [mutex](examples/10_concurrency/mutex/main.go)
- Любую переменную, доступную нескольким горутинам, защищай. Проверка — `go test -race` / `go build -race`.
  [race_detector](examples/13_testing/race_detector/race.go)

---

## 13. Стандартная библиотека

### io и друзья

```go
type Reader interface{ Read(p []byte) (n int, err error) }
type Writer interface{ Write(p []byte) (n int, err error) }
```

- `Read` пишет в **переданный буфер** (одна аллокация на весь цикл) и возвращает `n`; конец данных — `io.EOF`.
  Важно: последние байты могут прийти **вместе** с `io.EOF` — сначала обработай `buf[:n]`, потом проверяй ошибку.
- Декораторы: `gzip.NewReader(r)`, `io.MultiReader`, `io.LimitReader`, `io.MultiWriter`, `io.TeeReader`.
- Комбинированные интерфейсы: `io.ReadCloser`, `io.ReadSeeker`, `io.ReadWriteCloser` — используй самый узкий.
- `io.NopCloser(r)` — адаптер `Reader` → `ReadCloser` (редкий случай, когда стандартная библиотека возвращает интерфейс).
- `os.ReadFile`/`os.WriteFile`/`io.ReadAll` — только для небольших данных; иначе `os.Open` + `bufio.Scanner`.
- [io_friends](examples/11_stdlib/io_friends/main.go)

### time

- `time.Duration` — `int64` наносекунд: `2*time.Hour + 30*time.Minute`, `time.ParseDuration("2h30m")`, `Truncate`/`Round`.
- `time.Time` — момент времени: `Now`, `Parse(layout, s)`, `Format(layout)` (layout — эталон `2006-01-02 15:04:05`),
  `Add`, `Sub`, `Before`/`After`/`Equal`, `Year/Month/Day/Clock`.
- Периодика: `time.NewTicker` (есть `Stop`) вместо `time.Tick`; разовая задержка — `time.After` / `time.NewTimer`.
- Сравнивать `Time` через `==` нельзя (монотонные часы и локация) — только `Equal`.

### encoding/json

```go
type Order struct {
    ID          string    `json:"id"`
    DateOrdered time.Time `json:"date_ordered"`
    CustomerID  string    `json:"customer_id,omitempty"`
    Secret      string    `json:"-"`
}
data, err := json.Marshal(order)          // Go → JSON
err = json.Unmarshal(data, &order)        // JSON → Go (обязательно указатель)
json.NewEncoder(w).Encode(v)              // прямо в io.Writer
json.NewDecoder(r).Decode(&v)             // прямо из io.Reader
```

- Работают **только экспортируемые** поля. Без тега имя берется как есть, при разборе регистр игнорируется.
- `omitempty` считает пустыми `""`, `0`, `nil`, пустые срез/map — **но не нулевую структуру**.
- Свои правила разбора — методы `MarshalJSON`/`UnmarshalJSON` (обычно на типе-обертке).
- `json.MarshalIndent(v, "", "  ")` — читаемый вывод.
- [json](examples/11_stdlib/json/main.go) · [json_encode_decode](examples/11_stdlib/json_encode_decode/main.go) · [json_custom](examples/11_stdlib/json_custom/main.go) · [json_custom_wrapper](examples/11_stdlib/json_custom_wrapper/main.go)

### net/http — клиент

- `http.Client` безопасен для конкурентного использования, одного на приложение достаточно.
  **Всегда задавай `Timeout`** — у `http.DefaultClient` (и у `http.Get`/`Post`) его нет.
- Схема: `http.NewRequestWithContext(ctx, method, url, body)` → `req.Header.Add(...)` → `client.Do(req)` →
  проверить `res.StatusCode` → `defer res.Body.Close()` → `json.NewDecoder(res.Body).Decode(&v)`.
- [http_client](examples/11_stdlib/http_client/main.go)

### net/http — сервер

- `http.Server`: `Addr`, `Handler`, обязательно `ReadTimeout`/`WriteTimeout`/`IdleTimeout` (защита от зависших клиентов).
- Обработчик — интерфейс `http.Handler` с `ServeHTTP(w, r)`; порядок записи ответа:
  `w.Header()` → `w.WriteHeader(code)` → `w.Write(data)` (для 200 `WriteHeader` можно не вызывать).
- `http.NewServeMux()` + `mux.HandleFunc("GET /hello/{name}", ...)`; параметр пути — `r.PathValue("name")` (Go 1.22).
  `ServeMux` сам является `Handler`, поэтому mux'ы вкладываются; `http.StripPrefix` срезает префикс.
- Не используй `http.Handle`, `http.HandleFunc`, `http.ListenAndServe` — они пишут в глобальный `DefaultServeMux`.
- [http_server](examples/11_stdlib/http_server/main.go) · [http_server_mux](examples/11_stdlib/http_server_mux/main.go)

**Middleware** — `func(http.Handler) http.Handler`: обертка вокруг обработчика (логирование, авторизация, тайминги,
rate limiting, контекст). Цепочка: `mux → m1 → m2 → handler`. [http_middleware](examples/11_stdlib/http_middleware/main.go)

**`http.NewResponseController(w)`** — `Flush()` (потоковый ответ, `Transfer-Encoding: chunked`), `Hijack()` (доступ к TCP),
`SetReadDeadline`/`SetWriteDeadline`. [http_response_controller](examples/11_stdlib/http_response_controller/main.go)

### log/slog

Структурированное логирование: уровни, атрибуты `slog.String("k", v)`, JSON- и текстовый handler,
`logger.With(...)` для общих полей. Дешевле, чем форматировать строки, и пригодно для машинного разбора.
[structured_logging](examples/11_stdlib/structured_logging/main.go)

---

## 14. Контекст

- `context.Context` передается **первым параметром** и явно; не храни его в структурах.
- Корни: `context.Background()`, `context.TODO()` (заглушка).
- Порождение: `WithCancel`, `WithDeadline`, `WithTimeout` (+варианты `...Cause` с причиной), `WithValue`.
  Все возвращают новый контекст, отмена — только вниз по дереву.
- `cancel` **всегда** вызывается (`defer cancel()`) — иначе утекают таймеры и горутины.
  `ctx.Done()` — канал, закрывающийся при отмене; `ctx.Err()` — `context.Canceled` / `DeadlineExceeded`;
  `context.Cause(ctx)` — исходная причина.
- Значения: ключ — неэкспортируемый тип (`type userKey struct{}`), а не строка; наружу — пара функций
  `ContextWithUser(ctx, u)` / `UserFromContext(ctx)`. Класть в контекст стоит только то, что сквозное
  (trace id, пользователь), а не обязательные аргументы функции.
- В HTTP: `r.Context()`, `r.WithContext(ctx)` — так middleware добавляет значения (и дедлайны:
  [timeout_middleware](examples/12_context/timeout_middleware/main.go)); клиентский запрос строится через `NewRequestWithContext`.

[context_user](examples/12_context/context_user/main.go) · [context_guid: сервис](examples/12_context/context_guid/app1/main.go) и [обертки над контекстом](examples/12_context/context_guid/tracker/tracker.go) · [cancel](examples/12_context/cancel/main.go) · [cancel_cause](examples/12_context/cancel_cause/main.go) · [timeout_cause](examples/12_context/timeout_cause/main.go) · [nested_timers](examples/12_context/nested_timers/main.go) · [own_cancellation](examples/12_context/own_cancellation/main.go)

---

## 15. Тестирование

Тесты лежат рядом с кодом, в файлах `*_test.go`, функции `TestXxx(t *testing.T)`. Запуск — `go test ./...` (`-v`, `-run`).

- `t.Error`/`t.Errorf` — продолжить, `t.Fatal` — прервать тест. [basic_test](examples/13_testing/basic_test/adder_test.go)
- `TestMain(m *testing.M)` — общая настройка/очистка на пакет (`m.Run()` + `os.Exit`). [testmain](examples/13_testing/testmain/testmain_test.go)
- `t.Cleanup(func(){...})` — освобождение ресурсов конкретного теста, `t.Setenv` — переменная окружения на тест
  (несовместим с `t.Parallel`). [cleanup](examples/13_testing/cleanup/cleanup_test.go) · [env](examples/13_testing/env/env_test.go)
- Данные — в каталоге `testdata` (Go его игнорирует при сборке). [testdata_files](examples/13_testing/testdata_files/text_test.go)
- Тест публичного API — пакет `имя_test` с импортом тестируемого пакета. [public_api_test](examples/13_testing/public_api_test/adder_public_test.go)
- **Табличные тесты** — срез анонимных структур + `t.Run(name, ...)`. [table_test](examples/13_testing/table_test/table_test.go)
- Сравнение результатов — `github.com/google/go-cmp`: `cmp.Diff(expected, got)` печатает разницу;
  нестабильные поля глушатся `cmpopts.IgnoreFields` или своим `cmp.Comparer`. [go_cmp](examples/13_testing/go_cmp/cmp_test.go)
- `t.Parallel()` — тест исполняется конкурентно с другими параллельными. [parallel](examples/13_testing/parallel/code_test.go)
- Покрытие: `go test -cover -coverprofile=c.out` + `go tool cover -html=c.out`.

### Заглушки и имитации

- Абстрагировать зависимость можно интерфейсом или функциональным типом — тогда в тест подставляется заглушка.
- **Stub** возвращает заготовленный ответ; **mock** дополнительно проверяет, что вызовы были в нужном порядке
  с нужными аргументами. Библиотеки: `gomock`, `testify`.
- Приемы: встроить интерфейс в структуру-заглушку (реализовывать нужно только используемые методы);
  заглушка из полей-функций — гибкая настройка на каждый тест. [stub](examples/13_testing/stub/stub_test.go)
- HTTP-зависимости — `httptest.NewServer`. [httptest_and_integration](examples/13_testing/httptest_and_integration/remote_solver_test.go)

### Интеграционные тесты и теги сборки

```go
//go:build integration
```
```bash
go test -tags integration ./...
```

[remote_solver_integration_test.go](examples/13_testing/httptest_and_integration/remote_solver_integration_test.go) —
тег отделяет тесты, которым нужна внешняя среда (БД, сервис в Docker). Альтернатива `-short` + `testing.Short()`
хуже: она дает всего два уровня и говорит о «долготе», а не о зависимости.

### Бенчмарки

```go
func BenchmarkFileLen(b *testing.B) {
    for i := 0; i < b.N; i++ {   // b.N подбирается фреймворком
        result, _ := FileLen("testdata/data.txt", 1)
        blackhole = result       // присвоение в пакетную переменную, чтобы компилятор не выбросил вызов
    }
}
```

```bash
go test -bench=. -benchmem ./13_testing/benchmark   # ns/op, B/op, allocs/op
```

Подбенчмарки через `b.Run(name, ...)` удобно сравнивают варианты (например, размер буфера) —
вывод показывает, где кончается выигрыш. [benchmark](examples/13_testing/benchmark/bench_test.go)
Профилирование — `pprof` (CPU/память, в том числе через HTTP-эндпоинт `net/http/pprof`).

### Детектор гонок

`go test -race` (или `go run -race`) ловит несинхронизированный доступ к переменной из разных горутин.
Не гарантирует нахождение всех гонок и сильно замедляет программу, поэтому включается в тестах/CI, а не в проде.
[race_detector](examples/13_testing/race_detector/race.go) — `go test -tags racedemo -race ./13_testing/race_detector`.

### Фаззинг

```go
func FuzzXxx(f *testing.F) {
    f.Add("seed")                                   // корректные примеры
    f.Fuzz(func(t *testing.T, in string) { YourFunc(in) })
}
```

`go test -fuzz=FuzzXxx` генерирует входные данные и ищет паники/ошибки. Найденный вход сохраняется файлом
в `testdata/fuzz/<Тест>/` и дальше выполняется обычным `go test` как регрессионный тест.

---

## 16. reflect, unsafe, cgo

Все три пакета нужны **на границе** программы, а не в бизнес-логике: `reflect` — текстовые данные и БД,
`unsafe` — ОС и бинарные протоколы, `cgo` — C-библиотеки.

### reflect

Три концепции: **тип** (`reflect.TypeOf`), **разновидность** (`Kind`: Struct, Ptr, Slice, Map, Func…), **значение** (`reflect.ValueOf`).

- `Name()` — имя типа (у срезов и указателей пусто), `Kind()` — из чего он сделан (`type Foo struct{}` → Kind `Struct`, Type `Foo`).
- `Elem()` — тип/значение, на которое ссылается указатель, срез, map, канал, массив.
- Структуры: `NumField()` / `Field(i)` → `StructField` (имя, тип, **теги**). Так работают `encoding/json` и ORM.
  [reflect_struct_tag](examples/14_reflect_unsafe/reflect_struct_tag/main.go)
- Чтение значения: `Interface()` + утверждение типа, либо `Int()`, `String()`, `Bool()`, `Bytes()`.
- **Запись**: только через указатель — `reflect.ValueOf(&i).Elem().SetInt(20)`; иначе паника.
- Создание: `reflect.New(t)` (аналог `new`), `MakeSlice`, `MakeMap`, `MakeChan`, `MakeFunc`.
- Вызов метода, не подходящего разновидности, **паникует** — всегда проверяй `Kind`.
- Проверка «интерфейс содержит nil»: `IsValid()` + `IsNil()` для Ptr/Slice/Map/Func/Interface. [reflect_nil_check](examples/14_reflect_unsafe/reflect_nil_check/main.go)
- Свой маршалер (CSV ↔ структуры) — типовое применение. [reflect_csv_marshaler](examples/14_reflect_unsafe/reflect_csv_marshaler/main.go)
- `reflect.MakeFunc` оборачивает произвольную функцию (например, замером времени). [reflect_make_func](examples/14_reflect_unsafe/reflect_make_func/main.go)
- `reflect.StructOf` создает типы на лету — экзотика ([мемоизация любой функции](examples/14_reflect_unsafe/reflect_structof_memoizer/main.go)). **Методы рефлексией создать нельзя**, значит и реализовать интерфейс — нельзя.
- Рефлексия медленная: обобщенный `Filter` на рефлексии проигрывает обычному в разы. [reflect_filter_bench](examples/14_reflect_unsafe/reflect_filter_bench/filter_test.go)
- В стандартной библиотеке используется в `database/sql`, `text/template`, `fmt`, `errors.Is/As`, `sort.Slice`, `encoding/*`.
- `reflect.DeepEqual` — только для тестов (и то лучше `go-cmp`).

### unsafe

Три функции и один тип. `unsafe.Pointer` — универсальный указатель: в него и из него можно преобразовать
указатель любого типа и `uintptr` (над которым разрешена арифметика).

- `unsafe.Sizeof` — размер значения: указатель 8 байт, строка 16 (ptr+len), срез 24 (ptr+len+cap), map 8 (это указатель).
- `unsafe.Offsetof` — смещение поля. Размер структуры = сумма полей **+ выравнивание**: порядок полей влияет на размер
  (`bool,int64,bool` → 24 байта, а `int64,bool,bool` → 16). [unsafe_sizeof_offsetof](examples/14_reflect_unsafe/unsafe_sizeof_offsetof/main.go)
- Паттерн 1 — преобразование несовместимых типов цепочкой через `unsafe.Pointer`;
  паттерн 2 — чтение/запись байтов значения. Быстрый разбор бинарного протокола: `*(*Data)(unsafe.Pointer(&b))`
  (не забудь про порядок байтов). [unsafe_binary_data](examples/14_reflect_unsafe/unsafe_binary_data/main.go)
- `reflect` + `unsafe` дают доступ к неэкспортируемым полям — крайний случай. [unsafe_unexported_field](examples/14_reflect_unsafe/unsafe_unexported_field/main.go)
- Проверка корректности: `go test -gcflags=-d=checkptr`.

### cgo

```go
/*
#cgo LDFLAGS: -lm
#include <math.h>
int add(int a, int b) { return a + b; }
*/
import "C"

sum := C.add(3, 2)
```

- Комментарий перед `import "C"` — это C-код; `.c`/`.h` файлы рядом подхватываются автоматически.
- Нужен C-компилятор; cgo **замедляет** вызовы (переключение стеков) и ломает кросс-компиляцию — это средство
  интеграции, а не оптимизации. [cgo_call_c](examples/14_reflect_unsafe/cgo_call_c/main.go)

---

## 17. Идиомы

- Нулевое значение осмысленно: `var b bytes.Buffer`, `var mu sync.Mutex` и `var wg sync.WaitGroup` работают сразу.
- Каждая **локальная** переменная и каждый импорт должны использоваться — иначе ошибка компиляции
  (на переменные уровня пакета и константы это не распространяется; неиспользуемые константы просто не попадают в бинарник).
- Числовые литералы разделяются `_`: `1_000_000`. Исходники всегда в UTF-8.
- `camelCase` в именах; регистр первой буквы задает экспорт. Чем уже область видимости, тем короче имя (`i`, `k`, `v`, `p`).
- Передача всегда по значению; указатель — явный сигнал изменения.
- Функция высшего порядка = функция, принимающая функцию.
- Среда выполнения (планировщик, GC, сеть, аллокатор) вшита в каждый бинарник — отсюда ~2 МБ на «hello world»
  и отсутствие внешних зависимостей у результата сборки.
- Сложность кода растет в первую очередь от **глубины вложенности** и слабой структуризации — предпочитай ранний
  `return` и плоские функции.

### Что чаще всего спрашивают

| Вопрос | Ответ в одну строку |
| --- | --- |
| Чем срез отличается от массива | массив — значение фиксированного размера (размер в типе), срез — дескриптор `{ptr,len,cap}` на массив |
| Что будет при `append` в подсрез | перезапишет элементы родителя, пока хватает cap; при превышении — перевыделит и отвяжется |
| Почему `i != nil` для nil-указателя в интерфейсе | интерфейс = пара (тип, значение), тип не nil |
| Когда метод не попадает в набор методов | метод с указательным приемником у значимого типа |
| Чем `Is` отличается от `As` | `Is` ищет экземпляр/значение, `As` — тип и записывает результат |
| Как остановить горутину | контекст/канал отмены в `select`; запуск без плана завершения = утечка |
| Небуферизованный или буферизованный канал | по умолчанию небуферизованный; буфер — когда известно число результатов или нужен семафор |
| Каналы или мьютекс | передача владения данными — каналы, защита поля структуры — мьютекс |
| Где живет переменная | решает escape analysis: возвращаемый указатель и «переживание» кадра → куча |
| Зачем `defer cancel()` | освобождает таймер и дочерние горутины контекста |
