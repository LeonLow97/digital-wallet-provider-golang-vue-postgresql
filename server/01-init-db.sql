CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    first_name VARCHAR(50),
    last_name VARCHAR(50),
    email VARCHAR(255) NOT NULL UNIQUE,
    username VARCHAR(20) NOT NULL,
    password VARCHAR(60) NOT NULL,
    active INT NOT NULL DEFAULT 1,
    admin INT NOT NULL DEFAULT 0,
    mobile_number VARCHAR(255) NOT NULL UNIQUE,
    is_mfa_configured BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE
);

CREATE UNIQUE INDEX ON users(email);

CREATE TABLE IF NOT EXISTS user_totp_secrets (
    id SERIAL PRIMARY KEY,
    user_id INT NOT NULL,
    totp_encrypted_secret TEXT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE
);

CREATE TABLE IF NOT EXISTS user_beneficiary (
    user_id INT NOT NULL,
    beneficiary_id INT NOT NULL,
    is_deleted SMALLINT NOT NULL DEFAULT 0,
    PRIMARY KEY (user_id, beneficiary_id),
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (beneficiary_id) REFERENCES users(id) ON DELETE CASCADE
);

CREATE UNIQUE INDEX ON user_beneficiary(user_id);
CREATE UNIQUE INDEX ON user_beneficiary(beneficiary_id);

CREATE TABLE IF NOT EXISTS wallet_types (
    id SERIAL PRIMARY KEY,
    type VARCHAR(255) NOT NULL UNIQUE,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS wallets (
    id SERIAL PRIMARY KEY,
    balance DECIMAL(20,2) NOT NULL,
    currency CHAR(3) NOT NULL,
    wallet_type_id INT REFERENCES wallet_types(id) ON DELETE CASCADE,
    user_id INT REFERENCES users(id) ON DELETE CASCADE,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE
);

CREATE TABLE IF NOT EXISTS balances (
    id SERIAL PRIMARY KEY,
    balance DECIMAL(20,2) NOT NULL,
    currency CHAR(3) NOT NULL,
    user_id INT REFERENCES users(id) ON DELETE CASCADE,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE
);

CREATE TABLE IF NOT EXISTS transactions (
    id SERIAL PRIMARY KEY,
    user_id INT REFERENCES users(id) NOT NULL,
    sender_id INT REFERENCES users(id) NOT NULL,
    beneficiary_id INT REFERENCES users(id) NOT NULL,
    source_of_transfer VARCHAR(255) NOT NULL,
    sent_amount DECIMAL(20,2) NOT NULL,
    source_currency CHAR(3) NOT NULL,
    received_amount DECIMAL(20,2) NOT NULL,
    received_currency CHAR(3) NOT NULL,
    status VARCHAR(50) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);