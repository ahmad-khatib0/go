CREATE TABLE IF NOT EXISTS credentials(
  id         SERIAL NOT NULL PRIMARY KEY,
  created_at TIMESTAMP(3) NOT NULL,
  updated_at TIMESTAMP(3) NOT NULL,
  deleted_at TIMESTAMP(3),
  method     VARCHAR(16) NOT NULL,
  value      VARCHAR(128) NOT NULL,
  synthetic  VARCHAR(192) NOT NULL,
  user_id    BIGINT NOT NULL,
  response   VARCHAR(255),
  done       BOOLEAN NOT NULL DEFAULT FALSE,
  retries    INT NOT NULL DEFAULT 0,
  FOREIGN KEY(user_id) REFERENCES users(id)
);

CREATE UNIQUE INDEX IF NOT EXISTS credentials_uniqueness ON credentials(synthetic);

COMMENT ON COLUMN credentials.method IS 'Verification method (email, tel, captcha, etc)';
COMMENT ON COLUMN credentials.value IS 'Credential value like jdoe@example.com or +12345678901';
COMMENT ON COLUMN credentials.response IS 'Expected response';
COMMENT ON COLUMN credentials.done IS 'if credential was successfully confirmed';
