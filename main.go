package main

import (
	"context"
	"log"
	"time"

	"goNaverWorksBot/internal/config"
	"goNaverWorksBot/pkg/works"
)

func main() {
	log.Println("--- Starting goNaverWorksBot Initialization ---")
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Configuration Load Failed: %v", err)
	}
	targetID := "bb271bd3-2f4f-61e6-ab3b-d23f1e366c51"
	tokenManager, err := works.NewTokenManager(cfg)
	if err != nil {
		log.Fatalf("TokenManager Initialization Failed (Private Key Load): %v", err)
	}
	messageSender := works.NewMessageSender(tokenManager)
	testMessage := "안녕하세요! Go 언어 봇 개발 환경 구축 후 첫 번째 테스트 메시지입니다. 🎉"

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	err = messageSender.PostMessage(ctx, cfg, testMessage, targetID)
	if err != nil {
		log.Fatalf("Post Message Failed: %v", err)
	}
	log.Println("Initialization complete. Ready to start the server.")
}
