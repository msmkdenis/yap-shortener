package config

import (
	"flag"
	"os"
)

type Config struct {
	URLServer       string
	URLPrefix       string
	FileStoragePath string
}

func NewConfig() *Config {

	var config = Config{
		URLServer:       "8080",
		URLPrefix:       "http://localhost:8080",
		FileStoragePath: "/tmp/short-url-db.json",
	}

	config.parseFlags()
	config.parseEnv()

	return &config
}

func (c *Config) parseFlags() {
	var URLServer string
	flag.StringVar(&URLServer, "a", ":8080", "Enter URLServer as ip_address:port Or use SERVER_ADDRESS env")

	var URLPrefix string
	flag.StringVar(&URLPrefix, "b", "http://localhost:8080", "Enter URLPrefix as http://ip_address:port Or use BASE_URL env")

	var FileStoragePath string
	flag.StringVar(&FileStoragePath, "f", "/tmp/short-url-db.json", "Enter path for file Or use FILE_STORAGE_PATH env")

	flag.Parse()

	c.URLServer = URLServer
	c.URLPrefix = URLPrefix
	c.FileStoragePath = FileStoragePath
}

func (c *Config) parseEnv() {

	if envURLServer := os.Getenv("SERVER_ADDRESS"); envURLServer != "" {
		c.URLServer = envURLServer
	}

	if envURLPrefix := os.Getenv("BASE_URL"); envURLPrefix != "" {
		c.URLPrefix = envURLPrefix
	}

	if envFileStoragePath := os.Getenv("FILE_STORAGE_PATH"); envFileStoragePath != "" {
		c.FileStoragePath = envFileStoragePath
	}
}
