package spotify

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

type AuthResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int    `json:"expires_in"`
	Err         string `json:"error"`
}

func Authenticate(clientId, clientSecret string) (AuthResponse, error) {
	authRes := AuthResponse{}
	reqBody := fmt.Sprintf("grant_type=client_credentials&client_id=%s&client_secret=%s", clientId, clientSecret)
	req, err := http.NewRequest("POST", AuthUrl, strings.NewReader(reqBody))
	if err != nil {
		return authRes, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return authRes, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return authRes, err
	}

	err = json.Unmarshal(body, &authRes)
	if err != nil {
		return authRes, err
	}

	if authRes.Err != "" {
		return authRes, fmt.Errorf("%v", authRes.Err)
	}

	return authRes, nil
}
