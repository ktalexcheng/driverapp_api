package util

import (
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/rs/zerolog/pkgerrors"
)

// Initializes the global settings for zerolog
func InitLogger(debug bool) {
	// Global settings
	var logLevel int
	if debug {
		log.Logger = log.Output(zerolog.ConsoleWriter{
			Out:        os.Stdout,
			TimeFormat: zerolog.TimeFormatUnix,
		})

		logLevel = int(zerolog.DebugLevel)
	} else {
		logLevel = int(zerolog.InfoLevel)
	}

	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack
	zerolog.SetGlobalLevel(zerolog.Level(logLevel))
	log.Logger = log.With().Timestamp().Caller().Logger()
}
