CREATE TABLE group_relations (
  parent_id VARCHAR(36) NOT NULL,
  child_id  VARCHAR(36) NOT NULL,
  PRIMARY KEY (parent_id, child_id),
  FOREIGN KEY (parent_id) REFERENCES groups(id) ON DELETE CASCADE,
  FOREIGN KEY (child_id)  REFERENCES groups(id) ON DELETE CASCADE
);
