CREATE TABLE memberships (
  id        VARCHAR(36) PRIMARY KEY,
  user_id   VARCHAR(36) NOT NULL,
  group_id  VARCHAR(36) NOT NULL,
  role      VARCHAR(32) NOT NULL,
  joined_at TIMESTAMP   NOT NULL,
  UNIQUE (user_id, group_id),
  FOREIGN KEY (user_id)  REFERENCES users(id)  ON DELETE CASCADE,
  FOREIGN KEY (group_id) REFERENCES groups(id) ON DELETE CASCADE
);
