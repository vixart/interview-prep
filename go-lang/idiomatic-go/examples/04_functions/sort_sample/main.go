// Функция как параметр: sort.Slice принимает компаратор-замыкание,
// поэтому один и тот же срез сортируется по разным полям.
package main

import (
	"fmt"
	"sort"
)

func main() {
	type Person struct {
		FirstName string
		LastName  string
		Age       int
	}

	people := []Person{
		{"Pat", "Patterson", 37},
		{"Tracy", "Bobdaughter", 23},
		{"Fred", "Fredson", 18},
	}
	fmt.Println(people)

	// sort by last name
	sort.Slice(people, func(i, j int) bool {
		// функция-компаратор передается параметром — сортировка настраивается на месте
		return people[i].LastName < people[j].LastName
	})
	fmt.Println(people)

	// sort by age
	sort.Slice(people, func(i, j int) bool {
		return people[i].Age < people[j].Age
		// тот же срез, другой компаратор — другой порядок
	})
	fmt.Println(people)
}
