CREATE TABLE IF NOT EXISTS novels (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  prd_no TEXT NOT NULL UNIQUE,
  title TEXT NOT NULL,
  author TEXT NOT NULL,
  publisher TEXT NOT NULL,
  category TEXT NOT NULL,
  created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_novels_prd_no ON novels(prd_no);
CREATE INDEX IF NOT EXISTS idx_novels_category ON novels(category);
CREATE INDEX IF NOT EXISTS idx_novels_title ON novels(title);

CREATE TABLE IF NOT EXISTS tags (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  name TEXT NOT NULL UNIQUE
);

CREATE INDEX IF NOT EXISTS idx_tags_name ON tags(name);

CREATE TABLE IF NOT EXISTS novel_tags (
  novel_id INTEGER NOT NULL REFERENCES novels(id),
  tag_id INTEGER NOT NULL REFERENCES tags(id),
  PRIMARY KEY (novel_id, tag_id)
);

CREATE TABLE IF NOT EXISTS ignored_novels (
  novel_id INTEGER PRIMARY KEY REFERENCES novels(id)
);

CREATE TABLE IF NOT EXISTS novel_stats (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  novel_id INTEGER NOT NULL REFERENCES novels(id),
  rating DECIMAL(4,2) NOT NULL,
  download_count TEXT NOT NULL,
  comment_count TEXT NOT NULL,
  recorded_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_novel_stats_novel_id ON novel_stats(novel_id);
