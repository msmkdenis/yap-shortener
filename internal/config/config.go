package config

import (
	"flag"
	"os"
)

type Repository uint

const(
	DataBaseRepository Repository = iota + 1
	FileRepository
	MemoryRepostiory
)

type Config struct {
	URLServer          string
	URLPrefix          string
	FileStoragePath    string
	DataBaseDSN        string
	RepositoryType     Repository

}

func NewConfig() *Config {

	var config = Config{
		URLServer:       "8080",
		URLPrefix:       "http://localhost:8080",
	}

	config.parseFlags()
	config.parseEnv()

	config.RepositoryType = config.newRepositoryType()

	return &config
}

func (c *Config) parseFlags() {
	var URLServer string
	flag.StringVar(&URLServer, "a", ":8080", "Enter URLServer as ip_address:port Or use SERVER_ADDRESS env")

	var URLPrefix string
	flag.StringVar(&URLPrefix, "b", "http://localhost:8080", "Enter URLPrefix as http://ip_address:port Or use BASE_URL env")

	var FileStoragePath string
	flag.StringVar(&FileStoragePath, "f", "", "Enter path for file Or use FILE_STORAGE_PATH env")

	var DataBaseDSN string
	flag.StringVar(&DataBaseDSN, "d", "", "Enter url to connect database as host=host port=port user=postgres password=postgres dbname=dbname sslmode=disable Or use DATABASE_DSN env")

	flag.Parse()

	c.URLServer = URLServer
	c.URLPrefix = URLPrefix
	c.FileStoragePath = FileStoragePath
	c.DataBaseDSN = DataBaseDSN
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

	if envDataBaseDSN := os.Getenv("DATABASE_DSN"); envDataBaseDSN != "" {
		c.DataBaseDSN = envDataBaseDSN
	}
}

func (c *Config) newRepositoryType() Repository {
	if c.DataBaseDSN != "" {
		return DataBaseRepository
	}
	if c.FileStoragePath != "" {
		return FileRepository
	}
	return MemoryRepostiory
}
