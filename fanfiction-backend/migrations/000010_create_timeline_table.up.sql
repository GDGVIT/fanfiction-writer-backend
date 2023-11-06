CREATE TABLE IF NOT EXISTS timelines
(
    id bigserial,
    created_at timestamp with time zone,
    story_id bigint,
    name citext,
    version integer,
    PRIMARY KEY (id),
    UNIQUE (story_id, name),
    FOREIGN KEY (story_id) REFERENCES stories (id) MATCH SIMPLE ON DELETE CASCADE
);