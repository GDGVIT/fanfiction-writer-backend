CREATE TABLE IF NOT EXISTS stories
(
    id bigserial,
    created_at timestamp with time zone NOT NULL DEFAULT NOW(),
    user_id bigint,
    title citext,
    description text DEFAULT '',
    version integer NOT NULL DEFAULT 1,
    PRIMARY KEY (id),
    UNIQUE (user_id, title),
    FOREIGN KEY (user_id) REFERENCES users (id) MATCH SIMPLE ON DELETE CASCADE
);