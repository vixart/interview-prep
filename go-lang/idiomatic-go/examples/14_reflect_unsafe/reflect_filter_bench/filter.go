// Цена рефлексии: один и тот же Filter в трех версиях — на рефлексии,
// на дженериках и монолитом для конкретного типа. Бенчмарк показывает,
// что рефлексия на порядок медленнее и аллоцирует.
// Запуск: go test -bench=. -benchmem ./14_reflect_unsafe/reflect_filter_bench
package filter

import "reflect"

func FilterReflection(slice interface{}, filter interface{}) interface{} {
	sv := reflect.ValueOf(slice)
	fv := reflect.ValueOf(filter)

	sliceLen := sv.Len()
	out := reflect.MakeSlice(sv.Type(), 0, sliceLen)
	// аналог make() через рефлексию
	for i := 0; i < sliceLen; i++ {
		curVal := sv.Index(i)
		values := fv.Call([]reflect.Value{curVal})
		// каждый вызов фильтра идет через рефлексию — отсюда и тормоза
		if values[0].Bool() {
			out = reflect.Append(out, curVal)
		}
	}
	return out.Interface()
}

// FilterGeneric filters values from a slice using a filter function.
// It returns a new slice with only the elements of s
// for which f returned true.
func FilterGeneric[T any](s []T, f func(T) bool) []T {
	// то же самое на дженериках: типобезопасно и без рантайм-накладных
	out := make([]T, 0, len(s))
	for _, v := range s {
		if f(v) {
			out = append(out, v)
		}
	}
	return out
}

func FilterString(s []string, f func(string) bool) []string {
	out := make([]string, 0, len(s))
	for _, v := range s {
		if f(v) {
			out = append(out, v)
		}
	}
	return out
}

func FilterInt(s []int, f func(int) bool) []int {
	out := make([]int, 0, len(s))
	for _, v := range s {
		if f(v) {
			out = append(out, v)
		}
	}
	return out
}
