package transaction

type TransactionUsecase interface{}

type transactionUsecase struct {
	repo TransasctionRepository
}

func NewTransactionUsecse(repo TransasctionRepository) TransactionUsecase {
	return &transactionUsecase{repo: repo}
}
