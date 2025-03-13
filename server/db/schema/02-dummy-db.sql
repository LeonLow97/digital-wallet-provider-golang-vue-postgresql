INSERT INTO users (first_name, last_name, username, email, password, active, admin, mobile_country_code, mobile_number, created_at) 
VALUES
('Leon', 'Low', 'leonlow97', 'leonlow@email.com', '$2a$10$p444biF49.py2HOTVe5TSuUNAhSKqelEtlbLtZXghUh3o21Et7DNO', 1, 0, '+65', '87654321', NOW()),
('Admin', 'Admin', 'admin', 'adminuser@email.com', '$2a$10$GuPa.s8UqAVqDrC34uwRm.If7/vHAoAfGygZWPMk3UooywEsJivPu', 1, 1, '+65', '98765432', NOW());

-- INSERT INTO user_beneficiary (user_id, beneficiary_id)
-- VALUES
-- (1, 2);

INSERT INTO balances (balance, currency, user_id)
VALUES
(20000, 'AUD', 1),
(15000, 'SGD', 1),
(5000, 'USD', 1),
(30000, 'AUD', 2),
(4000, 'SGD', 2);

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

INSERT INTO wallets (wallet_type_id, user_id)
VALUES
(4, 1);

INSERT INTO wallet_balances (amount, currency, wallet_id, user_id)
VALUES
(10000, 'AUD', 1, 1);

-- Inserting dummy data for testing pagination
INSERT INTO transactions (user_id, sender_id, beneficiary_id, source_of_transfer, source_amount, source_currency, destination_amount, destination_currency, status, created_at)
VALUES
(1, 1, 2, 'Bank Transfer', 100.00, 'SGD', 70.00, 'USD', 'SUCCESS', '2022-01-01 10:00:00'),
(1, 1, 2, 'Credit Card', 500.00, 'SGD', 350.00, 'USD', 'PENDING', '2022-01-02 12:00:00'),
(1, 1, 2, 'PayPal', 200.00, 'SGD', 140.00, 'USD', 'FAILED', '2022-01-03 14:00:00'),
(1, 1, 2, 'Bank Transfer', 800.00, 'SGD', 560.00, 'USD', 'SUCCESS', '2022-01-04 16:00:00'),
(1, 1, 2, 'Credit Card', 300.00, 'SGD', 210.00, 'USD', 'PENDING', '2022-01-05 18:00:00'),
(1, 1, 2, 'Bank Transfer', 400.00, 'SGD', 280.00, 'USD', 'SUCCESS', '2022-01-06 20:00:00'),
(1, 1, 2, 'PayPal', 600.00, 'SGD', 420.00, 'USD', 'FAILED', '2022-01-07 22:00:00'),
(1, 1, 2, 'Credit Card', 900.00, 'SGD', 630.00, 'USD', 'PENDING', '2022-01-08 00:00:00'),
(1, 1, 2, 'Bank Transfer', 700.00, 'SGD', 490.00, 'USD', 'SUCCESS', '2022-01-09 02:00:00'),
(1, 1, 2, 'Bank Transfer', 950.00, 'SGD', 665.00, 'USD', 'SUCCESS', '2022-01-10 04:00:00'),
(1, 1, 2, 'PayPal', 250.00, 'SGD', 175.00, 'USD', 'FAILED', '2022-01-11 06:00:00'),
(1, 1, 2, 'Credit Card', 750.00, 'SGD', 525.00, 'USD', 'PENDING', '2022-01-12 08:00:00'),
(1, 1, 2, 'Bank Transfer', 650.00, 'SGD', 455.00, 'USD', 'SUCCESS', '2022-01-13 10:00:00'),
(1, 1, 2, 'Bank Transfer', 900.00, 'SGD', 630.00, 'USD', 'SUCCESS', '2022-01-14 12:00:00'),
(1, 1, 2, 'PayPal', 800.00, 'SGD', 560.00, 'USD', 'FAILED', '2022-01-15 14:00:00'),
(1, 1, 2, 'Credit Card', 950.00, 'SGD', 665.00, 'USD', 'PENDING', '2022-01-16 16:00:00'),
(1, 1, 2, 'Bank Transfer', 700.00, 'SGD', 490.00, 'USD', 'SUCCESS', '2022-01-17 18:00:00'),
(1, 1, 2, 'Bank Transfer', 500.00, 'SGD', 350.00, 'USD', 'SUCCESS', '2022-01-18 20:00:00'),
(1, 1, 2, 'PayPal', 650.00, 'SGD', 455.00, 'USD', 'FAILED', '2022-01-19 22:00:00'),
(1, 1, 2, 'Credit Card', 800.00, 'SGD', 560.00, 'USD', 'PENDING', '2022-01-20 00:00:00'),
(1, 1, 2, 'Bank Transfer', 950.00, 'SGD', 665.00, 'USD', 'SUCCESS', '2022-01-21 02:00:00'),
(1, 1, 2, 'PayPal', 750.00, 'SGD', 525.00, 'USD', 'FAILED', '2022-01-22 04:00:00'),
(1, 1, 2, 'Credit Card', 600.00, 'SGD', 420.00, 'USD', 'PENDING', '2022-01-23 06:00:00'),
(1, 1, 2, 'Bank Transfer', 400.00, 'SGD', 280.00, 'USD', 'SUCCESS', '2022-01-24 08:00:00'),
(1, 1, 2, 'Bank Transfer', 250.00, 'SGD', 175.00, 'USD', 'SUCCESS', '2022-01-25 10:00:00'),
(1, 1, 2, 'PayPal', 900.00, 'SGD', 630.00, 'USD', 'FAILED', '2022-01-26 12:00:00'),
(1, 1, 2, 'Credit Card', 700.00, 'SGD', 490.00, 'USD', 'PENDING', '2022-01-27 14:00:00'),
(1, 1, 2, 'Bank Transfer', 650.00, 'SGD', 455.00, 'USD', 'SUCCESS', '2022-01-28 16:00:00'),
(1, 1, 2, 'Bank Transfer', 500.00, 'SGD', 350.00, 'USD', 'SUCCESS', '2022-01-29 18:00:00'),
(1, 1, 2, 'PayPal', 800.00, 'SGD', 560.00, 'USD', 'FAILED', '2022-01-30 20:00:00');