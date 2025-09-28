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

func Search(apiKey, searchQuery string) error {
	searchUrl := fmt.Sprintf("%s?key=%v&q=%v&type=video", youtubeAPIBaseURl, url.QueryEscape(apiKey), url.QueryEscape(searchQuery))
	// url := "https://www.googleapis.com/youtube/v3/search?key=AIzaSyAgCdA5TnTjCRlVDBL4-1f-HMCzOrWo3do&q=Frank Ocean Lost&type=video"
	result := SearchResult{}
	response, err := http.Get(searchUrl)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	resBody, err := io.ReadAll(response.Body)
	if err != nil {
		return err
	}
	fmt.Println(string(resBody))
	err = json.Unmarshal(resBody, &result)
	if err != nil {
		return err
	}

	return fmt.Errorf("%v", result)
}
