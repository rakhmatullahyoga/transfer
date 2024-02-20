package transfer

import (
	"context"
	"encoding/json"
	"errors"
	"time"
)

var (
	ErrAccountNotFound          = errors.New("bank account not found")
	ErrInvalidTransactionStatus = errors.New("invalid transaction status")
	ErrInvalidStatusTransition  = errors.New("invalid state transition")
	ErrParseInput               = errors.New("invalid input format")
)

type TransactionStatus uint

const (
	PendingState TransactionStatus = iota
	ProcessedState
	FailedState
	SuccessState
)

func (ts TransactionStatus) String() (status string) {
	switch ts {
	case PendingState:
		status = "pending"
	case ProcessedState:
		status = "processed"
	case FailedState:
		status = "failed"
	case SuccessState:
		status = "success"
	default:
		status = "unknown"
	}
	return
}

func (ts TransactionStatus) MarshalJSON() ([]byte, error) {
	return json.Marshal(ts.String())
}

var (
	transactionStatusMap = map[string]TransactionStatus{
		"pending":   PendingState,
		"processed": ProcessedState,
		"failed":    FailedState,
		"success":   SuccessState,
	}
)

func ParseTransactionStatus(s string) (status TransactionStatus, err error) {
	status, ok := transactionStatusMap[s]
	if !ok {
		err = ErrInvalidTransactionStatus
	}
	return
}

type AccountDetail struct {
	AccountNumber string `json:"accountNumber"`
	Owner         string `json:"owner"`
}

type Transaction struct {
	ID                uint64            `json:"id"`
	BankTransactionID string            `json:"bankTransactionId"`
	AccountNumber     string            `json:"accountNumber"`
	AccountName       string            `json:"accountName"`
	Amount            float64           `json:"amount"`
	Status            TransactionStatus `json:"transactionStatus"`
	CreatedAt         time.Time         `json:"createdAt"`
	UpdatedAt         *time.Time        `json:"updatedAt"`
}

type CallbackRequest struct {
	TransactionID string `json:"transactionId"`
	Status        string `json:"status"`
}

type TransferRequest struct {
	Amount        float64 `json:"amount"`
	AccountNumber string  `json:"accountNumber"`
}

type BankTransferResponse struct {
	ID        string    `json:"id"`
	Amount    string    `json:"amount"`
	Success   bool      `json:"success"`
	AccountID string    `json:"accountId"`
	CreatedAt time.Time `json:"createdAt"`
}

type BankRepository interface {
	ValidateAccount(context.Context, string) (AccountDetail, error)
	Transfer(context.Context, string, *Transaction) error
}

type DbRepository interface {
	CreateTransaction(context.Context, *Transaction) error
	GetTransaction(context.Context, string) (Transaction, error)
	ProcessTransaction(context.Context, *Transaction) error
	SetTransactionStatus(context.Context, uint64, TransactionStatus) error
}

type TransferUsecase interface {
	Inquiry(context.Context, string) (AccountDetail, error)
	Transfer(context.Context, TransferRequest) (Transaction, error)
	ProcessCallback(context.Context, CallbackRequest) error
}
