package main

import (
	"fmt"
	_ "net/http/pprof"

	_ "github.com/golang-migrate/migrate/v4/database/postgres"

	"github.com/msmkdenis/yap-shortener/internal/app/shortener"
)

var (
	buildVersion = "N/A"
	buildDate    = "N/A"
	buildCommit  = "N/A"
)

func printGreeting() {
	fmt.Printf("Build version: %s\n", buildVersion)
	fmt.Printf("Build date: %s\n", buildDate)
	fmt.Printf("Build commit: %s\n\n", buildCommit)
}

func main() {
	printGreeting()
	shortener.URLShortenerRun()
}
