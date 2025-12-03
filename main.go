package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"time"

	_ "github.com/mattn/go-sqlite3"

	"goNaverWorksBot/internal/config"
	"goNaverWorksBot/internal/db"
	"goNaverWorksBot/internal/server"
	"goNaverWorksBot/pkg/works"
)

func main() {
	log.Println("[INFO] Application starting...")
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("[Fatal] Configuration Load Failed: %v", err)
	}
	log.Println("[INFO] Configuration loaded.")

	tokenManager, err := works.NewTokenManager(cfg)
	if err != nil {
		log.Fatalf("[Fatal] TokenManager Initialization Failed (Private Key Load): %v", err)
	}
	log.Println("[INFO] TokenManager initialized.")

	messageSender := works.NewMessageSender(tokenManager)
	log.Println("[INFO] Message Sender initialized.")

	dbFile := "chat_history.db"
	database, err := sql.Open("sqlite3", dbFile)
	if err != nil {
		log.Fatalf("[Fatal] Failed to open SQLite database %s: %v", dbFile, err)
	}
	defer database.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := db.InitializeSchema(ctx, database); err != nil {
		log.Fatalf("[Fatal] Failed to initialize database schema: %v", err)
	}
	log.Println("[Info] Database schema initialized.")

	historyRepo := db.NewHistoryRepository(database)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		server.WebhookHandler(w, r, cfg, messageSender, historyRepo)
	})
	addr := fmt.Sprintf(":%d", cfg.Port)
	log.Printf("[INFO] Server listening on address %s", addr)

	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatalf("[FATAL] Server failed: %v", err)
	}
}
