// Псевдоним импорта: два пакета называются rand (crypto/rand и math/rand),
// поэтому одному дают имя crand. Тот же прием нужен, когда имя пакета
// не совпадает с именем каталога или конфликтует с локальной переменной.
package main

import (
	crand "crypto/rand"
	// псевдоним: иначе два rand в одном файле не ужились бы
	"encoding/binary"
	"fmt"
	"math/rand"
)

func main() {
	r := seedRand()
	fmt.Println(r.Int())
}

func seedRand() *rand.Rand {
	var b [8]byte
	_, err := crand.Read(b[:])
	// криптостойкие байты...
	if err != nil {
		panic("cannot seed with cryptographic random number generator")
	}
	r := rand.New(rand.NewSource(int64(binary.LittleEndian.Uint64(b[:]))))
	// ...ими засеваем быстрый math/rand
	return r
}
