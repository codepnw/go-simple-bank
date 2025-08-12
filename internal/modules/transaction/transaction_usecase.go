package transaction

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/codepnw/simple-bank/internal/db"
	"github.com/codepnw/simple-bank/internal/modules/account"
)

type TransactionUsecase interface {
	Deposit(ctx context.Context, req *DepositReq) (*Transaction, error)
	Withdraw(ctx context.Context, req *WithdrawReq) (*Transaction, error)
	Transfer(ctx context.Context, req *TransferReq) (*Transaction, error)
	Transactions(ctx context.Context, userID int64) ([]*Transaction, error)
}

type transactionUsecase struct {
	tranRepo   TransasctionRepository
	accUsecase account.AccountUsecase
	txManager  *db.Tx
}

func NewTransactionUsecse(tranRepo TransasctionRepository, accUsecase account.AccountUsecase, txManager *db.Tx) TransactionUsecase {
	return &transactionUsecase{
		tranRepo:   tranRepo,
		accUsecase: accUsecase,
		txManager:  txManager,
	}
}

func (uc *transactionUsecase) Deposit(ctx context.Context, req *DepositReq) (*Transaction, error) {
	result := new(Transaction)

	err := uc.txManager.WithTx(ctx, func(tx *sql.Tx) error {
		if req.Amount <= 0 {
			return errors.New("amount must be greater than zaro")
		}

		// Find Account
		account, err := uc.accUsecase.GetAccountByIDWithTx(ctx, tx, req.ToAccount)
		if err != nil {
			return fmt.Errorf("account not found: %w", err)
		}

		// Update Account
		err = uc.accUsecase.UpdateBalanceWithTx(ctx, tx, account.ID, req.Amount)
		if err != nil {
			return fmt.Errorf("update balance failed: %w", err)
		}

		// Insert Transaction
		result, err = uc.tranRepo.DepositWithTx(ctx, tx, &Transaction{
			ToAccount: account.ID,
			Amount:    req.Amount,
		})
		if err != nil {
			return fmt.Errorf("insert transaction failed: %w", err)
		}
		
		return nil
	})

	return result, err
}

func (uc *transactionUsecase) Withdraw(ctx context.Context, req *WithdrawReq) (*Transaction, error) {
	panic("unimplemented")
}

func (uc *transactionUsecase) Transfer(ctx context.Context, req *TransferReq) (*Transaction, error) {
	panic("unimplemented")
}

func (uc *transactionUsecase) Transactions(ctx context.Context, userID int64) ([]*Transaction, error) {
	panic("unimplemented")
}
