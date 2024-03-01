// Package config contains the configuration for the application.
package config

import (
	"encoding/json"
	"flag"
	"os"

	"go.uber.org/zap"
)

// Repository represents the type of the repository.
type Repository uint

// RepositoryType represents the type of the repository.
const (
	DataBaseRepository Repository = iota + 1
	FileRepository
	MemoryRepostiory
)

type jsonConfig struct {
	URLServer       string `json:"url_server"`
	URLPrefix       string `json:"url_prefix"`
	FileStoragePath string `json:"file_storage_path"`
	DataBaseDSN     string `json:"database_dsn"`
	SecretKey       string `json:"secret_key"`
	TokenName       string `json:"token_name"`
	EnableHTTPS     string `json:"enable_https"`
	TrustedSubnet   string `json:"trusted_subnet"`
}

// Config represents the configuration for the application.
type Config struct {
	URLServer       string
	URLPrefix       string
	FileStoragePath string
	DataBaseDSN     string
	RepositoryType  Repository
	SecretKey       string
	TokenName       string
	EnableHTTPS     string
	ConfigFile      string
	TrustedSubnet   string
}

// NewConfig creates a new Config instance with default values and returns a pointer to it.
//
// No parameters.
// Returns a pointer to a Config instance.
func NewConfig(logger *zap.Logger) *Config {
	config := Config{
		URLServer: "8080",
		URLPrefix: "http://localhost:8080",
	}

	config.parseFlags()
	config.parseEnv()
	if config.ConfigFile != "" {
		if err := config.parseJSONConfig(); err != nil {
			logger.Error("unable to parse config file", zap.Error(err))
		}
	}

	config.RepositoryType = config.newRepositoryType()
	return &config
}

// user=postgres password=postgres host=localhost database=yap-shortener sslmode=disable
func (c *Config) parseFlags() {
	var URLServer string
	flag.StringVar(&URLServer, "a", "localhost:8080", "Enter URLServer as ip_address:port Or use SERVER_ADDRESS env")

	var URLPrefix string
	flag.StringVar(&URLPrefix, "b", "http://localhost:8080", "Enter URLPrefix as http://ip_address:port Or use BASE_URL env")

	var FileStoragePath string
	flag.StringVar(&FileStoragePath, "f", "", "Enter path for file Or use FILE_STORAGE_PATH env")

	var DataBaseDSN string
	flag.StringVar(&DataBaseDSN, "d", "", "Enter url to connect database as host=host port=port user=postgres password=postgres dbname=dbname sslmode=disable Or use DATABASE_DSN env")

	var SecretKey string
	flag.StringVar(&SecretKey, "k", "supersecretkey", "Enter secret key Or use SECRET_KEY env")

	var EnableHTTPS string
	flag.StringVar(&EnableHTTPS, "s", "false", "Enable HTTPS Or use ENABLE_HTTPS env")

	var TokenName string
	flag.StringVar(&TokenName, "n", "token", "Enter token name Or use TOKEN_NAME env")

	var ConfigFile string
	flag.StringVar(&ConfigFile, "c", "", "Enter path to config file Or use CONFIG env")

	var TrustedSubnet string
	flag.StringVar(&TrustedSubnet, "t", "", "Enter trusted subnet Or use TRUSTED_SUBNET env")

	flag.Parse()

	c.URLServer = URLServer
	c.URLPrefix = URLPrefix
	c.FileStoragePath = FileStoragePath
	c.DataBaseDSN = DataBaseDSN
	c.SecretKey = SecretKey
	c.EnableHTTPS = EnableHTTPS
	c.TokenName = TokenName
	c.ConfigFile = ConfigFile
	c.TrustedSubnet = TrustedSubnet
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

	if envSecretKey := os.Getenv("SECRET_KEY"); envSecretKey != "" {
		c.SecretKey = envSecretKey
	}

	if envEnableHTTPS := os.Getenv("ENABLE_HTTPS"); envEnableHTTPS != "" {
		c.EnableHTTPS = envEnableHTTPS
	}

	if envTokenName := os.Getenv("TOKEN_NAME"); envTokenName != "" {
		c.TokenName = envTokenName
	}

	if envConfigFile := os.Getenv("CONFIG"); envConfigFile != "" {
		c.ConfigFile = envConfigFile
	}

	if envTrustedSubnet := os.Getenv("TRUSTED_SUBNET"); envTrustedSubnet != "" {
		c.TrustedSubnet = envTrustedSubnet
	}
}

func (c *Config) parseJSONConfig() error {
	configFile, err := os.Open(c.ConfigFile)
	if err != nil {
		return err
	}

	var config jsonConfig
	err = json.NewDecoder(configFile).Decode(&config)
	if err != nil {
		return err
	}

	if c.URLServer == "" {
		c.URLServer = config.URLServer
	}

	if c.URLPrefix == "" {
		c.URLPrefix = config.URLPrefix
	}

	if c.SecretKey == "" {
		c.SecretKey = config.SecretKey
	}

	if c.EnableHTTPS == "" {
		c.EnableHTTPS = config.EnableHTTPS
	}

	if c.TokenName == "" {
		c.TokenName = config.TokenName
	}

	if c.TrustedSubnet == "" {
		c.TrustedSubnet = config.TrustedSubnet
	}

	return configFile.Close()
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
