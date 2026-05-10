package handler

import (
	"encoding/json"
	"net/http"

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
			WebURL  string `json:"web_url"`
			Project string `json:"project"`
			Level   string `json:"level"`
		} `json:"issue"`
	} `json:"data"`
}

func (h *WebhookHandler) SentryWebhook(w http.ResponseWriter, r *http.Request) {
	if r.Body == nil {
		http.Error(w, "empty body", http.StatusBadRequest)
		return
	}
	defer func() { _ = r.Body.Close() }()

	var payload SentryWebhook
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	if payload.Action != "created" {
		w.WriteHeader(http.StatusOK)
		return
	}

	issue := payload.Data.Issue
	if err := h.service.SendTelegram(issue.Title, issue.Project, issue.Level, issue.WebURL); err != nil {
		http.Error(w, "telegram error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
