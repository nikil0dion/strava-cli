package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

type Credentials struct {
	ClientID       string    `json:"client_id"`
	ClientSecret   string    `json:"client_secret"`
	AccessToken    string    `json:"access_token"`
	RefreshToken   string    `json:"refresh_token"`
	TokenExpiresAt time.Time `json:"token_expires_at"`
	Scope          string    `json:"scope"`
	FromEnv        bool      `json:"-"` // true if loaded from environment
}

// Environment variable names
const (
	EnvClientID       = "STRAVA_CLIENT_ID"
	EnvClientSecret   = "STRAVA_CLIENT_SECRET"
	EnvAccessToken    = "STRAVA_ACCESS_TOKEN"
	EnvRefreshToken   = "STRAVA_REFRESH_TOKEN"
	EnvTokenExpiresAt = "STRAVA_TOKEN_EXPIRES_AT" // Unix timestamp
)

func GetCredentialsPath() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return ""
	}
	return filepath.Join(home, ".config", "strava-cli", "credentials.json")
}

// LoadCredentials loads credentials from environment variables first,
// falls back to credentials.json if env is not set.
func LoadCredentials() (*Credentials, error) {
	// Try environment variables first
	if creds := loadFromEnv(); creds != nil {
		return creds, nil
	}

	// Fallback to file
	return loadFromFile()
}

func loadFromEnv() *Credentials {
	clientID := os.Getenv(EnvClientID)
	clientSecret := os.Getenv(EnvClientSecret)
	accessToken := os.Getenv(EnvAccessToken)
	refreshToken := os.Getenv(EnvRefreshToken)

	// Require at least access token
	if accessToken == "" {
		return nil
	}

	creds := &Credentials{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		FromEnv:      true,
	}

	// Parse expiry if provided
	if expiresStr := os.Getenv(EnvTokenExpiresAt); expiresStr != "" {
		if ts, err := strconv.ParseInt(expiresStr, 10, 64); err == nil {
			creds.TokenExpiresAt = time.Unix(ts, 0)
		}
	}

	return creds
}

func loadFromFile() (*Credentials, error) {
	path := GetCredentialsPath()

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("credentials not found in env (%s) or file (%s): %w",
			EnvAccessToken, path, err)
	}

	var creds Credentials
	if err := json.Unmarshal(data, &creds); err != nil {
		return nil, fmt.Errorf("failed to parse credentials: %w", err)
	}

	return &creds, nil
}

// SaveCredentials saves credentials to file, or prints to stderr if from env.
func SaveCredentials(creds *Credentials) error {
	if creds.FromEnv {
		// Print updated tokens to stderr for env update
		fmt.Fprintln(os.Stderr, "\n# Updated tokens (update your environment):")
		fmt.Fprintf(os.Stderr, "export %s=%s\n", EnvAccessToken, creds.AccessToken)
		fmt.Fprintf(os.Stderr, "export %s=%s\n", EnvRefreshToken, creds.RefreshToken)
		fmt.Fprintf(os.Stderr, "export %s=%d\n", EnvTokenExpiresAt, creds.TokenExpiresAt.Unix())
		return nil
	}

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
	// If no expiry set, assume not expired (for env-only setups)
	if c.TokenExpiresAt.IsZero() {
		return false
	}
	return time.Now().After(c.TokenExpiresAt)
}

// CanRefresh returns true if we have credentials to perform token refresh.
func (c *Credentials) CanRefresh() bool {
	return c.ClientID != "" && c.ClientSecret != "" && c.RefreshToken != ""
}
