package token

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt"
)

const minSecretKeySize = 32

var (
	ErrInvalidToken = fmt.Errorf("invalid token")
	ErrExpiredToken = fmt.Errorf("token has expired")
)

type JWTMaker struct {
	secretKey string
}

func NewJWTMaker(secretkey string) (Maker, error) {
	if len(secretkey) < minSecretKeySize {
		return nil, fmt.Errorf("invalid key size: must be at least %d characters", minSecretKeySize)
	}
	return &JWTMaker{secretKey: secretkey}, nil
}

func (jwtMaker *JWTMaker) CreateToken(username string, duration time.Duration) (string, error) {
	payload, err := NewPayload(username, duration)
	if err != nil {
		return "", nil
	}
	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)
	return jwtToken.SignedString([]byte(jwtMaker.secretKey))

}

func (jwtMaker *JWTMaker) VerifyToken(token string) (*Payload, error) {
	jwtToken, err := jwt.ParseWithClaims(token, &Payload{}, func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC) // HMAC is a type of signing method
		if !ok {
			return nil, ErrInvalidToken
		}
		return []byte(jwtMaker.secretKey), nil
	})
	if err != nil {
		verr, ok := err.(*jwt.ValidationError)
		if ok && errors.Is(verr.Inner, ErrExpiredToken) {
			return nil, ErrExpiredToken
		}
		return nil, ErrInvalidToken
	}
	payload, ok := jwtToken.Claims.(*Payload)
	if !ok {
		return nil, ErrInvalidToken
	}
	return payload, nil
}
