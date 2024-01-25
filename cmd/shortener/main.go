package main

import (
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/msmkdenis/yap-shortener/internal/app/shortener"
	_ "net/http/pprof"
)

func main() {
	shortener.URLShortenerRun()
}
