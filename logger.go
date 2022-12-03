package forge

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"time"
)

const (
	LevelCritical = "CRITICAL"
	LevelError    = "ERROR"
	LevelWarning  = "WARNING"
	LevelInfo     = "INFO"
	LevelDebug    = "DEBUG"
)

type LogSupplementer func(
	logEntry *LogEntry,
	ctx context.Context,
	r *http.Request,
)

func NewLogger(
	channel string,
	writer io.Writer,
	supplementer LogSupplementer,
) *Logger {
	return &Logger{
		channel:      channel,
		encoder:      json.NewEncoder(writer),
		supplementer: supplementer,
	}
}

type Logger struct {
	channel      string
	encoder      *json.Encoder
	supplementer LogSupplementer
}

func (logger *Logger) Copy(newChannel string) *Logger {
	return &Logger{
		channel:      newChannel,
		encoder:      logger.encoder,
		supplementer: logger.supplementer,
	}
}

func (logger *Logger) Log(
	ctx context.Context,
	level string,
	message string,
	extras map[string]any,
) {
	// Make sure extras is initialized
	if extras == nil {
		extras = map[string]any{}
	}

	r := getContextRequest(ctx)

	entry := LogEntry{
		Channel:   logger.channel,
		Timestamp: time.Now(),
		Level:     level,
		Message:   message,
		Extras:    extras,
		ContextID: getContextID(ctx),
	}

	// Supplement the logger if there is one
	if logger.supplementer != nil {
		logger.supplementer(&entry, ctx, r)
	}

	logger.encoder.Encode(entry)
}

func (logger *Logger) Critical(
	cxt context.Context,
	message string,
	extras map[string]any,
) {
	logger.Log(cxt, LevelCritical, message, extras)
}

func (logger *Logger) Error(
	cxt context.Context,
	message string,
	extras map[string]any,
) {
	logger.Log(cxt, LevelError, message, extras)
}

func (logger *Logger) Warning(
	cxt context.Context,
	message string,
	extras map[string]any,
) {
	logger.Log(cxt, LevelWarning, message, extras)
}

func (logger *Logger) Info(
	cxt context.Context,
	message string,
	extras map[string]any,
) {
	logger.Log(cxt, LevelInfo, message, extras)
}

func (logger *Logger) Debug(
	cxt context.Context,
	message string,
	extras map[string]any,
) {
	logger.Log(cxt, LevelDebug, message, extras)
}

func (logger *Logger) Write(b []byte) (int, error) {
	logger.Error(nil, "Standard Library Log", map[string]interface{}{
		"log": string(b),
	})

	return len(b), nil
}

func (logger *Logger) StandardLogger() *log.Logger {
	return log.New(logger, "", 0)
}

type LogEntry struct {
	Timestamp time.Time      `json:"@timestamp"`
	Channel   string         `json:"channel"`
	Level     string         `json:"level"`
	Message   string         `json:"message"`
	Extras    map[string]any `json:"extras"`
	ContextID string         `json:"context_id"`
}
