-- +migrate Up
ALTER TABLE ONLY blobs DROP COLUMN owner;

ALTER TABLE ONLY blobs ADD COLUMN owner_id text not null;
ALTER TABLE ONLY blobs ADD CONSTRAINT blobs_users_fk FOREIGN KEY (owner_id) REFERENCES users(id) ON UPDATE CASCADE ON DELETE CASCADE;

-- +migrate Down
ALTER TABLE ONLY blobs DROP CONSTRAINT blobs_users_fk;
ALTER TABLE ONLY blobs DROP COLUMN owner_id;

ALTER TABLE ONLY blobs ADD COLUMN owner text not null;
