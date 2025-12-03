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
	"goNaverWorksBot/internal/db"
	"goNaverWorksBot/pkg/works"
)

const UserRole = "user"
const AssistantRole = "assistant"

func WebhookHandler(w http.ResponseWriter, r *http.Request, cfg *config.Config, sender *works.MessageSender, history *db.HistoryRepository) {
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
		userID := webhookRequest.Source.UserID // Only respond to the user in 1-on-1 chats
		userMessage := webhookRequest.Content.Text
		go func() {
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			if err := history.SaveAndLimitChatHistory(ctx, userID, UserRole, userMessage); err != nil {
				log.Printf("ERROR: Failed to save user chat history for %s: %v", userID, err)
			}
			messages, err := history.GetRecentChatHistory(ctx, userID)
			if err != nil {
				log.Printf("ERROR: Failed to get recent chat history for %s: %v", userID, err)
			}
			finalResponseMessage := fmt.Sprintf(
				"chat history count: (%d) saved",
				len(messages),
			)

			if err := sender.PostMessage(ctx, cfg, finalResponseMessage, userID); err != nil {
				log.Printf("ERROR: Failed send messages: %v", err)
			}

			if err := history.SaveAndLimitChatHistory(ctx, userID, AssistantRole, finalResponseMessage); err != nil {
				log.Printf("ERROR: Failed to save user chat history for %s: %v", userID, err)
			}
		}()
	}
	w.WriteHeader(http.StatusOK)
}
