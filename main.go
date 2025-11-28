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
	log.Println("--- Starting goNaverWorksBot Initialization ---")
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("[Fatal] Configuration Load Failed: %v", err)
	}
	tokenManager, err := works.NewTokenManager(cfg)
	if err != nil {
		log.Fatalf("[Fatal] TokenManager Initialization Failed (Private Key Load): %v", err)
	}
	messageSender := works.NewMessageSender(tokenManager)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		server.WebhookHandler(w, r, cfg, messageSender)
	})
	addr := fmt.Sprintf(":%d", cfg.Port)
	log.Printf("goNaverWorks Server Started: %s port", addr)
	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatalf("goNaverWorks Server Failed: %v", err)
	}
}
