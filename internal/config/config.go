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

func NewConfig() *Config {

	var config = Config{
		URLServer: "8080",
		URLPrefix: "http://localhost:8080",
	}

	config.ParseFlags()
	config.ParseEnv()

	return &config
}

func (c *Config) ParseFlags() {
	var URLServer string
	flag.StringVar(&URLServer, "a", ":8080", "Enter URLServer as ip_address:port")

	var URLPrefix string
	flag.StringVar(&URLPrefix, "b", "http://localhost:8080", "Enter URLPrefix as http://ip_address:port")

	flag.Parse()

	c.URLServer = URLServer
	c.URLPrefix = URLPrefix
}

func (c *Config) ParseEnv() {

	if envURLServer := os.Getenv("SERVER_ADDRESS"); envURLServer != "" {
		c.URLServer = envURLServer
	}

	if envURLPrefix := os.Getenv("BASE_URL"); envURLPrefix != "" {
		c.URLPrefix = envURLPrefix
	}
}
