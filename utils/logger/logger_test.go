package logger

import (
	"bytes"
	"encoding/json"
	"errors"
	"testing"

	"github.com/krishnaZawar/distributed-logger/logger-sdk"
	"github.com/stretchr/testify/assert"
)

type testLogEntry struct {
	Level       string                 `json:"level"`
	Timestamp   string                 `json:"timestamp"`
	ServiceName string                 `json:"service"`
	Metadata    map[string]interface{} `json:"metadata,omitempty"`
	Message     string                 `json:"message"`
}

func Test_NewLogger(t *testing.T) {
	service := "test-srv"

	logger := New(service)

	assert.NotNil(t, logger)
}

func Test_Debug(t *testing.T) {
	var buf bytes.Buffer
	service := "test-srv"
	message := "debug log"

	ls := newWithOutput(service, &buf)

	ls.Debug().Msg(message)

	logEntry := &testLogEntry{}
	err := json.Unmarshal(buf.Bytes(), &logEntry)
	assert.Nil(t, err)

	assert.Equal(t, string(logger.LevelDebug), logEntry.Level)
	assert.Equal(t, service, logEntry.ServiceName)
	assert.Equal(t, message, logEntry.Message)
}

func Test_Info(t *testing.T) {
	var buf bytes.Buffer
	service := "test-srv"
	message := "info log"

	ls := newWithOutput(service, &buf)

	ls.Info().Msg(message)

	logEntry := &testLogEntry{}
	err := json.Unmarshal(buf.Bytes(), &logEntry)
	assert.Nil(t, err)

	assert.Equal(t, string(logger.LevelInfo), logEntry.Level)
	assert.Equal(t, service, logEntry.ServiceName)
	assert.Equal(t, message, logEntry.Message)
}

func Test_Warn(t *testing.T) {
	var buf bytes.Buffer
	service := "test-srv"
	message := "warn log"

	ls := newWithOutput(service, &buf)

	ls.Warn().Msg(message)

	logEntry := &testLogEntry{}
	err := json.Unmarshal(buf.Bytes(), &logEntry)
	assert.Nil(t, err)

	assert.Equal(t, string(logger.LevelWarn), logEntry.Level)
	assert.Equal(t, service, logEntry.ServiceName)
	assert.Equal(t, message, logEntry.Message)
}

func Test_Error(t *testing.T) {
	var buf bytes.Buffer
	service := "test-srv"
	message := "error log"

	ls := newWithOutput(service, &buf)

	ls.Error().Msg(message)

	logEntry := &testLogEntry{}
	err := json.Unmarshal(buf.Bytes(), &logEntry)
	assert.Nil(t, err)

	assert.Equal(t, string(logger.LevelError), logEntry.Level)
	assert.Equal(t, service, logEntry.ServiceName)
	assert.Equal(t, message, logEntry.Message)
}

func Test_ErrorWith(t *testing.T) {
	var buf bytes.Buffer
	service := "test-srv"
	message := "error log"
	expectedErr := errors.New("error")

	ls := newWithOutput(service, &buf)

	ls.ErrorWith(expectedErr).Msg(message)

	logEntry := &testLogEntry{}
	err := json.Unmarshal(buf.Bytes(), &logEntry)
	assert.Nil(t, err)

	assert.Equal(t, string(logger.LevelError), logEntry.Level)
	assert.Equal(t, service, logEntry.ServiceName)
	assert.Equal(t, message, logEntry.Message)
	assert.Equal(t, expectedErr.Error(), logEntry.Metadata[labelErr])
}
