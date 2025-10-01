-- name: CreateAlbum :one
INSERT INTO albums(id, name, numberOfTracks ,coverImageUrl, releaseDate, artists, spotifyUrl, jsonTrackList )
VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING *;

