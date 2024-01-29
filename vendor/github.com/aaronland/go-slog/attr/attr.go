package attr

import (
	"log/slog"
)

// Exported constants from a custom logging package.
const (
	LevelTrace     = slog.Level(-8)
	LevelDebug     = slog.LevelDebug
	LevelInfo      = slog.LevelInfo
	LevelNotice    = slog.Level(2)
	LevelWarning   = slog.LevelWarn
	LevelError     = slog.LevelError
	LevelEmergency = slog.Level(12)
)
