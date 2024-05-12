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
VALUES (500, 'SGD', 1, 1), (1500, 'SGD', 2, 1), (70000, 'CHN', 3, 1), (95600, 'IDR', 4, 1);