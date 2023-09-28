-- +migrate Up
ALTER TABLE users DROP COLUMN salt;

-- +migrate Down
ALTER TABLE users ADD COLUMN salt text not null;
