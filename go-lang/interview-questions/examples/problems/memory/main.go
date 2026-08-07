// Проблемы с памятью: четыре самых частых сюжета на собеседовании.
//
//  1. лишние аллокации из-за роста среза без предварительного make;
//  2. escape analysis: почему указатель уводит объект в кучу;
//  3. срез/подстрока держат ВЕСЬ исходный массив («утечка памяти без утечки»);
//  4. sync.Pool для переиспользования буферов.
//
// Посмотреть, что именно уезжает в кучу:
//
//	go build -gcflags='-m' ./problems/memory
//
// Посмотреть аллокации в цифрах:
//
//	go test -bench=. -benchmem ./problems/memory
package main

import (
	"fmt"
	"runtime"
	"strings"
	"sync"
)

// --- 1. Аллокации при росте среза ---

func growNaive(n int) []int {
	var out []int // cap 0
	for i := 0; i < n; i++ {
		out = append(out, i) // при каждом исчерпании cap: новый массив + копирование
	}
	return out
}

func growPrealloc(n int) []int {
	out := make([]int, 0, n) // одна аллокация на весь цикл
	for i := 0; i < n; i++ {
		out = append(out, i)
	}
	return out
}

// --- 2. Escape analysis ---

type point struct{ X, Y int }

// stays on stack: значение не покидает функцию.
func onStack() int {
	p := point{1, 2}
	return p.X + p.Y
}

// escapes to heap: возвращаем указатель на локальную переменную,
// компилятор обязан разместить ее в куче.
func onHeap() *point {
	p := point{1, 2}
	return &p
}

// --- 3. Срез держит исходный массив ---

// bad возвращает 10 элементов, но не дает освободить массив на 10 миллионов:
// подсрез ссылается на тот же backing array.
func keepsWholeArray() []int {
	big := make([]int, 10_000_000)
	return big[:10]
}

// good копирует нужное — исходный массив становится мусором и будет собран.
func copiesWhatIsNeeded() []int {
	big := make([]int, 10_000_000)
	out := make([]int, 10)
	copy(out, big[:10])
	return out
}

// То же самое со строками: s[:10] от мегабайтной строки держит весь буфер.
// Лечение — strings.Clone(s[:10]).

// --- 4. sync.Pool ---

var bufPool = sync.Pool{
	New: func() any { return new(strings.Builder) },
}

func withPool(n int) string {
	b := bufPool.Get().(*strings.Builder)
	defer func() {
		b.Reset() // обязательно чистим перед возвратом
		bufPool.Put(b)
	}()
	for i := 0; i < n; i++ {
		b.WriteString("x")
	}
	return b.String()
}

func heapMB() float64 {
	var ms runtime.MemStats
	runtime.ReadMemStats(&ms)
	return float64(ms.HeapAlloc) / 1024 / 1024
}

func main() {
	fmt.Printf("старт, куча: %.1f МБ\n", heapMB())

	_ = growNaive(1_000_000)
	fmt.Printf("после growNaive:    %.1f МБ (много промежуточных массивов)\n", heapMB())
	runtime.GC()

	_ = growPrealloc(1_000_000)
	fmt.Printf("после growPrealloc: %.1f МБ\n", heapMB())
	runtime.GC()

	small := keepsWholeArray()
	runtime.GC()
	mb := heapMB()
	runtime.KeepAlive(small)
	// KeepAlive обязателен: иначе компилятор посчитает small мертвым до GC,
	// массив соберется и демонстрация развалится
	fmt.Printf("держим срез из %d элементов: %.1f МБ — жив весь массив на 10 млн\n",
		len(small), mb)

	copied := copiesWhatIsNeeded()
	runtime.GC()
	mb = heapMB()
	runtime.KeepAlive(copied)
	fmt.Printf("после copy той же длины %d:  %.1f МБ — исходный массив собран\n",
		len(copied), mb)

	fmt.Println("sync.Pool:", len(withPool(5)), "символов")
	fmt.Println("escape analysis:", onStack(), onHeap().X)

	// Чек-лист на собеседовании:
	// - знаешь размер — делай make с cap (срезы, map, strings.Builder);
	// - не возвращай указатель на локальную структуру без нужды;
	// - подсрез/подстрока держат весь буфер → copy или strings.Clone;
	// - sync.Pool только для ДОЛГОЖИВУЩИХ дорогих объектов и обязательно с Reset;
	// - профиль памяти: go tool pprof -alloc_space (сколько выделили всего)
	//   и -inuse_space (сколько занято сейчас).
}
