-- +goose Up
CREATE TABLE refresh_token(token TEXT PRIMARY KEY, created_on TIMESTAMP NOT NULL DEFAULT NOW(), updated_on TIMESTAMP NOT NULL DEFAULT NOW(), revokedAt TIMESTAMP DEFAULT NULL, 
userId UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE);

-- +goose Down
DROP TABLE refresh_token;
