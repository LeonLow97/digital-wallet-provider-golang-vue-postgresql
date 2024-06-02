/*
failed to create beneficiary with error: ERROR: duplicate key value violates unique constraint "user_beneficiary_user_id_idx" (SQLSTATE 23505)

CREATE TABLE IF NOT EXISTS user_beneficiary (
    user_id INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    beneficiary_id INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    is_deleted SMALLINT NOT NULL DEFAULT 0,
    PRIMARY KEY (user_id, beneficiary_id)
);

CREATE UNIQUE INDEX ON user_beneficiary(user_id);
CREATE UNIQUE INDEX ON user_beneficiary(beneficiary_id);

When creating 2 beneficiaries with the same user, we get the above error because
`CREATE UNIQUE INDEX ON user_beneficiary(user_id);` enforces that each user_id must be unique across the table.
However, we allow a user to have multiple beneficiaries, so user_id will appear in multiple rows but with a different beneficiary_id
*/

-- 1. Find Index Names
SELECT indexname
FROM pg_indexes
WHERE tablename = 'user_beneficiary';

-- 2. Drop the Indexes
BEGIN;

DROP INDEX IF EXISTS user_beneficiary_user_id_idx;
DROP INDEX IF EXISTS user_beneficiary_beneficiary_id_idx;

COMMIT;
