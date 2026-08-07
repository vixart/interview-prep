// go-cmp вместо ручных сравнений: cmp.Diff печатает читаемую разницу.
// Нестабильное поле глушится либо cmpopts.IgnoreFields, либо своим cmp.Comparer.
package cmp

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func TestCreatePerson(t *testing.T) {
	expected := Person{
		Name: "Dennis",
		Age:  37,
	}
	result := CreatePerson("Dennis", 37)
	// Без опций cmp.Diff сравнил бы и DateAdded — тест всегда падал бы.
	if diff := cmp.Diff(expected, result, cmpopts.IgnoreFields(Person{}, "DateAdded")); diff != "" {
		t.Error(diff)
	}
}

func TestCreatePersonIgnoreDate(t *testing.T) {
	expected := Person{
		Name: "Dennis",
		Age:  37,
	}
	result := CreatePerson("Dennis", 37)
	comparer := cmp.Comparer(func(x, y Person) bool {
		// второй способ: своя функция сравнения вместо игнорирования поля
		return x.Name == y.Name && x.Age == y.Age
	})
	if diff := cmp.Diff(expected, result, comparer); diff != "" {
		t.Error(diff)
	}
	if result.DateAdded.IsZero() {
		// поле все равно проверяем — просто отдельно и по смыслу
		t.Error("DateAdded wasn't assigned")
	}
}
