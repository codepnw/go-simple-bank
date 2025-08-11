package transaction

import "time"

type Transaction struct {
	ID          int64
	FromAccount int64
	ToAccount   int64
	Amount      float64
	Type        transactionType
	CreatedAt   time.Time
}
