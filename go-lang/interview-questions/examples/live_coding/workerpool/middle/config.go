package main

// Config задаёт параметры пула. Валидируется в New.
type Config struct {
	Workers   int
	QueueSize int
	// QueueSize 0 = небуферизованный канал: Submit ждет свободного воркера.
	// Это честное противодавление, но Submit блокируется — см. senior/OnFull.
	ErrBuf int
}

func (c Config) validate() {
	// Паника, а не error: невалидный конфиг — это ошибка программиста,
	// она обязана падать на старте, а не превращаться в рантайм-ветку.
	if c.Workers <= 0 {
		panic("workerpool: Workers должен быть > 0")
	}
	if c.QueueSize < 0 {
		panic("workerpool: QueueSize должен быть >= 0")
	}
	if c.ErrBuf < 0 {
		panic("workerpool: ErrBuf должен быть >= 0")
	}
}
