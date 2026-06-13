// Package app is the main crawler app
package app

import (
	"context"
	"fmt"

	"series-analytics/internal/crawler"
	"series-analytics/internal/model"
	"series-analytics/internal/store"
)

type App struct {
	store   *store.Store
	crawler *crawler.Crawler
}

func New(dbStore *store.Store) (*App, error) {
	c, err := crawler.New()
	if err != nil {
		return nil, err
	}
	return &App{
		store:   dbStore,
		crawler: c,
	}, nil
}

func (a *App) RecordTop100Tags(ctx context.Context, rankingType model.RankingType, isDryRun bool) error {
	rankList, err := a.crawler.QueryTop100(ctx, rankingType)
	if err != nil {
		return err
	}

	for i, ranking := range rankList {
		if ranking.PrdNo == "" {
			break
		}
		novel, nerr := a.getNovel(ctx, ranking.PrdNo)
		if nerr != nil {
			return nerr
		}
		if !isDryRun {
			ID, nerr := a.store.StoreNovel(ctx, novel, true)
			if nerr != nil {
				return nerr
			}
			rankList[i].NovelID = ID
		}
	}

	if !isDryRun {
		err = a.store.StoreNovelRank(ctx, rankList, rankingType)
	} else {
		fmt.Printf("type: %s\n", rankingType)
		fmt.Printf("%v\n", rankList)
	}

	return err
}

func (a *App) RecordNovelStat(ctx context.Context, prdNo string, isDryRun bool) error {
	novel, err := a.crawler.QueryNovel(ctx, prdNo)
	if err != nil {
		return err
	}

	if !isDryRun {
		_, err = a.store.StoreNovel(ctx, novel, false)
	} else {
		fmt.Printf("novel: %v\n", novel)
	}

	return err 
}


func (a *App) getNovel(ctx context.Context, prdNo string) (novel model.Novel, err error) {
	ID, detail, err := a.store.GetNovelDetailByPrdNo(ctx, prdNo)
	if err != nil || detail.Title == "" {
		return a.crawler.QueryNovel(ctx, prdNo)
	}

	tags, err := a.store.GetNovelTagsByNovelID(ctx, ID)
	if err != nil {
		return
	}
	stats, err := a.store.GetNovelStatByNovelID(ctx, ID, 1)
	if err != nil {
		return
	}

	novel.Detail = detail
	if len(stats.NovelStatRecord) > 0 {
		novel.Stat.CommentCount = stats.NovelStatRecord[0].Stat.CommentCount
		novel.Stat.DownloadCount = stats.NovelStatRecord[0].Stat.DownloadCount
		novel.Stat.Rating = stats.NovelStatRecord[0].Stat.Rating
	}
	novel.Tags = tags
	return
}

func (a *App) Close() {
	a.crawler.Close()
}
