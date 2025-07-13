package logger

import (
	"os"

	"github.com/rs/zerolog"
)

func New() zerolog.Logger {
	logger := zerolog.New(os.Stdout).With().Timestamp().Logger()

	if os.Getenv("LOG_FORMAT") != "json" {
		logger = logger.Output(zerolog.ConsoleWriter{Out: os.Stdout})
	}

	return logger
}
