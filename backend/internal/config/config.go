package config

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	_ "github.com/lib/pq"
	"database/sql"
)

type Config struct {
	ServerAddress string
	PostgresConn  string
}

func Load() (*Config, error) {
	serverAddr := os.Getenv("SERVER_ADDRESS")
	if serverAddr == "" {
		return nil, fmt.Errorf("SERVER_ADDRESS environment variable is not set")
	}

	postgresConn := os.Getenv("POSTGRES_CONN")
	if postgresConn == "" {
		return nil, fmt.Errorf("POSTGRES_CONN environment variable is not set")
	}

	cfg := &Config{
		ServerAddress: serverAddr,
		PostgresConn:  postgresConn,
	}

	logDatabaseConnection(cfg)

	return cfg, nil
}

func RunMigrations(cfg *Config) error {
	db, err := sql.Open("postgres", cfg.PostgresConn)
	if err != nil {
		return fmt.Errorf("failed to connect to database: %v", err)
	}
	defer db.Close()

	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		return fmt.Errorf("failed to get current file path")
	}
	migrationsDir := filepath.Join(filepath.Dir(filename), "..", "bd", "migrations")

	files, err := filepath.Glob(filepath.Join(migrationsDir, "*.sql"))
	if err != nil {
		return fmt.Errorf("failed to read migration files: %v", err)
	}

	for _, file := range files {
		content, err := os.ReadFile(file)
		if err != nil {
			return fmt.Errorf("failed to read migration file %s: %v", file, err)
		}

		_, err = db.Exec(string(content))
		if err != nil {
			return fmt.Errorf("failed to execute migration %s: %v", file, err)
		}

		log.Printf("Executed migration: %s", file)
	}

	return nil
}

func logDatabaseConnection(cfg *Config) {
	log.Printf("Подключение к базе данных установлено. Строка подключения: %s", cfg)
}
