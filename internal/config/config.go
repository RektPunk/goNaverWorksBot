package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	Port           int
	ClientID       string
	ClientSecret   string
	ServiceAccount string
	BotID          string
	BotSecret      string
	PrivateKeyPath string
}

func Load() (*Config, error) {
	if err := godotenv.Load(); err != nil {
		return nil, fmt.Errorf("FATAL: Failed to load .env file. Please ensure .env file exists in the root directory. Error: %w", err)
	}

	conf := &Config{
		ClientID:       os.Getenv("CLIENT_ID"),
		ClientSecret:   os.Getenv("CLIENT_SECRET"),
		ServiceAccount: os.Getenv("SERVICE_ACCOUNT"),
		BotID:          os.Getenv("BOT_ID"),
		BotSecret:      os.Getenv("BOT_SECRET"),
		PrivateKeyPath: os.Getenv("PRIVATE_KEY_PATH"),
	}

	portStr := os.Getenv("PORT")
	if portStr == "" {
		portStr = "8080" // if PORT is not exists in .env
	}

	port, err := strconv.Atoi(portStr)
	if err != nil {
		return nil, fmt.Errorf("invalid port number: %s", portStr)
	}
	conf.Port = port

	if conf.BotID == "" ||
		conf.ClientID == "" ||
		conf.ClientSecret == "" ||
		conf.ServiceAccount == "" ||
		conf.BotSecret == "" ||
		conf.PrivateKeyPath == "" {

		return nil, fmt.Errorf("missing one or more required environment variables (WORKS_BOT_ID, WORKS_CLIENT_ID, WORKS_CLIENT_SECRET, SERVICE_ACCOUNT, NAVERWORKS_BOT_SECRET, PRIVATE_KEY_PATH) must be set")
	}

	return conf, nil
}
