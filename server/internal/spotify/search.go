package spotify

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

type SearchResponse struct {
	Albums struct {
		Items []struct {
			AlbumType   string `json:"album_type"`
			TotalTracks int    `json:"total_tracks"`
			AlbumUrl    string `json:"href"`
			Id          string `json:"id"`
			Images      []struct {
				Url    string `json:"url"`
				Height int    `json:"height"`
				Width  int    `json:"width"`
			} `json:"images"`
			Name                 string `json:"name"`
			ReleaseDate          string `json:"release_date"`
			ReleaseDatePrecision string `json:"release_date_precision"`
			Type                 string `json:"type"`
			Artists              []struct {
				Id           string `json:"id"`
				Name         string `json:"name"`
				Type         string `json:"type"`
				ArtistApiUrl string `json:"href"`
			} `json:"artists"`
		} `json:"items"`
	} `json:"albums"`
}

func Search(accessToken, albumName string) (SearchResponse, error) {
	searchRes := SearchResponse{}
	searchParams := fmt.Sprintf("?offset=0&limit=20&query=%s&type=album", url.QueryEscape(albumName))
	req, err := http.NewRequest("GET", SearchUrl+searchParams, nil)
	if err != nil {
		return searchRes, err
	}
	req.Header.Set("Authorization", "Bearer "+accessToken)

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return searchRes, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return searchRes, err
	}

	err = json.Unmarshal(body, &searchRes)
	if err != nil {
		return searchRes, err
	}

	fmt.Println(searchRes)

	return searchRes, nil
}
