package server

import (
	"errors"
	"fmt"

	"github.com/codepnw/simple-bank/config"
	"github.com/codepnw/simple-bank/internal/db"
	"github.com/gin-gonic/gin"
)

func Run(cfg *config.EnvConfig) error {
	if cfg == nil {
		return errors.New("config is nil")
	}

	// Init Postgres DB
	pg, err := db.PostgresConnect(cfg)
	if err != nil {
		return fmt.Errorf("database connect failed: %w", err)
	}
	defer pg.Close()

	// Init TX
	tx := db.InitTx(pg)

	// Init gin router
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()

	// Init Routes
	routes := setupRoutes(&routeConfig{
		router: r,
		db:     pg,
		tx:     tx,
		cfg:    cfg,
	})

	routes.authRoutes()
	routes.userRoutes()
	routes.accountRoutes()

	return r.Run(cfg.APP.Port)
}
