// Package app is the main crawler app
package app

import (
	"context"
	"fmt"

	"series-analytics/internal/crawler"
	"series-analytics/internal/store"
)

type App struct {
	db	*store.DB
	crawler *crawler.Crawler
}

func New(dbStore *store.DB) (*App, error) {
	c, err := crawler.New();
	if err != nil {
		return nil, err
	}
	return &App{
		db:	dbStore,
		crawler: c,
	}, nil
}

func (a *App) RecordNovelStat(ctx context.Context) error {
	// 14143381
	_, err := a.crawler.QueryNovelStat(ctx, "14207846")
	if err != nil {
		return err
	}

	return nil;
}

func (a *App) Close() error {
	println("Close")
	return a.db.Close()
}

