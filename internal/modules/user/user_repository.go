package user

import (
	"context"
	"database/sql"

	"github.com/codepnw/simple-bank/internal/utils/errs"
)

type UserRepository interface {
	Create(ctx context.Context, u *User) (*User, error)
	FindByID(ctx context.Context, id int64) (*User, error)
	FindByEmail(ctx context.Context, email string) (*User, error)
	List(ctx context.Context) ([]*User, error)
	Update(ctx context.Context, u *User) error
	Delete(ctx context.Context, id int64) error
}

type userRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) Create(ctx context.Context, u *User) (*User, error) {
	query := `
		INSERT INTO users (email, password, first_name, last_name, phone)
		VALUES ($1, $2, $3, $4, $5) RETURNING id, created_at;
	`
	err := r.db.QueryRowContext(
		ctx,
		query,
		u.Email,
		u.Password,
		u.FirstName,
		u.LastName,
		u.Phone,
	).Scan(&u.ID, &u.CreatedAt)

	if err != nil {
		return nil, err
	}

	return u, nil
}

func (r *userRepository) FindByID(ctx context.Context, id int64) (*User, error) {
	query := `
		SELECT id, email, password, first_name, last_name, phone, created_at, updated_at
		FROM users WHERE id = $1 LIMIT 1;
	`
	var u User

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&u.ID,
		&u.Email,
		&u.Password,
		&u.FirstName,
		&u.LastName,
		&u.Phone,
		&u.CreatedAt,
		&u.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &u, nil
}

func (r *userRepository) FindByEmail(ctx context.Context, email string) (*User, error) {
	query := `
		SELECT id, email, password, first_name, last_name, phone, created_at, updated_at
		FROM users WHERE email = $1 LIMIT 1;
	`
	var u User

	err := r.db.QueryRowContext(ctx, query, email).Scan(
		&u.ID,
		&u.Email,
		&u.Password,
		&u.FirstName,
		&u.LastName,
		&u.Phone,
		&u.CreatedAt,
		&u.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &u, nil
}

func (r *userRepository) List(ctx context.Context) ([]*User, error) {
	query := `
		SELECT id, email, password, first_name, last_name, phone, created_at, updated_at
		FROM users;
	`
	var users []*User

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var u User
		err = rows.Scan(
			&u.ID,
			&u.Email,
			&u.Password,
			&u.FirstName,
			&u.LastName,
			&u.Phone,
			&u.CreatedAt,
			&u.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		users = append(users, &u)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

func (r *userRepository) Update(ctx context.Context, u *User) error {
	query := `
		UPDATE users SET first_name = $1, last_name = $2, phone = $3
		WHERE id = $4
	`
	res, err := r.db.ExecContext(ctx, query, u.FirstName, u.LastName, u.Phone, u.ID)
	if err != nil {
		return err
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return errs.ErrUserNotFound
	}

	return nil
}

func (r *userRepository) Delete(ctx context.Context, id int64) error {
	res, err := r.db.ExecContext(ctx, "DELETE FROM users WHERE id = $1", id)
	if err != nil {
		return err
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return errs.ErrUserNotFound
	}

	return nil
}
