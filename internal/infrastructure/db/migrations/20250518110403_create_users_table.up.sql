CREATE TABLE users (
  id VARCHAR(36) PRIMARY KEY,
  display_name VARCHAR(255) NOT NULL,
  birth_date TIMESTAMP NOT NULL,
  created_at TIMESTAMP NOT NULL,
  updated_at TIMESTAMP NOT NULL
);