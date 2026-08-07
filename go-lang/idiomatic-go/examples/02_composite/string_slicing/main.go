// Индекс и срез строки работают в БАЙТАХ: s[6] дает byte (число), s[4:7] — подстроку.
// Для ASCII это безопасно, но на многобайтовой руне такой срез разрежет символ пополам.
package main

import "fmt"

func main() {
	var s string = "Hello there"
	var b byte = s[6]
	// индексация строки дает БАЙТ (число), а не символ
	fmt.Println(b)
	var s2 string = s[4:7]
	// срез строки тоже в байтах — на многобайтовой руне разрезал бы символ
	var s3 string = s[:5]
	var s4 string = s[6:]
	fmt.Println(s2)
	fmt.Println(s3)
	fmt.Println(s4)
}
