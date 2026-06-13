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

type NovelRank struct {
	NovelID int
	Rank int
	PrdNo string
}

type RankingType string

const (
	HourlyRank RankingType = "HourlyRank"
	DailyRank   RankingType = "DailyRank"
	WeeklyRank  RankingType = "WeeklyRank"
	MonthlyRank RankingType = "MonthlyRank"
)

func (r RankingType) ToQuery() string {
	var q string
	switch r {
	case HourlyRank:
		q = "HOURLY"
	case DailyRank:
		q = "DAILY"
	case WeeklyRank:
		q = "WEEKLY"
	case MonthlyRank:
		q = "MONTHLY"
	} 
	return q
}

