package server

import (
	"fmt"

	"github.com/codepnw/simple-bank/config"
	"github.com/codepnw/simple-bank/internal/db"
	"github.com/gin-gonic/gin"
)

func Run(cfg *config.EnvConfig) error {
	db, err := db.PostgresConnect(cfg)
	if err != nil {
		return fmt.Errorf("database connect failed: %w", err)
	}
	defer db.Close()

	r := gin.Default()

	return r.Run(cfg.APP.Port)
}
