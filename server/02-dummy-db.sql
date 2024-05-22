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

INSERT INTO balances_history (amount, currency, type, user_id, balance_id, created_at)
VALUES
(4000, 'AUD', 'deposit', 1, 1, '2024-01-01 12:00:00'),
(16000, 'AUD', 'deposit', 1, 1, '2024-01-05 14:30:00'),
(11000, 'SGD', 'deposit', 1, 2, '2024-02-10 10:45:00'),
(4000, 'SGD', 'deposit', 1, 2, '2024-03-15 16:15:00'),
(4300, 'USD', 'deposit', 1, 3, '2024-04-20 12:30:00'),
(700, 'USD', 'deposit', 1, 3, '2024-05-25 14:00:00');

INSERT INTO wallet_types (type)
VALUES ('personal'), ('savings'), ('investment'), ('business');

INSERT INTO wallets (balance, currency, wallet_type_id, user_id)
VALUES
(8000, 'AUD', 4, 1);
