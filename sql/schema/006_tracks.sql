-- +goose Up
ALTER TABLE tracks ADD COLUMN spotifyId TEXT;
ALTER TABLE tracks ADD COLUMN trackNumber INT;

-- +goose Down
ALTER TABLE tracks DROP COLUMN spotifyId;
ALTER TABLE tracks DROP COLUMN trackNumber;
