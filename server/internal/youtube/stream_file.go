package youtube

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"regexp"
)

func SearchVideoUrl() {

}

func GetStreamUrl(videoUrl string) error {

	htmlText, err := downloadWebPage(videoUrl)
	if err != nil {
		return err
	}
	_, err = GetDataFromHTMLThroughRegex(htmlText, ytInitialPlayerResponseRegexString)
	if err != nil {
		return err
	}

	_, err = GetDataFromHTMLThroughRegex(htmlText, ytInitialDataRegexString)
	if err != nil {
		return err
	}

	ytCfg, err := GetDataFromHTMLThroughRegex(htmlText, ytCfgRegexString)
	if err != nil {
		return err
	}

	fmt.Println("youtube config: ", ytCfg)
	playerUrl := GetPlayerURL(ytCfg)
	if playerUrl != "" {
		return fmt.Errorf("unable to get the player url for the video")
	}
	_, err = downloadWebPage(playerUrl)
	if err != nil {
		return fmt.Errorf("%v", err)
	}
	fmt.Println("Player URL result: ", playerUrl)
	return nil
}

func downloadWebPage(videoUrl string) (string, error) {
	res, err := http.Get(videoUrl)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	pageContent, err := io.ReadAll(res.Body)
	if err != nil {
		return "", err
	}

	return string(pageContent), nil
}

func GetDataFromHTMLThroughRegex(htmlText, regexString string) (map[string]any, error) {
	// putting the JSON into an empty struct
	var prResponse map[string]interface{}
	ytVarRegex, err := regexp.Compile(regexString)
	if err != nil {
		return prResponse, err
	}

	if !ytVarRegex.MatchString(htmlText) {
		return prResponse, fmt.Errorf("unable to get the video informations from the provided url")
	}

	regexResult := ytVarRegex.FindStringSubmatch(htmlText)

	fmt.Println(regexResult)
	if len(regexResult) <= 1 {
		return prResponse, fmt.Errorf("unable to get the video information from the provided url")
	}
	regexJsonResult := regexResult[1] //strings.ReplaceAll(ytPlayerResponse[1], ";", "")
	// fmt.Println(len(ytPlayerResponse))
	err = json.Unmarshal([]byte(regexJsonResult), &prResponse)
	if err != nil {
		return prResponse, err
	}

	return prResponse, nil
}

func GetPlayerURL(cfgMap map[string]any) string {
	// playerBaseUrl := "https://youtube.com/"
	// find the key in the maps recursively
	playerURL, isPlayerInMap := cfgMap["PLAYER_JS_URL"]
	if isPlayerInMap {
		return fmt.Sprintf("%s%s", ytPlayerBaseUrl, playerURL)
	}
	for key := range cfgMap {
		innerMap, isValueMap := cfgMap[key].(map[string]any)
		playerUrl := ""
		if isValueMap {
			playerUrl = GetPlayerURL(innerMap)
		}
		if playerUrl != "" {
			return playerUrl
		}
	}

	return ""
}

// func GetYtCfg(webHtmlText string, videoId string) error {
// 	clients := map[string]string{
// 		"web": "https://www.youtube.com",
// 		"tv":  "https://www.youtube.com/tv",
// 	}
// 	for client := range clients {
// 		clients[client]
// 	}
// 	return nil
// }

// get ytInitialPlayerResponse
// get ytcfg

// get playerJS URL
// get playerJS
