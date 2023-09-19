-- +migrate Up
CREATE TABLE blobs (
    id text primary key,
    value text not null
);

-- +migrate Down

DROP TABLE blobs;
