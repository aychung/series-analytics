// Package store implements DB store for crawler
package store

import (
	"context"
	"database/sql"
	"fmt"
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

func (s *Store) StoreNovel(ctx context.Context, novel model.Novel, skipStat bool) (id int, err error) {
	tx, err := s.DB.BeginTx(ctx, nil)
	if err != nil {
		return
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback()
		}
	}()

	id, err = storeNovelDetail(ctx, tx, novel.Detail)
	if err != nil {
		return
	}
	if !skipStat {
		err = storeNovelStat(ctx, tx, id, novel.Stat)
		if err != nil {
			return
		}
	}
	err = storeNovelTags(ctx, tx, id, novel.Tags)
	if err != nil {
		return
	}

	err = tx.Commit()
	return
}

func storeNovelDetail(ctx context.Context, tx *sql.Tx, detail model.NovelDetail) (id int, err error) {
	err = tx.QueryRowContext(ctx, `
			INSERT INTO novels (prd_no, title, author, publisher, category)
			VALUES (?, ?, ?, ?, ?)
			ON CONFLICT(prd_no) DO UPDATE SET updated_at = CURRENT_TIMESTAMP
			RETURNING id
		`, detail.PrdNo, detail.Title, detail.Author, detail.Publisher, detail.Category,
	).Scan(&id)
	return
}

func storeNovelStat(ctx context.Context, tx *sql.Tx, novelID int, stat model.NovelStat) error {
	_, err := tx.ExecContext(ctx, `
		INSERT INTO novel_stats (novel_id, rating, comment_count, download_count)
		VALUES (?, ?, ?, ?)
	`, novelID, stat.Rating, stat.CommentCount, stat.DownloadCount)

	return err
}

func storeNovelTags(ctx context.Context, tx *sql.Tx, novelID int, tags []string) (err error) {
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
			INSERT INTO novel_tags (novel_id, tag_id, tag_name)
			VALUES (?, ?, ?)
			ON CONFLICT(novel_id, tag_id) DO NOTHING
			`, novelID, tagID, tag)
		if err != nil {
			return
		}
	}

	return
}

func (s *Store) StoreNovelRank(ctx context.Context, ranking [100]model.NovelRank, rankType model.RankingType) error {
	tx, err := s.DB.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback()
		}
	}()

	tableName := ""
	switch rankType {
	case model.HourlyRank:
		tableName = "novel_hourly_ranking"
	case model.DailyRank:
		tableName = "novel_daily_ranking"
	case model.WeeklyRank:
		tableName = "novel_weekly_ranking"
	case model.MonthlyRank:
		tableName = "novel_monthly_ranking"
	default:
		return fmt.Errorf("unknown ranking type: %s", rankType)
	}
	for _, rank := range ranking {
		if rank.PrdNo == "" {
			break
		}
		_, err = tx.ExecContext(ctx, fmt.Sprintf(`
			INSERT INTO %s (novel_id, ranking)
			VALUES (?, ?)
			`, tableName), rank.NovelID, rank.Rank)
		if err != nil {
			return err
		}
	}
	err = tx.Commit()

	return nil
}

func (s *Store) GetNovelStatByNovelID(ctx context.Context, novelID int, limit int) (statRL NovelStatRecordList, err error) {
	rows, err := s.DB.QueryContext(ctx, `
		SELECT rating, comment_count, download_count, recorded_at
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

func (s *Store) GetNovelTagsByNovelID(ctx context.Context, novelID int) (tags []string, err error) {
	rows, err := s.DB.QueryContext(ctx, `
		SELECT tag_name FROM novel_tags WHERE novel_id = ?;
	`, novelID)
	if err != nil {
		return
	}
	defer rows.Close()

	for rows.Next() {
		var tag string
		err = rows.Scan(&tag)
		if err != nil {
			return
		}
		tags = append(tags, tag)
	}
	return
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

func (s *Store) GetNovelDetailByPrdNo(ctx context.Context, prdNo string) (ID int, detail model.NovelDetail, err error) {
	row := s.DB.QueryRowContext(ctx, `
		SELECT id, prd_no, title, author, publisher, category FROM novels WHERE prd_no = ?;
	`, prdNo)

	err = row.Scan(
		&ID,
		&detail.PrdNo,
		&detail.Title,
		&detail.Author,
		&detail.Publisher,
		&detail.Category,
	)
	return
}

func (s *Store) Close() error {
	// println("DB CLOSE")
	return s.DB.Close()
}
