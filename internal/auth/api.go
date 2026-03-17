package auth

import (
	"errors"
	"net/http"
	"strings"
)

func GetAPIKey(headers http.Header) (string, error) {
	apiKey := headers.Get("Authorization")
	if apiKey == "" {
		return apiKey, errors.New("auth header empty, no api found")
	}

	key := strings.Fields(apiKey)
	if len(key) != 2 || key[0] != "ApiKey" || key[1] == "" {
		return "", errors.New("wrong auth apiKey")
	}

	return key[1], nil
}
