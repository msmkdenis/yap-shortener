package main

import (
	"github.com/msmkdenis/yap-shortener/internal/config"
	"github.com/msmkdenis/yap-shortener/internal/server"
	"github.com/msmkdenis/yap-shortener/internal/storage"
)

func main() {
	config.AppConfig = *config.NewConfig()
	storage.GlobalRepository = storage.NewMemoryRepository()
	server.InitServer(config.AppConfig.URLServer)
}
