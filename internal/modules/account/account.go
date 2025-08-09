package account

type Account struct {
	ID       int64         `json:"id"`
	UserID   int64         `json:"user_id"`
	Name     string        `json:"name"`
	Balance  int           `json:"balance"`
	Currency string        `json:"currency"`
	Status   accountStatus `json:"status"`
}
