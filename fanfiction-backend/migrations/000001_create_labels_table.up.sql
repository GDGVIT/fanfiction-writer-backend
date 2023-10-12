-- Creates the labels table
CREATE TABLE IF NOT EXISTS labels
(
    id bigserial,
    created_at timestamp with time zone NOT NULL DEFAULT NOW(),
    name citext NOT NULL,
    version integer NOT NULL DEFAULT 1,
    PRIMARY KEY (id),
    UNIQUE (name)
);