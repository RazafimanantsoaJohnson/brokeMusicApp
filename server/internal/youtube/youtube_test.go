package youtube

import (
	"fmt"
	"os"
	"testing"

	"github.com/joho/godotenv"
)

// func TestGetStreamUrl(t *testing.T) {
// 	err := GetStreamUrl("https://www.youtube.com/watch?v=PVRbKHXwM58")
// 	if err != nil {
// 		t.Error(err)
// 	}
// }

func getCredentials() (apiKey string) {
	godotenv.Load("../../.env")
	apiKey = os.Getenv("YOUTUBE_APIKEY")

	return
}

func TestSearchYoutubeVideo(t *testing.T) {
	apiKey := getCredentials()
	cases := []string{
		"Frank Ocean - Godspeed",
		"Frank Ocean - Wiseman",
	}
	for c := range cases {
		_, err := Search(apiKey, cases[c])

		if err != nil {
			t.Error(err)
		}
	}
}

func TestYtDlpCmd(t *testing.T) {
	cases := [][]string{
		[]string{"https://www.youtube.com/watch?v=KmqrB3I-26Y"},
		[]string{"https://www.youtube.com/watch?v=oCnotRXfR_M", "https://www.youtube.com/watch?v=xjoBP7SDgaY", "https://www.youtube.com/watch?v=KmqrB3I-26Y"},
	}

	for c := range cases {
		_, err := CallYtDlpCmd(cases[c])
		if err != nil {
			fmt.Println(cases[c])
			t.Error(err)
		}
	}
}
