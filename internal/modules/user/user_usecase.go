package user

import (
	"context"
	"time"

	"github.com/codepnw/simple-bank/internal/utils/security"
)

const queryTimeout = time.Second * 5

type UserUsecase interface {
	Create(ctx context.Context, req *UserRequest) (*User, error)
	GetUserByID(ctx context.Context, id int64) (*User, error)
	GetUserByEmail(ctx context.Context, email string) (*User, error)
	GetUsers(ctx context.Context) ([]*User, error)
	Update(ctx context.Context, id int64, req *UserUpdateRequest) (*User, error)
	Delete(ctx context.Context, id int64) error
}

type userUsecase struct {
	repo UserRepository
}

func NewUserUsecase(repo UserRepository) UserUsecase {
	return &userUsecase{repo: repo}
}

func (uc *userUsecase) Create(ctx context.Context, req *UserRequest) (*User, error) {
	ctx, cancel := context.WithTimeout(ctx, queryTimeout)
	defer cancel()

	hashPassword, err := security.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	user := &User{
		Email:     req.Email,
		Password:  hashPassword,
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Phone:     req.Phone,
	}

	created, err := uc.repo.Create(ctx, user)
	if err != nil {
		return nil, err
	}

	return created, nil
}

func (uc *userUsecase) GetUserByID(ctx context.Context, id int64) (*User, error) {
	ctx, cancel := context.WithTimeout(ctx, queryTimeout)
	defer cancel()

	return uc.repo.FindByID(ctx, id)
}

func (uc *userUsecase) GetUserByEmail(ctx context.Context, email string) (*User, error) {
	ctx, cancel := context.WithTimeout(ctx, queryTimeout)
	defer cancel()

	return uc.repo.FindByEmail(ctx, email)
}

func (uc *userUsecase) GetUsers(ctx context.Context) ([]*User, error) {
	ctx, cancel := context.WithTimeout(ctx, queryTimeout)
	defer cancel()

	return uc.repo.List(ctx)
}

func (uc *userUsecase) Update(ctx context.Context, id int64, req *UserUpdateRequest) (*User, error) {
	ctx, cancel := context.WithTimeout(ctx, queryTimeout)
	defer cancel()

	user, err := uc.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if req.FirstName != nil {
		user.FirstName = *req.FirstName
	}

	if req.LastName != nil {
		user.LastName = *req.LastName
	}

	if req.Phone != nil {
		user.Phone = req.Phone
	}

	now := time.Now()
	user.UpdatedAt = &now

	if err = uc.repo.Update(ctx, user); err != nil {
		return nil, err
	}

	return user, nil
}

func (uc *userUsecase) Delete(ctx context.Context, id int64) error {
	ctx, cancel := context.WithTimeout(ctx, queryTimeout)
	defer cancel()

	return uc.repo.Delete(ctx, id)
}
