
-- records of uploaded files.
-- don't add foreign key on user_id. it's not needed and it will break user deletion.
-- using INDEX rather than FK on topic because it's either 'topics' or 'users' reference.
CREATE TABLE IF NOT EXISTS file_uploads(
  id         BIGINT NOT NULL PRIMARY KEY,
  created_at TIMESTAMP(3) NOT NULL,
  updated_at TIMESTAMP(3) NOT NULL,
  user_id    BIGINT,
  status     INT NOT NULL,
  mime_type  VARCHAR(255) NOT NULL,
  size       BIGINT NOT NULL,
  location   VARCHAR(2048) NOT NULL
);

CREATE INDEX IF NOT EXISTS fileup_loads_status ON file_uploads(status);

COMMENT ON COLUMN file_uploads.status    IS 'Status of upload';
COMMENT ON COLUMN file_uploads.mime_type IS 'Type of the file';
COMMENT ON COLUMN file_uploads.size      IS 'Size of the file in bytes';
COMMENT ON COLUMN file_uploads.location  IS 'Internal file location, i.e. path on disk or an S3 blob address';
