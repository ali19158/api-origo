package logger

import (
	"context"
	"log"

	"github.com/getsentry/sentry-go"
)

type Logger struct {
	sentry sentry.Logger
}

var Default *Logger

func Init(ctx context.Context) {
	Default = New(ctx)
}

func New(ctx context.Context) *Logger {
	return &Logger{
		sentry: sentry.NewLogger(ctx),
	}
}

func (l *Logger) Info(msg string, args ...any) {
	log.Printf("[INFO] "+msg, args...)
	l.sentry.Info().Emitf(msg, args...)
}

func (l *Logger) Warn(msg string, args ...any) {
	log.Printf("[WARN] "+msg, args...)
	l.sentry.Warn().Emitf(msg, args...)
}

func (l *Logger) Error(msg string, args ...any) {
	log.Printf("[ERROR] "+msg, args...)
	l.sentry.Error().Emitf(msg, args...)
}
