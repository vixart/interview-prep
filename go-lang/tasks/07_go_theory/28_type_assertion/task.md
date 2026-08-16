# 4.5.1. Приведение типов: type assertion и type switch

Тип задачи: вопрос с собеседования «скомпилируется ли код и что произойдёт».

## Условие

```go
package main

type User struct{}

func (u *User) Create() {}
func (u *User) Get()    {}
func (u *User) List()   {}
func (u *User) Delete() {}

type Reader interface {
    Get()
    List()
}

type Writer interface {
    Create()
    Delete()
}

func main() {
    var userReader Reader = &User{}
    userWriter := userReader.(Writer)
    userWriter.Get()
    _ = userWriter
}
```

**Задание:** определи, что выведет программа (и скомпилируется ли она), и объясни почему.
