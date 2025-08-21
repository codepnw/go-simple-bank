package transaction

import (
	"context"
	"testing"

	"github.com/codepnw/simple-bank/internal/db"
	"github.com/codepnw/simple-bank/internal/modules/account"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type testCase struct {
	name        string
	method      string
	req         any
	mockSetup   func(acc *account.AccountUsecaseMock, tranRepo *transactionRepositoryMock)
	expected    any
	expectedErr bool
}

func TestTransactionUsecaseAll(t *testing.T) {
	var tests []*testCase

	// Mock Account
	account1 := &account.Account{ID: int64(1), Balance: 0}
	account2 := &account.Account{ID: int64(2), Balance: 200}
	// Mock Request
	depositReq := &DepositReq{ToAccount: account1.ID, Amount: float64(100)}
	withdrawReq := &WithdrawReq{FromAccount: account2.ID, Amount: float64(50)}
	transferReq := &TransferReq{FromAccount: account2.ID, ToAccount: account1.ID, Amount: float64(50)}
	// Mock Expected
	depositExpected := &Transaction{ToAccount: &account1.ID, Amount: float64(100)}
	withdrawExpected := &Transaction{FromAccount: &account2.ID, Amount: float64(150)}
	transferExpected := &Transaction{FromAccount: &account2.ID, ToAccount: &account1.ID, Amount: float64(150)}

	// Deposit Method
	deposit := &testCase{
		name:   "Deposit success",
		method: "Deposit",
		req:    depositReq,
		mockSetup: func(accUsecase *account.AccountUsecaseMock, tranRepo *transactionRepositoryMock) {
			accUsecase.On("GetAccountByID", mock.Anything, account1.ID).Return(account1, nil)
			accUsecase.On("UpdateBalanceWithTx", mock.Anything, mock.Anything, account1.ID, depositReq.Amount).Return(nil)

			tranResponse := &Transaction{ToAccount: &account1.ID, Amount: float64(100)}
			tranRepo.On("DepositWithTx", mock.Anything, mock.Anything, mock.Anything).Return(tranResponse, nil)
		},
		expected: depositExpected,
	}
	tests = append(tests, deposit)

	// Withdraw Method
	withdraw := &testCase{
		name:   "Withdraw success",
		method: "Withdraw",
		req:    withdrawReq,
		mockSetup: func(acc *account.AccountUsecaseMock, tranRepo *transactionRepositoryMock) {
			acc.On("GetAccountByID", mock.Anything, account2.ID).Return(account2, nil)
			acc.On("UpdateBalanceWithTx", mock.Anything, mock.Anything, account2.ID, -withdrawReq.Amount).Return(nil)

			tranResponse := &Transaction{FromAccount: &account2.ID, Amount: float64(150)}
			tranRepo.On("WithdrawWithTx", mock.Anything, mock.Anything, mock.Anything).Return(tranResponse, nil)
		},
		expected: withdrawExpected,
	}
	tests = append(tests, withdraw)

	// Transfer Method
	transfer := &testCase{
		name:   "Transfer success",
		method: "Transfer",
		req:    transferReq,
		mockSetup: func(acc *account.AccountUsecaseMock, tranRepo *transactionRepositoryMock) {
			acc.On("GetAccountByID", mock.Anything, account1.ID).Return(account1, nil)
			acc.On("GetAccountByID", mock.Anything, account2.ID).Return(account2, nil)

			acc.On("UpdateBalanceWithTx", mock.Anything, mock.Anything, account1.ID, transferReq.Amount).Return(nil)
			acc.On("UpdateBalanceWithTx", mock.Anything, mock.Anything, account2.ID, -transferReq.Amount).Return(nil)

			tranResponse := &Transaction{FromAccount: &account2.ID, ToAccount: &account1.ID, Amount: float64(150)}
			tranRepo.On("TransferWithTx", mock.Anything, mock.Anything, mock.Anything).Return(tranResponse, nil)
		},
		expected: transferExpected,
	}
	tests = append(tests, transfer)

	// Transactions Method
	trans := &testCase{
		name:   "Transactions success",
		method: "Transactions",
		req:    int64(1), // UserID
		mockSetup: func(acc *account.AccountUsecaseMock, tranRepo *transactionRepositoryMock) {
			tranRepo.On("Transactions", mock.Anything, int64(1)).Return([]*Transaction{
				{ID: 1, Amount: 100},
				{ID: 2, Amount: 200},
			}, nil)
		},
		expected: []*Transaction{
			{ID: 1, Amount: 100},
			{ID: 2, Amount: 200},
		},
	}
	tests = append(tests, trans)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tranRepo := NewtransactionRepositoryMockMock()
			accUsecase := account.NewAccountUsecaseMock()
			tx := db.TxMock{}
			uc := NewTransactionUsecse(tranRepo, accUsecase, &tx)

			if tt.mockSetup != nil {
				tt.mockSetup(accUsecase, tranRepo)
			}

			var result any
			var err error

			switch tt.method {
			case "Deposit":
				result, err = uc.Deposit(context.Background(), tt.req.(*DepositReq))
			case "Withdraw":
				result, err = uc.Withdraw(context.Background(), tt.req.(*WithdrawReq))
			case "Transfer":
				result, err = uc.Transfer(context.Background(), tt.req.(*TransferReq))
			case "Transactions":
				result, err = uc.Transactions(context.Background(), tt.req.(int64))
			}

			if tt.expectedErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)

				switch tt.method {
				case "Transactions":
					assert.EqualValues(t, tt.expected.([]*Transaction), result)
				default:
					assert.EqualValues(t, tt.expected.(*Transaction), result)
				}
			}
		})
	}
}
