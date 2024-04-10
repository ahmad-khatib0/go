
-- Links between uploaded files and the topics, users or messages they are attached to
CREATE TABLE IF NOT EXISTS file_message_links(
  id         SERIAL NOT NULL PRIMARY KEY,
  created_at TIMESTAMP(3) NOT NULL,
  file_id    BIGINT NOT NULL,
  message_id INT,
  topic      VARCHAR(25),
  user_id    BIGINT,
  FOREIGN KEY(file_id)    REFERENCES file_uploads(id) ON DELETE CASCADE,
  FOREIGN KEY(message_id) REFERENCES messages(id)     ON DELETE CASCADE,
  FOREIGN KEY(topic)      REFERENCES topics(name)     ON DELETE CASCADE,
  FOREIGN KEY(user_id)    REFERENCES users(id)        ON DELETE CASCADE
);
