package account

import (
	"context"
	"time"
)

const queryTimeout = time.Second * 5

type AccountUsecase interface {
	CreateAccount(ctx context.Context, req *accountRequest) (*Account, error)
	GetAccountByID(ctx context.Context, id int64) (*Account, error)
	ListAccounts(ctx context.Context, userID int64) ([]*Account, error)
	UpdateStatusPending(ctx context.Context, id int64) error
	UpdateStatusApproved(ctx context.Context, id int64) error
	UpdateStatusRejected(ctx context.Context, id int64) error
}

type accountUsecase struct {
	repo AccountRepository
}

func NewAccountUsecse(repo AccountRepository) AccountUsecase {
	return &accountUsecase{repo: repo}
}

func (uc *accountUsecase) CreateAccount(ctx context.Context, req *accountRequest) (*Account, error) {
	ctx, cancel := context.WithTimeout(ctx, queryTimeout)
	defer cancel()

	acc := &Account{
		UserID: req.UserID,
		Name:   req.Name,
		Status: StatusPending,
	}

	result, err := uc.repo.Create(ctx, acc)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (uc *accountUsecase) GetAccountByID(ctx context.Context, id int64) (*Account, error) {
	ctx, cancel := context.WithTimeout(ctx, queryTimeout)
	defer cancel()

	return uc.repo.FindByID(ctx, id)
}

func (uc *accountUsecase) ListAccounts(ctx context.Context, userID int64) ([]*Account, error) {
	ctx, cancel := context.WithTimeout(ctx, queryTimeout)
	defer cancel()

	return uc.repo.List(ctx, userID)
}

func (uc *accountUsecase) UpdateStatusPending(ctx context.Context, id int64) error {
	ctx, cancel := context.WithTimeout(ctx, queryTimeout)
	defer cancel()

	return uc.repo.UpdateStatus(ctx, id, string(StatusApproved))
}

func (uc *accountUsecase) UpdateStatusApproved(ctx context.Context, id int64) error {
	ctx, cancel := context.WithTimeout(ctx, queryTimeout)
	defer cancel()

	return uc.repo.UpdateStatus(ctx, id, string(StatusApproved))
}

func (uc *accountUsecase) UpdateStatusRejected(ctx context.Context, id int64) error {
	ctx, cancel := context.WithTimeout(ctx, queryTimeout)
	defer cancel()

	return uc.repo.UpdateStatus(ctx, id, string(StatusRejected))
}
