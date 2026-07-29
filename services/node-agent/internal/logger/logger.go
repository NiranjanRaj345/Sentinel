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

func (l *Logger) logMessage(level Level, format string, args ...any) {
	if level < l.level {
		return
	}

	prefix := "[" + level.String() + "] "
	l.log.Printf(prefix+format, args...)
}

func (l *Logger) Debug(format string, args ...any) {
	l.logMessage(Debug, format, args...)
}

func (l *Logger) Info(format string, args ...any) {
	l.logMessage(Info, format, args...)
}

func (l *Logger) Warn(format string, args ...any) {
	l.logMessage(Warn, format, args...)
}

func (l *Logger) Error(format string, args ...any) {
	l.logMessage(Error, format, args...)
}
