package server

import (
	"database/sql"

	"github.com/codepnw/simple-bank/config"
	"github.com/codepnw/simple-bank/internal/db"
	"github.com/codepnw/simple-bank/internal/modules/account"
	"github.com/codepnw/simple-bank/internal/modules/auth"
	"github.com/codepnw/simple-bank/internal/modules/transaction"
	"github.com/codepnw/simple-bank/internal/modules/user"
	"github.com/gin-gonic/gin"
)

type routeConfig struct {
	router *gin.Engine
	db     *sql.DB
	tx     *db.Tx
	cfg    *config.EnvConfig
}

func setupRoutes(params *routeConfig) *routeConfig {
	return &routeConfig{
		router: params.router,
		db:     params.db,
		tx:     params.tx,
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

func (r *routeConfig) accountRoutes() {
	accRepo := account.NewAccountRepository(r.db)
	accUsecase := account.NewAccountUsecse(accRepo)
	accHandler := account.NewAccountHandler(accUsecase)

	account := r.router.Group("/accounts")

	account.POST("/", accHandler.CreateAccount)
	account.GET("/", accHandler.ListAccounts)
	account.GET("/:id", accHandler.GetAccountByID)
	account.GET("/:id/pending", accHandler.UpdateStatusPending)
	account.GET("/:id/approved", accHandler.UpdateStatusApproved)
	account.GET("/:id/rejected", accHandler.UpdateStatusRejected)
}

func (r *routeConfig) transactionRoutes() {
	accRepo := account.NewAccountRepository(r.db)
	accUsecase := account.NewAccountUsecse(accRepo)

	tranRepo := transaction.NewTransactionRepository(r.db)
	tranUsecase := transaction.NewTransactionUsecse(tranRepo, accUsecase, r.tx)
	tranHandler := transaction.NewTransactionHandler(tranUsecase)

	route := r.router.Group("/transactions")

	route.POST("/deposit", tranHandler.Deposit)
	route.POST("/withdraw", tranHandler.Withdraw)
	route.POST("/transfer", tranHandler.Transfer)
	route.GET("/:userID", tranHandler.Transactions)
}
