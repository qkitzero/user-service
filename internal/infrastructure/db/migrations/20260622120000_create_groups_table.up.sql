CREATE TABLE groups (
  id         VARCHAR(36) PRIMARY KEY,
  name       VARCHAR(255) NOT NULL,
  created_at TIMESTAMP NOT NULL,
  updated_at TIMESTAMP NOT NULL
);
