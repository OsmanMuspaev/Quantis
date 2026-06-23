package clients

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

type githubTokenResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	Scope       string `json:"scope"`
}

func ExchangeGithubCodeForToken(code string) (string, error) {
	data := url.Values{}
	data.Set("client_id", os.Getenv("GITHUB_CLIENT_ID"))
	data.Set("client_secret", os.Getenv("GITHUB_CLIENT_SECRET"))
	data.Set("code", code)

	redirectURI := os.Getenv("GITHUB_REDIRECT_URI")
	if redirectURI == "" {
		redirectURI = "http://localhost:8081/auth/github/callback"
	}
	data.Set("redirect_uri", redirectURI)

	req, err := http.NewRequest(
		"POST",
		"https://github.com/login/oauth/access_token",
		strings.NewReader(data.Encode()),
	)
	if err != nil {
		return "", err
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var tokenResp githubTokenResponse

	if err := json.Unmarshal(bodyBytes, &tokenResp); err != nil {
		return "", err
	}

	if tokenResp.AccessToken == "" {
		return "", errors.New("github did not return access token")
	}

	return tokenResp.AccessToken, nil
}
