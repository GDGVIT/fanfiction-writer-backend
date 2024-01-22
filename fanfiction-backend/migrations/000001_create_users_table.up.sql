CREATE TABLE IF NOT EXISTS users
(
    id bigserial,
    created_at timestamp with time zone NOT NULL DEFAULT NOW(),
    name text NOT NULL,
    email citext UNIQUE NOT NULL,
    password_hash bytea NOT NULL,
    activated boolean NOT NULL DEFAULT false,
    version integer NOT NULL DEFAULT 1,
    PRIMARY KEY (id)
);