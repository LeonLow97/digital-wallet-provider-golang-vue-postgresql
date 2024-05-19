INSERT INTO users (first_name, last_name, username, email, password, active, admin, mobile_number, created_at) 
VALUES
('Leon', 'Low', 'leonlow97', 'leonlow@email.com', '$2a$10$p444biF49.py2HOTVe5TSuUNAhSKqelEtlbLtZXghUh3o21Et7DNO', 1, 1, '+65 1234567890', NOW());

-- INSERT INTO user_beneficiary (user_id, beneficiary_id)
-- VALUES
-- (1, 2);

INSERT INTO balances (balance, currency, user_id)
VALUES
(20000, 'AUD', 1),
(15000, 'SGD', 1),
(5000, 'USD', 1);

INSERT INTO balances_history (amount, currency, type, user_id, balance_id)
VALUES
(4000, 'AUD', 'deposit', 1, 1),
(16000, 'AUD', 'deposit', 1, 1),
(1000, 'SGD', 'deposit', 1, 2),
(4000, 'SGD', 'deposit', 1, 2),
(4300, 'USD', 'deposit', 1, 3),
(700, 'USD', 'deposit', 1, 3);

INSERT INTO wallet_types (type)
VALUES ('personal'), ('savings'), ('investment'), ('business');

INSERT INTO wallets (balance, currency, wallet_type_id, user_id)
VALUES
(8000, 'AUD', 4, 1);
