package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"series-analytics/internal/app"
	"series-analytics/internal/config"
	"series-analytics/internal/store"

	"github.com/joho/godotenv"
)

func main() {
	println("> Starting new context")
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	println("> Loading configs")
	if err := godotenv.Load(); err != nil {
		log.Fatal(err)
	}
	config, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	println("> Opening DB")
	dbStore, err := store.Open(config.DBPath)
	if err != nil {
		log.Fatal(err)
	}

	println("> Running DB migration")
	if err = store.RunMigration(dbStore.DB, config.DBPath); err != nil {
		dbStore.Close()
		log.Fatal(err)
	}
	defer dbStore.Close()

	println("> Starting new app")
	a, err := app.New(dbStore)
	if err != nil {
		log.Fatal(err)
	}

	println("> Running app")
	start := time.Now()
	err = a.RecordNovelStat(ctx)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("> Job done in %s", time.Since(start))
}

