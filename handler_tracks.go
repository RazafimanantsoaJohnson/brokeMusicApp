package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/RazafimanantsoaJohnson/brokeMusicApp/internal/database"
	"github.com/RazafimanantsoaJohnson/brokeMusicApp/internal/spotify"
	"github.com/RazafimanantsoaJohnson/brokeMusicApp/internal/youtube"
	"github.com/google/uuid"
)

type trackResponse struct {
	DiscNumber  int    `json:"disc_number"` // most album are just 1 disc (some might be 2)
	Duration    int    `json:"duration_ms"`
	Explicit    bool   `json:"explicit"`
	Id          string `json:"id"`
	Name        string `json:"name"`
	TrackNumber int    `json:"track_number"`
	TrackUri    string `json:"uri"`
	IsAvailable bool   `json:"is_available"`
	YoutubeURL  string `json:"youtube_url"`
	FileURL     string `json:"file_url"`
}

func HandleGetAlbumTracks(cfg *ApiConfig, curUserId uuid.UUID, w http.ResponseWriter, r *http.Request) {
	albumId := r.PathValue("albumId")

	queriedAlbum, err := cfg.db.GetAlbumFromSpotifyId(context.Background(), albumId)
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte(err.Error()))
	}
	result, areTracksLoaded, err := fetchAlbumTracks(cfg, albumId)
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte(err.Error()))
	}
	if areTracksLoaded {
		fmt.Println("the tracks are loaded from DB")
		jsonResult, err := json.Marshal(&result)
		if err != nil {
			w.WriteHeader(500)
			w.Write([]byte(err.Error()))
			return
		}
		w.WriteHeader(200)
		w.Write(jsonResult)
		return
	}

	albumTracks, err := spotify.GetAlbumTracks(cfg.spotifyAccessToken.AccessToken, albumId)

	if err != nil && err.Error() == spotify.UnvalidAuthErrorMessage {
		err = cfg.renewSpotifyAuth() // we renew the auth and reset the err
		albumTracks, _ = spotify.GetAlbumTracks(cfg.spotifyAccessToken.AccessToken, albumId)
	}
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte(err.Error()))
		return
	}

	if !areTracksLoaded {
		go saveAlbumTracksInDB(cfg, albumId, albumTracks, queriedAlbum)
	}

	jsonTracks, err := json.Marshal(&albumTracks.Tracks.Items)
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte(err.Error()))
		return
	}

	err = cfg.db.SaveUserVisitedAlbum(context.Background(), database.SaveUserVisitedAlbumParams{
		Userid:  uuid.NullUUID{Valid: true, UUID: curUserId},
		Albumid: sql.NullString{Valid: true, String: albumId},
	})
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte(err.Error()))
		return
	}

	fmt.Println("we are returning tracks from spotify")
	w.WriteHeader(200)
	w.Write(jsonTracks)
}

func HandleGetTrack(cfg *ApiConfig, curUserId uuid.UUID, w http.ResponseWriter, r *http.Request) {
	trackId := r.PathValue("trackId")
	queryParams := r.URL.Query()
	retryParam := queryParams.Get("retry")
	isRetry := !(retryParam == "")
	id, err := uuid.Parse(trackId)
	if err != nil {
		w.WriteHeader(400)
		w.Write([]byte("the provided track ID is not valid"))
		return
	}
	dbTrack, err := cfg.db.FetchTrack(context.Background(), id)
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte(err.Error()))
		return
	}
	// add a forced refresh
	if !dbTrack.Youtubeurl.Valid || isRetry {
		resultChan := make(chan YtDlpTaskResult)
		mutex.Lock()
		pushTask(&Tasks, YtDlpTask{
			YoutubeId:  dbTrack.Youtubeid.String,
			AlbumId:    dbTrack.Albumid.String,
			Priority:   1,
			ResultChan: resultChan,
		})
		mutex.Unlock()

		extracted := <-resultChan
		audioFormat := youtube.GetAudioStreamingUrl(extracted.result)

		dbTrack.Youtubeurl.Valid = true
		dbTrack.Youtubeurl.String = audioFormat.Url
	}

	result := trackResponse{
		Id:          dbTrack.ID.String(),
		Name:        dbTrack.Name,
		Duration:    int(dbTrack.Spotifyduration.Int32),
		Explicit:    dbTrack.Isexplicit.Bool,
		TrackNumber: int(dbTrack.Tracknumber.Int32),
		IsAvailable: dbTrack.Isavailable,
		YoutubeURL:  dbTrack.Youtubeurl.String,
		FileURL:     dbTrack.Fileurl.String,
	}

	err = cfg.db.SaveUserVisitedAlbum(context.Background(), database.SaveUserVisitedAlbumParams{
		Userid:  uuid.NullUUID{Valid: true, UUID: curUserId},
		Albumid: sql.NullString{Valid: true, String: dbTrack.Albumid.String},
	})
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte(err.Error()))
		return
	}

	returnJson(w, result)
	cfg.db.InsertTrackYoutubeUrl(context.Background(), database.InsertTrackYoutubeUrlParams{
		ID:         dbTrack.ID,
		Youtubeurl: dbTrack.Youtubeurl,
	})
}

func fetchAlbumTracks(cfg *ApiConfig, albumId string) ([]trackResponse, bool, error) {
	// it doesn't return errors when there is no tracks in DB
	result := make([]trackResponse, 0)
	existingAlbumTracks, err := cfg.db.FetchAlbumTracks(context.Background(), sql.NullString{String: albumId, Valid: true})
	if err != nil {
		fmt.Println(result)
		return result, false, err
	}
	if len(existingAlbumTracks) <= 0 {
		return result, false, nil
	}

	for _, tmpTrack := range existingAlbumTracks {
		result = append(result, trackResponse{
			Id:          tmpTrack.ID.String(),
			Name:        tmpTrack.Name,
			Duration:    int(tmpTrack.Spotifyduration.Int32),
			Explicit:    tmpTrack.Isexplicit.Bool,
			TrackNumber: int(tmpTrack.Tracknumber.Int32),
			IsAvailable: tmpTrack.Isavailable,
			YoutubeURL:  tmpTrack.Youtubeurl.String,
			FileURL:     tmpTrack.Fileurl.String,
		})
	}
	return result, true, nil
}

func saveAlbumTracksInDB(cfg *ApiConfig, albumId string, tracks spotify.AlbumResponse, album database.GetAlbumFromSpotifyIdRow) error {
	fmt.Println("We are saving the tracks in DB")
	for i := range tracks.Tracks.Items {
		track := tracks.Tracks.Items[i]
		searchQuery := fmt.Sprintf("%s %s", album.Artists.String, track.Name)
		ytSearchResult, err := youtube.Search(cfg.ytApiKey, searchQuery)
		ytTrackId := sql.NullString{String: "", Valid: false}
		if err != nil {
			return err
		}
		if len(ytSearchResult.Items) > 0 {
			ytTrackId = sql.NullString{String: ytSearchResult.Items[0].Id.VideoId, Valid: true}
		}

		newTrack, err := cfg.db.InsertAlbumTrack(context.Background(), database.InsertAlbumTrackParams{
			Name:            track.Name,
			Tracknumber:     sql.NullInt32{Int32: int32(track.TrackNumber), Valid: true},
			Isexplicit:      sql.NullBool{Bool: track.Explicit, Valid: true},
			Albumid:         sql.NullString{String: albumId, Valid: true},
			Spotifyid:       sql.NullString{String: track.Id, Valid: true},
			Spotifyduration: sql.NullInt32{Int32: int32(track.Duration), Valid: true},
			Spotifyuri:      sql.NullString{String: track.TrackUri, Valid: true},
			Youtubeid:       ytTrackId,
		})
		if err != nil {
			return err
		}

		if !ytTrackId.Valid {
			err = cfg.db.SetTrackAsUnavailable(context.Background(), newTrack.ID)
			if err != nil {
				return err
			}
			return fmt.Errorf("unable to find track")
		}
		// launching yt-dlp task
		mutex.Lock()
		pushTask(&Tasks, YtDlpTask{
			YoutubeId: ytSearchResult.Items[0].Id.VideoId,
			AlbumId:   albumId,
			TrackId:   newTrack.ID,
			Priority:  0,
		})
		mutex.Unlock()
	}
	return nil
}
