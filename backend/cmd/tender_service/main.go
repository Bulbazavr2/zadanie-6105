package main

import (
	"database/sql"
	"log"
	"tender_srevice/internal/config"
	"tender_srevice/internal/app/server"
	"tender_srevice/internal/repository"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Printf("Error loading .env file: %v", err)
	}

	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}
	
	err = config.RunMigrations(cfg)
	if err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}
	
	// Создаем подключение к базе данных
	db, err := sql.Open("postgres", cfg.PostgresConn)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Создаем репозиторий
	repo := repository.NewPostgresRepository(db)
	
	// Создаем сервер, передавая конфигурацию и репозиторий
	srv := server.New(cfg, repo)
	if err := srv.Run(); err != nil {
		log.Fatalf("Error running server: %v", err)
	}
}