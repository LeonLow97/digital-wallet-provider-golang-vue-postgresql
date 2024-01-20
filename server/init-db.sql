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
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

CREATE UNIQUE INDEX ON users(email);

CREATE TABLE IF NOT EXISTS user_beneficiary (
    user_id INT NOT NULL,
    beneficiary_id INT NOT NULL,
    PRIMARY KEY (user_id, beneficiary_id),
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (beneficiary_id) REFERENCES users(id) ON DELETE CASCADE
);

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
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS balances (
    id SERIAL PRIMARY KEY,
    balance DECIMAL(20,2) NOT NULL,
    currency CHAR(3) NOT NULL,
    user_id INT REFERENCES users(id) ON DELETE CASCADE,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS transactions (
    id SERIAL PRIMARY KEY,
    sender_wallet_id INT REFERENCES wallets(id) ON DELETE CASCADE,
    beneficiary_wallet_id INT REFERENCES wallets(id) ON DELETE CASCADE,
    type VARCHAR(50) NOT NULL,
    status VARCHAR(50) NOT NULL,
    amount DECIMAL(20,2) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

-- temporary table (remove when connected to real Exchange Rate API)
CREATE TABLE IF NOT EXISTS exchange_rates (
    currency_from VARCHAR(5) NOT NULL,
    currency_to VARCHAR(5) NOT NULL,
    rate DECIMAL(18, 6) NOT NULL,
    PRIMARY KEY (currency_from, currency_to)
);

INSERT INTO users (first_name, last_name, username, email, password, active, admin, mobile_number, created_at) 
VALUES
('Alice', 'Tan', 'alice', 'alicetan@email.com', '$2a$10$CerQd299qowq2ck8k/EqQeB7Jpjd/4Cut/Df.f8jnq9kYsuG0W7zG', 1, 1, '+65 90399012', NOW()),
('Bob', 'Smith', 'bob', 'bobsmith@email.com', '$2a$10$MVLL5BT/nIQKk6OYbgzK7.fbT0XKMBtNdeoy64ihYUUhr8Ag6358u', 1, 1, '+65 89230122', NOW()),
('Charlie', 'Brown', 'charlie', 'charliebrown@email.com', '$2a$10$yKz0rguTzykTec4Bgke7LempFl/GQVTw9w9qEXfGUpI/XGK97VHFq', 1, 1, '+1 5551234567', NOW()),
('David', 'Johnson', 'david', 'davidjohnson@email.com', '$2a$10$p444biF49.py2HOTVe5TSuUNAhSKqelEtlbLtZXghUh3o21Et7DNO', 1, 1, '+49 1234567890', NOW()),
('Leon', 'Low', 'leon', 'leonlow@email.com', '$2a$10$p444biF49.py2HOTVe5TSuUNAhSKqelEtlbLtZXghUh3o21Et7DNO', 1, 1, '+65 1234567890', NOW());

INSERT INTO user_beneficiary (user_id, beneficiary_id)
VALUES
(1,2),(1,3),(1,4),(2,1),(2,3),(3,2),(3,4),(4,1);

INSERT INTO wallet_types (type)
VALUES ('personal'), ('savings'), ('investment'), ('business');
