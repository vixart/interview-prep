// Фильтрация продавцов по городам — пересечение множеств.
//
// Строим множество интересующих городов (map[string]struct{}), затем для
// каждого продавца оставляем только города из множества. Продавцы без
// пересечения в результат не попадают. Время O(len(cities) + сумм. длина
// списков городов) вместо двойного цикла.
package main

import "fmt"

func filterSellersByCities(sellers map[int][]string, cities []string) map[int][]string {
	citySet := make(map[string]struct{}, len(cities))
	for _, c := range cities {
		citySet[c] = struct{}{}
	}

	result := make(map[int][]string)
	for id, sellerCities := range sellers {
		var matched []string
		for _, c := range sellerCities {
			if _, ok := citySet[c]; ok {
				matched = append(matched, c)
			}
		}
		if len(matched) > 0 {
			result[id] = matched
		}
	}

	return result
}

func main() {
	sellers := map[int][]string{
		1: {"Москва", "Самара", "Ростов"},
		2: {"Москва", "Самара", "Ростов", "Казань", "Курган", "Пенза"},
		3: {"Самара", "Ростов"},
		4: {"Москва", "Казань", "Тула"},
	}
	cities := []string{"Москва", "Казань", "Тула"}

	fmt.Println(filterSellersByCities(sellers, cities))
	// map[1:[Москва] 2:[Москва Казань] 4:[Москва Казань Тула]]
}
