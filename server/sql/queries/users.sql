-- name: CreateUser :one
INSERT INTO users(id, email, password, authType)
VALUES (GEN_RANDOM_UUID(), $1, $2, $3) RETURNING *;

-- name: FetchUserByEmail :one
SELECT * FROM users WHERE email= $1 LIMIT 1;
