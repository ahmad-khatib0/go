CREATE TABLE IF NOT EXISTS subscriptions (
  id               SERIAL NOT NULL PRIMARY KEY,
  created_at       TIMESTAMP(3) NOT NULL,
  updated_at       TIMESTAMP(3) NOT NULL,
  deleted_at       TIMESTAMP(3),
  user_id          BIGINT NOT NULL,
  topic            VARCHAR(25) NOT NULL,
  del_id           INT DEFAULT 0,
  received_seq_id  INT DEFAULT 0,
  read_seq_id      INT DEFAULT 0,
  mode_want        VARCHAR(8),
  mode_given       VARCHAR(8),
  private          JSON,
  FOREIGN KEY(user_id) REFERENCES users(id)
);
  
CREATE UNIQUE INDEX IF NOT EXISTS subscriptions_topic_uid ON subscriptions(topic, user_id);
CREATE INDEX IF NOT EXISTS subscriptions_topic ON subscriptions(topic);
CREATE INDEX IF NOT EXISTS subscriptions_deletedat ON subscriptions(deleted_at);

COMMENT ON COLUMN subscriptions.del_id IS 
  'ID of the latest Soft-delete operation';

COMMENT ON COLUMN subscriptions.received_seq_id IS 
  'Last SeqId reported by user as received by at least one of his sessions';

COMMENT ON COLUMN subscriptions.read_seq_id IS
  'Last SeqID reported read by the user';

COMMENT ON COLUMN subscriptions.mode_want IS 
  'Access mode requested by this user';

COMMENT ON COLUMN subscriptions.mode_given IS 
  'Access mode granted to this user';

COMMENT ON COLUMN subscriptions.private IS
  'user private data associated with the subscription to topic';
