
INSERT INTO  topics(
  created_at,
  updated_at,
  state,
  touched_at,
  name,
  access,
  public
) VALUES ( 
  CURRENT_TIMESTAMP(3),
  CURRENT_TIMESTAMP(3),
  0,
  CURRENT_TIMESTAMP(3),
  'sys',
  '{"Auth": "N","Anon": "N"}',
  '{"fn": "System"}'
);
