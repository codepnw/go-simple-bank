package transaction

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/codepnw/simple-bank/internal/db"
	"github.com/codepnw/simple-bank/internal/modules/account"
	"github.com/codepnw/simple-bank/internal/utils/errs"
)

const queryTimeout = time.Second * 5

type TransactionUsecase interface {
	Deposit(ctx context.Context, req *DepositReq) (*Transaction, error)
	Withdraw(ctx context.Context, req *WithdrawReq) (*Transaction, error)
	Transfer(ctx context.Context, req *TransferReq) (*Transaction, error)
	Transactions(ctx context.Context, userID int64) ([]*Transaction, error)
}

type transactionUsecase struct {
	tranRepo   TransasctionRepository
	accUsecase account.AccountUsecase
	txManager  db.TxManager
}

func NewTransactionUsecse(tranRepo TransasctionRepository, accUsecase account.AccountUsecase, txManager db.TxManager) TransactionUsecase {
	return &transactionUsecase{
		tranRepo:   tranRepo,
		accUsecase: accUsecase,
		txManager:  txManager,
	}
}

func (uc *transactionUsecase) Deposit(ctx context.Context, req *DepositReq) (*Transaction, error) {
	ctx, cancel := context.WithTimeout(ctx, queryTimeout)
	defer cancel()

	result := new(Transaction)

	// Find Account
	account, err := uc.accUsecase.GetAccountByID(ctx, req.ToAccount)
	if err != nil {
		return nil, errs.ErrAccountNotFound
	}

	// Tx Transaction
	err = uc.txManager.WithTx(ctx, func(tx *sql.Tx) error {
		// Update Account
		err = uc.accUsecase.UpdateBalanceWithTx(ctx, tx, account.ID, req.Amount)
		if err != nil {
			return fmt.Errorf("update balance failed: %w", err)
		}

		// Insert Transaction
		result, err = uc.tranRepo.DepositWithTx(ctx, tx, &Transaction{
			ToAccount: &account.ID,
			Amount:    req.Amount,
		})
		if err != nil {
			return fmt.Errorf("insert transaction failed: %w", err)
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (uc *transactionUsecase) Withdraw(ctx context.Context, req *WithdrawReq) (*Transaction, error) {
	ctx, cancel := context.WithTimeout(ctx, queryTimeout)
	defer cancel()

	result := new(Transaction)

	// Find Account
	account, err := uc.accUsecase.GetAccountByID(ctx, req.FromAccount)
	if err != nil {
		return nil, errs.ErrAccountNotFound
	}

	if req.Amount > float64(account.Balance) {
		return nil, errs.ErrAmountGreaterAccBalance
	}

	// Tx Transaction
	err = uc.txManager.WithTx(ctx, func(tx *sql.Tx) error {
		// Update Account
		err = uc.accUsecase.UpdateBalanceWithTx(ctx, tx, account.ID, -req.Amount)
		if err != nil {
			return fmt.Errorf("update balance failed: %w", err)
		}

		// Insert Transaction
		result, err = uc.tranRepo.WithdrawWithTx(ctx, tx, &Transaction{
			FromAccount: &account.ID,
			Amount:      req.Amount,
		})
		if err != nil {
			return fmt.Errorf("insert transaction failed: %w", err)
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (uc *transactionUsecase) Transfer(ctx context.Context, req *TransferReq) (*Transaction, error) {
	ctx, cancel := context.WithTimeout(ctx, queryTimeout)
	defer cancel()

	result := new(Transaction)

	// Find From Account
	fromAcc, err := uc.accUsecase.GetAccountByID(ctx, req.FromAccount)
	if err != nil {
		return nil, errs.ErrAccountNotFound
	}

	// Find To Account
	toAcc, err := uc.accUsecase.GetAccountByID(ctx, req.ToAccount)
	if err != nil {
		return nil, errs.ErrAccountNotFound
	}

	if req.FromAccount == req.ToAccount {
		return nil, errs.ErrTranSameAccount
	}

	if req.Amount <= 0 {
		return nil, errs.ErrAmountGreaterThanZero
	}

	// Check Account Balance
	if req.Amount > float64(fromAcc.Balance) {
		return nil, errs.ErrAmountGreaterAccBalance
	}

	// Tx Transaction
	err = uc.txManager.WithTx(ctx, func(tx *sql.Tx) error {
		// Update From Account
		err = uc.accUsecase.UpdateBalanceWithTx(ctx, tx, fromAcc.ID, -req.Amount)
		if err != nil {
			return fmt.Errorf("update from account failed: %w", err)
		}

		// Update To Account
		err = uc.accUsecase.UpdateBalanceWithTx(ctx, tx, toAcc.ID, req.Amount)
		if err != nil {
			return fmt.Errorf("update to account failed: %w", err)
		}

		// Insert Transaction
		result, err = uc.tranRepo.TransferWithTx(ctx, tx, &Transaction{
			FromAccount: &req.FromAccount,
			ToAccount:   &req.ToAccount,
			Amount:      req.Amount,
		})
		if err != nil {
			return fmt.Errorf("insert transaction failed: %w", err)
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (uc *transactionUsecase) Transactions(ctx context.Context, userID int64) ([]*Transaction, error) {
	ctx, cancel := context.WithTimeout(ctx, queryTimeout)
	defer cancel()

	return uc.tranRepo.Transactions(ctx, userID)
}
