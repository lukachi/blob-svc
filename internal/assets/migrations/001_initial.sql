-- +migrate Up
CREATE TABLE blobs (
    id text primary key,
    data bytea not null
);

-- +migrate Down

DROP TABLE blobs;
