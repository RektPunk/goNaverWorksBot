package works

import (
	"bytes"
	"crypto/rsa"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"io"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"goNaverWorksBot/internal/config"
)

type TokenManager struct {
	cfg         *config.Config
	client      *http.Client
	privateKey  *rsa.PrivateKey
	accessToken string
	tokenExpiry int64
	mu          sync.RWMutex
}

type TokenResponse struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   string `json:"expires_in"`
	TokenType   string `json:"token_type"`
}

func NewTokenManager(cfg *config.Config) (*TokenManager, error) {
	keyBytes, err := os.ReadFile(cfg.PrivateKeyPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read private key file at %s: %w", cfg.PrivateKeyPath, err)
	}
	block, _ := pem.Decode(keyBytes)
	if block == nil {
		return nil, fmt.Errorf("failed to decode private key: key is not in PEM format or file is invalid")
	}
	var privateKey *rsa.PrivateKey
	key, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse PKCS#8 private key: %w", err)
	}
	var ok bool
	privateKey, ok = key.(*rsa.PrivateKey)
	if !ok {
		return nil, fmt.Errorf("parsed PKCS#8 key is not an RSA private key. Got type %T", key)
	}
	return &TokenManager{
		cfg:        cfg,
		client:     &http.Client{Timeout: 10 * time.Second},
		privateKey: privateKey,
	}, nil
}

func (m *TokenManager) GetToken() (string, error) {
	m.mu.RLock()
	now := time.Now().Unix()
	if m.accessToken != "" && now < m.tokenExpiry-60 {
		m.mu.RUnlock()
		return m.accessToken, nil
	}
	m.mu.RUnlock()

	m.mu.Lock()
	defer m.mu.Unlock()
	if m.accessToken != "" && now < m.tokenExpiry-60 {
		return m.accessToken, nil
	}

	exp := now + 3600
	claims := jwt.MapClaims{
		"iss": m.cfg.ClientID,
		"sub": m.cfg.ServiceAccount,
		"iat": now,
		"exp": exp,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	encodedJWT, err := token.SignedString(m.privateKey)
	if err != nil {
		return "", fmt.Errorf("failed to sign JWT: %w", err)
	}
	tokenURL := "https://auth.worksmobile.com/oauth2/v2.0/token"
	data := fmt.Sprintf(
		"assertion=%s&grant_type=urn:ietf:params:oauth:grant-type:jwt-bearer&client_id=%s&client_secret=%s&scope=bot bot.message bot.read",
		encodedJWT,
		m.cfg.ClientID,
		m.cfg.ClientSecret,
	)
	req, err := http.NewRequest("POST", tokenURL, bytes.NewBufferString(data))
	if err != nil {
		return "", fmt.Errorf("failed to create token request: %w", err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")

	resp, err := m.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("token request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("failed to get token (status: %d). response body: %s", resp.StatusCode, string(bodyBytes))
	}

	var tokenData TokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&tokenData); err != nil {
		return "", fmt.Errorf("failed to decode token response: %w", err)
	}
	m.accessToken = tokenData.AccessToken
	m.tokenExpiry = exp
	return m.accessToken, nil
}

func (m *TokenManager) SetHeaders() (map[string]string, error) {
	token, err := m.GetToken()
	if err != nil {
		return nil, err
	}
	return map[string]string{
		"Authorization": fmt.Sprintf("Bearer %s", token),
		"Content-Type":  "application/json",
	}, nil
}
