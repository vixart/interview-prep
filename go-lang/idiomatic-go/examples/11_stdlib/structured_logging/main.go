// log/slog: уровни, пары ключ-значение, JSON-handler со своим уровнем,
// LogAttrs (быстрый путь без аллокаций на интерфейсах) и мост к старому log.Logger.
package main

import (
	"context"
	"log/slog"
	"os"
	"time"
)

func main() {
	slog.Debug("debug log message")
	slog.Info("info log message")
	slog.Warn("warning log message")
	slog.Error("error log message")

	userID := "fred"
	loginCount := 20
	slog.Info("user login",
		// после сообщения идут пары ключ-значение — их можно искать в логах машинно
		"id", userID,
		"login_count", loginCount)

	options := &slog.HandlerOptions{Level: slog.LevelDebug}
	handler := slog.NewJSONHandler(os.Stderr, options)
	// handler решает формат и уровень; для прода обычно JSON
	mySlog := slog.New(handler)
	lastLogin := time.Date(2023, 01, 01, 11, 50, 00, 00, time.UTC)
	mySlog.Debug("debug message", "id", userID, "last_login", lastLogin)

	ctx := context.Background()
	mySlog.LogAttrs(ctx, slog.LevelInfo, "faster logging", slog.String("id", userID), slog.Time("last_login", lastLogin))
	// быстрый путь: типизированные атрибуты без упаковки в any

	myLog := slog.NewLogLogger(mySlog.Handler(), slog.LevelDebug)
	myLog.Println("using the mySlog Handler")
}
