package security

import (
	"errors"
	"fmt"
	"time"

	"github.com/codepnw/simple-bank/config"
	"github.com/dgrijalva/jwt-go"
)

type TokenUser struct {
	ID    int64
	Email string
	Role  string
}

type Token struct {
	cfg config.EnvConfig
}

func InitJWT(cfg config.EnvConfig) *Token {
	return &Token{cfg: cfg}
}

func (t *Token) GenerateAccessToken(user *TokenUser) (string, error) {
	exp := time.Hour * 24
	return t.generateToken(t.cfg.JWT.SecretKey, exp, user)
}

func (t *Token) GenerateRefreshToken(user *TokenUser) (string, error) {
	exp := time.Hour * 24 * 7
	return t.generateToken(t.cfg.JWT.RefreshKey, exp, user)
}

func (t *Token) VerifyAccessToken(token string) (*TokenUser, error) {
	return t.verifyToken(token, t.cfg.JWT.SecretKey)
}

func (t *Token) VerifyRefreshToken(token string) (*TokenUser, error) {
	return t.verifyToken(token, t.cfg.JWT.RefreshKey)
}

// Generate Token String
func (t *Token) generateToken(key string, exp time.Duration, user *TokenUser) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"email":   user.Email,
		"role":    user.Role,
		"exp":     time.Now().Add(exp).Unix(),
	})

	tokenString, err := token.SignedString([]byte(key))
	if err != nil {
		return "", fmt.Errorf("sign token failed: %w", err)
	}

	return tokenString, nil
}

// Verify Token String
func (t *Token) verifyToken(tokenStr, key string) (*TokenUser, error) {
	token, err := jwt.Parse(tokenStr, func(tt *jwt.Token) (any, error) {
		if _, ok := tt.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unknow signing method: %v", tt.Header)
		}
		return []byte(key), nil
	})
	if err != nil {
		return nil, err
	}
	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("verification failed")
	}

	u := new(TokenUser)

	u.ID = int64(claims["user_id"].(float64))
	u.Email = claims["email"].(string)
	u.Role = claims["role"].(string)

	return u, nil
}
