// t.Parallel: подтесты выполняются конкурентно друг с другом.
// Закомментированная строка `d := d` была обязательна до Go 1.22 — иначе все
// параллельные подтесты видели последнее значение переменной цикла. Сейчас не нужна.
// Осторожно: параллельные тесты нельзя совмещать с t.Setenv.
package parallel

import (
	"fmt"
	"testing"
)

func TestParallelTable(t *testing.T) {
	data := []struct {
		name   string
		input  int
		output int
	}{
		{"a", 10, 20},
		{"b", 30, 40},
		{"c", 50, 60},
	}
	for _, d := range data {
		// d := d //UNCOMMENT THIS LINE SO THAT d IS SHADOWED AND THE TEST WORKS!
		t.Run(d.name, func(t *testing.T) {
			t.Parallel()
			// подтест сообщает: меня можно выполнять параллельно с другими такими же
			fmt.Println(d.input, d.output)
			out := toTest(d.input)
			if out != d.output {
				t.Error("didn't match", out, d.output)
			}
		})
	}
}
