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
    user_id SERIAL NOT NULL UNIQUE,
    balance DECIMAL(20,2) NOT NULL,
    currency CHAR(3) NOT NULL,
    country_iso_code CHAR(3) NOT NULL,
    PRIMARY KEY (user_id),
    FOREIGN KEY (user_id) REFERENCES users(id)
);

CREATE TABLE IF NOT EXISTS user_beneficiary (
    user_id SERIAL NOT NULL,
    beneficiary_id BIGSERIAL NOT NULL,
    PRIMARY KEY (user_id, beneficiary_id),
    FOREIGN KEY (user_id) REFERENCES users(id),
    FOREIGN KEY (beneficiary_id) REFERENCES users(id)
);