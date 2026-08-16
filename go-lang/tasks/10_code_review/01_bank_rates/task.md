# 4.8.1. Ревью кода: сбор курсов валют из банков

Раздел: `code_review / 01_bank`

Тип задачи: ревью кода — найти проблемы, объяснить и предложить рефакторинг.

## Что делает программа

CLI-утилита `./currency update`, которая:

1. Разбирает аргументы командной строки (`os.Args`): команды `help` и `update`.
2. Держит захардкоженный список банков со ссылками на курсы:

```go
urlsBank := []struct {
    bankName string
    curFrom  string
    curTo    string
    url      string
}{
    {bankName: "Bank 1", curFrom: "RUB", curTo: "USD", url: "http://bank.example.com/rates/rub-usd"},
    {bankName: "Bank 2", curFrom: "RUB", curTo: "USD", url: "http://bank2.example.com/rates?..."},
}

clientBank := &http.Client{}
```

3. Последовательно ходит по HTTP за курсом каждого банка, парсит ответ.
4. Для каждого курса вызывает `updateCurrency`, которая пишет значение в PostgreSQL:

```go
const (
    host     = "localhost"
    port     = 5432
    user     = "postgres"
    password = ""
    dbname   = ""
)

func updateCurrency(bank, from, to string, value float64) error {
    psqlconn := fmt.Sprintf("host=%s port=%d "+
        "user=%s password=%s dbname=%s sslmode=disable",
        host, port, user, password, dbname)

    db, err := sql.Open("postgres", psqlconn)
    if err != nil { panic(err) }

    defer db.Close()

    err = db.Ping()
    if err != nil { panic(err) }

    fmt.Println("Connected!")

    insertStmt := fmt.Sprintf("insert into currency_rates (bank, ...) values ('%s', ...)", ...)
    _, err = db.Exec(insertStmt)
    return err
}
```

**Задание:** сделай ревью: найди проблемы, объясни их и предложи исправленный вариант.
