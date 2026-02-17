package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

type Credentials struct {
	ClientID       string    `json:"client_id"`
	ClientSecret   string    `json:"client_secret"`
	AccessToken    string    `json:"access_token"`
	RefreshToken   string    `json:"refresh_token"`
	TokenExpiresAt time.Time `json:"token_expires_at"`
	Scope          string    `json:"scope"`
}

func GetCredentialsPath() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return ""
	}
	return filepath.Join(home, ".config", "strava-cli", "credentials.json")
}

func LoadCredentials() (*Credentials, error) {
	path := GetCredentialsPath()
	
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read credentials: %w", err)
	}

	var creds Credentials
	if err := json.Unmarshal(data, &creds); err != nil {
		return nil, fmt.Errorf("failed to parse credentials: %w", err)
	}

	return &creds, nil
}

func SaveCredentials(creds *Credentials) error {
	path := GetCredentialsPath()
	
	data, err := json.MarshalIndent(creds, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal credentials: %w", err)
	}

	if err := os.WriteFile(path, data, 0600); err != nil {
		return fmt.Errorf("failed to write credentials: %w", err)
	}

	return nil
}

func (c *Credentials) IsExpired() bool {
	return time.Now().After(c.TokenExpiresAt)
}
