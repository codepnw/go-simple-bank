package server

import (
	"database/sql"

	"github.com/codepnw/simple-bank/config"
	"github.com/codepnw/simple-bank/internal/db"
	"github.com/codepnw/simple-bank/internal/middleware"
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
	mid    middleware.Auth
}

func setupRoutes(params *routeConfig) *routeConfig {
	return &routeConfig{
		router: params.router,
		db:     params.db,
		tx:     params.tx,
		cfg:    params.cfg,
		mid:    middleware.AuthMiddleware(params.cfg),
	}
}

// Route: Auth
func (r *routeConfig) authRoutes() {
	userRepo := user.NewUserRepository(r.db)
	userUsecase := user.NewUserUsecase(userRepo)

	authUsecase := auth.NewAuthUsecase(r.cfg, userUsecase)
	authHandler := auth.NewAuthHandler(authUsecase)

	// Public
	public := r.router.Group("/auth")
	{
		public.POST("/register", authHandler.Register)
		public.POST("/login", authHandler.Login)
	}
}

// Route: Users
func (r *routeConfig) userRoutes() {
	userRepo := user.NewUserRepository(r.db)
	userUsecase := user.NewUserUsecase(userRepo)
	userHandler := user.NewUserHandler(userUsecase)

	// Group: All Role
	authorized := r.router.Group("/users/profile", r.mid.Authorized())
	{
		authorized.GET("/", userHandler.GetProfile)
		authorized.PATCH("/", userHandler.UpdateProfile)
	}

	// Group: Admin Role
	permission := r.router.Group("/users", r.mid.Authorized(), r.mid.Permissions(user.RoleAdmin))
	{
		permission.POST("/", userHandler.CreateUser)
		permission.GET("/", userHandler.GetUsers)
		permission.GET("/:id", userHandler.GetUser)
		permission.DELETE("/:id", userHandler.DeleteUser)
	}
}

// Route: Accounts
func (r *routeConfig) accountRoutes() {
	accRepo := account.NewAccountRepository(r.db)
	accUsecase := account.NewAccountUsecse(accRepo)
	accHandler := account.NewAccountHandler(accUsecase)

	authorized := r.router.Group("/accounts", r.mid.Authorized())
	{
		authorized.POST("/", accHandler.CreateAccount)
		authorized.GET("/user/:userID", accHandler.ListAccounts)
		authorized.GET("/:id", accHandler.GetAccountByID)
	}

	// Group: Staff, Admin
	permission := r.router.Group("/accounts", r.mid.Authorized(), r.mid.Permissions(user.RoleStaff, user.RoleAdmin))
	{
		permission.GET("/:id/pending", accHandler.UpdateStatusPending)
		permission.GET("/:id/approved", accHandler.UpdateStatusApproved)
		permission.GET("/:id/rejected", accHandler.UpdateStatusRejected)
	}
}

// Route: Transactions
func (r *routeConfig) transactionRoutes() {
	accRepo := account.NewAccountRepository(r.db)
	accUsecase := account.NewAccountUsecse(accRepo)

	tranRepo := transaction.NewTransactionRepository(r.db)
	tranUsecase := transaction.NewTransactionUsecse(tranRepo, accUsecase, r.tx)
	tranHandler := transaction.NewTransactionHandler(tranUsecase)

	// Public
	public := r.router.Group("/transactions")
	{
		public.POST("/deposit", tranHandler.Deposit)
	}

	// Authorized
	authorized := r.router.Group("/transactions", r.mid.Authorized())
	{
		authorized.POST("/withdraw", tranHandler.Withdraw)
		authorized.POST("/transfer", tranHandler.Transfer)
		authorized.GET("/", tranHandler.TransactionsByCurrentUser)
	}

	// Group: Staff, Admin
	permission := r.router.Group("/transactions/user", r.mid.Authorized(), r.mid.Permissions(user.RoleStaff, user.RoleAdmin))
	{
		permission.GET("/:id", tranHandler.TransactionsByUserID)
	}
}
