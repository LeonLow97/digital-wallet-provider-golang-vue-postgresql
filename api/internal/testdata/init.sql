CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    first_name VARCHAR(50),
    last_name VARCHAR(50),
    username VARCHAR(20) NOT NULL UNIQUE,
    password VARCHAR(60) NOT NULL,
    active INT NOT NULL DEFAULT 1,
    admin INT NOT NULL DEFAULT 0,
    mobile_number VARCHAR(50) NOT NULL UNIQUE,
    creation_date TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

CREATE UNIQUE INDEX ON users(username);

CREATE TABLE IF NOT EXISTS user_balance (
    id SERIAL PRIMARY KEY,
    user_id INT NOT NULL,
    balance DECIMAL(20,2) NOT NULL,
    currency CHAR(3) NOT NULL,
    country_iso_code CHAR(3) NOT NULL,
    is_primary INT NOT NULL DEFAULT 0,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS user_beneficiary (
    user_id INT NOT NULL,
    beneficiary_id INT NOT NULL,
    PRIMARY KEY (user_id, beneficiary_id),
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (beneficiary_id) REFERENCES users(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS transactions (
    id SERIAL PRIMARY KEY,
    user_id SERIAL NOT NULL,
    sender_id SERIAL NOT NULL,
    beneficiary_id SERIAL NOT NULL,
    transferred_amount DECIMAL(20,2) NOT NULL,
    transferred_amount_currency CHAR(3) NOT NULL,
    received_amount DECIMAL(20,2) NOT NULL,
    received_amount_currency CHAR(3) NOT NULL,
    status VARCHAR(50) NOT NULL,
    transferred_date TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    received_date TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (sender_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (beneficiary_id) REFERENCES users(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS transaction_fee (
    id SERIAL PRIMARY KEY,
    transaction_id SERIAL NOT NULL,
    fee DECIMAL(10,2) NOT NULL,
    FOREIGN KEY (transaction_id) REFERENCES transactions(id) ON DELETE CASCADE
);
