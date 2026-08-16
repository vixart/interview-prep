// Определение чемпионов: участвовал во ВСЕХ днях + максимальная сумма шагов.
//
// Одним проходом копим по каждому участнику количество дней участия и сумму
// шагов, затем среди участвовавших во всех днях выбираем максимум суммы.
// Время O(n), память O(участников).
package main

import (
	"fmt"
	"sort"
)

type UserSteps struct {
	UserId int
	Steps  int
}

type Champions struct {
	UserIds []int
	Steps   int
}

func findChampions(statistics [][]UserSteps) Champions {
	type agg struct {
		days  int
		steps int
	}

	users := make(map[int]*agg)
	for _, day := range statistics {
		for _, us := range day {
			a, ok := users[us.UserId]
			if !ok {
				a = &agg{}
				users[us.UserId] = a
			}
			a.days++
			a.steps += us.Steps
		}
	}

	best := Champions{}
	for id, a := range users {
		if a.days != len(statistics) {
			continue // участвовал не во всех днях
		}
		switch {
		case a.steps > best.Steps:
			best = Champions{UserIds: []int{id}, Steps: a.steps}
		case a.steps == best.Steps:
			best.UserIds = append(best.UserIds, id)
		}
	}

	sort.Ints(best.UserIds) // детерминированный порядок (map не упорядочен)
	return best
}

func main() {
	fmt.Println(findChampions([][]UserSteps{
		{{UserId: 1, Steps: 1000}, {UserId: 2, Steps: 1500}},
		{{UserId: 2, Steps: 1000}},
	})) // {[2] 2500}

	fmt.Println(findChampions([][]UserSteps{
		{{UserId: 1, Steps: 2000}, {UserId: 2, Steps: 1500}},
		{{UserId: 1, Steps: 3500}, {UserId: 2, Steps: 4000}},
	})) // {[1 2] 5500}
}
