CREATE TABLE IF NOT EXISTS messages(
  id         SERIAL NOT NULL PRIMARY KEY,
  created_at TIMESTAMP(3) NOT NULL,
  updated_at TIMESTAMP(3) NOT NULL,
  deleted_at TIMESTAMP(3),
  del_id     INT DEFAULT 0,
  seq_id     INT NOT NULL,
  topic      VARCHAR(25) NOT NULL,
  "from"     BIGINT NOT NULL,
  head       JSON,
  content    JSON,
  FOREIGN KEY(topic) REFERENCES topics(name)
);

CREATE UNIQUE INDEX IF NOT EXISTS messages_topic_seqid ON messages(topic, seq_id);

COMMENT ON COLUMN messages.del_id IS 'ID of the hard-delete operation';
