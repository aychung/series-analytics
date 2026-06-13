// Package store implements DB store for crawler
package store

import (
	"context"
	"database/sql"
	"series-analytics/internal/model"
	"time"

	_ "modernc.org/sqlite"
)

type NovelDetailRecord struct {
	NovelInfo model.NovelDetail
	CreatedAt time.Time
	UpdatedAt time.Time
}

type NovelStatRecord struct {
	Stat       model.NovelStat
	RecordedAt time.Time
}

type NovelStatRecordList struct {
	PrdNo           string
	NovelID         int
	NovelStatRecord []NovelStatRecord
}

type Store struct {
	DB *sql.DB
}

func (s *Store) StoreNovelDetail(ctx context.Context, detail model.NovelDetail) (id int, err error) {
	err = s.DB.QueryRowContext(ctx, `
			INSERT INTO novels (prd_no, title, author, publisher, category)
			VALUES (?, ?, ?, ?, ?)
			ON CONFLICT(prd_no) DO UPDATE SET updated_at = CURRENT_TIMESTAMP
			RETURNING id
		`, detail.PrdNo, detail.Title, detail.Author, detail.Publisher, detail.Category,
	).Scan(&id)
	return
}

func (s *Store) StoreNovelStat(ctx context.Context, novelID int, stat model.NovelStat) error {
	_, err := s.DB.ExecContext(ctx, `
		INSERT INTO novel_stats (novel_id, rating, comment_count, download_count)
		VALUES (?, ?, ?, ?)
	`, novelID, stat.Rating, stat.CommentCount, stat.DownloadCount)

	return err
}

func (s *Store) StoreNovelTags(ctx context.Context, novelID int, tags []string) (err error) {
	tx, err := s.DB.BeginTx(ctx, nil)
	if err != nil {
		return
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback()
		}
	}()

	for _, tag := range tags {
		var tagID int
		_, err = tx.ExecContext(ctx, `
			INSERT INTO tags (name) VALUES (?)
			ON CONFLICT(name) DO NOTHING
			`, tag,
		)
		if err != nil {
			return
		}

		err = tx.QueryRowContext(ctx, `
			SELECT id FROM tags WHERE name = ?
			`, tag,
		).Scan(&tagID)
		if err != nil {
			return
		}

		_, err = tx.ExecContext(ctx, `
			INSERT INTO novel_tags (novel_id, tag_id)
			VALUES (?, ?)
			ON CONFLICT(novel_id, tag_id) DO NOTHING
			`, novelID, tagID)
		if err != nil {
			return
		}
	}

	err = tx.Commit()
	return
}

func (s *Store) GetAllNovelStatByNovelID(ctx context.Context, novelID int, limit int) (statRL NovelStatRecordList, err error) {
	rows, err := s.DB.QueryContext(ctx, `
		SELECT 
			rating, 
			comment_count
			download_count,
			recorded_at
		FROM novel_stats WHERE novel_id = ? ORDER BY recorded_at DESC LIMIT ?;
	`, novelID, limit)
	if err != nil {
		return
	}
	defer rows.Close()

	statRL = NovelStatRecordList{
		PrdNo:           "",
		NovelID:         novelID,
		NovelStatRecord: make([]NovelStatRecord, 0, limit),
	}

	for rows.Next() {
		var record NovelStatRecord

		err = rows.Scan(
			&record.Stat.Rating,
			&record.Stat.CommentCount,
			&record.Stat.DownloadCount,
			&record.RecordedAt,
		)
		if err != nil {
			return
		}

		statRL.NovelStatRecord = append(statRL.NovelStatRecord, record)
	}

	return
}

func (s *Store) GetAllNovelStatByPrdNo(ctx context.Context, prdNo string) error {
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
