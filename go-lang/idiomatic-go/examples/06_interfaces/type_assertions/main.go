// Утверждение типа: правильный вариант, две паники (неверный тип и «похожий, но не идентичный»
// тип MyInt vs int) и идиома «запятая-ok». Паники перехвачены recover, чтобы увидеть все случаи.
// Вывод: утверждение проверяется в рантайме, всегда используй запятая-ok.
package main

import "fmt"

func main() {
	typeAssert()
	typeAssertPanicWrongType()
	typeAssertPanicTypeNotIdentical()
	err := typeAssertCommaOK()
	if err != nil {
		fmt.Println(err)
	}
}

func typeAssertCommaOK() error {
	var i any
	var mine MyInt = 20
	i = mine
	i2, ok := i.(int)
	// идиома «запятая-ok»: вместо паники получаем ok == false
	if !ok {
		// we are constructing a new error with fmt.Errorf.
		// fmt.Errorf is covered in chapter 9.
		return fmt.Errorf("unexpected type for %v", i)
	}
	fmt.Println(i2 + 1)
	return nil
}

func typeAssertPanicTypeNotIdentical() {
	// we are using recover to allow us to run through the
	// failing type assertions. recover is explained in chapter 9.
	defer func() {
		if m := recover(); m != nil {
			fmt.Println(m) // prints out because a panic happens
		}
	}()
	var i any
	var mine MyInt = 20
	i = mine
	i2 := i.(int)
	// ПАНИКА: базовый тип MyInt, а не int, хотя они «одинаковые» по смыслу
	fmt.Println(i2 + 1)
}

func typeAssertPanicWrongType() {
	// we are using recover to allow us to run through the
	// failing type assertions. recover is explained in chapter 9.
	defer func() {
		if m := recover(); m != nil {
			fmt.Println(m) // prints out because a panic happens
		}
	}()
	var i any
	var mine MyInt = 20
	i = mine
	i2 := i.(string)
	// ПАНИКА: тип вообще другой
	fmt.Println(i2)
}

type MyInt int

func typeAssert() {
	var i any
	var mine MyInt = 20
	i = mine
	i2 := i.(MyInt)
	// тип совпал — все хорошо
	fmt.Println(i2 + 1)
}
