package youtube

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"regexp"
	"strings"
)

type YtDlpExtractedJson struct {
	Title     string        `json:"title"`
	Formats   []YtDlpFormat `json:"formats"`
	Thumbnail string        `json:"thumbnail"`
	Channel   string        `json:"channel"`
}

type YtDlpFormat struct {
	FormatId   string `json:"format_id"`
	Url        string `json:"url"`
	Ext        string `json:"ext"`
	Format     string `json:"format"`
	FormatNote string `json:"format_note"`
	// FileSize int    `json:"filesize"`
	// Quality  int    `json:"quality"`
	// 	Height   int    `json:"height"`
	// 	Width    int    `json:"width"`
	// 	AudioExt string `json:"audio_ext"`
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

func DownloadVideo(url, outputPath string) error {
	fileName := fmt.Sprintf("%s.%%(ext)s", outputPath)

	params := []string{"-f ba", url, "-o", fileName}
	cmd := exec.Command("yt-dlp", params...)
	cmdOutput, err := cmd.Output()
	if err != nil {
		return err
	}
	stringOutput := string(cmdOutput)
	downloadSuccessReg, err := regexp.Compile(downloadSuccessRegex)
	if err != nil {
		return err
	}
	fmt.Println(downloadSuccessReg.MatchString(stringOutput))
	if !downloadSuccessReg.MatchString(stringOutput) {
		return fmt.Errorf("something went wrong, not able to download video")
	}
	return nil
}

func GetAudioStreamingUrl(json YtDlpExtractedJson) YtDlpFormat {
	result := YtDlpFormat{}
	for _, format := range json.Formats {
		if format.FormatNote == "medium" && format.FormatId == "140" { // m4a
			result = format
			break
		}
		if format.FormatId == "249" || format.Format == "250" || format.Format == "251" {
			result = format
		}
	}
	return result
}
