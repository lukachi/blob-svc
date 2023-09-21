-- +migrate Up
CREATE TABLE blobs (
    id text primary key not null,
    value text not null
);

CREATE TABLE users (
    id text primary key not null,
    login text not null,
    password text not null,
    salt text not null,
    username text not null
);

-- +migrate Down

DROP TABLE blobs;

DROP TABLE users;
