ALTER TABLE users
  ADD COLUMN username VARCHAR(255) NOT NULL UNIQUE default '',
  ADD COLUMN last_login_at TIMESTAMP;

CREATE INDEX users_username_idx ON users(username);

ALTER TABLE users
  DROP COLUMN first_name,
  DROP COLUMN last_name;

ALTER TABLE users
  ALTER COLUMN password TYPE VARCHAR(255) 