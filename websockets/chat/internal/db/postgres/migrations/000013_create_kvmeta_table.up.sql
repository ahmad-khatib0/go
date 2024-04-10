
-- key value of metadata information
CREATE TABLE IF NOT EXISTS kvmeta(
  "key"      VARCHAR(64) NOT NULL PRIMARY KEY,
	created_at TIMESTAMP(3),
	"value"    TEXT
);

CREATE INDEX kvmeta_created_at_key ON kvmeta(created_at, "key");

