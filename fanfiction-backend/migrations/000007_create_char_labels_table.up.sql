CREATE TABLE IF NOT EXISTS characters_labels
(
    characters_id bigint NOT NULL,
    labels_id bigint NOT NULL,
    PRIMARY KEY (characters_id, labels_id),
    FOREIGN KEY (characters_id) REFERENCES characters(id) MATCH SIMPLE ON DELETE CASCADE,
    FOREIGN KEY (labels_id) REFERENCES labels(id) MATCH SIMPLE ON DELETE CASCADE
);