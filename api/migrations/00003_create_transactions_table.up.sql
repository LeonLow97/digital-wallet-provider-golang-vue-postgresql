CREATE TABLE IF NOT EXISTS transactions (
    transaction_id BIGSERIAL PRIMARY KEY,
    user_id BIGSERIAL NOT NULL,
    sender_id BIGSERIAL NOT NULL,
    beneficiary_id BIGSERIAL NOT NULL,
    amount_transferred DECIMAL(20,2) NOT NULL,
    amount_transferred_currency CHAR(3) NOT NULL,
    amount_received DECIMAL(20,2) NOT NULL,
    amount_received_currency CHAR(3) NOT NULL,
    status VARCHAR(50) NOT NULL,
    date_transferred TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    date_received TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);