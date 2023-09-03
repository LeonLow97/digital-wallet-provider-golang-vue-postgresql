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


INSERT INTO users (first_name, last_name, username, password, active, admin, mobile_number, creation_date) 
VALUES
('Alice', 'Tan', 'alice', '$2a$10$CerQd299qowq2ck8k/EqQeB7Jpjd/4Cut/Df.f8jnq9kYsuG0W7zG', 1, 1, '+65 90399012', NOW()),
('Bob', 'Smith', 'bob', '$2a$10$MVLL5BT/nIQKk6OYbgzK7.fbT0XKMBtNdeoy64ihYUUhr8Ag6358u', 1, 1, '+65 89230122', NOW()),
('Charlie', 'Brown', 'charlie', '$2a$10$yKz0rguTzykTec4Bgke7LempFl/GQVTw9w9qEXfGUpI/XGK97VHFq', 1, 1, '+1 5551234567', NOW()),
('David', 'Johnson', 'david', '$2a$10$p444biF49.py2HOTVe5TSuUNAhSKqelEtlbLtZXghUh3o21Et7DNO', 1, 1, '+49 1234567890', NOW());

INSERT INTO user_balance (user_id, balance, currency, country_iso_code, is_primary)
VALUES
(1, 70000.00, 'SGD', 'SG', 1),
(2, 15000.00, 'SGD', 'SG', 1),
(3, 78000.00, 'USD', 'US', 1),
(4, 4700.00, 'EUR', 'FR', 1);

INSERT INTO user_beneficiary (user_id, beneficiary_id)
VALUES
(1,2),(1,3),(1,4),(2,1),(2,3),(3,2),(3,4),(4,1);
