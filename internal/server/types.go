package server

import (
	"github.com/openai/openai-go/v3"

	"goNaverWorksBot/internal/config"
	"goNaverWorksBot/internal/db"
	"goNaverWorksBot/pkg/works"
)

const userRole string = "user"
const assistantRole string = "assistant"

type Handler struct {
	Cfg          *config.Config
	Sender       *works.MessageSender
	HistoryRepo  *db.HistoryRepository
	OpenAIClient *openai.Client
}
