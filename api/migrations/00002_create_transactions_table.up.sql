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
    FOREIGN KEY (user_id) REFERENCES users(id),
    FOREIGN KEY (sender_id) REFERENCES users(id),
    FOREIGN KEY (beneficiary_id) REFERENCES users(id)
);

CREATE TABLE IF NOT EXISTS transaction_fee (
    id SERIAL PRIMARY KEY,
    transaction_id SERIAL NOT NULL,
    fee DECIMAL(10,2) NOT NULL,
    FOREIGN KEY (transaction_id) REFERENCES transactions(id)
);

INSERT INTO transactions (user_id, sender_id, beneficiary_id, transferred_amount, transferred_amount_currency, received_amount, received_amount_currency, status)
SELECT
    1, 1, 2,
    gs.num,
    'SGD',
    6000,
    'SGD',
    'COMPLETED'
FROM
    generate_series(1, 1000) AS gs(num);

