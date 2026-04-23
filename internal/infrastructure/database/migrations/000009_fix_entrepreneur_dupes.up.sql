-- Expand free-text fields that can exceed 255 chars when coming from birdarcha.
ALTER TABLE entrepreneurs
    ALTER COLUMN legal_name TYPE TEXT,
    ALTER COLUMN registration_authority TYPE TEXT,
    ALTER COLUMN registration_number TYPE TEXT,
    ALTER COLUMN legal_form TYPE TEXT,
    ALTER COLUMN founders TYPE TEXT,
    ALTER COLUMN address TYPE TEXT,
    ALTER COLUMN director_name TYPE TEXT,
    ALTER COLUMN mhobt_code TYPE TEXT;

-- De-duplicate existing rows before adding the unique constraint.
-- Keep the oldest row per (inn_id, registration_authority) pair.
DELETE FROM entrepreneurs e
USING entrepreneurs dup
WHERE e.inn_id = dup.inn_id
  AND e.registration_authority = dup.registration_authority
  AND e.created_at > dup.created_at;

-- Prevent re-inserting the same entity from the same source.
CREATE UNIQUE INDEX entrepreneurs_inn_source_uniq
    ON entrepreneurs (inn_id, registration_authority);
