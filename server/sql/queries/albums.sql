-- name: CreateAlbum :one
INSERT INTO albums(id, name, numberOfTracks ,coverImageUrl, releaseDate, artists, spotifyUrl, jsonTrackList )
RETURNING *;

