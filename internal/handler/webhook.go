package handler

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/online-shop/internal/logger"
	"github.com/online-shop/internal/service"
)

type WebhookHandler struct {
	service *service.WebhookService
}

func NewWebhookHandler(service *service.WebhookService) *WebhookHandler {
	return &WebhookHandler{service: service}
}

type SentryWebhook struct {
	Action string `json:"action"`
	Data   struct {
		Issue struct {
			Title   string `json:"title"`
			WebURL  string `json:"permalink"`
			Project struct {
				Name string `json:"name"`
				Slug string `json:"slug"`
			} `json:"project"`
			Level string `json:"level"`
		} `json:"issue"`
	} `json:"data"`
}

func (h *WebhookHandler) SentryWebhook(w http.ResponseWriter, r *http.Request) {
	if r.Body == nil {
		logger.Default.Warn("SentryWebhook Empty body")
		http.Error(w, "empty body", http.StatusBadRequest)
		return
	}
	defer func() { _ = r.Body.Close() }()

	var payload SentryWebhook
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		logger.Default.Warn("SentryWebhook payload %v, %s", err.Error(), r.Body)
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	log.Printf("[INFO] sentry webhook payload: %+v", payload)

	if payload.Action != "created" && payload.Action != "resolved" && payload.Action != "unresolved" {
		w.WriteHeader(http.StatusOK)
		return
	}

	issue := payload.Data.Issue
	if err := h.service.SendTelegram(issue.Title, issue.Project.Name, issue.Level, issue.WebURL); err != nil {
		http.Error(w, "telegram error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
