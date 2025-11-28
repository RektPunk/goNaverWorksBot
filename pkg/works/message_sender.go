package works

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"goNaverWorksBot/internal/config"
)

type MessagePayload struct {
	Content Content `json:"content"`
}

type MessageSender struct {
	Client *http.Client
	TM     *TokenManager
}

func NewMessageSender(tm *TokenManager) *MessageSender {
	return &MessageSender{
		Client: &http.Client{Timeout: 15 * time.Second},
		TM:     tm,
	}
}

func _setMessagePayload(message string) MessagePayload {
	return MessagePayload{
		Content: Content{
			Type: "text",
			Text: message,
		},
	}
}

func (s *MessageSender) postPayloadWithRetry(ctx context.Context, payload MessagePayload, url string) error {
	const maxTries = 3
	delay := time.Second * 1
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal payload: %w", err)
	}
	for i := 0; i < maxTries; i++ {
		headers, err := s.TM.SetHeaders()
		if err != nil {
			return fmt.Errorf("failed to get token headers: %w", err)
		}
		req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(jsonPayload))
		if err != nil {
			return fmt.Errorf("failed to create request: %w", err)
		}
		for k, v := range headers {
			req.Header.Set(k, v)
		}
		resp, err := s.Client.Do(req)

		if err != nil {
			fmt.Printf("Attempt %d failed: HTTP request error (network/timeout): %v\n", i+1, err)
			if i == maxTries-1 {
				return fmt.Errorf("failed to post payload after %d network errors: %w", maxTries, err)
			}
			continue
		}
		defer resp.Body.Close()
		if resp.StatusCode >= 200 && resp.StatusCode < 300 {
			return nil
		}

		bodyBytes, _ := io.ReadAll(resp.Body)
		err = fmt.Errorf("API failed (status: %d). body: %s", resp.StatusCode, string(bodyBytes))
		fmt.Printf("Attempt %d failed: %v\n", i+1, err)

		if i == maxTries-1 {
			return fmt.Errorf("failed to post message after %d API errors: %w", maxTries, err)
		} else {
			fmt.Printf("Retrying in %v (Attempt %d of %d)...\n", delay, i+1, maxTries)
			time.Sleep(delay)
			delay *= 2
		}
		continue
	}
	return fmt.Errorf("FATAL: postPayloadWithRetry loop terminated unexpectedly")
}

func (s *MessageSender) PostMessage(ctx context.Context, cfg *config.Config, message string, id string) error {
	payload := _setMessagePayload(message)
	url := fmt.Sprintf(UserPostURL, cfg.BotID, id)
	return s.postPayloadWithRetry(ctx, payload, url)
}
