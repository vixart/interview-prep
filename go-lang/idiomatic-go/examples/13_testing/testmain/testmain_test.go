// TestMain: единственная точка настройки и очистки НА ПАКЕТ.
// Все, что до m.Run() — подготовка, после — уборка; код обязан завершиться os.Exit(exitVal).
// Нужен, когда тестам требуется общий ресурс (БД, переменные уровня пакета).
package testmain

import (
	"fmt"
	"os"
	"testing"
	"time"
)

var testTime time.Time

func TestMain(m *testing.M) {
	fmt.Println("Set up stuff for tests here")
	testTime = time.Now()
	exitVal := m.Run()
	// m.Run() запускает ВСЕ тесты пакета; до него — настройка, после — уборка
	fmt.Println("Clean up stuff after tests here")
	os.Exit(exitVal)
	// код возврата обязателен — иначе go test не узнает результат
}

func TestFirst(t *testing.T) {
	fmt.Println("TestFirst uses stuff set up in TestMain", testTime)
}

func TestSecond(t *testing.T) {
	fmt.Println("TestSecond also uses stuff set up in TestMain", testTime)
}
