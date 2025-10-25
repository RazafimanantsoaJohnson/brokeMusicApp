-- name: SaveUserVisitedAlbum :exec
INSERT INTO user_album (id, userId, albumId, isSaved)
VALUES (GEN_RANDOM_UUID(), $1, $2, TRUE);

-- name: GetUserRecentlyVisitedAlbums :many
SELECT albums.name, albums.id, albums.coverimageurl, albums.releasedate, albums.numberoftracks FROM user_album 
INNER JOIN albums ON user_album.albumId=albums.id WHERE userId = $1 ORDER BY user_album.created_on DESC LIMIT 10;

