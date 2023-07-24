CREATE TABLE IF NOT EXISTS users (
    id BIGSERIAL PRIMARY KEY,
    username TEXT NOT NULL UNIQUE,
    password TEXT NOT NULL,
    mobile_number VARCHAR(255) NOT NULL UNIQUE,
    creation_date TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

CREATE UNIQUE INDEX ON users(username);

CREATE TABLE IF NOT EXISTS user_balance (
    user_id BIGSERIAL NOT NULL,
    balance DECIMAL(20,2) NOT NULL,
    currency CHAR(3) NOT NULL,
    country_iso_code CHAR(3) NOT NULL,
    PRIMARY KEY (user_id),
    FOREIGN KEY (user_id) REFERENCES users(id)
);

CREATE TABLE IF NOT EXISTS beneficiaries (
    beneficiary_id BIGSERIAL PRIMARY KEY,
    beneficiary_name VARCHAR(255) NOT NULL,
    mobile_number VARCHAR(255) NOT NULL UNIQUE,
    currency CHAR(3) NOT NULL,
    is_internal INT NOT NULL
);

CREATE TABLE IF NOT EXISTS user_beneficiary (
    user_id BIGSERIAL NOT NULL,
    beneficiary_id BIGSERIAL NOT NULL,
    PRIMARY KEY (user_id, beneficiary_id),
    FOREIGN KEY (user_id) REFERENCES users(id),
    FOREIGN KEY (beneficiary_id) REFERENCES beneficiaries(beneficiary_id)
);

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