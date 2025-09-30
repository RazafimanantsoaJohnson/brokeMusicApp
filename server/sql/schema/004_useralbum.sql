-- +goose Up
CREATE TABLE user_album(id UUID PRIMARY KEY, userId UUID REFERENCES users(id), albumId TEXT REFERENCES albums(id), isSaved BOOLEAN, isDownloaded BOOLEAN );

-- +goose Down
DROP TABLE user_album;
