-- +goose Up
ALTER TABLE user_album ADD COLUMN created_on TIMESTAMP DEFAULT NOW();
ALTER TABLE user_album ADD COLUMN updated_on TIMESTAMP DEFAULT NOW();

-- +goose Down
ALTER TABLE user_album DROP COLUMN created_on, updated_on;
