CREATE EXTENSION IF NOT EXISTS citext;

-- Creates the labels table
CREATE TABLE IF NOT EXISTS labels
(
    id bigserial,
    created_at timestamp with time zone NOT NULL DEFAULT NOW(),
    story_id bigint NOT NULL,
    name citext NOT NULL,
    version integer NOT NULL DEFAULT 1,
    PRIMARY KEY (id),
    FOREIGN KEY (story_id) REFERENCES stories (id) MATCH SIMPLE ON DELETE CASCADE,
    UNIQUE (story_id, name)
);
