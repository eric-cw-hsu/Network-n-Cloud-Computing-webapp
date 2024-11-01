CREATE TABLE
  user_pic (
    -- meta data
    filename VARCHAR(255) NOT NULL,
    uploaded_at TIMESTAMP NOT NULL,
    url VARCHAR(255) NOT NULL,
    s3_key VARCHAR(255) NOT NULL,
    etag VARCHAR(255) NOT NULL,
    encryption VARCHAR(255) NOT NULL,
    encryption_key VARCHAR(255) NOT NULL,
    user_id VARCHAR(36) NOT NULL,
    foreign key (user_id) references users (id) on delete cascade,
    primary key (user_id)
  );