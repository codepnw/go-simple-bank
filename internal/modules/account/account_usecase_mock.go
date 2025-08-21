package account

import (
	"context"
	"database/sql"

	"github.com/stretchr/testify/mock"
)

type AccountUsecaseMock struct {
	mock.Mock
}

func NewAccountUsecaseMock() *AccountUsecaseMock {
	return &AccountUsecaseMock{}
}

func (m *AccountUsecaseMock) CreateAccount(ctx context.Context, req *AccountRequest) (*Account, error) {
	args := m.Called(ctx, req)
	return args.Get(0).(*Account), args.Error(1)
}

func (m *AccountUsecaseMock) GetAccountByID(ctx context.Context, id int64) (*Account, error) {
	args := m.Called(ctx, id)

	res, ok := args.Get(0).(*Account)
	if !ok {
		return nil, args.Error(1)
	}

	return res, args.Error(1)
}

func (m *AccountUsecaseMock) ListAccounts(ctx context.Context, userID int64) ([]*Account, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).([]*Account), args.Error(1)
}

func (m *AccountUsecaseMock) UpdateStatusPending(ctx context.Context, id int64) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *AccountUsecaseMock) UpdateStatusApproved(ctx context.Context, id int64) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *AccountUsecaseMock) UpdateStatusRejected(ctx context.Context, id int64) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *AccountUsecaseMock) UpdateBalanceWithTx(ctx context.Context, tx *sql.Tx, id int64, balance float64) error {
	args := m.Called(ctx, tx, id, balance)
	return args.Error(0)
}
