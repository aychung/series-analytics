CREATE TABLE IF NOT EXISTS novel_hourly_ranking (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  novel_id INTEGER NOT NULL REFERENCES novels(id),
  ranking INTEGER NOT NULL,
  recorded_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_novel_hourly_ranking_novel_id ON novel_hourly_ranking(novel_id);

CREATE TABLE IF NOT EXISTS novel_daily_ranking (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  novel_id INTEGER NOT NULL REFERENCES novels(id),
  ranking INTEGER NOT NULL,
  recorded_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_novel_daily_ranking_novel_id ON novel_daily_ranking(novel_id);

CREATE TABLE IF NOT EXISTS novel_weekly_ranking (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  novel_id INTEGER NOT NULL REFERENCES novels(id),
  ranking INTEGER NOT NULL,
  recorded_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_novel_weekly_ranking_novel_id ON novel_weekly_ranking(novel_id);

CREATE TABLE IF NOT EXISTS novel_monthly_ranking (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  novel_id INTEGER NOT NULL REFERENCES novels(id),
  ranking INTEGER NOT NULL,
  recorded_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_novel_monthly_ranking_novel_id ON novel_monthly_ranking(novel_id);
