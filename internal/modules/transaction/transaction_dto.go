package transaction

type transactionType string

const (
	TypeDeposit  transactionType = "DEPOSIT"
	TypeTransfer transactionType = "TRANSFER"
	TypeWithdraw transactionType = "WITHDRAW"
)

type DepositReq struct {
	FromAccount int64   `json:"from_account"`
	Amount      float64 `json:"amount"`
}

type WithdrawReq struct {
	ToAccount int64   `json:"to_account"`
	Amount    float64 `json:"amount"`
}

type TransferReq struct {
	FromAccount int64   `json:"from_account"`
	ToAccount   int64   `json:"to_account"`
	Amount      float64 `json:"amount"`
}
