package log

import (
	"os"

	"github.com/rs/zerolog"
)

var Log = MustInitLogger()

type Logger struct {
	logger     *zerolog.Logger
	infoLevel  zerolog.Level
	errorLevel zerolog.Level
}

func (l *Logger) Info(s string) {
	l.logger.WithLevel(l.infoLevel).Caller(1).Msg(s)
}

func (l *Logger) Error(s error) {
	l.logger.WithLevel(l.errorLevel).Caller(1).Msg(s.Error())
}

func MustInitLogger() *Logger {
	zerolog.TimeFieldFormat = "15:04:05 02-01-2006"

	infoLevel := zerolog.InfoLevel
	errorLevel := zerolog.ErrorLevel

	logger := zerolog.New(os.Stdout).With().Timestamp().Logger()

	log := &Logger{
		logger:     &logger,
		infoLevel:  infoLevel,
		errorLevel: errorLevel,
	}

	return log
}
