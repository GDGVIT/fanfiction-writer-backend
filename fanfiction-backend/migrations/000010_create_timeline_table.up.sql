CREATE TABLE IF NOT EXISTS timelines
(
    id bigserial,
    created_at timestamp with time zone NOT NULL DEFAULT NOW(),
    story_id bigint,
    name citext,
    version integer NOT NULL DEFAULT 1,
    PRIMARY KEY (id),
    UNIQUE (story_id, name),
    FOREIGN KEY (story_id) REFERENCES stories (id) MATCH SIMPLE ON DELETE CASCADE
);