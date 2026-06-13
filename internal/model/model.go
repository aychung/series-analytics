// Package model contains type definitions
package model

type NovelDetail struct {
	PrdNo     string
	Title     string
	Author    string
	Publisher string
	Category  string
}

type NovelStat struct {
	CommentCount  string
	DownloadCount string
	Rating        float64
}

type Novel struct {
	Detail NovelDetail
	Tags   []string
	Stat   NovelStat
}
