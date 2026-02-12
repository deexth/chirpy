package auth

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func GetBearerToken(headers http.Header) (string, error) {
	tokenString := headers.Get("Authorization")
	if tokenString == "" {
		return tokenString, errors.New("auth header empty, no token found")
	}

	token := strings.Fields(tokenString)
	if len(token) != 2 || token[0] != "Bearer" || token[1] == "" {
		return "", errors.New("wrong auth token")
	}

	return token[1], nil
}

func MakeRefreshToken() (string, error) {
	key := make([]byte, 32)
	_, err := rand.Read(key)
	if err != nil {
		return "", fmt.Errorf("issue creating a random key: %v", err)
	}

	encodedStr := hex.EncodeToString(key)
	return encodedStr, nil
}

func MakeJWT(userID uuid.UUID, tokenSecret string, expiresIn time.Duration) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Issuer:    "chirpy-access",
		IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(expiresIn * time.Hour)),
		Subject:   userID.String(),
	})

	ss, err := token.SignedString([]byte(tokenSecret))
	if err != nil {
		return "", fmt.Errorf("issue signing token: %v", err)
	}

	return ss, nil
}

func ValidateJWT(tokenString, tokenSecret string) (uuid.UUID, error) {
	token, err := jwt.ParseWithClaims(tokenString, &jwt.RegisteredClaims{}, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return []byte(tokenSecret), nil
	})

	if err != nil {
		return uuid.Nil, fmt.Errorf("couldn't parse the token string: %v", err)
	} else if claims, ok := token.Claims.(*jwt.RegisteredClaims); ok {
		if id, err := uuid.Parse(claims.Subject); err != nil {
			return uuid.Nil, fmt.Errorf("unable to parse user_id: %v", err)
		} else {
			return id, nil
		}
	} else {
		return uuid.Nil, errors.New("wrong claims")
	}
}
