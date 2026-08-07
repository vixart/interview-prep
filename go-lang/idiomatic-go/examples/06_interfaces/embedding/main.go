// Встраивание = композиция: поля и методы Employee доступны прямо на Manager.
// Но это НЕ наследование: присвоить Manager переменной типа Employee нельзя
// (строка закомментирована), нужно явно взять m.Employee.
package main

import (
	"fmt"
)

type Employee struct {
	Name string
	ID   string
}

func (e Employee) Description() string {
	return fmt.Sprintf("%s (%s)", e.Name, e.ID)
}

type Manager struct {
	Employee
	// встроенное поле без имени — так делается композиция:
	// поля и методы Employee становятся доступны прямо на Manager
	Reports []Employee
}

func (m Manager) FindNewEmployees() []Employee {
	// do business logic
	return nil
}

func main() {
	m := Manager{
		Employee: Employee{
			Name: "Bob Bobson",
			ID:   "12345",
		},
		Reports: []Employee{},
	}
	// Встраивание — не наследование: Manager нельзя присвоить Employee.
	// var eFail Employee = m // ошибка компиляции        // compilation error!
	var eOK Employee = m.Employee // ok!
	// к встроенной структуре обращаемся по имени типа
	fmt.Println(eOK.Description())
	// метод Description поднят на Manager: можно звать и m.Description()
}
