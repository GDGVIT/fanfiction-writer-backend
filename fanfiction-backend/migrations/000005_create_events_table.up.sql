CREATE TABLE IF NOT EXISTS events
(
    id uuid NOT NULL DEFAULT uuid_generate_v4(),
    created_at timestamp with time zone NOT NULL DEFAULT NOW(),
    character_id uuid,
    event_time timestamp without time zone,
    title citext,
    description text DEFAULT '',
    details text,
    version integer NOT NULL DEFAULT 1,
    PRIMARY KEY (id),
    UNIQUE (character_id, title),
    FOREIGN KEY (character_id) REFERENCES characters (id) MATCH SIMPLE ON DELETE CASCADE
);