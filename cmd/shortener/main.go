package main

import (
	_ "net/http/pprof"
	"os"

	_ "github.com/golang-migrate/migrate/v4/database/postgres"

	"github.com/msmkdenis/yap-shortener/internal/app/shortener"
)

func main() {
	shortener.URLShortenerRun()
	os.Exit(1)
}
