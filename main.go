package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/option"

	"goNaverWorksBot/internal/config"
	"goNaverWorksBot/internal/db"
	"goNaverWorksBot/internal/server"
	"goNaverWorksBot/pkg/works"
)

const dbFile string = "chat_history.db"

func main() {
	log.Println("[INFO] Application starting...")

	// load config
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("[FATAL] Configuration Load Failed: %v", err)
	}
	log.Println("[INFO] Configuration loaded.")

	// Naver Works initialization
	// access token manager initialization
	tokenManager, err := works.NewTokenManager(cfg)
	if err != nil {
		log.Fatalf("[FATAL] TokenManager Initialization Failed (Private Key Load): %v", err)
	}
	log.Println("[INFO] TokenManager initialized.")

	// message sender initialization
	messageSender := works.NewMessageSender(tokenManager)
	log.Println("[INFO] Message Sender initialized.")

	// DB initialization
	database, err := sql.Open("sqlite3", dbFile)
	if err != nil {
		log.Fatalf("[FATAL] Failed to open SQLite database %s: %v", dbFile, err)
	}
	defer database.Close()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := db.InitializeSchema(ctx, database); err != nil {
		log.Fatalf("[FATAL] Failed to initialize database schema: %v", err)
	}
	log.Println("[INFO] Database schema initialized.")

	// chat history repository initialization
	historyRepo := db.NewHistoryRepository(database)

	// openai client initialization
	client := openai.NewClient(option.WithAPIKey(cfg.OpenAIAPIKey))

	// handler initialization
	handler := server.NewHandler(
		cfg,
		messageSender,
		historyRepo,
		&client,
	)
	log.Println("[INFO] Webhook Handler initialized with all dependencies.")

	// set http routing
	http.HandleFunc("/", handler.WebhookHandler)

	// server start
	addr := fmt.Sprintf(":%d", cfg.Port)
	log.Printf("[INFO] Server listening on address %s", addr)

	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatalf("[FATAL] Server failed: %v", err)
	}
}
