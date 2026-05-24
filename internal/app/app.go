// Package app is the main crawler app
package app

import (
	"context"

	"series-analytics/internal/store"
	// "series-analytics/internal/crawler"
)

type App struct {
	db	*store.DB
//	crawler *crawler.Crawler
}

func New(dbStore *store.DB) (*App, error) {
	return &App{
		db:	dbStore,
//		crawler: crawler.New(dbStore),
	}, nil
}

func (a *App) Run(ctx context.Context) error {
	println("Run")
	return nil;
//	return a.crawler.Run(ctx)
}

func (a *App) Close() error {
	println("Close")
	return a.db.Close()
}

