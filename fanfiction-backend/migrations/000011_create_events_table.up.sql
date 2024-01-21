CREATE TABLE IF NOT EXISTS events
(
    id bigserial,
    created_at timestamp with time zone NOT NULL DEFAULT NOW(),
    timeline_id bigint,
    event_time timestamp without time zone,
    title citext,
    description text DEFAULT '',
    details text,
    version integer NOT NULL DEFAULT 1,
    PRIMARY KEY (id),
    UNIQUE (timeline_id, title),
    FOREIGN KEY (timeline_id) REFERENCES timelines (id) MATCH SIMPLE ON DELETE CASCADE
);