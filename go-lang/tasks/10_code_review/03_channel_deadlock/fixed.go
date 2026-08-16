// Исправленный вариант: 5 горутин пишут в канал, main читает все сообщения.
//
// Что изменилось относительно исходника:
//   - убран мьютекс вокруг отправки в канал (канал сам потокобезопасен;
//     лок вокруг блокирующей отправки только мешал);
//   - канал закрывается в отдельной горутине после wg.Wait() —
//     это сигнал читателю «сообщений больше не будет»;
//   - чтение через for range вместо вечного for/select -> нет deadlock,
//     wg.Wait() больше не недостижим;
//   - нормальные имена (wg вместо wc, messages вместо m).
package main

import (
	"fmt"
	"sync"
)

func main() {
	var wg sync.WaitGroup
	messages := make(chan string, 3)

	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			messages <- fmt.Sprintf("Gorutine %d", i)
		}(i)
	}

	go func() {
		wg.Wait()
		close(messages)
	}()

	for msg := range messages {
		fmt.Println(msg)
	}
}
