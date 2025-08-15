package middleware

import (
	"strings"

	"github.com/codepnw/simple-bank/config"

	"github.com/codepnw/simple-bank/internal/modules/user"
	"github.com/codepnw/simple-bank/internal/utils"
	"github.com/codepnw/simple-bank/internal/utils/response"
	"github.com/codepnw/simple-bank/internal/utils/security"
	"github.com/gin-gonic/gin"
)

type Auth interface {
	Authorized() gin.HandlerFunc
	Permissions(roles ...user.UserRole) gin.HandlerFunc
}

type auth struct {
	cfg   *config.EnvConfig
	token *security.Token
}

func AuthMiddleware(cfg *config.EnvConfig) Auth {
	return &auth{
		cfg:   cfg,
		token: security.InitJWT(cfg),
	}
}

func (a *auth) Authorized() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authHeader := ctx.Request.Header.Get("Authorization")
		if authHeader == "" {
			response.Unauthorized(ctx, "auth header is missing")
			ctx.Abort()
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			response.Unauthorized(ctx, "invalid token format")
			ctx.Abort()
			return
		}

		token := parts[1]
		u, err := a.token.VerifyAccessToken(token)
		if err != nil {
			response.Unauthorized(ctx, err.Error())
			ctx.Abort()
			return
		}

		ctx.Set(utils.ContextKeyUser, u)
		ctx.Next()
	}
}

func (a *auth) Permissions(roles ...user.UserRole) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		u, err := user.CurrentUser(ctx)
		if err != nil {
			response.Unauthorized(ctx, err.Error())
			ctx.Abort()
			return
		}

		for _, role := range roles {
			if u.Role == role {
				ctx.Next()
				return
			}
		}

		response.Unauthorized(ctx, "permission denied")
		ctx.Abort()
	}
}
