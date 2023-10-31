package forge

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/kyberbits/forge/forgeutils"
)

type LogSupplementerFunc func(
	logEntry *Log,
	ctx context.Context,
	r *http.Request,
)

type LoggerErrorExtrasAdder interface {
	error
	AddExtras(extras map[string]any)
}

func NewLogger(
	channel string,
	writer io.Writer,
	supplementer LogSupplementerFunc,
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
	supplementer LogSupplementerFunc
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

	r := forgeutils.ContextGetRequest(ctx)

	entry := Log{
		Channel:   logger.channel,
		Timestamp: time.Now(),
		Level:     level,
		Message:   message,
		Extras:    extras,
		ContextID: forgeutils.ContextGetID(ctx),
	}

	for key, extra := range extras {
		err, isErr := extra.(error)
		if !isErr {
			continue
		}

		// Convert the error to a string in the "extras"
		extras[key] = err.Error()

		extrasAdder, isExtrasLogger := err.(LoggerErrorExtrasAdder)
		if !isExtrasLogger {
			continue
		}

		extrasAdder.AddExtras(extras)
	}

	// Supplement the logger if there is one
	if logger.supplementer != nil {
		logger.supplementer(&entry, ctx, r)
	}

	if err := logger.encoder.Encode(entry); err != nil {
		panic(err)
	}
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
	logger.Error(context.Background(), "Standard Library Log", map[string]interface{}{
		"log": string(b),
	})

	return len(b), nil
}

func (logger *Logger) StandardLogger() *log.Logger {
	return log.New(logger, "", 0)
}

type Log struct {
	Timestamp time.Time      `json:"@timestamp"`
	Channel   string         `json:"channel"`
	Level     string         `json:"level"`
	Message   string         `json:"message"`
	Extras    map[string]any `json:"extras"`
	ContextID string         `json:"context_id"`
}
