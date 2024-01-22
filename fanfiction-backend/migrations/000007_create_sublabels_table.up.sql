-- Creates the sublabels table
CREATE TABLE IF NOT EXISTS sublabels
(
    label_id bigint,
    sublabel_id bigint,
    PRIMARY KEY (label_id, sublabel_id),
    FOREIGN KEY (label_id) REFERENCES labels (id) MATCH SIMPLE ON DELETE CASCADE,
    FOREIGN KEY (sublabel_id) REFERENCES labels (id) MATCH SIMPLE ON DELETE CASCADE
);
