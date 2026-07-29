package logger

import (
	"io"
	"log"
)

type Logger struct {
	level Level
	log   *log.Logger
}

func New(level Level, out io.Writer) *Logger {
	return &Logger{
		level: level,
		log:   log.New(out, "", log.LstdFlags|log.LUTC),
	}
}
