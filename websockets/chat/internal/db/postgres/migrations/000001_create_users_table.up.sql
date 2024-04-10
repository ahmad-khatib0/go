
CREATE TABLE IF NOT EXISTS users (
  id         BIGINT NOT NULL,
  created_at TIMESTAMP(3) NOT NULL,
  updated_at TIMESTAMP(3) NOT NULL,
  state      SMALLINT NOT NULL DEFAULT 0,
  state_at   TIMESTAMP(3),
  access     JSON,
  last_seen  TIMESTAMP,
  user_agent VARCHAR(255) DEFAULT '',
  public     JSON,
  trusted    JSON,
  tags       JSON,
  PRIMARY KEY(id)
);

CREATE INDEX IF NOT EXISTS users_state_state_at ON users(state, state_at);
CREATE INDEX IF NOT EXISTS users_last_seen_updated_at ON users(last_seen, updated_at);

