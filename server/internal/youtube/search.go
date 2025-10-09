package youtube

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

type SearchResult struct {
	Etag          string `json:"etag"`
	NextPageToken string `json:"nextPageToken"`
	Items         []struct {
		Id struct {
			VideoId string `json:"videoId"`
		} `json:"id"`
	} `json:"items"`
}

func Search(apiKey, searchQuery string) (SearchResult, error) {
	searchUrl := fmt.Sprintf("%s?key=%v&q=%v&type=video", youtubeAPIBaseURl, url.QueryEscape(apiKey), url.QueryEscape(searchQuery))
	result := SearchResult{}
	response, err := http.Get(searchUrl)
	if err != nil {
		return result, err
	}
	defer response.Body.Close()

	resBody, err := io.ReadAll(response.Body)
	if err != nil {
		return result, err
	}
	err = json.Unmarshal(resBody, &result)
	if err != nil {
		return result, err
	}

	return result, nil
}
