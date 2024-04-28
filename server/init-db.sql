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
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
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

INSERT INTO users (first_name, last_name, username, email, password, active, admin, mobile_number, created_at) 
VALUES
('Leon', 'Low', 'leonlow97', 'leonlow@email.com', '$2a$10$p444biF49.py2HOTVe5TSuUNAhSKqelEtlbLtZXghUh3o21Et7DNO', 1, 1, '+65 1234567890', NOW()),
('Bob', 'Smith', 'bobsmith97', 'bobsmith@email.com', '$2a$10$MVLL5BT/nIQKk6OYbgzK7.fbT0XKMBtNdeoy64ihYUUhr8Ag6358u', 1, 1, '+65 89230122', NOW()),
('Charlie', 'Brown', 'charliebrown97', 'charliebrown@email.com', '$2a$10$yKz0rguTzykTec4Bgke7LempFl/GQVTw9w9qEXfGUpI/XGK97VHFq', 1, 1, '+1 5551234567', NOW()),
('David', 'Johnson', 'davidjohnson97', 'davidjohnson@email.com', '$2a$10$p444biF49.py2HOTVe5TSuUNAhSKqelEtlbLtZXghUh3o21Et7DNO', 1, 1, '+49 1234567890', NOW()),
('Alice', 'Tan', 'alicetan97', 'alicetan@email.com', '$2a$10$CerQd299qowq2ck8k/EqQeB7Jpjd/4Cut/Df.f8jnq9kYsuG0W7zG', 1, 1, '+65 90399012', NOW());

INSERT INTO user_beneficiary (user_id, beneficiary_id)
VALUES
(1, 5);

INSERT INTO balances (balance, currency, user_id)
VALUES (70000, 'SGD', 1), (20000, 'SGD', 5);

INSERT INTO wallet_types (type)
VALUES ('personal'), ('savings'), ('investment'), ('business');

INSERT INTO wallets (balance, currency, wallet_type_id, user_id)
VALUES (500, 'SGD', 1, 1), (1500, 'SGD', 2, 1);