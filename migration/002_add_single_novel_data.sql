CREATE TABLE IF NOT EXISTS novel_data (
  id INTEGER PRIMARY KEY AUTOINCREMENT
  novel_id INTEGER NOT NULL references novels(id)
  prd_no TEXT NOT NULL references novels(prd_no)
  recorded_at datetime NOT NULL
  rating DECIMAL(4,2) NOT NULL
  download_count TEXT NOT NULL
  comment_count TEXT NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_novel_data_novel_id ON novel_data(name);

CREATE INDEX IF NOT EXISTS idx_novels_title ON novels(title);

