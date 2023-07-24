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