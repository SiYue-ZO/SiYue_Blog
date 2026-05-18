package logger

import (
	"os"
	"time"

	"github.com/rs/zerolog"
)

func NewLogger(level, format string) zerolog.Logger {
	var l zerolog.Level
	switch level {
	case "debug":
		l = zerolog.DebugLevel
	case "info":
		l = zerolog.InfoLevel
	case "warn":
		l = zerolog.WarnLevel
	case "error":
		l = zerolog.ErrorLevel
	default:
		l = zerolog.InfoLevel
	}

	zerolog.SetGlobalLevel(l)

	var logger zerolog.Logger
	if format == "json" {
		logger = zerolog.New(os.Stderr).With().Timestamp().Logger()
	} else {
		output := zerolog.ConsoleWriter{
			Out:        os.Stderr,
			TimeFormat: time.RFC3339,
		}
		logger = zerolog.New(output).With().Timestamp().Logger()
	}

	return logger
}
