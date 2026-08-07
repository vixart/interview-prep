package main

import (
	"errors"
	"fmt"
)

var (
	ErrPoolClosed      = errors.New("workerpool: пул останавливается")
	ErrAlreadyShutdown = errors.New("workerpool: shutdown уже был вызван")
	ErrNilTask         = errors.New("workerpool: пустая функция задачи")
)

// PanicError - паника задачи, восстановленная воркером. Доходит до
// потребителя через Errors(); извлекается через errors.As(err, &pe).
type PanicError struct {
	// Паника задачи доезжает до потребителя как обычная ошибка со стеком —
	// так ее видно в логах и алертах, а не только в упавшем процессе.
	Recovered any
	Stack     []byte
}

func (e *PanicError) Error() string {
	return fmt.Sprintf("workerpool: паника в задаче: %v", e.Recovered)
}
