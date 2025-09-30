-- +goose Up
CREATE TABLE albums(id TEXT PRIMARY KEY, name TEXT NOT NULL, coverImageUrl TEXT, releaseDate TEXT, artists VARCHAR(255), spotifyUrl TEXT, 
jsonTrackList VARCHAR(255),created_on TIMESTAMP NOT NULL DEFAULT NOW(), updated_on TIMESTAMP NOT NULL DEFAULT NOW());

-- +goose Down
DROP TABLE albums;