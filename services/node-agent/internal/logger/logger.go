package logger

import (
	"io"
	"log"
)

type Logger struct {
	level     Level
	log       *log.Logger
	component string
}

func New(level Level, out io.Writer) *Logger {
	return &Logger{
		level: level,
		log:   log.New(out, "", log.LstdFlags|log.LUTC),
	}
}

// Component returns a new logger that shares the same configuration
// but prefixes every log message with the given component name.
func (l *Logger) Component(name string) *Logger {
	return &Logger{
		level:     l.level,
		log:       l.log,
		component: name,
	}
}

func (l *Logger) logMessage(level Level, format string, args ...any) {
	if level < l.level {
		return
	}

	prefix := "[" + level.String() + "]"

	if l.component != "" {
		prefix += " [" + l.component + "]"
	}

	prefix += " "

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

func NewFromString(level string, out io.Writer) (*Logger, error) {
	parsed, err := ParseLevel(level)
	if err != nil {
		return nil, err
	}

	return New(parsed, out), nil
}
