ALTER TABLE users
  ADD COLUMN first_name VARCHAR(255) NOT NULL DEFAULT '',
  ADD COLUMN last_name VARCHAR(255) NOT NULL DEFAULT '';

DROP INDEX users_username_idx;

ALTER TABLE users
  DROP COLUMN username,
  DROP COLUMN last_login_at;

ALTER TABLE users
  ALTER COLUMN password TYPE VARCHAR(60);
