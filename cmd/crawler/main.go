package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"series-analytics/internal/app"
	"series-analytics/internal/config"
	"series-analytics/internal/model"
	"series-analytics/internal/store"

	"github.com/joho/godotenv"
)

func main() {
	mode := flag.String("mode", "",
		"'stat': for hourly stat recording\n"+
			"'top100': for top100 tag trends recording")
	isDryRun := flag.Bool("dryrun", false, "set dry-run to true to run without actually recording into DB")
	rankType := flag.String("rankType", "", "hourly, daily, weekly, monthly")

	flag.Parse()

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	if err := godotenv.Load(); err != nil {
		log.Fatal(err)
	}
	config, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	dbStore, err := store.Open(config.DBPath)
	if err != nil {
		log.Fatal(err)
	}

	if err = store.RunMigration(dbStore.DB, config.MigrationPath); err != nil {
		dbStore.Close()
		log.Fatal(err)
	}
	defer dbStore.Close()

	a, err := app.New(dbStore)
	if err != nil {
		log.Fatal(err)
	}
	defer a.Close()

	log.Printf("> Running %s - %s", *mode, *rankType)
	start := time.Now()
	switch *mode {
	case "stat":
		novelPrdNo := []string{
			"14253468",
			"14143381",
		}
		errs := []error{}
		for _, prdNo := range novelPrdNo {
			runErr := a.RecordNovelStat(ctx, prdNo, *isDryRun)
			if runErr != nil {
				errs = append(errs, runErr)
			}
		}

		if len(errs) == len(novelPrdNo) {
			log.Fatal(errs)
		} else if len(errs) > 0 {
			log.Print(errs)
		}

	case "top100":
		isValidRankType := false
		var rType model.RankingType
		switch *rankType {
		case "hourly":
			isValidRankType = true
			rType = model.HourlyRank
		case "daily":
			isValidRankType = true
			rType = model.DailyRank
		case "weekly":
			isValidRankType = true
			rType = model.WeeklyRank
		case "monthly":
			isValidRankType = true
			rType = model.MonthlyRank
		}
		if !isValidRankType {
			log.Fatal("> valid rankType option is required for running top100 mode")
		}
		err := a.RecordTop100Tags(ctx, rType, *isDryRun)
		if err != nil {
			log.Fatal(err)
		}

	default:
		log.Printf("> Invalid mode provided. ('stat' or 'top100')");
	}
	log.Printf("> Job %s - %s done in %s", *mode, *rankType, time.Since(start))
}
