-- name: CreateRefreshToken :one
INSERT INTO refresh_token(token, userId) 
VALUES ($1, $2) RETURNING *;

-- name: GetTokenById :one
SELECT * FROM refresh_token WHERE token=$1 LIMIT 1;

-- name: RevokeToken :exec
UPDATE refresh_token SET revokedAt=NOW() WHERE token=$1;
