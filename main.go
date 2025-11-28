package main

import (
	"log"

	"goNaverWorksBot/internal/config"
	"goNaverWorksBot/pkg/works"
)

func main() {
	log.Println("--- Starting goNaverWorksBot Initialization ---")
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Configuration Load Failed: %v", err)
	}
	log.Println("Configuration Loaded Successfully!")

	tokenManager, err := works.NewTokenManager(cfg)
	if err != nil {
		log.Fatalf("TokenManager Initialization Failed (Private Key Load): %v", err)
	}
	log.Println("TokenManager Initialized.")

	log.Println("--- Attempting to retrieve Access Token ---")
	accessToken, err := tokenManager.GetToken()
	if err != nil {
		log.Fatalf("Failed to get Access Token: %v", err)
	}
	log.Printf("Access Token successfully retrieved (Partial): %s...", accessToken[:20])

	log.Println("Initialization complete. Ready to start the server.")
}
