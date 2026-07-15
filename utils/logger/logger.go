package logger

import (
	"io"
	"os"

	"github.com/krishnaZawar/distributed-logger/logger-sdk"
)

const (
	labelErr = "error"
)

// Logger is used for logging across both the editor and the builder
type Logger struct {
	ls *logger.Logger
}

func New(serviceName string) *Logger {
	return &Logger{
		logger.New(serviceName, os.Stdout),
	}
}

func newWithOutput(serviceName string, writer io.Writer) *Logger {
	return &Logger{
		logger.New(serviceName, writer),
	}
}

// initiates a debug level log
func (logger *Logger) Debug() *logger.LogEvent {
	return logger.ls.Debug()
}

// initiates a info level log
func (logger *Logger) Info() *logger.LogEvent {
	return logger.ls.Info()
}

// initiates a warn level log
func (logger *Logger) Warn() *logger.LogEvent {
	return logger.ls.Warn()
}

// initiates a error level log
func (logger *Logger) Error() *logger.LogEvent {
	return logger.ls.Error()
}

// initiates an error log with the error as a separate field by default
func (logger *Logger) ErrorWith(err error) *logger.LogEvent {
	return logger.ls.Error().WithMetadata(labelErr, err.Error())
}
