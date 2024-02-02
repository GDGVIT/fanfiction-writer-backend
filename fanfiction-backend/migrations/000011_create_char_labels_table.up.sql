CREATE TABLE IF NOT EXISTS characters_labels
(
    character_id uuid NOT NULL,
    label_id bigint NOT NULL,
    PRIMARY KEY (character_id, label_id),
    FOREIGN KEY (character_id) REFERENCES characters(id) MATCH SIMPLE ON DELETE CASCADE,
    FOREIGN KEY (label_id) REFERENCES labels(id) MATCH SIMPLE ON DELETE CASCADE
);