package transaction

type transactionHandler struct {
	uc TransactionUsecase
}

func NewTransactionHandler(uc TransactionUsecase) *transactionHandler {
	return &transactionHandler{uc: uc}
}
