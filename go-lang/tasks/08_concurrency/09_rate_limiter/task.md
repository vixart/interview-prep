# 4.6.9. Rate Limiter — ограничение RPS

Раздел: `concurrency / 09_rate_limiter`

## Условие

Реализовать rate limiter для ограничения RPS запросов к API.

### Требования

- Функция вызывается при каждом запросе к API.
- Возвращает `true`, если запрос можно обработать, `false`, если нужно отклонить.
- Гарантировать, что ни за какой промежуток времени длиной в 1 секунду
  не будет обработано более `RPS` запросов.
- `RPS` может быть от 1 до 100000.
- Учитывать конкурентный доступ (метод вызывается из нескольких горутин).

Напишите код и несколько тестов, демонстрирующих работоспособность.

## Заготовка

```go
package main

type RateLimiter struct{}

// NewRateLimiter создает новый rate limiter с заданным RPS.
func NewRateLimiter(rps int) *RateLimiter { panic("implement me") }

// Allow возвращает true, если запрос можно пропустить.
func (rl *RateLimiter) Allow() bool { panic("implement me") }
```
