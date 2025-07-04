package logger

import (
	"io"
	"os"
	"silent-sort/internal/config"
	"time"

	"github.com/rs/zerolog"
	"gopkg.in/natefinch/lumberjack.v2"
)

var log zerolog.Logger

func Init(cfg *config.Config) {
	var writers []io.Writer

	if cfg.LogToConsole {
		writers = append(writers, zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339})
	}

	if cfg.LogFile != "" {
		fileLogger := &lumberjack.Logger{
			Filename: cfg.LogFile,
			Compress: true,
		}
		writers = append(writers, fileLogger)
	}

	mw := io.MultiWriter(writers...)

	log = zerolog.New(mw).With().Timestamp().Logger()
	log.Fatal()
}

func GetLogger() zerolog.Logger {
	return log
}

func Error() *zerolog.Event {
	return log.Error()
}

func Fatal() *zerolog.Event {
	return log.Fatal()
}

func Info() *zerolog.Event {
	return log.Info()
}

func Debug() *zerolog.Event {
	return log.Debug()
}
