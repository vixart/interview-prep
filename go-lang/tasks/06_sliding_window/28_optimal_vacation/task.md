# 6.3. Оптимальное планирование отпуска

Раздел: `06_sliding_window / 28_optimal_vacation`

## Условие

Необходимо определить день `X` начала отпуска длиной в `V` дней так, чтобы отгулять его
в ближайшие `P` дней и пропустить минимум `Y` встреч.
Считаем, что уже завтра первый возможный день отпуска (`day: 1`).

## Входные данные

- `daysWithMeetings: []DayMeetings` — дни со встречами (уже упорядочены по дню)
  - `day: int` — номер дня
  - `meetings: int` — количество встреч в этот день
- `periodLength: int` — в какой срок надо отгулять ВЕСЬ отпуск (в ближайшие `P` дней)
- `vacationLength: int` — длительность отпуска в днях

## Выходные данные

- `[]int` — `[день X начала отпуска, сколько встреч Y пропустим]`

## Заготовка

```go
type DayMeetings struct {
    Day      int
    Meetings int
}

func optimalVacation(daysWithMeetings []DayMeetings, periodLength, vacationLength int) []int {
    // ...
}
```

## Пример

```
daysWithMeetings = [{day: 3, meetings: 1}, {day: 4, meetings: 3}, {day: 14, meetings: 3}, {day: 21, ...}, ...]
periodLength = 30
vacationLength = 7

Результат: [5, 0] - начать отпуск с 5 дня, пропустить 0 встреч
```
