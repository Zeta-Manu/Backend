package utils

import (
	"fmt"
	"net/http"
	"strings"
)

func ParseToken(r *http.Request) (string, error) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return "", fmt.Errorf("authorization header not found")
	}

	splitToken := strings.Split(authHeader, "Bearer ")
	if len(splitToken) != 2 {
		return "", fmt.Errorf("bearer token not in proper format")
	}

	token := strings.TrimSpace(splitToken[1])
	return token, nil
}
