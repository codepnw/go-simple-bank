package account

import (
	"context"
	"database/sql"

	"github.com/codepnw/simple-bank/internal/utils/errs"
)

type AccountRepository interface {
	Create(ctx context.Context, acc *Account) (*Account, error)
	FindByID(ctx context.Context, id int64) (*Account, error)
	List(ctx context.Context, userID int64) ([]*Account, error)
	UpdateStatus(ctx context.Context, id int64, status string) error
	UpdateBalanceWithTx(ctx context.Context, tx *sql.Tx, id int64, balance float64) error
	GetAccountBalance(ctx context.Context, accountID int64) (float64, error)
	GetAccountBalanceByUserID(ctx context.Context, accountID, userID int64) (float64, error)
}

type accountRepository struct {
	db *sql.DB
}

func NewAccountRepository(db *sql.DB) AccountRepository {
	return &accountRepository{db: db}
}

func (r *accountRepository) Create(ctx context.Context, acc *Account) (*Account, error) {
	query := `
		INSERT INTO accounts (user_id, name, balance)
		VALUES ($1, $2, $3) RETURNING id, currency, status;
	`
	err := r.db.QueryRowContext(
		ctx,
		query,
		acc.UserID,
		acc.Name,
		acc.Balance,
	).Scan(
		&acc.ID,
		&acc.Currency,
		&acc.Status,
	)
	if err != nil {
		return nil, err
	}

	return acc, nil
}

func (r *accountRepository) FindByID(ctx context.Context, id int64) (*Account, error) {
	query := `
		SELECT id, user_id, name, balance, currency, status
		FROM accounts WHERE id = $1 LIMIT 1;
	`
	acc := new(Account)

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&acc.ID,
		&acc.UserID,
		&acc.Name,
		&acc.Balance,
		&acc.Currency,
		&acc.Status,
	)
	if err != nil {
		return nil, err
	}

	return acc, nil
}

func (r *accountRepository) List(ctx context.Context, userID int64) ([]*Account, error) {
	query := `
		SELECT id, user_id, name, balance, currency, status
		FROM accounts WHERE user_id = $1
	`
	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var accs []*Account

	for rows.Next() {
		acc := new(Account)
		err = rows.Scan(
			&acc.ID,
			&acc.UserID,
			&acc.Name,
			&acc.Balance,
			&acc.Currency,
			&acc.Status,
		)
		if err != nil {
			return nil, err
		}
		accs = append(accs, acc)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return accs, nil
}

func (r *accountRepository) UpdateStatus(ctx context.Context, id int64, status string) error {
	query := `UPDATE accounts SET status = $1 WHERE id = $2`

	res, err := r.db.ExecContext(ctx, query, status, id)
	if err != nil {
		return err
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return errs.ErrAccountNotFound
	}

	return nil
}

func (r *accountRepository) UpdateBalanceWithTx(ctx context.Context, tx *sql.Tx, id int64, balance float64) error {
	query := `
		UPDATE accounts SET balance = balance + $1 
		WHERE id = $2 AND balance + $1 >= 0
	`
	res, err := tx.ExecContext(ctx, query, balance, id)
	if err != nil {
		return err
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return errs.ErrAccountNotFound
	}

	return nil
}

func (r *accountRepository) GetAccountBalance(ctx context.Context, accountID int64) (float64, error) {
	var balance float64
	query := `SELECT balance FROM accounts WHERE id = $1`

	err := r.db.QueryRowContext(ctx, query, accountID).Scan(&balance)
	if err != nil {
		return 0, err
	}

	return balance, nil
}

func (r *accountRepository) GetAccountBalanceByUserID(ctx context.Context, accountID, userID int64) (float64, error) {
	var balance float64
	query := `SELECT balance FROM accounts WHERE id = $1 AND user_id = $2`

	err := r.db.QueryRowContext(ctx, query, accountID, userID).Scan(&balance)
	if err != nil {
		return 0, err
	}

	return balance, nil
}
