package spotify

import (
	"os"
	"testing"

	"github.com/joho/godotenv"
)

func getCredentials() (appId, appSecret string) {
	godotenv.Load("../../.env")
	appId = os.Getenv("SPOTIFY_CLIENTID")
	appSecret = os.Getenv("SPOTIFY_CLIENTSECRET")

	return
}

func TestAuth(t *testing.T) {
	appId, appSecret := getCredentials()
	_, err := Authenticate(appId, appSecret)

	if err != nil {
		t.Errorf("%s", err)
	}
}

func TestFetchAlbum(t *testing.T) {
	appId, appSecret := getCredentials()
	auth, err := Authenticate(appId, appSecret)

	if err != nil {
		t.Errorf("%s", err)
	}
	_, err = Search(auth.AccessToken, "Never Enough")
	if err != nil {
		t.Errorf("%s", err)
	}

	t.Errorf("e")
}

func TestFetchAlbumTracks(t *testing.T) {
	appId, appSecret := getCredentials()
	auth, err := Authenticate(appId, appSecret)

	if err != nil {
		t.Errorf("%s", err)
	}
	_, err = GetAlbumTracks(auth.AccessToken, "3mH6qwIy9crq0I9YQbOuDf")
	if err != nil {
		t.Errorf("%s", err)
	}

	t.Errorf("e")
}
