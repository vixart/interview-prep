package main

import (
	"context"
	"log/slog"
	"os"
	"time"
)

func main() {

	// Использование глобального логгера slog.
	// По умолчанию пишет в stderr и уровень INFO,
	// поэтому Debug может не отображаться.
	slog.Debug("debug log message")
	slog.Info("info log message")
	slog.Warn("warning log message")
	slog.Error("error log message")

	userID := "fred"
	loginCount := 20

	// Структурированное логирование:
	// ключ-значение добавляются как поля события.
	slog.Info("user login",
		"id", userID,
		"login_count", loginCount)

	// Настройка Handler:
	// Handler отвечает за формат и вывод логов.
	// LevelDebug включает вывод debug сообщений.
	options := &slog.HandlerOptions{Level: slog.LevelDebug}

	// JSONHandler — форматирует логи в JSON
	// и пишет их в указанный io.Writer (stderr).
	handler := slog.NewJSONHandler(os.Stderr, options)

	// Создание собственного логгера с кастомным handler.
	mySlog := slog.New(handler)

	// Пример значения времени для логирования.
	lastLogin := time.Date(2023, 01, 01, 11, 50, 00, 00, time.UTC)

	// Логирование через кастомный logger.
	// Поля автоматически сериализуются в JSON.
	mySlog.Debug("debug message", "id", userID, "last_login", lastLogin)

	// Context передается для возможной интеграции
	// с tracing, cancellation или middleware.
	ctx := context.Background()

	// LogAttrs — более производительный вариант логирования:
	// атрибуты создаются напрямую (без аллокаций map/interface).
	mySlog.LogAttrs(
		ctx,
		slog.LevelInfo,
		"faster logging",
		slog.String("id", userID),
		slog.Time("last_login", lastLogin),
	)

	// Адаптер для стандартного пакета log.
	// Позволяет использовать slog handler
	// через привычный интерфейс log.Logger.
	myLog := slog.NewLogLogger(mySlog.Handler(), slog.LevelDebug)

	// Сообщение пройдет через тот же handler,
	// поэтому будет форматировано так же (JSON).
	myLog.Println("using the mySlog Handler")
}
