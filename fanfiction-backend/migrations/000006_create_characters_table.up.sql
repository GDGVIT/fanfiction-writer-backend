CREATE TABLE IF NOT EXISTS characters
(
    id bigserial NOT NULL,
    created_at timestamp with time zone NOT NULL DEFAULT NOW(),
    story_id bigint NOT NULL,
    name citext NOT NULL,
    description text NOT NULL DEFAULT '',
    version integer NOT NULL DEFAULT 1,
    PRIMARY KEY (id),
    UNIQUE (story_id, name)
);