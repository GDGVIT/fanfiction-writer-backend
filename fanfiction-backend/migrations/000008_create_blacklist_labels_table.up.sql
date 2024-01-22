-- Creates the blacklist labels table
CREATE TABLE IF NOT EXISTS blacklist_labels
(
    label_id bigint,
    blacklist_id bigint,
    PRIMARY KEY (label_id, blacklist_id),
    FOREIGN KEY (label_id) REFERENCES labels (id) MATCH SIMPLE ON DELETE CASCADE,
    FOREIGN KEY (blacklist_id) REFERENCES labels (id) MATCH SIMPLE ON DELETE CASCADE
);
