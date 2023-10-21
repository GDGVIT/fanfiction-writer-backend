CREATE TABLE IF NOT EXISTS characters
(
    id bigserial NOT NULL,
    created_at timestamp with time zone NOT NULL DEFAULT NOW(),
    name citext NOT NULL,
    version integer NOT NULL DEFAULT 1,
    PRIMARY KEY (id)
);