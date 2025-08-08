package server

import (
	"database/sql"

	"github.com/codepnw/simple-bank/config"
	"github.com/codepnw/simple-bank/internal/modules/auth"
	"github.com/codepnw/simple-bank/internal/modules/user"
	"github.com/gin-gonic/gin"
)

type routeConfig struct {
	router *gin.Engine
	db     *sql.DB
	cfg    *config.EnvConfig
}

func setupRoutes(params *routeConfig) *routeConfig {
	return &routeConfig{
		router: params.router,
		db:     params.db,
		cfg:    params.cfg,
	}
}

func (r *routeConfig) authRoutes() {
	userRepo := user.NewUserRepository(r.db)
	userUsecase := user.NewUserUsecase(userRepo)

	authUsecase := auth.NewAuthUsecase(r.cfg, userUsecase)
	authHandler := auth.NewAuthHandler(authUsecase)

	pub := r.router.Group("/auth")
	pub.POST("/register", authHandler.Register)
	pub.POST("/login", authHandler.Login)
}

func (r *routeConfig) userRoutes() {
	userRepo := user.NewUserRepository(r.db)
	userUsecase := user.NewUserUsecase(userRepo)
	userHandler := user.NewUserHandler(userUsecase)

	user := r.router.Group("/users")
	user.POST("/", userHandler.CreateUser)
	user.GET("/:id", userHandler.GetUser)
	user.GET("/", userHandler.GetUsers)
	user.PATCH("/", userHandler.UpdateUser)
	user.DELETE("/", userHandler.DeleteUser)
}
