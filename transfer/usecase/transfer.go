package usecase

import (
	"context"
	"transfer/transfer"
)

type transferUsecase struct {
	bankRepo transfer.BankRepository
	dbRepo   transfer.DbRepository
}

func NewUsecase(bankRepo transfer.BankRepository, dbRepo transfer.DbRepository) *transferUsecase {
	return &transferUsecase{
		bankRepo: bankRepo,
		dbRepo:   dbRepo,
	}
}

func (u *transferUsecase) Inquiry(ctx context.Context, accountNumber string) (acc transfer.AccountDetail, err error) {
	acc, err = u.bankRepo.ValidateAccount(ctx, accountNumber)
	return
}

func (u *transferUsecase) Transfer(ctx context.Context, params transfer.TransferRequest) (trx transfer.Transaction, err error) {
	acc, err := u.bankRepo.ValidateAccount(ctx, params.AccountNumber)
	if err != nil {
		return
	}

	trx = transfer.Transaction{
		AccountNumber: acc.AccountNumber,
		AccountName:   acc.Owner,
		Amount:        params.Amount,
	}
	err = u.dbRepo.CreateTransaction(ctx, &trx)
	if err != nil {
		return
	}

	err = u.bankRepo.Transfer(ctx, acc.AccountNumber, &trx)
	if err != nil {
		return
	}

	err = u.dbRepo.ProcessTransaction(ctx, &trx)
	return
}

func (u *transferUsecase) ProcessCallback(ctx context.Context, params transfer.CallbackRequest) (err error) {
	trx, err := u.dbRepo.GetTransaction(ctx, params.TransactionID)
	if err != nil {
		return
	}

	status, err := validateStateTransition(trx, params.Status)
	if err != nil {
		return
	}

	err = u.dbRepo.SetTransactionStatus(ctx, trx.ID, status)
	return
}

func validateStateTransition(trx transfer.Transaction, newState string) (status transfer.TransactionStatus, err error) {
	status, err = transfer.ParseTransactionStatus(newState)
	if err != nil {
		return
	}

	var valid bool
	switch trx.Status {
	case transfer.PendingState:
		valid = status == transfer.ProcessedState || status == transfer.FailedState
	case transfer.ProcessedState:
		valid = status == transfer.FailedState || status == transfer.SuccessState
	default:
		valid = false
	}
	if !valid {
		err = transfer.ErrInvalidStatusTransition
	}
	return
}
