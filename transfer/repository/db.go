package repository

import (
	"context"

	"transfer/transfer"

	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	getTransactionQuery       = `SELECT id, status FROM transactions WHERE bank_transaction_id = $1`
	insertTransactionQuery    = `INSERT INTO transactions (account_number, account_name, amount, status) VALUES ( $1, $2, $3, $4 ) RETURNING id, created_at`
	processTransactionQuery   = `UPDATE transactions SET bank_transaction_id = $1, status = $2, updated_at = NOW() WHERE id = $3 RETURNING updated_at`
	setTransactionStatusQuery = `UPDATE transactions SET status = $1, updated_at = NOW() WHERE id = $2 RETURNING updated_at`
)

type dbRepository struct {
	pool *pgxpool.Pool
}

func NewDbRepository(pool *pgxpool.Pool) *dbRepository {
	return &dbRepository{
		pool: pool,
	}
}

func (r *dbRepository) CreateTransaction(ctx context.Context, trx *transfer.Transaction) (err error) {
	err = r.pool.QueryRow(ctx, insertTransactionQuery, trx.AccountNumber, trx.AccountName, trx.Amount, trx.Status).Scan(&trx.ID, &trx.CreatedAt)
	return
}

func (r *dbRepository) GetTransaction(ctx context.Context, bankTrxID string) (trx transfer.Transaction, err error) {
	err = r.pool.QueryRow(ctx, getTransactionQuery, bankTrxID).Scan(&trx.ID, &trx.Status)
	return
}

func (r *dbRepository) ProcessTransaction(ctx context.Context, trx *transfer.Transaction) (err error) {
	err = r.pool.QueryRow(ctx, processTransactionQuery, trx.BankTransactionID, trx.Status, trx.ID).Scan(&trx.UpdatedAt)
	return
}

func (r *dbRepository) SetTransactionStatus(ctx context.Context, ID uint64, status transfer.TransactionStatus) (err error) {
	_, err = r.pool.Exec(ctx, setTransactionStatusQuery, status, ID)
	return
}
