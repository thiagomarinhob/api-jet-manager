package auth

import (
	"errors"
	"time"

	"api-jet-manager/internal/domain/models"

	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
)

type JWTService struct {
	secretKey     string
	tokenDuration time.Duration
}

type JWTClaims struct {
	UserID       string          `json:"user_id"`
	Email        string          `json:"email"`
	UserType     models.UserType `json:"user_type"`
	RestaurantID *uuid.UUID      `json:"restaurant_id,omitempty"`
	jwt.RegisteredClaims
}

func NewJWTService(secretKey string, tokenDuration time.Duration) *JWTService {
	return &JWTService{
		secretKey:     secretKey,
		tokenDuration: tokenDuration,
	}
}

func (j *JWTService) GenerateToken(user *models.User) (string, error) {
	claims := JWTClaims{
		UserID:       user.ID.String(),
		Email:        user.Email,
		UserType:     user.Type,
		RestaurantID: user.RestaurantID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(j.tokenDuration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(j.secretKey))
}

func (j *JWTService) ValidateToken(tokenString string) (*JWTClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(j.secretKey), nil
	})

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*JWTClaims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token")
	}

	return claims, nil
}
