package nk

import (
	"os"
	"time"

	"github.com/rs/zerolog"
)

var Log zerolog.Logger

func init() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	output := zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339}

	Log = zerolog.New(output).With().Timestamp().Caller().Logger()
}

func SetLogger(logger zerolog.Logger) {
	Log = logger
}
