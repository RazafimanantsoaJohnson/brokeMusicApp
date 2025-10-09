-- name: FetchAlbumTracks :many
SELECT * FROM tracks WHERE  (albumId = $1) ORDER BY trackNumber;

-- name: InsertAlbumTrack :exec
INSERT INTO tracks (id, isAvailable,name, trackNumber,spotifyId,spotifyDuration, spotifyUri, isExplicit, albumId, youtubeid)
VALUES (GEN_RANDOM_UUID(), TRUE, $1, $2, $3, $4, $5, $6, $7, $8);

-- name: FetchTrack :one
SELECT * FROM tracks WHERE (id = $1);

-- name: InsertTrackYoutubeUrl :one
UPDATE tracks SET youtubeUrl= $2 WHERE (id=$1) RETURNING *;
