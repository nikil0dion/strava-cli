package auth

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/nikilodion/strava-cli/internal/config"
)

const tokenURL = "https://www.strava.com/oauth/token"

type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresAt    int64  `json:"expires_at"`
	ExpiresIn    int    `json:"expires_in"`
}

func RefreshToken(creds *config.Credentials) error {
	data := url.Values{}
	data.Set("client_id", creds.ClientID)
	data.Set("client_secret", creds.ClientSecret)
	data.Set("grant_type", "refresh_token")
	data.Set("refresh_token", creds.RefreshToken)

	resp, err := http.Post(tokenURL, "application/x-www-form-urlencoded", bytes.NewBufferString(data.Encode()))
	if err != nil {
		return fmt.Errorf("failed to refresh token: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("token refresh failed (status %d): %s", resp.StatusCode, body)
	}

	var tokenResp TokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&tokenResp); err != nil {
		return fmt.Errorf("failed to parse token response: %w", err)
	}

	// Update credentials
	creds.AccessToken = tokenResp.AccessToken
	creds.RefreshToken = tokenResp.RefreshToken
	creds.TokenExpiresAt = time.Unix(tokenResp.ExpiresAt, 0)

	// Save updated credentials
	if err := config.SaveCredentials(creds); err != nil {
		return fmt.Errorf("failed to save refreshed credentials: %w", err)
	}

	return nil
}

func EnsureValidToken(creds *config.Credentials) error {
	if creds.IsExpired() {
		fmt.Println("Token expired, refreshing...")
		return RefreshToken(creds)
	}
	return nil
}
