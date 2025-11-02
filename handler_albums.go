package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"log"
	"net/http"

	"github.com/RazafimanantsoaJohnson/brokeMusicApp/internal/database"
	"github.com/RazafimanantsoaJohnson/brokeMusicApp/internal/spotify"
	"github.com/google/uuid"
)

type UserVisitedAlbumResponse struct {
	AlbumId        string `json:"album_id"`
	Artist         string `json:"artist"`
	Name           string `json:"album_name"`
	CoverImageUrl  string `json:"album_cover_url"`
	ReleaseDate    string `json:"release_date"`
	NumberOfTracks int    `json:"number_of_tracks"`
	CreatedOn      string `json:"created_on"`
}

func HandleSearchAlbum(cfg *ApiConfig, curUserId uuid.UUID, w http.ResponseWriter, r *http.Request) {
	queryParams := r.URL.Query()
	searchQuery := queryParams.Get("query")
	// We will probably want to execute a request -> get the answer, check if the response is 400; authenticate, save the token, redo the request
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
			Releasedate:    sql.NullString{String: album.ReleaseDate, Valid: true},
			Spotifyurl:     sql.NullString{String: album.AlbumUrl, Valid: true},
			Coverimageurl:  sql.NullString{String: album.Images[1].Url, Valid: true},
			Artists:        sql.NullString{String: album.Artists[0].Name, Valid: true},
		})
		if err != nil && err.Error() == "pq: duplicate key value violates unique constraint \"albums_pkey\"" {
			continue
		} else if err != nil {
			log.Fatal(err.Error()) // should keep the log somewhere instead of crash the system
		}
	}
}

func (cfg *ApiConfig) renewSpotifyAuth() error {
	authResponse, err := spotify.Authenticate(cfg.spotifyClientId, cfg.spotifyClientSecret)
	if err != nil {
		return err
	}
	cfg.spotifyAccessToken = authResponse

	return nil
}

func HandleGetUserVisitedAlbums(cfg *ApiConfig, curUserId uuid.UUID, w http.ResponseWriter, r *http.Request) {
	userVisitedAlbums, err := cfg.db.GetUserRecentlyVisitedAlbums(context.Background(), uuid.NullUUID{Valid: true, UUID: curUserId})
	result := []UserVisitedAlbumResponse{}
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte(err.Error()))
		return
	}

	for _, userAlbum := range userVisitedAlbums {
		result = append(result, UserVisitedAlbumResponse{
			AlbumId:        userAlbum.ID,
			Artist:         userAlbum.Artists.String,
			Name:           userAlbum.Name,
			CoverImageUrl:  userAlbum.Coverimageurl.String,
			ReleaseDate:    userAlbum.Releasedate.String,
			NumberOfTracks: int(userAlbum.Numberoftracks),
		})
	}

	returnJson(w, result)
}
