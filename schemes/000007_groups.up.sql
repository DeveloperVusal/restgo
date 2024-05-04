CREATE TABLE IF NOT EXISTS groups (
  id SERIAL,
  name VARCHAR(32) NOT NULL,
  is_draft BOOLEAN DEFAULT true,
  updated_at TIMESTAMP(0),
  created_at TIMESTAMP(0) NOT NULL,
  CONSTRAINT groups_pkey PRIMARY KEY (id)
);