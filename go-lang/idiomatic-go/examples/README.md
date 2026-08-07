# examples — примеры к конспекту

Один Go-модуль (`interviewprep/examples`), каждый подкаталог — отдельная запускаемая программа
или пакет с тестами. Ссылки из [`../conspect.md`](../conspect.md) ведут сюда.

```bash
cd examples
go build ./...            # собрать все
go vet ./...              # статический анализ
go test ./...             # все тесты зеленые
go run ./10_concurrency/backpressure    # запустить конкретный пример
go test -bench=. ./13_testing/benchmark # бенчмарки
```

Модуль вложен в репозиторий, но независим: у него свой `go.mod`, поэтому корневой модуль его не видит.
Внешние зависимости: `go-cmp` (тесты), `uuid` и `chi` (примеры с контекстом), `x/exp` (бенчмарк дженериков).

## Примеры, которые ведут себя «неправильно» специально

| Пример | Что происходит |
| --- | --- |
| `02_composite/array_conversion` | паникует в конце: срез короче массива |
| `07_generics/non_generic_tree` | паникует: в дерево из чисел кладут строку (до дженериков это ловилось только в рантайме) |
| `06_interfaces/interface_comparable`, `07_generics/comparable_pitfall` | паникуют: сравнение интерфейсов с несравнимым базовым типом |
| `10_concurrency/deadlock` | `fatal error: all goroutines are asleep` |
| `10_concurrency/goroutine_leak` | показывает утекшую горутину |
| `01_basics/const_declaration`, `06_interfaces/embedding`, `06_interfaces/method_set`, `07_generics/stack`, `07_generics/impossible_constraint` | строки с ошибками компиляции закомментированы, ошибка приведена рядом |
| `13_testing/race_detector` | гоночный тест спрятан за тегом: `go test -tags racedemo -race ./13_testing/race_detector` |
| `13_testing/httptest_and_integration` | интеграционный тест за тегом: `go test -tags integration ./...` (нужен сервер на :8080) |
| серверы и демо с сетью | работают, пока не остановишь: `11_stdlib/http_*`, `12_context/*`, `10_concurrency/backpressure`, `06_interfaces/dependency_injection` |
| `14_reflect_unsafe/reflect_structof_memoizer` | крутит демонстрацию кэша несколько секунд |
| `04_functions/defer_closer`, `10_concurrency/pipeline` | требуют аргументов: `go run . <файл>` / `go run . <строкаA> <строкаB>` |

## Содержание

### 01_basics — типы, переменные, константы
- `const_declaration` — типизированные/нетипизированные константы, попытка присвоить константе
- `iota_enum` — перечисления, повтор выражения, битовые флаги
- `type_conversion` — явные преобразования числовых типов
- `complex_numbers` — комплексные числа
- `shadow_variables`, `shadow_multiple_assignment`, `shadow_universe_block` — затенение переменных, `:=` и имен из universe block

### 02_composite — массивы, срезы, отображения, строки
- `slice_len_cap` — рост длины и емкости при `append`
- `slice_share_storage` — подсрез разделяет память с родителем
- `slice_append_storage` — `append` в подсрез затирает элементы родителя
- `full_slice_expression` — `x[low:high:max]` ограничивает емкость
- `slicing_slices`, `copy_slice` — срезание и `copy`
- `array_conversion`, `slice_array_memory` — срез ↔ массив (копия) и `(*[N]T)` (общая память)
- `string_slicing`, `string_to_slice` — байты, руны, срезы строк
- `map_read_write` — «запятая-ok», `delete`
- `map_set` — множество на `map[T]bool`
- `confusing_slices` — тот же сценарий без `:max`: append дерется за общую память

### 03_control_flow — управляющие конструкции
- `for_range`, `iterate_map` (случайный порядок), `iterate_string` (руны, а не байты)
- `for_label` — `continue`/`break` с меткой
- `switch`, `blank_switch` — обычный и пустой switch
- `goto_valid` — редкий валидный случай `goto`

### 04_functions — функции, замыкания, defer
- `variadic`, `named_optional_parameters` (структура опций), `named_return_values`
- `func_value`, `anon_func`, `closure`, `closure_shadow` (`:=` вместо `=` в замыкании), `closure_factory`, `sort_sample`
- `defer_order` — LIFO и вычисление аргументов в момент `defer`
- `defer_db_tx` — commit/rollback через именованную ошибку
- `defer_closer` — функция возвращает closure для освобождения ресурса
- `pass_value_type`, `pass_map_slice` — что реально копируется при передаче

### 05_pointers — указатели, память
- `pointer_primer` — `&`, `*`, `new`, nil-указатель
- `update_via_pointer` — почему «заполнение структуры по указателю» плохо
- `reusable_buffer` — один буфер вместо аллокаций в цикле
- `pointer_perf_bench` — бенчмарк: значение vs указатель (`go test -bench=. ./05_pointers/pointer_perf_bench`)

### 06_interfaces — типы, методы, интерфейсы
- `value_vs_pointer_receiver` — приемник по значению не меняет оригинал
- `auto_address_of` — Go сам берет адрес переменной при вызове указательного метода
- `method_set` — набор методов: почему значение не реализует интерфейс с указательными методами
- `embedding` — встраивание как композиция; встраивание ≠ наследование
- `no_dynamic_dispatch` — метод встроенного типа не видит переопределение
- `interface_nil` — типизированный nil ≠ nil-интерфейс
- `interface_comparable` — сравнение интерфейсов и паника
- `type_assertions`, `type_switch` — утверждение и переключатель типа
- `dependency_injection` — DI на неявных интерфейсах + функциональный тип как реализация интерфейса (поднимает сервер на :8080)
- `nil_receiver_tree` — методы, работающие на nil-приемнике

### 07_generics — обобщенные типы
- `stack`, `comparable_stack` — обобщенный стек, ограничение `comparable`
- `map_filter_reduce` — обобщенные функции высшего порядка
- `generic_tree`, `generic_linked_list` — обобщенные структуры данных
- `non_generic_tree` — как то же дерево писали до дженериков (типобезопасности нет)
- `generic_interface` — интерфейс с параметром типа как ограничение
- `type_terms` — списки типов и `~`
- `type_inference` — когда параметр типа приходится указывать явно
- `impossible_constraint` — нереализуемое ограничение (`int` + метод)
- `comparable_pitfall` — `comparable` пропускает интерфейсы, паника в рантайме
- `generics_perf_bench` — цена дженериков (`go test -bench=. ./07_generics/generics_perf_bench`)

### 08_errors — ошибки, panic/recover
- `error_basics`, `sentinel_error`, `custom_error`, `custom_error_nil_trap` (почему возвращать `error`, а не свой тип)
- `wrap_error`, `custom_wrapped_error` — `%w`, `Unwrap`
- `defer_wrap_error` — один `defer` вместо трех `fmt.Errorf`
- `join_error`, `multi_error` — `errors.Join`, `Unwrap() []error`
- `errors_is`, `errors_is_custom`, `errors_is_pattern_match`, `errors_as`
- `panic_recover`

### 09_tools — модули и инструменты
- `embed_file`, `embed_fs` — `//go:embed` в строку и в `embed.FS`
- `embed_hidden` — `parent_dir` vs `parent_dir/*` vs `all:parent_dir` (скрытые файлы)
- `import_alias` — псевдоним импорта при конфликте имен пакетов
- `go_generate_stringer` — `//go:generate` + сгенерированный `direction_string.go`
- `godoc_example` — правила Go Doc (`go doc ./09_tools/godoc_example`)

### 10_concurrency — конкурентность
- `goroutine_basics`, `buffered_channel`, `loop_var_capture`
- `deadlock` → `select_avoids_deadlock` — как select снимает взаимоблокировку
- `select_disable_case` — отключение ветви через `nil`-канал
- `timeout` — таймаут через `context.WithTimeout` + буфер на 1
- `waitgroup`, `waitgroup_gather` — ожидание и сбор результатов, закрытие канала следящей горутиной
- `goroutine_leak` — утечка и лечение контекстом
- `context_cancel` — остановка горутины по `ctx.Done()`
- `backpressure` — ограничение нагрузки буферизованным каналом
- `sync_once`, `sync_oncevalue` — ленивая инициализация
- `mutex` — `Mutex`/`RWMutex`
- `channel_vs_mutex/channel` и `channel_vs_mutex/mutex` — одна задача двумя способами
- `pipeline` — конвейер из этапов с контекстом (`go run . abc def`)

### 11_stdlib — стандартная библиотека
- `io_friends` — `io.Reader`/`io.Writer`, декораторы, `MultiReader`, `LimitReader`
- `json`, `json_encode_decode`, `json_custom`, `json_custom_wrapper` — теги, потоковый Encoder/Decoder, свои маршалеры
- `http_client` — `http.Client`, `NewRequestWithContext`, декодирование ответа
- `http_server`, `http_server_mux` — сервер с таймаутами, маршруты Go 1.22, вложенные mux
- `http_middleware` — `func(http.Handler) http.Handler`
- `http_response_controller` — потоковый ответ и `Flush`
- `structured_logging` — `log/slog`

### 12_context — контекст
- `context_user`, `context_guid` — значения в контексте, идиоматические обертки (поднимают серверы)
- `cancel`, `cancel_cause`, `timeout_cause` — отмена, отмена с причиной, таймаут с причиной
- `timeout_middleware` — дедлайн на каждый HTTP-запрос + различие 504/500
- `nested_timers` — вложенные дедлайны
- `own_cancellation` — собственная отмена долгих вычислений

### 13_testing — тестирование
- `basic_test`, `table_test` — обычный и табличный тест
- `testmain` — `TestMain`, `cleanup` — `t.Cleanup`, `env` — `t.Setenv`
- `testdata_files` — данные в каталоге `testdata`
- `public_api_test` — пакет `..._test` для тестирования публичного API
- `go_cmp` — сравнение через `cmp.Diff` и `cmpopts`
- `parallel` — `t.Parallel`
- `stub` — заглушки на интерфейсах
- `httptest_and_integration` — `httptest` + интеграционный тест с тегом сборки
- `benchmark` — `Benchmark*`, `b.N`, `-benchmem`
- `race_detector` — гонка и ее исправление

### 14_reflect_unsafe — reflect, unsafe, cgo
- `reflect_struct_tag` — `Type`, `Kind`, `Field`, теги структур
- `reflect_nil_check` — проверка «интерфейс содержит nil»
- `reflect_csv_marshaler` — свой маршалер CSV ↔ структуры
- `reflect_make_func` — обертка функции через `reflect.MakeFunc`
- `reflect_filter_bench` — цена рефлексии (`go test -bench=. ./14_reflect_unsafe/reflect_filter_bench`)
- `reflect_structof_memoizer` — мемоизация любой функции: тип ключа создается в рантайме (`reflect.StructOf`)
- `unsafe_sizeof_offsetof` — размеры, выравнивание, порядок полей
- `unsafe_binary_data` — разбор бинарного протокола через `unsafe.Pointer`
- `unsafe_unexported_field` — доступ к неэкспортируемому полю
- `cgo_call_c` — вызов C из Go (нужен gcc)
