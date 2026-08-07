// switch: несколько значений в case, пустой case (ничего не делать), default.
// missingLabel показывает ловушку: break внутри case выходит из switch, а не из цикла;
// labeledBreak — правильное решение через метку над for.
package main

import "fmt"

func main() {
	basicSwitch()
	fmt.Println()
	missingLabel()
	fmt.Println()
	labeledBreak()
}

func basicSwitch() {
	fmt.Println("basic switch statement example")
	words := []string{"a", "cow", "smile", "gopher",
		"octopus", "anthropologist"}
	for _, word := range words {
		switch size := len(word); size {
		case 1, 2, 3, 4:
			// несколько значений в одной ветви
			fmt.Println(word, "is a short word!")
		case 5:
			wordLen := len(word)
			fmt.Println(word, "is exactly the right length:", wordLen)
		case 6, 7, 8, 9:
		// пустая ветвь = «ничего не делать» (провала на default не будет)
		default:
			fmt.Println(word, "is a long word!")
		}
	}
}

func missingLabel() {
	fmt.Println("the case of the missing label...")
	for i := 0; i < 10; i++ {
		switch {
		case i%2 == 0:
			fmt.Println(i, "is even")
		case i%3 == 0:
			fmt.Println(i, "is divisible by 3 but not 2")
		case i%7 == 0:
			fmt.Println("exit the loop!")
			break
			// ЛОВУШКА: этот break выходит из switch, цикл продолжается
		default:
			fmt.Println(i, "is boring")
		}
	}
}

func labeledBreak() {
	fmt.Println("the label has been added!")
	// метка ставится перед циклом
loop:
	for i := 0; i < 10; i++ {
		switch {
		case i%2 == 0:
			fmt.Println(i, "is even")
		case i%3 == 0:
			fmt.Println(i, "is divisible by 3 but not 2")
		case i%7 == 0:
			fmt.Println("exit the loop!")
			break loop
			// ...и break с меткой — теперь выходим действительно из цикла
		default:
			fmt.Println(i, "is boring")
		}
	}
}
