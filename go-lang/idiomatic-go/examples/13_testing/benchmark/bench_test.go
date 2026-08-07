// Бенчмарки: функция Benchmark*(b *testing.B) с циклом на b.N (число итераций
// подбирает сам фреймворк), результат кладется в пакетную переменную blackhole,
// чтобы компилятор не выбросил вызов. b.Run сравнивает размеры буфера между собой.
// Запуск: go test -bench=. -benchmem ./13_testing/benchmark
package bench

import (
	"fmt"
	"math/rand"
	"os"
	"testing"
)

var dataSize int64

func TestMain(m *testing.M) {
	dataSize = makeData()
	exitVal := m.Run()
	os.Remove("testdata/data.txt")
	os.Exit(exitVal)
}

// makeData makes our data file for us. Rather than checking in a large file, we recreate it for the test.
// By setting the random seed to the same value every time, we ensure that we generate the same file every time.
// This random seed generates a file that's 65,204 bytes long.
func makeData() int64 {
	file, err := os.Create("testdata/data.txt")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	r := rand.New(rand.NewSource(1))
	var total int64
	for i := 0; i < 10000; i++ {
		data := makeWord(r, r.Intn(10)+1)
		n, _ := file.Write(data)
		total += int64(n)
	}
	return total
}

func makeWord(r *rand.Rand, l int) []byte {
	out := make([]byte, l+1)
	for i := 0; i < l; i++ {
		out[i] = 'a' + byte(r.Intn(26))
	}
	out[l] = '\n'
	return out
}

func TestFileLen(t *testing.T) {
	result, err := FileLen("testdata/data.txt", 1)
	if err != nil {
		t.Fatal(err)
	}
	if int64(result) != dataSize {
		t.Error("Expected", dataSize, "got", result)
	}
}

var blackhole int

func BenchmarkFileLen1(b *testing.B) {
	for i := 0; i < b.N; i++ {
		// b.N подбирает фреймворк, пока замер не станет статистически значимым
		result, err := FileLen("testdata/data.txt", 1)
		if err != nil {
			b.Fatal(err)
		}
		blackhole = result
		// результат «сливаем» в пакетную переменную, чтобы компилятор не выбросил вызов
	}
}

func BenchmarkFileLen(b *testing.B) {
	for _, v := range []int{1, 10, 100, 1000, 10000, 100000} {
		b.Run(fmt.Sprintf("FileLen-%d", v), func(b *testing.B) {
			// подбенчмарки: сравниваем размеры буфера между собой
			for i := 0; i < b.N; i++ {
				result, err := FileLen("testdata/data.txt", v)
				if err != nil {
					b.Fatal(err)
				}
				blackhole = result
			}
		})
	}
}
