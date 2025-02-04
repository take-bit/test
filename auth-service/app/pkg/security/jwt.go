package security

import (
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"os"
	"time"
)

var (
	ErrInvalidToken = errors.New("invalid token")
	ErrExpiredToken = errors.New("token has expired")
)

var secretKey = []byte(os.Getenv("SECRET_KEY"))

func CreateToken(uuid, username, jit string, isVerifiedEmail bool, ttl time.Duration) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.MapClaims{
			"uuid":            uuid,
			"username":        username,
			"isVerifiedEmail": isVerifiedEmail,
			"iat":             time.Now().Unix(),
			"exp":             time.Now().Add(ttl).Unix(),
			"jit":             jit,
		})

	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func ValidateToken(tokenString string) (string, bool, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrInvalidToken
		}
		return secretKey, nil
	})

	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return "", false, ErrExpiredToken
		}
		return "", false, ErrInvalidToken
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return "", false, ErrInvalidToken
	}

	username, ok1 := claims["username"].(string)
	isVerified, ok2 := claims["isVerified"].(bool)
	if !ok1 || !ok2 {
		return "", false, ErrInvalidToken
	}

	return username, isVerified, nil
}
