package youtube

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"regexp"
	"strings"
)

type YtDlpExtractedJson struct {
	Title   string `json:"title"`
	Formats []struct {
		FormatId string `json:"format_id"`
		Url      string `json:"url"`
		Ext      string `json:"ext"`
		Format   string `json:"format"`
		// FileSize int    `json:"filesize"`
		// Quality  int    `json:"quality"`
		// 	Height   int    `json:"height"`
		// 	Width    int    `json:"width"`
		// 	AudioExt string `json:"audio_ext"`
	} `json:"formats"`
	Thumbnail string `json:"thumbnail"`
	Channel   string `json:"channel"`
}

func CallYtDlpCmd(urls []string) ([]YtDlpExtractedJson, error) {
	params := []string{"--skip-download", "--dump-json"}
	params = append(params, urls...)
	cmd := exec.Command("yt-dlp", params...)
	cmdOutput, err := cmd.Output()
	result := []YtDlpExtractedJson{}
	if err != nil {
		return result, err
	}
	txtCmdOutput := string(cmdOutput)
	jsonStringOutput := ""

	// fmt.Println(txtCmdOutput)
	jsonRegex, err := regexp.Compile(ytDlpGreedyJsonRegex)
	if err != nil {
		return result, err
	}

	// fmt.Println("Index new line: ", strings.Index(txtCmdOutput, "\n"))
	// fmt.Println("Number of lines: ", strings.Count(txtCmdOutput, "\n"))
	if !strings.Contains(txtCmdOutput, "\n") {
		jsonStringOutput = fmt.Sprintf("[%v]", txtCmdOutput)
		err = json.Unmarshal([]byte(jsonStringOutput), &result)
		if err != nil {
			return result, err
		}
	} else {
		lines := strings.Split(txtCmdOutput, "\n")
		for line := range lines {
			if !jsonRegex.MatchString(lines[line]) {
				continue
			}
			jsonInLine := jsonRegex.FindString(lines[line])
			newExtractedVideo := YtDlpExtractedJson{}
			err = json.Unmarshal([]byte(jsonInLine), &newExtractedVideo)
			if err != nil {
				continue
			}
			result = append(result, newExtractedVideo)
		}
	}

	return result, nil
}
