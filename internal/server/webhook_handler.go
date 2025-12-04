package server

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/openai/openai-go/v3"

	"goNaverWorksBot/internal/config"
	"goNaverWorksBot/internal/db"
	"goNaverWorksBot/pkg/works"
)

func NewHandler(
	cfg *config.Config,
	sender *works.MessageSender,
	historyRepo *db.HistoryRepository,
	openaiClient *openai.Client,
) *Handler {
	return &Handler{
		Cfg:          cfg,
		Sender:       sender,
		HistoryRepo:  historyRepo,
		OpenAIClient: openaiClient,
	}
}

func (h *Handler) WebhookHandler(w http.ResponseWriter, r *http.Request) {
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
		go h.processMessage(webhookRequest.Source.UserID, webhookRequest.Content.Text)
	}
	w.WriteHeader(http.StatusOK)
}

func (h *Handler) processMessage(userID string, message string) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	userTurn := db.ChatTurn{
		Role: userRole,
		Text: message,
	}
	chatTurns, err := h.HistoryRepo.GetRecentChatHistory(ctx, userID)
	if err != nil {
		log.Printf("ERROR: Failed to get recent chat history for %s: %v", userID, err)
		return
	}
	messageParams := h.convertHistoryToMessagesParams(chatTurns, userTurn)
	chatCompletion, err := h.OpenAIClient.Chat.Completions.New(ctx, messageParams)

	var finalResponseMessage string
	if err != nil {
		log.Printf("ERROR: Failed to get openai completion for %s: %v", userID, err)
		finalResponseMessage = "Sorry, failed to communicate with the AI server."
	} else if len(chatCompletion.Choices) == 0 {
		finalResponseMessage = "The AI model failed to generate a valid response."
	} else {
		finalResponseMessage = chatCompletion.Choices[0].Message.Content
	}

	if err := h.Sender.PostMessage(ctx, h.Cfg, finalResponseMessage, userID); err != nil {
		log.Printf("ERROR: Failed send messages to %s: %v", userID, err)
	}

	assistantTurn := db.ChatTurn{
		Role: assistantRole,
		Text: finalResponseMessage,
	}
	if err := h.HistoryRepo.SaveAndLimitChatHistory(ctx, userID, userTurn, assistantTurn); err != nil {
		log.Printf("ERROR: Failed to save full chat history for %s: %v", userID, err)
	}
}
