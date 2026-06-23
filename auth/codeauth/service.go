package codeauth

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"time"

	"auth/domain"
	"auth/storage"
)

// GenerateCode creates a cryptographically secure 6-digit verification code.
func GenerateCode(entryToken string) (string, error) {
	n, err := rand.Int(rand.Reader, big.NewInt(900000))
	if err != nil {
		return "", fmt.Errorf("failed to generate random code: %w", err)
	}
	code := fmt.Sprintf("%06d", n.Int64()+100000)

	codeState := domain.CodeState{
		EntryToken: entryToken,
		ExpiresAt:  time.Now().Add(1 * time.Minute),
	}
	storage.SaveCode(codeState, code)

	return code, nil
}
