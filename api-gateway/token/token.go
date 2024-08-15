package token

import (
	"Api_Gateway/config"
	"fmt"
	"strings"

	"github.com/golang-jwt/jwt"
)

func ValidateToken(tokenstr string) (bool, error) {
	_, err := ExtractClaims(tokenstr)
	if err != nil {
		return false, err
	}
	return true, nil
}

func ExtractClaims(tokenstr string) (jwt.MapClaims, error) {
	// Tokenni tozalash
	tokenstr = strings.TrimPrefix(tokenstr, "\"")
	tokenstr = strings.TrimSuffix(tokenstr, "\"")
	token, err := jwt.ParseWithClaims(tokenstr, jwt.MapClaims{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return []byte(config.Load().ACCESS_TOKEN), nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %w", err)
	}

	if !token.Valid {
		return nil, fmt.Errorf("invalid token: %s", tokenstr)
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("failed to parse token claims")
	}
	return claims, nil
}
