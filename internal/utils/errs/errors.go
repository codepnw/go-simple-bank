package errs

import "errors"

var (
	// Error Account
	ErrAccountNotFound         = errors.New("account not found")
	ErrAccountAmountNotZero    = errors.New("amount mest not zero")
	ErrAmountGreaterThanZero   = errors.New("amount must be greater than zero")
	ErrAmountGreaterAccBalance = errors.New("amount must be greater than account balance")

	// Error Transaction
	ErrTranSameAccount = errors.New("cant transfer to the same account")

	// Error Users
	ErrUserNotFound = errors.New("user not found")
)
