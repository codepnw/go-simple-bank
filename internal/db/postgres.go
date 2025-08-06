package db

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/codepnw/simple-bank/config"
	_ "github.com/lib/pq"
)

func PostgresConnect(cfg *config.EnvConfig) (*sql.DB, error) {
	connectStr := fmt.Sprintf(
		"user=%s password=%s dbname=%s port=%s sslmode=%s",
		cfg.DB.User,
		cfg.DB.Pass,
		cfg.DB.DB,
		cfg.DB.Port,
		cfg.DB.SSL,
	)
	db, err := sql.Open("postgres", connectStr)
	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}

	log.Println("database connected....")

	return db, nil
}
