// Функция, производительность которой измеряется: размер буфера чтения — параметр,
// и бенчмарк показывает, как он влияет на скорость.
package bench

import "os"

func FileLen(f string, bufsize int) (int, error) {
	file, err := os.Open(f)
	if err != nil {
		return 0, err
	}
	defer file.Close()
	count := 0
	for {
		buf := make([]byte, bufsize)
		// аллокация ВНУТРИ цикла — ее и покажет -benchmem как allocs/op
		num, err := file.Read(buf)
		// чем меньше bufsize, тем больше системных вызовов — это и меряет бенчмарк
		count += num
		if err != nil {
			break
		}
	}
	return count, nil
}
