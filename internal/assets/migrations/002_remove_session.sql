-- +migrate Up
DROP TABLE sessions;

-- +migrate Down

CREATE TABLE sessions (
    id text primary key not null,
    user_id text not null,
    created_at timestamp not null,
    expires_at timestamp not null
);
