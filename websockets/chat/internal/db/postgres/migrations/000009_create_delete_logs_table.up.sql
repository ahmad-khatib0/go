
CREATE TABLE IF NOT EXISTS delete_logs(
  id          SERIAL NOT NULL PRIMARY KEY,
  topic       VARCHAR(25) NOT NULL,
  deleted_for BIGINT NOT NULL DEFAULT 0,
  del_id      INT NOT NULL,
  low         INT NOT NULL,
  hi          INT NOT NULL,
  FOREIGN KEY(topic) REFERENCES topics(name)
);

-- For getting the list of deleted message ranges
CREATE INDEX IF NOT EXISTS dellog_topic_delid_deletedfor  ON delete_logs(topic, del_id, deleted_for);

-- used when getting not-yet-deleted messages(messages LEFT JOIN delete_logs)
CREATE INDEX IF NOT EXISTS dellog_topic_deletedfor_low_hi ON delete_logs(topic, deleted_for, low, hi);

-- Used when deleting a user
CREATE INDEX IF NOT EXISTS dellog_deletedfor              ON delete_logs(deleted_for);

