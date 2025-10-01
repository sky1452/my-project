package auth

import (
	"errors"
	"time"
	

	"github.com/golang-jwt/jwt/v5"
)

type JWTManager struct {
	secretKey     string
	tokenDuration time.Duration
}

type UserClaims struct {
	jwt.RegisteredClaims
	UserID int64 `json:"user_id"`
	RoleID int   `json:"role_id"`
}

func NewJWTManager(secretKey string, duration time.Duration) *JWTManager {
	return &JWTManager{secretKey, duration}
}

// 🔹 Генерируем токен JWT для пользователя
func (j *JWTManager) Generate(userID int64, roleID int) (string, error) {
	claims := UserClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(j.tokenDuration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
		UserID: userID,
		RoleID: roleID,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(j.secretKey))
}

// 🔹 Проверяем токен: подпись + срок действия
func (j *JWTManager) Verify(accessToken string) (*UserClaims, error) {
	token, err := jwt.ParseWithClaims(
		accessToken,
		&UserClaims{},
		func(token *jwt.Token) (interface{}, error) {
			return []byte(j.secretKey), nil
		},
	)

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*UserClaims)
	if !ok {
		return nil, errors.New("invalid token claims")
	}

	// JWT сам проверит exp, если используешь jwt.RegisteredClaims
	if !token.Valid {
		return nil, errors.New("invalid or expired token")
	}

	return claims, nil
}
