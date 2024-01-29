package main

import (
	_ "net/http/pprof"

	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/msmkdenis/yap-shortener/internal/app/shortener"
)

func main() {
	shortener.URLShortenerRun()
}
