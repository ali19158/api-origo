package service

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
)

type WebhookService struct {
	telegramToken  string
	telegramChatID string
}

func NewWebhookService(telegramToken, telegramChatID string) *WebhookService {
	return &WebhookService{
		telegramToken:  telegramToken,
		telegramChatID: telegramChatID,
	}
}

func (s *WebhookService) SendTelegram(title, project, level, webURL string) error {
	token := s.telegramToken
	chatID := s.telegramChatID

	if token == "" || chatID == "" {
		return fmt.Errorf("TELEGRAM_TOKEN or TELEGRAM_CHAT_ID not set")
	}

	msg := fmt.Sprintf(
		"🚨 *%s*\n\nПроект: %s\nУровень: %s\n\n🔗 %s",
		title, project, level, webURL,
	)

	apiURL := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", token)
	resp, err := http.PostForm(apiURL, url.Values{
		"chat_id":    {chatID},
		"text":       {msg},
		"parse_mode": {"Markdown"},
	})
	if err != nil {
		return err
	}
	defer func() { _ = resp.Body.Close() }()

	body, _ := io.ReadAll(resp.Body)
	log.Printf("[INFO] telegram response: %s", string(body))

	return nil
}
