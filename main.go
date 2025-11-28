package main

import (
	"fmt"
	"goNaverWorksBot/internal/config"
	"goNaverWorksBot/internal/server"
	"goNaverWorksBot/pkg/works"
	"log"
	"net/http"
)

func main() {
	log.Println("[INFO] Application starting...")
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("[Fatal] Configuration Load Failed: %v", err)
	}
	log.Println("[INFO] Configuration loaded successfully.")

	tokenManager, err := works.NewTokenManager(cfg)
	if err != nil {
		log.Fatalf("[Fatal] TokenManager Initialization Failed (Private Key Load): %v", err)
	}
	log.Println("[INFO] TokenManager initialization successfully.")

	messageSender := works.NewMessageSender(tokenManager)
	log.Println("[INFO] Message Sender initialized.")

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		server.WebhookHandler(w, r, cfg, messageSender)
	})
	addr := fmt.Sprintf(":%d", cfg.Port)
	log.Printf("[INFO] Server listening on address %s", addr)

	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatalf("[FATAL] Server failed: %v", err)
	}
}
