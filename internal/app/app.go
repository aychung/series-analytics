// Package app is the main crawler app
package app

import (
	"context"
	"fmt"

	"series-analytics/internal/crawler"
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

func (a *App) RecordTop100Tags(ctx context.Context, isDryRun bool) error {
	fmt.Printf("> isDryRun: %v\n", isDryRun)
	a.crawler.QueryTop100(ctx)
	return nil
}

func (a *App) RecordNovelStat(ctx context.Context, prdNo string, isDryRun bool) error {
	detail, stat, tags, err := a.crawler.QueryNovelInfo(ctx, prdNo)
	if err != nil {
		return err
	}

	// println("> a.store.StoreNovelDetail")
	if !isDryRun {
		ID, err := a.store.StoreNovelDetail(ctx, detail)
		if err != nil {
			return err
		}
		// println("> a.store.StoreNovelStat")
		err = a.store.StoreNovelStat(ctx, ID, stat)
		if err != nil {
			return err
		}
		// println("> a.store.StoreNovelTags")
		err = a.store.StoreNovelTags(ctx, ID, tags)
		if err != nil {
			return err
		}
	} else {
		fmt.Printf("title: %s\nauthor: %s\npublisher: %s\ncategory: %s\n", detail.Title, detail.Author, detail.Publisher, detail.Category)
		fmt.Printf("rating: %f\ndownload_count: %s\ncomment_count: %s\n", stat.Rating, stat.DownloadCount, stat.CommentCount)
		fmt.Printf("tags: %v\n", tags)
	}

	return nil
}

func (a *App) Close() error {
	// println("Close")
	return a.store.Close()
}
