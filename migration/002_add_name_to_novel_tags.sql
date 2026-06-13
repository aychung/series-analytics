ALTER TABLE novel_tags ADD COLUMN tag_name TEXT NULL REFERENCES tags(name);

UPDATE novel_tags AS nt
SET tag_name = t.name
FROM tags AS t
WHERE nt.tag_id = t.id;

ALTER TABLE novel_tags ALTER COLUMN tag_name SET NOT NULL;

CREATE INDEX IF NOT EXISTS idx_novel_tags_tag_name ON novel_tags(tag_name);

