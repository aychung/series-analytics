CREATE TABLE IF NOT EXISTS novels (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  prd_no TEXT NOT NULL UNIQUE,
  title TEXT,
  author TEXT,
  publisher TEXT,
  category TEXT,
  detail_url text NOT NULL
  first_seen_at datetime NOT NULL
  first_seen_at datetime NOT NULL
  created_at datetime NOT NULL
  updated_at datetime NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_prd_no ON novels(prd_no);
CREATE INDEX IF NOT EXISTS idx_category ON novels(category);

CREATE TABLE IF NOT EXISTS tags (
  id INTEGER PRIMARY KEY AUTOINCREMENT
  name TEXT NOT NULL UNIQUE
);

CREATE INDEX IF NOT EXISTS idx_tag_name ON tags(name);

CREATE TABLE IF NOT EXISTS novel_tags (
  novel_id INTEGER NOT NULL references novels(id)
  tag_id INTEGER NOT NULL references tags(id)
  created_at datetime not NULL
  PRIMARY KEY (novel_id, tag_id)
);

CREATE TABLE IF NOT EXISTS ignored_novels (
  novel_id INTEG
)

