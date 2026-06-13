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
	"series-analytics/internal/store"

	"github.com/joho/godotenv"
)

func main() {
	mode := flag.String("mode", "",
		"'hourlyStat': for hourly stat recording\n"+
			"'daily100': for daily top100 tag trends recording")
	isDryRun := flag.Bool("dryrun", false, "set dry-run to true to run without actually recording into DB")

	flag.Parse()

	// println("> Starting new context")
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	// println("> Loading configs")
	if err := godotenv.Load(); err != nil {
		log.Fatal(err)
	}
	config, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	// println("> Opening DB")
	dbStore, err := store.Open(config.DBPath)
	if err != nil {
		log.Fatal(err)
	}

	// println("> Running DB migration")
	if err = store.RunMigration(dbStore.DB, config.MigrationPath); err != nil {
		dbStore.Close()
		log.Fatal(err)
	}
	defer dbStore.Close()

	// println("> Starting new app")
	a, err := app.New(dbStore)
	if err != nil {
		log.Fatal(err)
	}

	// println("> Running app")
	start := time.Now()
	switch *mode {
	case "hourlyStat":
		// println("> Running hourly stat")
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

	case "daily100":
		println("> Running daily 100 tags")
		err := a.RecordTop100Tags(ctx, *isDryRun)
		if err != nil {
			log.Fatal(err)
		}
	}
	log.Printf("> Crawler job %s done in %s", *mode, time.Since(start))
}
