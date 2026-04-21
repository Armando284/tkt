package logger

import (
	"log/slog"
	"os"
	"strings"

	"github.com/charmbracelet/log"
)

var L *slog.Logger // logger global (igual que en JS)

func Init() {
	level := parseLevel(os.Getenv("LOG_LEVEL"))

	// Charmbracelet/log es un handler de slog con salida hermosa para terminal
	handler := log.NewWithOptions(os.Stderr, log.Options{
		Level:           log.Level(level),
		ReportTimestamp: true,
		ReportCaller:    false, // opcional: true si quieres archivo:line
	})

	L = slog.New(handler)

	// Lo ponemos como default para que cualquier paquete pueda usar slog.Default()
	slog.SetDefault(L)
}

func parseLevel(lvl string) slog.Level {
	switch strings.ToLower(strings.TrimSpace(lvl)) {
	case "debug":
		return slog.LevelDebug
	case "warn", "warning":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}