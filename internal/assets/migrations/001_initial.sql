-- +migrate Up
CREATE TABLE blobs (
    id text primary key not null,
    value jsonb not null,
    owner_id text not null
);

CREATE TABLE users (
    id text primary key not null,
    login text not null,
    password text not null,
    username text not null
);

ALTER TABLE ONLY blobs ADD CONSTRAINT blobs_users_fk FOREIGN KEY (owner_id) REFERENCES users(id) ON UPDATE CASCADE ON DELETE CASCADE;

-- +migrate Down
ALTER TABLE ONLY blobs DROP CONSTRAINT blobs_users_fk;

DROP TABLE blobs;

DROP TABLE users;
