-- name: CreateRefreshToken :one
INSERT INTO refresh_tokens(token, userId, revokedAt) 
VALUES ($1, $2, $3) RETURNING *;

-- name: GetTokenById :one
SELECT * FROM refresh_tokens WHERE token=$1 LIMIT 1;

-- name: RevokeToken :exec
UPDATE refresh_tokens SET revokedAt=NOW() WHERE token=$1;
