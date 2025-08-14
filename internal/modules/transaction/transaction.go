package transaction

import "time"

type Transaction struct {
	ID          int64           `json:"id"`
	FromAccount *int64          `json:"from_account"`
	ToAccount   *int64          `json:"to_account"`
	Amount      float64         `json:"amount"`
	Type        transactionType `json:"type"`
	Role        *string         `json:"role"`
	CreatedAt   time.Time       `json:"created_at"`
}
