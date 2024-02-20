CREATE TABLE IF NOT EXISTS transactions (
    id BIGSERIAL PRIMARY KEY,
    bank_transaction_id VARCHAR(20),
    account_number VARCHAR(20) NOT NULL,
    account_name VARCHAR(50) NOT NULL,
    amount DECIMAL(10,2) NOT NULL,
    status SMALLINT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP
);
CREATE INDEX transactions_on_bank_transaction_id ON transactions USING btree (bank_transaction_id);
CREATE INDEX transactions_on_bank_account_number ON transactions USING btree (account_number);
