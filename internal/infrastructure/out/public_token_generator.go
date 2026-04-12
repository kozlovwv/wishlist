package out

import (
	"crypto/rand"
	"encoding/hex"
)

type PublicTokenGenerator struct{}

func NewPublicTokenGenerator() *PublicTokenGenerator {
	return &PublicTokenGenerator{}
}

func (g *PublicTokenGenerator) Generate() (string, error) {
	bytes := make([]byte, 16)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}
