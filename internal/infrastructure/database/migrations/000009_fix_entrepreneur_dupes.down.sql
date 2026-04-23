DROP INDEX IF EXISTS entrepreneurs_inn_source_uniq;

ALTER TABLE entrepreneurs
    ALTER COLUMN legal_name TYPE VARCHAR(255),
    ALTER COLUMN registration_authority TYPE VARCHAR(255),
    ALTER COLUMN registration_number TYPE VARCHAR(255),
    ALTER COLUMN legal_form TYPE VARCHAR(255),
    ALTER COLUMN founders TYPE VARCHAR(255),
    ALTER COLUMN address TYPE VARCHAR(255),
    ALTER COLUMN director_name TYPE VARCHAR(255),
    ALTER COLUMN mhobt_code TYPE VARCHAR(255);
