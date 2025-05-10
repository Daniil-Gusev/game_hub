package core

import (
	"io"
	"log"
)

type Logger interface {
	Printf(format string, v ...any)
	Error(err error)
}

type StdLogger struct {
	logger       *log.Logger
	errorHandler ErrorHandler
}

func NewStdLogger(output io.Writer, errorHandler ErrorHandler) *StdLogger {
	return &StdLogger{
		logger:       log.New(output, "", 0),
		errorHandler: errorHandler,
	}
}

func (l *StdLogger) Printf(format string, v ...interface{}) {
	l.logger.Printf(format, v...)
}

func (l *StdLogger) Error(err error) {
	if err != nil {
		text := l.errorHandler.Handle(err)
		if text != "" {
			l.logger.Println(text)
		}
	}
}
