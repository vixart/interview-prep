# 4.7.1. Fan-in: объединение нескольких каналов в один

## Условие

> Напиши функцию, которая merge'ит N каналов в один.

```go
package main

func main() {
    channels := make([]chan int64, 10)
    for i := range channels {
        channels[i] = make(chan int64)
    }

    for i := range channels {
        go func(i int) {
            channels[i] <- int64(i)
            close(channels[i])
        }(i)
    }

    for v := range merge(channels...) {
        println(v)
    }
}
```

## Что нужно реализовать

```go
func merge(channels ...chan int64) <-chan int64
```

Функция должна вернуть канал, в который попадут значения из всех входных каналов,
и который закроется, когда закроются все входные каналы.
