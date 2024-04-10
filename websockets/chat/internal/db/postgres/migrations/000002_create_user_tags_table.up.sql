
-- indexed user tags.

CREATE TABLE IF NOT EXISTS user_tags (
  id       SERIAL NOT NULL PRIMARY KEY,
  user_id  BIGINT NOT NULL,
  tag      VARCHAR(96) NOT NULL,
  FOREIGN KEY(user_id) REFERENCES users(id)
);

CREATE INDEX IF NOT EXISTS user_tags_tag ON user_tags(tag);
CREATE UNIQUE INDEX IF NOT EXISTS user_tags_uid_tag ON user_tags(user_id, tag);
