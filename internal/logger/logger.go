package logger

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

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
	full := fmt.Sprintf(msg, args...)
	log.Printf("[INFO] %s", full)
	l.sentry.Info().Emit(full)
}

func (l *Logger) Warn(msg string, args ...any) {
	full := fmt.Sprintf(msg, args...)
	log.Printf("[WARN] %s", full)
	l.sentry.Warn().Emit(full)
}

func (l *Logger) Error(msg string, args ...any) {
	full := fmt.Sprintf(msg, args...)
	log.Printf("[ERROR] %s", full)
	l.sentry.Error().Emit(full)
	sentry.CaptureMessage(full)
}

func (l *Logger) Fatal(msg string, args ...any) {
	full := fmt.Sprintf(msg, args...)
	log.Printf("[FATAL] %s", full)
	l.sentry.Fatal().Emit(full)
	sentry.CaptureMessage(full)
	sentry.Flush(2 * time.Second)
	os.Exit(1)
}
