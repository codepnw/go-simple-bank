package transaction

import (
	"context"
	"database/sql"
)

type TransasctionRepository interface {
	DepositWithTx(ctx context.Context, tx *sql.Tx, input *Transaction) (*Transaction, error)
	WithdrawWithTx(ctx context.Context, tx *sql.Tx, input *Transaction) (*Transaction, error)
	TransferWithTx(ctx context.Context, tx *sql.Tx, input *Transaction) (*Transaction, error)
	Transactions(ctx context.Context, userID int64) ([]*Transaction, error)
}

type transactionRepository struct {
	db *sql.DB
}

func NewTransactionRepository(db *sql.DB) TransasctionRepository {
	return &transactionRepository{db: db}
}

func (r *transactionRepository) DepositWithTx(ctx context.Context, tx *sql.Tx, input *Transaction) (*Transaction, error) {
	query := `
		INSERT INTO transactions (to_account, amount, type)
		VALUES ($1, $2, $3)
		RETURNING id, type, created_at
	`
	err := tx.QueryRowContext(
		ctx,
		query,
		input.ToAccount,
		input.Amount,
		TypeDeposit,
	).Scan(
		&input.ID,
		&input.Type,
		&input.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	return input, nil
}

func (r *transactionRepository) WithdrawWithTx(ctx context.Context, tx *sql.Tx, input *Transaction) (*Transaction, error) {
	query := `
		INSERT INTO transactions (from_account, amount, type)
		VALUES ($1, $2, $3)
		RETURNING id, type, created_at
	`
	err := tx.QueryRowContext(
		ctx,
		query,
		input.FromAccount,
		input.Amount,
		TypeWithdraw,
	).Scan(
		&input.ID,
		&input.Type,
		&input.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	return input, nil
}

func (r *transactionRepository) TransferWithTx(ctx context.Context, tx *sql.Tx, input *Transaction) (*Transaction, error) {
	query := `
		INSERT INTO transactions (from_account, to_account, amount, type)
		VALUES ($1, $2, $3, $4)
		RETURNING id, type, created_at
	`
	err := tx.QueryRowContext(
		ctx,
		query,
		input.FromAccount,
		input.ToAccount,
		input.Amount,
		TypeTransfer,
	).Scan(
		&input.ID,
		&input.Type,
		&input.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	return input, nil
}

func (r *transactionRepository) Transactions(ctx context.Context, userID int64) ([]*Transaction, error) {
	query := `
		SELECT t.id, t.from_account, t.to_account, t.amount, t.type, t.created_at,
			CASE 
				WHEN t.from_account = a.id THEN 'SENDER'
				WHEN t.to_account = a.id THEN 'RECEIVER'
			END AS role
		FROM transactions t 
		JOIN accounts a 
			ON (a.id = t.from_account OR a.id = t.to_account)
		WHERE a.user_id = $1
	`
	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var transactions []*Transaction

	for rows.Next() {
		t := new(Transaction)

		err = rows.Scan(
			&t.ID,
			&t.FromAccount,
			&t.ToAccount,
			&t.Amount,
			&t.Type,
			&t.CreatedAt,
			&t.Role,
		)
		if err != nil {
			return nil, err
		}

		transactions = append(transactions, t)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return transactions, nil
}
