package spotify

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

type AlbumResponse struct {
	Tracks struct {
		Limit    int    `json:"limit"`
		Next     string `json:"next"`
		Previous string `json:"previous"`
		Total    int    `json:"total"`
		Items    []struct {
			DiscNumber  int    `json:"disc_number"` // most album are just 1 disc (some might be 2)
			Duration    int    `json:"duration_ms"`
			Explicit    bool   `json:"explicit"`
			Id          string `json:"id"`
			Name        string `json:"name"`
			TrackNumber int    `json:"track_number"`
			TrackUri    string `json:"uri"`
		}
	} `json:"tracks"`
}

func GetAlbumTracks(accessToken, albumId string) (AlbumResponse, error) {
	albumRes := AlbumResponse{}
	req, err := http.NewRequest("GET", AlbumUrl+url.PathEscape(albumId), nil)
	if err != nil {
		return albumRes, err
	}
	req.Header.Set("Authorization", "Bearer "+accessToken)

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return albumRes, err
	}
	defer res.Body.Close()
	if res.StatusCode >= 400 && res.StatusCode <= 499 {
		return albumRes, fmt.Errorf("%v", UnvalidAuthErrorMessage)
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return albumRes, err
	}

	err = json.Unmarshal(body, &albumRes)
	if err != nil {
		return albumRes, err
	}

	return albumRes, nil
}
