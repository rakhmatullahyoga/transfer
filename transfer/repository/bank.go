package repository

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"transfer/transfer"
)

type bankRepository struct {
	baseURL string
}

func NewBankRepository(baseUrl string) *bankRepository {
	return &bankRepository{
		baseURL: baseUrl,
	}
}

func (r *bankRepository) ValidateAccount(ctx context.Context, accountNumber string) (acc transfer.AccountDetail, err error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, r.baseURL+"/accounts", nil)
	if err != nil {
		return
	}

	q := req.URL.Query()
	q.Add("accountNumber", accountNumber)
	req.URL.RawQuery = q.Encode()
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return
	}

	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return
	}

	accounts := []transfer.AccountDetail{}
	err = json.Unmarshal(body, &accounts)
	if err != nil {
		if strings.Contains(string(body), "Not found") {
			err = transfer.ErrAccountNotFound
		}
		return
	}

	if len(accounts) > 0 {
		acc = accounts[0]
	} else {
		err = transfer.ErrAccountNotFound
	}
	return
}

func (r *bankRepository) Transfer(ctx context.Context, accountID string, trx *transfer.Transaction) (err error) {
	var trfRes transfer.BankTransferResponse
	payload := struct {
		Amount    string
		AccountID string
	}{
		Amount:    fmt.Sprintf("%.2f", trx.Amount),
		AccountID: accountID,
	}
	jsonPayload, _ := json.Marshal(payload)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, r.baseURL+"/transfers", bytes.NewReader(jsonPayload))
	if err != nil {
		return
	}

	req.Header.Set("Content-Type", "application/json")
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return
	}

	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return
	}

	err = json.Unmarshal(body, &trfRes)
	if err != nil {
		return
	}

	if trfRes.Success {
		trx.Status = transfer.ProcessedState
	} else {
		trx.Status = transfer.FailedState
	}
	trx.BankTransactionID = trfRes.ID
	return
}
