
-- indexed devices. Normalized into a separate table.
CREATE TABLE IF NOT EXISTS devices (
  id        SERIAL NOT NULL PRIMARY KEY,
  user_id   BIGINT NOT NULL,
  hash      CHAR(16) NOT NULL,
  device_id TEXT NOT NULL,
  platform  VARCHAR(32),
  last_seen TIMESTAMP NOT NULL,
  lang      VARCHAR(8),
  FOREIGN KEY(user_id) REFERENCES users(id)
);

CREATE UNIQUE INDEX IF NOT EXISTS devices_hash ON devices(hash);
