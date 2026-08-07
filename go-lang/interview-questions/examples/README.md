# examples — код к разбору вопросов

Отдельный Go-модуль (`interviewprep/interview`), не связанный с модулем конспекта.
Ссылки из [`../answers.md`](../answers.md) ведут сюда.

```bash
cd examples
make check                              # fmt + vet + build + test
go run ./live_coding/workerpool/senior  # запустить конкретный пример
go run -race ./problems/data_race       # то же с детектором гонок
```

## live_coding — то, что просят написать руками

| Пример | Что показывает |
| --- | --- |
| [parallel_processing](live_coding/parallel_processing/main.go) | распараллеливание с сохранением порядка: каждый пишет в свой индекс, синхронизация не нужна |
| [sync_goroutines](live_coding/sync_goroutines/main.go) | четыре способа синхронизации, включая свою мини-реализацию `errgroup` |
| [fan_in](live_coding/fan_in/main.go) | объединение нескольких источников в один канал |
| [limit_concurrency](live_coding/limit_concurrency/main.go) | семафор на канале: ждать слот или сразу отказать (противодавление) |

### workerpool — одна задача на трех уровнях

Полные реализации из [`go-lang/workerpool`](../../workerpool/README.md), разбор — в
[answers.md](../answers.md#сделать-воркер-пул). Комментарии в коде помечают,
как это работает и где именно закопаны баги.

| Пример | Что показывает |
| --- | --- |
| [junior](live_coding/workerpool/junior/pool.go) | базовые 40 строк и пять багов: `send on closed`, нет recover, нет ctx, потеря задач, мьютекс рядом с каналом |
| [middle](live_coding/workerpool/middle/pool.go) | канал задач не закрывается никогда, две фазы остановки, recover на задачу, `Submit(ctx, fn) error`, `Config` |
| [senior](live_coding/workerpool/senior/pool.go) | политики `OnFull`, `Resize` через `quitOne`, счетчики дропов, `worker_id` в контексте |

```bash
go run  ./live_coding/workerpool/junior         # счастливый путь: баги не видны
go test ./live_coding/workerpool/junior         # тест воспроизводит панику send on closed channel
go test -race ./live_coding/workerpool/junior   # ПАДАЕТ намеренно: детектор показывает гонку
go test -race ./live_coding/workerpool/middle ./live_coding/workerpool/senior   # а эти зеленые
```

## theory — вопросы про поведение языка

| Пример | Что показывает |
| --- | --- |
| [buffered_vs_unbuffered](theory/buffered_vs_unbuffered/main.go) | рандеву против очереди, что паникует у закрытых и nil-каналов |
| [map_internals](theory/map_internals/main.go) | адрес элемента, nil-map, рост, `delete`, конкурентная запись, `sync.Map` |
| [defer_pitfalls](theory/defer_pitfalls/main.go) | порядок, момент вычисления аргументов, defer в цикле, именованный результат |

## problems — что ломается в проде

| Пример | Что показывает |
| --- | --- |
| [data_race](problems/data_race/main.go) | гонка и четыре способа ее убрать: шардирование, atomic, mutex, канал |
| [goroutine_leak](problems/goroutine_leak/main.go) | четыре источника утечки горутин и как их считать |
| [deadlock](problems/deadlock/main.go) | повторный `Lock`, разный порядок замков, канал без пары, `RLock`→`Lock` |
| [mutex_copy](problems/mutex_copy/main.go) | что происходит при копировании мьютекса |
| [memory](problems/memory/main.go) | аллокации, escape analysis, удержание массива подсрезом, `sync.Pool` |

## runtime — как это работает внутри

| Пример | Что показывает |
| --- | --- |
| [scheduler](runtime/scheduler/main.go) | G-M-P, `GOMAXPROCS`, цена горутины, вытеснение при `GOMAXPROCS=1` |
| [gc](runtime/gc/main.go) | циклы и паузы GC, влияние `GOGC` и `GOMEMLIMIT`, живые данные против мусора |

## Примеры, которые запускаются отдельно

Часть демонстраций спрятана за тегами сборки, чтобы обычный `make check` оставался зеленым:

```bash
go run -tags copylocks ./problems/mutex_copy   # копирование мьютекса + go vet ругается
go vet  -tags copylocks ./problems/mutex_copy
go run -tags mapfatal   ./theory/map_internals # fatal error: concurrent map writes (процесс умрет)
go run -race            ./problems/data_race   # WARNING: DATA RACE
```

Полезные переменные окружения:

```bash
GODEBUG=schedtrace=1000 go run ./runtime/scheduler   # состояние очередей планировщика
GODEBUG=gctrace=1       go run ./runtime/gc          # каждая сборка мусора одной строкой
go build -gcflags='-m'  ./problems/memory            # решения escape analysis
```

`problems/deadlock` намеренно оставляет четыре заблокированные горутины — каждый сценарий
запускается с таймаутом, поэтому программа завершается сама.
