package server

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"goNaverWorksBot/internal/config"
	"goNaverWorksBot/pkg/works"
)

func WebhookHandler(w http.ResponseWriter, r *http.Request, cfg *config.Config, sender *works.MessageSender) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()

	var webhookRequest works.WebhookRequest
	if err := json.Unmarshal(body, &webhookRequest); err != nil {
		log.Printf("ERROR: Failed to unmarshal JSON: %v. Body: %s", err, string(body))
		http.Error(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}

	if webhookRequest.Type == "message" && webhookRequest.Content.Type == "text" {
		responseMessage := fmt.Sprintf("received: %s", webhookRequest.Content.Text)

		targetID := webhookRequest.Source.UserID // Only respond to the user in 1-on-1 chats
		go func() {
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			if err := sender.PostMessage(ctx, cfg, responseMessage, targetID); err != nil {
				log.Printf("ERROR: Failed send messages: %v", err)
			}
		}()
	}
	w.WriteHeader(http.StatusOK)
}
