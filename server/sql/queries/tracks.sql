-- name: FetchAlbumTracks :many
SELECT * FROM tracks WHERE  (albumId = $1);

-- name: InsertAlbumTrack :exec
INSERT INTO tracks (id, isAvailable,name, trackNumber,spotifyId,spotifyDuration, spotifyUri, isExplicit, albumId)
VALUES (GEN_RANDOM_UUID(), TRUE, $1, $2, $3, $4, $5, $6, $7);
