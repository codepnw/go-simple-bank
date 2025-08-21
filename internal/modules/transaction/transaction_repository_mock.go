package transaction

import (
	"context"
	"database/sql"

	"github.com/stretchr/testify/mock"
)

type transactionRepositoryMock struct {
	mock.Mock
}

func NewtransactionRepositoryMockMock() *transactionRepositoryMock {
	return &transactionRepositoryMock{}
}

func (m *transactionRepositoryMock) DepositWithTx(ctx context.Context, tx *sql.Tx, input *Transaction) (*Transaction, error) {
	args := m.Called(ctx, tx, input)

	res, ok := args.Get(0).(*Transaction)
	if !ok {
		return nil, args.Error(1)
	}

	return res, args.Error(1)
}

func (m *transactionRepositoryMock) Transactions(ctx context.Context, userID int64) ([]*Transaction, error) {
	args := m.Called(ctx, userID)

	res, ok := args.Get(0).([]*Transaction)
	if !ok {
		return nil, args.Error(1)
	}

	return res, args.Error(1)
}

func (m *transactionRepositoryMock) TransferWithTx(ctx context.Context, tx *sql.Tx, input *Transaction) (*Transaction, error) {
	args := m.Called(ctx, tx, input)

	res, ok := args.Get(0).(*Transaction)
	if !ok {
		return nil, args.Error(1)
	}

	return res, args.Error(1)
}

func (m *transactionRepositoryMock) WithdrawWithTx(ctx context.Context, tx *sql.Tx, input *Transaction) (*Transaction, error) {
	args := m.Called(ctx, tx, input)

	res, ok := args.Get(0).(*Transaction)
	if !ok {
		return nil, args.Error(1)
	}

	return res, args.Error(1)
}
