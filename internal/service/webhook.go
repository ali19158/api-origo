package service

import (
	"fmt"
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

func (s *WebhookService) SendTelegram(title, project, level, status, shortID, count, permalink string) error {
	token := s.telegramToken
	chatID := s.telegramChatID

	if token == "" || chatID == "" {
		return fmt.Errorf("TELEGRAM_TOKEN or TELEGRAM_CHAT_ID not set")
	}

	msg := fmt.Sprintf(
		"🔔 *%s*\n\nПроект: %s\nУровень: %s\nСтатус: %s\nID: %s\nКол-во: %s\n\n🔗 %s",
		title, project, level, status, shortID, count, permalink,
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

	return nil
}
