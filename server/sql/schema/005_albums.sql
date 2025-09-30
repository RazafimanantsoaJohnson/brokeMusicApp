-- +goose Up
ALTER TABLE albums ADD COLUMN numberOfTracks INT NOT NULL;

-- +goose Down
UPDATE TABLE albums DROP COLUMN numberOfTracks;
