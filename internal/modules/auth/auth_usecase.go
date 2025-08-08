package auth

import (
	"context"
	"log"
	"time"

	"github.com/codepnw/simple-bank/config"
	"github.com/codepnw/simple-bank/internal/modules/user"
	"github.com/codepnw/simple-bank/internal/utils/security"
)

const queryTimeout = time.Second * 5

type AuthUsecase interface {
	Login(ctx context.Context, req *authRequest) (*JWTTokenResponse, error)
	Register(ctx context.Context, req *user.UserRequest) (*JWTTokenResponse, error)
}

type authUsecase struct {
	userUsecase user.UserUsecase
	jwt         *security.Token
}

func NewAuthUsecase(cfg *config.EnvConfig, userUsecase user.UserUsecase) AuthUsecase {
	return &authUsecase{
		userUsecase: userUsecase,
		jwt:         security.InitJWT(cfg),
	}
}

func (uc *authUsecase) Login(ctx context.Context, req *authRequest) (*JWTTokenResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, queryTimeout)
	defer cancel()

	user, err := uc.userUsecase.GetUserByEmail(ctx, req.Email)
	if err != nil {
		return nil, err
	}

	err = security.ComparePassword(user.Password, req.Password)
	if err != nil {
		return nil, err
	}

	token, err := uc.jwtTokenResponse(&security.TokenUser{
		ID:    user.ID,
		Email: user.Email,
		Role:  "user", // TODO: change later
	})
	if err != nil {
		return nil, err
	}

	return token, nil
}

func (uc *authUsecase) Register(ctx context.Context, req *user.UserRequest) (*JWTTokenResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, queryTimeout)
	defer cancel()

	created, err := uc.userUsecase.Create(ctx, req)
	if err != nil {
		return nil, err
	}

	user := &security.TokenUser{
		ID:    created.ID,
		Email: created.Email,
		Role:  "user", // TODO: change later
	}

	return uc.jwtTokenResponse(user)
}

func (uc *authUsecase) jwtTokenResponse(user *security.TokenUser) (*JWTTokenResponse, error) {
	accessToken, err := uc.jwt.GenerateAccessToken(user)
	log.Println("acc", accessToken)
	if err != nil {
		return nil, err
	}

	refreshToken, err := uc.jwt.GenerateRefreshToken(user)
	if err != nil {
		return nil, err
	}

	token := &JWTTokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}

	return token, nil
}
