
-- Indexed topic tags.
CREATE TABLE IF NOT EXISTS topic_tags (
  id    SERIAL NOT NULL PRIMARY KEY,
	topic VARCHAR(25) NOT NULL,
	tag   VARCHAR(96) NOT NULL,
  FOREIGN KEY(topic) REFERENCES topics(name)
);

CREATE INDEX IF NOT EXISTS        topic_tags_tag        ON topic_tags(tag);
CREATE UNIQUE INDEX IF NOT EXISTS topic_tags_userid_tag ON topic_tags(topic, tag);
