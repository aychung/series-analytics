// Package store implements DB store for crawler
package store

import (
	"context"
	"database/sql"
	"series-analytics/internal/model"

	_ "modernc.org/sqlite"
)

type Store struct {
	DB *sql.DB
}

func (s Store) StoreNovelDetail(ctx context.Context, detail *model.NovelDetail) (id int, err error) {
	err = s.DB.QueryRowContext(ctx, `
			INSERT INTO novels (prd_no, title, author, publisher, category)
			VALUES (?, ?, ?, ?, ?)
			ON CONFLICT(prd_no) DO UPDATE SET updated_at = CURRENT_TIMESTAMP
			RETURNING id
		`, detail.PrdNo, detail.Title, detail.Author, detail.Publisher, detail.Category,
	).Scan(&id)
	return
}

func (s Store) StoreNovelStat(ctx context.Context, novelID int, stat *model.NovelStat) error {
	_, err := s.DB.ExecContext(ctx, `
		INSERT INTO novel_stats (novel_id, rating, comment_count, download_count)
		VALUES (?, ?, ?, ?)
	`, novelID, stat.Rating, stat.CommentCount, stat.DownloadCount)

	return err
}

func (s Store) GetAllNovelStatByNovelID(ctx context.Context, novelID int, limit int) (*model.NovelStatRecordList, error) {
	rows, err := s.DB.QueryContext(ctx, `
		SELECT 
			rating, 
			comment_count
			download_count,
			recorded_at
		FROM novel_stats WHERE novel_id = ? ORDER BY recorded_at DESC LIMIT ?;
	`, novelID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	novelStatRecordList := model.NovelStatRecordList{
		PrdNo:           "",
		NovelID:         novelID,
		NovelStatRecord: make([]model.NovelStatRecord, 0, limit),
	}

	for rows.Next() {
		var record model.NovelStatRecord

		err := rows.Scan(
			&record.Stat.Rating,
			&record.Stat.CommentCount,
			&record.Stat.DownloadCount,
			&record.RecordedAt,
		)
		if err != nil {
			return nil, err
		}

		novelStatRecordList.NovelStatRecord = append(novelStatRecordList.NovelStatRecord, record)
	}

	return &novelStatRecordList, nil
}

func (s Store) GetAllNovelStatByPrdNo(ctx context.Context, prdNo string) error {
	_, err := s.DB.ExecContext(ctx, `
		SELECT * FROM novel_stats WHERE prd_no = ? ORDER BY recorded_at DESC;
	`, prdNo)
	return err
}

func Open(dbPath string) (*Store, error) {
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, err
	}

	settings := []string{
		"PRAGMA journal_mode = WAL;",
		"PRAGMA busy_timeout = 5000;",
		"PRAGMA foreign_keys = ON;",
	}

	for _, stmt := range settings {
		if _, err := db.Exec(stmt); err != nil {
			db.Close()
			return nil, err
		}
	}

	return &Store{DB: db}, nil
}

func (s *Store) Close() error {
	// println("DB CLOSE")
	return s.DB.Close()
}
