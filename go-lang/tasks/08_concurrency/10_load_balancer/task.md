# 4.6.10. Client-side балансировщик нагрузки

Раздел: `concurrency / 10_balancer`

## Условие

Реализовать client-side балансировщик нагрузки для микросервисов.

### Контекст

- Есть микросервис с интерфейсом `Backend`.
- Для каждого микросервиса запущено несколько десятков экземпляров.
- Каждый экземпляр доступен по своему адресу.
- Экземпляры ненадёжны: могут падать, быть недоступными или перегруженными.

### Требования

- Реализовать тип `Balancer`, который также реализует `Backend`.
- Использовать алгоритм Round Robin для балансировки.
- Учитывать конкурентный доступ (метод может вызываться из горутин).
- Балансировщик должен равномерно распределять нагрузку.

## Заготовка

```go
package main

import "context"

type Request interface{}

type Response interface{}

type Backend interface {
    Invoke(ctx context.Context, req Request) (Response, error)
}

type BackendImpl struct {
    addr string
}

var _ Backend = &BackendImpl{}

// TODO: реализовать
type Balancer struct {
    // backends []Backend
    // next     atomic.Uint64  (или mutex)
}

var _ Backend = &Balancer{}
```
