package server

import (
	"goNaverWorksBot/internal/db"

	"github.com/openai/openai-go/v3"
)

const systemPrompt string = `You are RektPunkBot, an insightful, encouraging assistant. 
Do not explain your role. Answer the user's questions directly, concisely, and exclusively. 
Refer to the previous conversation history to continue the dialogue naturally.`
const maxTokens int64 = 500

func (h *Handler) convertHistoryToMessagesParams(history []db.ChatTurn, userTurn db.ChatTurn) openai.ChatCompletionNewParams {
	messages := []openai.ChatCompletionMessageParamUnion{
		openai.SystemMessage(systemPrompt),
	}
	for _, turn := range history {
		switch turn.Role {
		case userRole:
			messages = append(messages, openai.UserMessage(turn.Text))
		case assistantRole:
			messages = append(messages, openai.AssistantMessage(turn.Text))
		}
	}
	messages = append(messages, openai.UserMessage(userTurn.Text))
	messageParams := openai.ChatCompletionNewParams{
		Messages:            messages,
		Model:               openai.ChatModelGPT4_1,
		MaxCompletionTokens: openai.Int(maxTokens),
	}
	return messageParams
}
