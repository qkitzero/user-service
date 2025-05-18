CREATE TABLE user (
  id VARCHAR(36) PRIMARY KEY,
  display_name VARCHAR(255) NOT NULL,
  birth_date DATETIME NOT NULL,
  created_at DATETIME NOT NULL,
  updated_at DATETIME NOT NULL
);