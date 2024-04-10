
-- Authentication records for the basic authentication scheme.
CREATE TABLE IF NOT EXISTS auth (
  id        SERIAL NOT NULL PRIMARY KEY,
  user_name VARCHAR(32) NOT NULL,
  user_id   BIGINT NOT NULL,
  scheme    VARCHAR(16) NOT NULL,
  level     INT NOT NULL,
  secret    VARCHAR(255) NOT NULL,
  expires   TIMESTAMP,
  FOREIGN KEY(user_id) REFERENCES users(id)
);

CREATE UNIQUE INDEX IF NOT EXISTS auth_user_id_scheme ON auth(user_id, scheme);
CREATE UNIQUE INDEX IF NOT EXISTS auth_user_name ON auth(user_name);
