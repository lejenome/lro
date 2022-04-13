package auth

import (
	"errors"
	"log"
	"net/http"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/dgrijalva/jwt-go/request"
	"github.com/google/uuid"
	"github.com/lejenome/lro/pkg/config"
	"github.com/lejenome/lro/services/process-api/providers/context"
)

var (
	ErrInvalidToken = errors.New("token is invalid")
	ErrExpiredToken = errors.New("token has expired")
)

const COOKIE_NAME = "access_token"

type AuthPayload struct {
	*context.Context
	UUID uuid.UUID `json:"uuid"`
	// UserID    uint      `json:"user_id"`
	// UserRoles []string  `json:"user_roles"`
	// TenantRef  string    `json:"tenant_ref"`
	IssuedAt time.Time `json:"issued_at"`
	ExpireAt time.Time `json:"expire_at"`
}

func NewAuthPayload(userID uint, duration time.Duration) (*AuthPayload, time.Time) {
	expireAt := time.Now().Add(duration)
	tokenID := uuid.New()
	payload := &AuthPayload{
		&context.Context{
			UserID:    userID,
			UserRoles: []string{},
			TenantRef: "public",
		},
		tokenID,
		time.Now(),
		expireAt,
	}
	return payload, expireAt
}

func (payload *AuthPayload) Valid() error {
	if time.Now().After(payload.ExpireAt) {
		return ErrExpiredToken
	}
	/*
		if res := redis.Exists(fmt.Sprintf("jwt:token:%s", payload.UUID)).Val(); res != 1 {
			return ErrExpiredToken
		}
	*/
	return nil
}

func (payload *AuthPayload) Map() map[string]interface{} {
	return map[string]interface{}{
		"uuid":       payload.UUID.String(),
		"user_id":    payload.UserID,
		"user_roles": payload.UserRoles,
		"tenant_ref": payload.TenantRef,
	}
}

type Maker interface {
	CreateToken(userID uint) (string, time.Time, error)
	VerifyToken(token string) (*AuthPayload, error)
	ExtractToken(req *http.Request) (string, error)
}

const minSecretKeySize = 32

type JWTMaker struct {
	Config *config.JWTConfig
}

func NewJWTMaker(config *config.JWTConfig) Maker {
	if len(config.SecretKey) < minSecretKeySize {
		log.Fatalf("Invalid JWT Secret Key: must be at least %d characters", minSecretKeySize)
	}
	return &JWTMaker{config}
}

func (m *JWTMaker) CreateToken(userID uint) (string, time.Time, error) {
	payload, expireAt := NewAuthPayload(userID, m.Config.Duration)
	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)
	token, err := jwtToken.SignedString([]byte(m.Config.SecretKey))
	/*
		res := redis.Set(fmt.Sprintf("jwt:token:%s", payload.UUID), token, m.Config.Duration)
		log.Printf("Redis Set: %v", res)
	*/
	return token, expireAt, err
}

func (m *JWTMaker) VerifyToken(token string) (*AuthPayload, error) {
	keyFunc := func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, ErrInvalidToken
		}
		return []byte(m.Config.SecretKey), nil
	}

	jwtToken, err := jwt.ParseWithClaims(token, &AuthPayload{}, keyFunc)
	if err != nil {
		verr, ok := err.(*jwt.ValidationError)
		if ok && errors.Is(verr.Inner, ErrExpiredToken) {
			return nil, ErrExpiredToken
		}
		return nil, ErrInvalidToken
	}

	payload, ok := jwtToken.Claims.(*AuthPayload)
	if !ok {
		return nil, ErrInvalidToken
	}
	err = payload.Valid()
	if err != nil {
		return nil, err
	}

	return payload, nil
}

func (m *JWTMaker) ExtractToken(req *http.Request) (string, error) {
	token, err := request.OAuth2Extractor.ExtractToken(req)
	if err == nil {
		return token, err
	}
	cookie, err := req.Cookie(COOKIE_NAME)
	if err == nil && cookie.Value != "" {
		return cookie.Value, nil
	}
	return "", err
}
