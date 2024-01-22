-- Drops the sublabels check constraint
ALTER TABLE sublabels DROP CONSTRAINT sublabel_check;

-- Drops the blacklist labels check constraint
ALTER TABLE blacklist_labels DROP CONSTRAINT blacklist_check;