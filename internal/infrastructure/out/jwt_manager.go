package out

import (
	"fmt"
	"time"

	"wishlist/internal/domain"

	"github.com/golang-jwt/jwt/v5"
)

type JWTManager struct {
	secretKey []byte
	ttl       time.Duration
}

type claims struct {
	jwt.RegisteredClaims
	UserID int64 `json:"user_id"`
}

func NewJWTManager(secretKey string, ttl time.Duration) *JWTManager {
	return &JWTManager{
		secretKey: []byte(secretKey),
		ttl:       ttl,
	}
}

func (m *JWTManager) Generate(userID int64) (string, error) {
	c := claims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(m.ttl)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
		UserID: userID,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, c)

	signed, err := token.SignedString(m.secretKey)
	if err != nil {
		return "", fmt.Errorf("sign jwt token: %w", err)
	}

	return signed, nil
}

func (m *JWTManager) Parse(tokenStr string) (int64, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, domain.ErrInvalidToken
		}
		return m.secretKey, nil
	})
	if err != nil {
		return 0, domain.ErrInvalidToken
	}

	c, ok := token.Claims.(*claims)
	if !ok || !token.Valid {
		return 0, domain.ErrInvalidToken
	}

	return c.UserID, nil
}
