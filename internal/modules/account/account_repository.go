package account

import (
	"context"
	"database/sql"
	"errors"
)

type AccountRepository interface {
	Create(ctx context.Context, acc *Account) (*Account, error)
	FindByID(ctx context.Context, id int64) (*Account, error)
	List(ctx context.Context, userID int64) ([]*Account, error)
	UpdateStatus(ctx context.Context, id int64, status string) error
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
	query := `UPDATE FROM accounts SET status = $1 WHERE id = $2`

	res, err := r.db.ExecContext(ctx, query, status, id)
	if err != nil {
		return err
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return errors.New("account not found")
	}

	return nil
}
