CREATE TABLE IF NOT EXISTS reasons_bans (
  id SERIAL,
  name VARCHAR(64) NOT NULL,
  description VARCHAR(255) DEFAULT NULL,
  is_draft BOOLEAN DEFAULT true,
  updated_at TIMESTAMP(0),
  created_at TIMESTAMP(0) NOT NULL,
  CONSTRAINT reasons_bans_pkey PRIMARY KEY (id)
);