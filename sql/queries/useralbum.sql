-- name: SaveUserVisitedAlbum :exec
INSERT INTO user_album (id, userId, albumId, isSaved)
VALUES (GEN_RANDOM_UUID(), $1, $2, TRUE);

-- name: GetUserRecentlyVisitedAlbums :many
SELECT * FROM albums WHERE id IN (SELECT DISTINCT user_album.albumId FROM user_album WHERE user_album.userId = $1);

