ALTER TABLE users ALTER COLUMN id DROP DEFAULT;
ALTER TABLE users ALTER COLUMN id TYPE VARCHAR(36) USING id::VARCHAR(36);

UPDATE users SET id = gen_random_uuid();