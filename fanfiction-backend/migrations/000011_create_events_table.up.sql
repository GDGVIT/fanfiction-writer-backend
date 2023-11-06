CREATE TABLE IF NOT EXISTS events
(
    id bigserial,
    created_at timestamp with time zone,
    timeline_id bigint,
    event_time timestamp without time zone,
    title citext,
    description text,
    details text,
    version integer,
    PRIMARY KEY (id),
    UNIQUE (timeline_id, title),
    FOREIGN KEY (timeline_id) REFERENCES timelines (id) MATCH SIMPLE ON DELETE CASCADE
);