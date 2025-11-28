package main

import (
	"fmt"
	"goNaverWorksBot/internal/config"
	"log"
)

func main() {
	log.Println("--- Starting goNaverWorksBot Initialization ---")

	// 1. Config 로딩 시도
	cfg, err := config.Load()
	if err != nil {
		// config.Load()에서 반환된 모든 치명적인 오류(FATAL, Missing variables)를 처리합니다.
		log.Fatalf("❌ Configuration Load Failed: %v", err)
	}

	log.Println("✅ Configuration Loaded Successfully!")

	// 2. 로드된 핵심 설정 값 출력 및 확인
	fmt.Println("\n--- Loaded Configuration Details ---")
	fmt.Printf("✅ Server Port: %d (Default if not specified in .env)\n", cfg.Port)
	fmt.Printf("✅ Works Bot ID: %s\n", cfg.BotID)

	// 민감 정보는 출력하지 않거나 마스킹하여 보안을 유지합니다.
	fmt.Printf("✅ Works Client ID (Partial): %s...\n", cfg.ClientID[:4])
	fmt.Printf("✅ Service Account: %s...\n", cfg.ServiceAccount[:4])
	fmt.Printf("✅ Private Key Path: %s\n", cfg.PrivateKeyPath)
	fmt.Println("------------------------------------")

	// 테스트 성공 후, 이제 HTTP 서버 구동 로직이 여기에 들어갈 예정입니다.
	log.Println("Initialization complete. Ready to start the server.")
}
