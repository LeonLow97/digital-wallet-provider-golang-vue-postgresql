-- INSERT INTO `users` table
INSERT INTO users (first_name, last_name, username, password, active, admin, mobile_number)
VALUES
    ('Leon', 'Low', 'leonlow', '$2a$10$kCM6Ou5nwmzbiyrKIqVpCekuDBvuI0Swv1e8v51klmDPqMGcugO2O', 1, 1, '+65 98765432'),
    ('Alice', 'Johnson', 'alicejohnson', '$2a$10$6pZitxHzYR21Q08bwDo0zOt.U3C54p8nBDKRWkoC.waIAFmyYFWYG', 1, 0, '+1 1234567890'),
    ('Bob', 'Smith', 'bobsmith', '$2a$10$Pha2.FaRlHItL.NRzzm7Oe6kftIawkm4joGbrhU9TSWMUNeR8o3bC', 1, 0, '+1 9876543210'),
    ('Charlie', 'Brown', 'charliebrown', '$2a$10$y8bjN8Xp7U/2RXrXEmzoWeI/jj2n0bz7RFADFTXZUpG8NzgPbb0Q6', 0, 0, '+44 7712345678')
;
