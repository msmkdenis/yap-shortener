package db

import (
	"embed"
	"fmt"
	"strings"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

//go:embed migration/*.sql
var migrationsFS embed.FS

type Migrations struct {
	migrations *migrate.Migrate
	logger     *zap.Logger
}

func NewMigrations(connection string, logger *zap.Logger) *Migrations {
	dbConfig, err := pgxpool.ParseConfig(connection)
	if err != nil {
		logger.Fatal("Unable to parse connection string", zap.Error(err))
	}

	dbURL := dbURL(dbConfig, sslMode(connection))

	driver, err := iofs.New(migrationsFS, "migration")
	if err != nil {
		logger.Fatal("Unable to create iofs driver", zap.Error(err))
	}

	migrations, err := migrate.NewWithSourceInstance("iofs", driver, dbURL)
	if err != nil {
		logger.Fatal("Unable to create new migrations", zap.Error(err))
	}

	return &Migrations{
		migrations: migrations,
		logger:     logger,
	}
}

func (m *Migrations) MigrateUp() {	
	err := m.migrations.Up()
	if err != nil && err.Error() != "no change" {
		m.logger.Fatal("Unable to up migrations", zap.Error(err))
	}
}

func dbURL(config *pgxpool.Config, sslMode string) string{
	var dbURL strings.Builder

	dbURL.WriteString("postgres://")
	dbURL.WriteString(string(config.ConnConfig.User))
	dbURL.WriteString(":")
	dbURL.WriteString(config.ConnConfig.Password)
	dbURL.WriteString("@")
	dbURL.WriteString(config.ConnConfig.Host)
	dbURL.WriteString(":")
	dbURL.WriteString(fmt.Sprint(config.ConnConfig.Port))
	dbURL.WriteString("/")
	dbURL.WriteString(config.ConnConfig.Database)
	dbURL.WriteString("?sslmode=")
	dbURL.WriteString(sslMode)

	return dbURL.String()
}

func sslMode(connection string) string {
	var con []string = strings.Split(connection, " ")
	var sslMode string
	for _, v := range con {
		pair := strings.Split(v, "=")
		if pair[0] == "sslmode" {
			if pair[1] == "disable" {
				sslMode = pair[1]
			}
		}
	}
	return sslMode	
}


