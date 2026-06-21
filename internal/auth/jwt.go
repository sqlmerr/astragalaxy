package core_auth

import (
	"context"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	core_errors "github.com/sqlmerr/astragalaxy/internal/errors"
)

const ClaimsContextKey = "claims"

type JWTProcessor interface {
	GenerateToken(userID uuid.UUID) (string, error)
	ValidateToken(tokenString string) (uuid.UUID, error)
}

type JWTProcessorImpl struct {
	config Config
}

func NewJWTProcessor(config Config) *JWTProcessorImpl {
	return &JWTProcessorImpl{config}
}

func (m JWTProcessorImpl) GenerateToken(userID uuid.UUID) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": userID.String(),
		"exp": time.Now().Add(m.config.TokenDuration).Unix(),
	})
	signedString, err := token.SignedString([]byte(m.config.JWTSecret))
	if err != nil {
		return "", fmt.Errorf("generate jwt token: %w", err)
	}
	return signedString, nil
}
func (m JWTProcessorImpl) ValidateToken(tokenString string) (uuid.UUID, error) {
	claims := &jwt.MapClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (any, error) {
		return []byte(m.config.JWTSecret), nil
	})
	if err != nil {
		return uuid.Nil, fmt.Errorf("validate jwt token: %w: %w", core_errors.ErrUnauthorized, err)
	}

	errInvalidJwtToken := fmt.Errorf("invalid jwt token: %w", core_errors.ErrUnauthorized)

	if !token.Valid {
		return uuid.Nil, errInvalidJwtToken
	}

	exp, err := claims.GetExpirationTime()
	if err != nil {
		return uuid.Nil, errInvalidJwtToken
	}

	if time.Until(exp.Time).Seconds() <= 0 {
		return uuid.Nil, fmt.Errorf("jwt token expired: %w", errInvalidJwtToken)
	}

	sub, err := claims.GetSubject()
	if err != nil {
		return uuid.Nil, errInvalidJwtToken
	}
	userID, err := uuid.Parse(sub)
	if err != nil {
		return uuid.Nil, errInvalidJwtToken
	}
	return userID, nil
}

func GetUserIDFromContext(ctx context.Context) uuid.UUID {
	return ctx.Value(UserIDContextKey).(uuid.UUID)
}
