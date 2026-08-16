// Оптимальное планирование отпуска: окно длиной vacationLength по дням
// 1..periodLength, минимизируем количество встреч в окне.
//
// daysWithMeetings разрежен и отсортирован по дню, поэтому сумму встреч
// в окне поддерживаем двумя указателями по массиву встреч: при сдвиге
// старта на 1 день убираем день, вышедший слева, и добавляем день,
// вошедший справа. Время O(P + M), память O(1).
package main

import "fmt"

type DayMeetings struct {
	Day      int
	Meetings int
}

func optimalVacation(daysWithMeetings []DayMeetings, periodLength, vacationLength int) []int {
	lastStart := periodLength - vacationLength + 1
	if lastStart < 1 {
		return nil // отпуск не помещается в период
	}

	sum := 0
	l, r := 0, 0 // [l, r) — встречи внутри текущего окна

	// начальное окно [1, vacationLength]
	for r < len(daysWithMeetings) && daysWithMeetings[r].Day <= vacationLength {
		sum += daysWithMeetings[r].Meetings
		r++
	}

	bestStart, bestSum := 1, sum
	for start := 2; start <= lastStart; start++ {
		// день start-1 вышел из окна
		if l < len(daysWithMeetings) && daysWithMeetings[l].Day == start-1 {
			sum -= daysWithMeetings[l].Meetings
			l++
		}
		// день start+vacationLength-1 вошёл в окно
		if r < len(daysWithMeetings) && daysWithMeetings[r].Day == start+vacationLength-1 {
			sum += daysWithMeetings[r].Meetings
			r++
		}

		if sum < bestSum { // строгое < даёт самый ранний из лучших стартов
			bestStart, bestSum = start, sum
		}
	}

	return []int{bestStart, bestSum}
}

func main() {
	meetings := []DayMeetings{
		{Day: 3, Meetings: 1},
		{Day: 4, Meetings: 3},
		{Day: 14, Meetings: 3},
		{Day: 21, Meetings: 2},
	}

	fmt.Println(optimalVacation(meetings, 30, 7)) // [5 0]
	fmt.Println(optimalVacation(meetings, 10, 7)) // [4 3] — все старты 1..4 задевают встречи
}
