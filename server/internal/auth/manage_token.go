package auth

import (
	"fmt"
	"net/http"
	"strings"
)

func GetBearerToken(headers http.Header) (string, error) {
	authorization := headers.Get("Authorization")
	if authorization == "" {
		return "", fmt.Errorf("no authorization was provided")
	}
	return strings.Replace(authorization, "Bearer ", "", 1), nil
}
