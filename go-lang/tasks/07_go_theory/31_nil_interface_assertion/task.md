# 4.5.4. Ещё один баг с интерфейсами: type assertion к `interface{}`

Тип задачи: найти и объяснить баг, исправить код.

## Условие

```go
package main

import (
    "fmt"
    "sync"
)

type Cache struct {
    data map[string]interface{}
    mu   sync.RWMutex
}

func (c *Cache) Store(key string, value interface{}) {
    c.mu.Lock()
    defer c.mu.Unlock()

    c.data[key] = value
}

func (c *Cache) Load(key string) interface{} {
    c.mu.RLock()
    defer c.mu.RUnlock()

    return c.data[key]
}

func main() {
    cache := &Cache{
        data: make(map[string]interface{}),
    }

    cache.Store("name", "Alice")
    cache.Store("age", 25)

    name := cache.Load("name").(string)
    fmt.Println("Name:", name)

    age := cache.Load("age").(int)
    fmt.Println("Age:", age)

    height := cache.Load("height")
    if height == nil {
        fmt.Println("Height not found")
    }

    width := cache.Load("width").(float64)
    fmt.Println("Width:", width)
}
```

**Задание:** определи, что выведет программа (и скомпилируется ли она), и объясни почему.
