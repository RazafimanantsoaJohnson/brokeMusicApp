package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"log"
	"net/http"

	"github.com/RazafimanantsoaJohnson/brokeMusicApp/internal/database"
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
	go saveAlbumsInDB(cfg, foundAlbums)
}

func saveAlbumsInDB(cfg *ApiConfig, searchResponse spotify.SearchResponse) {
	albums := searchResponse.Albums.Items
	for i := range albums {
		album := albums[i]
		_, err := cfg.db.CreateAlbum(context.Background(), database.CreateAlbumParams{
			ID:   album.Id,
			Name: album.Name,
			// Artists: album.Artists,
			Numberoftracks: int32(album.TotalTracks),
			Releasedate:    sql.NullString{String: album.ReleaseDate},
			Spotifyurl:     sql.NullString{String: album.AlbumUrl},
			Coverimageurl:  sql.NullString{String: album.Images[1].Url},
		})
		if err.Error() == "pq: duplicate key value violates unique constraint \"albums_pkey\"" {
			// fmt.Println(err)
			continue
		} else {
			log.Fatalf(err.Error()) // should keep the log somewhere instead of crash the system
		}
	}
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
