package token

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const minSecretKeySize = 32

type JWTMarker struct {
	secretKey string
}

func NewJWTMarker(secretKey string) (Marker, error) {
	if len(secretKey) < minSecretKeySize {
		return nil, fmt.Errorf("invalid key size: must be at least %d chars", minSecretKeySize)

	}
	return &JWTMarker{secretKey}, nil
}

func (marker *JWTMarker) CreateToken(username string, duration time.Duration) (string, error) {
	payload, err := NewPayload(username, duration)
	if err != nil {
		return "", err
	}
	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)
	return jwtToken.SignedString([]byte(marker.secretKey))
}
func (marker *JWTMarker) VerifyToken(token string) (*Payload, error) {
	keyFunc := func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, ErrorInvalidToken
		}
		return []byte(marker.secretKey), nil
	}
	jwtToken, err := jwt.ParseWithClaims(token, &Payload{}, keyFunc)

	//Todo
	if err != nil {
		if errors.Is(err, jwt.ErrTokenInvalidClaims) {
			return nil, ErrorExpiredToken
		}
		return nil, ErrorInvalidToken
	}
	payload, ok := jwtToken.Claims.(*Payload)
	fmt.Println(payload)
	if !ok {
		return nil, ErrorInvalidToken
	}
	return payload, nil
}
