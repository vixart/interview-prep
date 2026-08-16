# 4.1.1. Указатели: как работают и где ловят новичков

Тип задачи: вопрос с собеседования «что выведет код и почему».

## Условие

Дан код. Нужно сказать, что выведет программа на каждом `Println`, и объяснить почему.

```go
package main

import "fmt"

type User struct {
    Name string
}

func main() {
    user := User{Name: "Олег"}
    fmt.Println("Имя до обновления:", user.Name)

    updateUser(user)
    fmt.Println("Имя после обновления:", user.Name)
}

func updateUser(u User) {
    u.Name = "Таненбаум"
    fmt.Println("Имя внутри функции [updateUser]:", u.Name)

    resetUser(&u)
    fmt.Println("Имя после вызова функции [resetUser] внутри функции [updateUser]:", u.Name)
}

func resetUser(u *User) {
    u = &User{Name: "Безымянный"}
}
```

**Задание:** определи, что выведет программа (и скомпилируется ли она), и объясни почему.
