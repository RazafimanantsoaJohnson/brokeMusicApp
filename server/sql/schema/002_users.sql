-- +goose Up
CREATE TABLE users( id UUID PRIMARY KEY, username TEXT, email TEXT NOT NULL UNIQUE, password TEXT NOT NULL, authType TEXT NOT NULL, created_on TIMESTAMP
NOT NULL DEFAULT NOW(), updated_on TIMESTAMP NOT NULL DEFAULT NOW());

-- +goose Down
DROP TABLE users;
