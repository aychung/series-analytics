// Package model contains type definitions
package model

import "time"

type NovelDetail struct {
	PrdNo     string
	Title     string
	Author    string
	Publisher string
	Category  string
}

type NovelDetailRecord struct {
	NovelInfo NovelDetail
	CreatedAt time.Time
	UpdatedAt time.Time
}

type NovelStat struct {
	CommentCount  string
	DownloadCount string
	Rating        float64
}

type NovelStatRecord struct {
	Stat       NovelStat
	RecordedAt time.Time
}

type NovelStatRecordList struct {
	PrdNo           string
	NovelID         int
	NovelStatRecord []NovelStatRecord
}
