// Метка над циклом: `continue outer` переходит к следующей итерации ВНЕШНЕГО цикла.
// Без метки continue управлял бы только внутренним циклом. Так же работает `break outer`.
package main

import "fmt"

func main() {
	samples := []string{"hello", "apple_π!"}
outer:
	// метка ставится ПЕРЕД циклом, на том же уровне отступа, что и скобки
	for _, sample := range samples {
		for i, r := range sample {
			fmt.Println(i, r, string(r))
			if r == 'l' {
				continue outer
				// без метки continue перешел бы к следующей руне, с меткой — к следующему слову
			}
		}
		fmt.Println()
	}
}
