package main

import (
	"encoding/json"
	"net/http"

	"github.com/RazafimanantsoaJohnson/brokeMusicApp/internal/spotify"
)

func (cfg *ApiConfig) HandleSearchAlbum(w http.ResponseWriter, r *http.Request) {
	queryParams := r.URL.Query()
	searchQuery := queryParams.Get("query")
	// We will probably want to execute a request -> get the answer, check if the response is 400; authenticate, save the token, redo the request
	// authResponse, err := spotify.Authenticate(cfg.spotifyClientId, cfg.spotifyClientSecret)
	// if err != nil {
	// 	w.WriteHeader(500)
	// 	w.Write([]byte(authResponse.Err))
	// 	return
	// }
	// cfg.renewSpotifyAuth()
	foundAlbums, err := spotify.Search(cfg.spotifyAccessToken.AccessToken, searchQuery)
	if err != nil && err.Error() == spotify.UnvalidAuthErrorMessage {
		err = cfg.renewSpotifyAuth() // we renew the auth and reset the err
		foundAlbums, _ = spotify.Search(cfg.spotifyAccessToken.AccessToken, searchQuery)
	}
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte(err.Error()))
		return
	}

	jsonValue, err := json.Marshal(&foundAlbums)
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte(err.Error()))
		return
	}
	w.WriteHeader(200)
	w.Write(jsonValue)
}

func (cfg *ApiConfig) HandleGetAlbumTracks(w http.ResponseWriter, r *http.Request) {
	albumId := r.PathValue("albumId")
	albumTracks, err := spotify.GetAlbumTracks(cfg.spotifyAccessToken.AccessToken, albumId)

	if err != nil && err.Error() == spotify.UnvalidAuthErrorMessage {
		err = cfg.renewSpotifyAuth() // we renew the auth and reset the err
		albumTracks, err = spotify.GetAlbumTracks(cfg.spotifyAccessToken.AccessToken, albumId)
	}
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte(err.Error()))
		return
	}

	jsonTracks, err := json.Marshal(&albumTracks)
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte(err.Error()))
		return
	}

	w.WriteHeader(200)
	w.Write(jsonTracks)
}

func (cfg *ApiConfig) renewSpotifyAuth() error {
	authResponse, err := spotify.Authenticate(cfg.spotifyClientId, cfg.spotifyClientSecret)
	if err != nil {
		return err
	}
	cfg.spotifyAccessToken = authResponse

	return nil
}
