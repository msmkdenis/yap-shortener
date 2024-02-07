package main

import (
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "net/http/pprof"

	"github.com/msmkdenis/yap-shortener/internal/app/shortener"
)

func main() {
	shortener.URLShortenerRun()
}
