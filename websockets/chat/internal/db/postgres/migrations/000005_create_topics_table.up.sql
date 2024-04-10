
CREATE TABLE IF NOT EXISTS topics (
  id          serial not null primary key,
  created_at  timestamp(3) not null,
  updated_at  timestamp(3) not null,
  state       smallint not null default 0,
  state_at    timestamp(3),
  touched_at  timestamp(3),
  name        varchar(25) not null,
  use_bt      boolean default false,
  owner       bigint not null default 0,
  access      json,
  seq_id      int not null default 0,
  del_id      int default 0,
  public      json,
  trusted     json,
  tags        json
);

CREATE UNIQUE INDEX IF NOT EXISTS topics_name ON topics(name);
CREATE INDEX IF NOT EXISTS        topics_owner ON topics(owner);
CREATE INDEX IF NOT EXISTS        topics_state_state_at ON topics(state, state_at);
