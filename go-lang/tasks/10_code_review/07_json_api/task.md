# 4.8.7. Ревью кода: JSON API для управления пользователями

Раздел: `code_review / 07_json_api`

Тип задачи: ревью кода — сделайте ревью кода и исправьте проблемы.

## ТЗ

REST API для управления пользователями.
Программа предоставляет HTTP endpoints для создания и получения пользователей.

> Этот код НАМЕРЕННО содержит ошибки для учебных целей! Не запускайте в production!

## Исходный код

```go
package main

import (
    "encoding/json"
    "fmt"
    "io"
    "net/http"
)

type User struct {
    ID    int
    Name  string
    Email string
    Age   int
}

var users = make(map[int]User)
var nextID = 1

func main() {
    http.HandleFunc("/users", handleUsers)
    http.HandleFunc("/user/", handleUser)

    fmt.Println("Server starting on :8080")
    http.ListenAndServe(":8080", nil)
}

func handleUsers(w http.ResponseWriter, r *http.Request) {
    if r.Method == "POST" {
        createUser(w, r)
    } else if r.Method == "GET" {
        listUsers(w, r)
    }
}

func createUser(w http.ResponseWriter, r *http.Request) {
    // Читаем тело запроса
    body, _ := io.ReadAll(r.Body)

    // Парсим JSON
    var user User
    json.Unmarshal(body, &user)

    // Создаем пользователя
    user.ID = nextID
    nextID++
    users[user.ID] = user

    // Возвращаем ответ
    json.NewEncoder(w).Encode(user)
}

func listUsers(w http.ResponseWriter, r *http.Request) {
    // Собираем всех пользователей
    // ...
    json.NewEncoder(w).Encode(allUsers)
}

func handleUser(w http.ResponseWriter, r *http.Request) {
    // Получаем ID из URL
    var id int
    fmt.Sscanf(r.URL.Path, "/user/%d", &id)

    if r.Method == "GET" {
        user := users[id]
        json.NewEncoder(w).Encode(user)
    } else if r.Method == "DELETE" {
        delete(users, id)
        w.WriteHeader(http.StatusOK)
    }
}
```

**Задание:** сделай ревью: найди проблемы, объясни их и предложи исправленный вариант.
