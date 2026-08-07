// Отключение ветви select: закрытый канал всегда «готов» и вернет нулевое значение,
// поэтому цикл крутился бы вхолостую. Присваивание `in = nil` делает ветвь вечно неготовой,
// а счетчик закрытых каналов дает условие выхода.
package main

import "fmt"

func main() {
	in1 := make(chan int)
	in2 := make(chan int)
	go func() {
		for i := 10; i < 100; i += 10 {
			in1 <- i
		}
		close(in1)
	}()
	go func() {
		for i := 20; i >= 0; i-- {
			in2 <- i
		}
		close(in2)
	}()
	result := readFromTwoChannels(in1, in2)
	fmt.Println(result)
}

func readFromTwoChannels(in, in2 chan int) []int {
	var out []int
	// in and in2 are channels
	for count := 0; count < 2; {
		// выходим, когда закрылись оба канала
		select {
		case v, ok := <-in:
			// ok == false означает «канал закрыт», а не «пришел 0»
			if !ok {
				in = nil // the case will never succeed again!
				// nil-канал блокируется навсегда → ветвь больше никогда не будет готова
				count++
				continue
			}
			// process the v that was read from in
			out = append(out, v)
		case v, ok := <-in2:
			if !ok {
				in2 = nil // the case will never succeed again!
				count++
				continue
			}
			// process the v that was read from in2
			out = append(out, v)
		}
	}
	return out
}
