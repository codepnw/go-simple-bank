package transaction

type transactionType string

const (
	TypeDeposit  transactionType = "DEPOSIT"
	TypeTransfer transactionType = "TRANSFER"
	TypeWithdraw transactionType = "WITHDRAW"
)

type DepositReq struct {
	ToAccount int64   `json:"to_account" validate:"required"`
	Amount    float64 `json:"amount" validate:"required,gt=0"`
}

type WithdrawReq struct {
	FromAccount int64   `json:"from_account" validate:"required"`
	Amount      float64 `json:"amount" validate:"required,gt=0"`
}

type TransferReq struct {
	FromAccount int64   `json:"from_account" validate:"required"`
	ToAccount   int64   `json:"to_account" validate:"required"`
	Amount      float64 `json:"amount" validate:"required,gt=0"`
}
