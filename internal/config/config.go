package config

import (
	"flag"
	"os"
)

var AppConfig Config

type Config struct {
	URLServer string
	URLPrefix string
}

func InitConfig() Config {

	var URLServer string
	flag.StringVar(&URLServer, "a", ":8080", "Enter URLServer as ip_address:port")

	var URLPrefix string
	flag.StringVar(&URLPrefix, "b", "http://localhost:8080", "Enter URLPrefix as http://ip_address:port")

	flag.Parse()

	var configuration Config

	if envURLServer := os.Getenv("SERVER_ADDRESS"); envURLServer != "" {
		URLServer = envURLServer
	}

	if envURLPrefix := os.Getenv("BASE_URL"); envURLPrefix != "" {
		URLPrefix = envURLPrefix
	}

	configuration.URLServer = URLServer
	configuration.URLPrefix = URLPrefix

	return configuration
}
