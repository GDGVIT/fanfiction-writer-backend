-- Adds a constraint to make sure sublabel_id is not label_id
ALTER TABLE IF EXISTS sublabels
    ADD CONSTRAINT sublabel_check CHECK (label_id <> sublabel_id);

-- Adds a constraint to make sure blacklist_id is not label_id
ALTER TABLE IF EXISTS blacklist_labels
    ADD CONSTRAINT blacklist_check CHECK (label_id <> blacklist_id);