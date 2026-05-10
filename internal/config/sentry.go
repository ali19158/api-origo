package config

import (
	"log"
	"time"

	"github.com/getsentry/sentry-go"
	sentryhttp "github.com/getsentry/sentry-go/http"
)

type SentryConfig struct {
	Handler       *sentryhttp.Handler
	WebhookSecret string
}

func InitSentry() *SentryConfig {
	dsn := getEnv("SENTRY_DSN", "")
	if dsn == "" {
		log.Println("SENTRY_DSN not set, skipping Sentry initialization")
		return nil
	}

	err := sentry.Init(sentry.ClientOptions{
		Dsn:              dsn,
		Environment:      getEnv("APP_ENV", "production"),
		TracesSampleRate: 0.5,
		EnableLogs:       true,
	})
	if err != nil {
		log.Fatalf("sentry.Init: %s", err)
	}

	handler := sentryhttp.New(sentryhttp.Options{
		Repanic: true,
	})

	return &SentryConfig{
		Handler:       handler,
		WebhookSecret: getEnv("SENTRY_WEBHOOK_SECRET", ""),
	}
}

func FlushSentry() {
	sentry.Flush(2 * time.Second)
}
