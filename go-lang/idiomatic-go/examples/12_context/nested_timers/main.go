// Вложенные дедлайны: дочерний контекст просит 3 секунды, родительский дает 2 —
// побеждает более короткий. Дочерний контекст не может продлить жизнь родительского.
package main

import (
	"context"
	"fmt"
	"time"
)

func main() {
	ctx := context.Background()
	parent, cancel := context.WithTimeout(ctx, 2*time.Second)
	// родитель дает 2 секунды
	defer cancel()
	child, cancel2 := context.WithTimeout(parent, 3*time.Second)
	// ребенок просит 3 — но продлить жизнь родителя не может
	defer cancel2()
	start := time.Now()
	<-child.Done()
	// сработает через 2 секунды: побеждает более короткий дедлайн
	end := time.Now()
	fmt.Println(end.Sub(start).Truncate(time.Second))
}
