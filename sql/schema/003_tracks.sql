-- +goose Up
CREATE TABLE tracks(id UUID PRIMARY KEY, youtubeId TEXT, name TEXT NOT NULL, spotifyDuration INT, spotifyUri TEXT, isExplicit BOOLEAN , isAvailable BOOLEAN NOT NULL DEFAULT FALSE,
youtubeUrlType TEXT, youtubeUrl TEXT, fileUrl TEXT,albumId TEXT REFERENCES albums(id) ,created_on TIMESTAMP NOT NULL DEFAULT NOW(), updated_on TIMESTAMP NOT NULL DEFAULT NOW());

-- +goose Down
DROP TABLE tracks;