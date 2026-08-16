# 4.8.2. Ревью кода: аналитика посещений

Раздел: `code_review / 02_visits_analytics`

Тип задачи: ревью кода — найти проблемы, объяснить и исправить.

## ТЗ

Веб-сервис для анализа посещений пользователями различных локаций.
Программа предоставляет HTTP endpoint `/visits`, который:

- Принимает параметр `period` (`week` или `day`).
- Запрашивает данные из БД за указанный период.
- Группирует посещения по дням.
- Возвращает JSON с количеством посещений и топ локацией за каждый день.

### Требования

- Код должен быть безопасным (нет SQL injection).
- Ошибки должны обрабатываться корректно.
- Типы данных должны быть правильными.

> Этот код НАМЕРЕННО содержит ошибки для учебных целей! Не запускайте в production!

## Исходный код (ключевые фрагменты)

```go
type Visit struct { /* ... */ }

type DayVisit struct {
    Day         string // день, в формате "2006-01-02"
    Count       int    // количество посещений в этот день
    TopLocation string // наиболее посещаемая локация в этот день
}

func main() {
    connStr := "postgres://admin:password123@localhost:5432/analytics?sslmode=disable"
    db, _ := sql.Open("postgres", connStr)

    r := chi.NewRouter()
    r.Use(middleware.Logger)
    r.Get("/visits", handleVisits(db))

    _ = http.ListenAndServe(":3000", r)
}

func handleVisits(db *sql.DB) func(w http.ResponseWriter, r *http.Request) {
    return func(w http.ResponseWriter, r *http.Request) {
        ctx := context.WithValue(
            context.Background(),
            "period", Period(r.URL.Query().Get("period")))

        visits, _ := getVisitsFromDB(db, ctx)
        dayVisits := dayVisitsFromVisits(visits)

        bytes, _ := json.Marshal(dayVisits)
        _ = w.Write(bytes)
    }
}

func getVisitsFromDB(db *sql.DB, ctx context.Context) ([]Visit, error) {
    // ... собирает SQL-запрос из значения period через конкатенацию/Sprintf
}
```

**Задание:** сделай ревью: найди проблемы, объясни их и предложи исправленный вариант.
