package config

import (
	"log"
	"time"

	"github.com/getsentry/sentry-go"
	sentryhttp "github.com/getsentry/sentry-go/http"
)

func InitSentry() *sentryhttp.Handler {
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

	return sentryhttp.New(sentryhttp.Options{
		Repanic: true,
	})
}

func FlushSentry() {
	sentry.Flush(2 * time.Second)
}
