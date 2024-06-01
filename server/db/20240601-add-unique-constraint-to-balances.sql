-- add_unique_constraint_to_balances.sql

BEGIN;

-- Add unique constraint on (user_id, currency)
ALTER TABLE balances
ADD CONSTRAINT unique_user_currency UNIQUE (user_id, currency);

COMMIT;
