-- +goose Up
ALTER TABLE albums ADD COLUMN created_by UUID 

-- +goose Down

