DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'users_confirm_status') THEN
        CREATE TYPE users_confirm_status AS ENUM ('quest', 'success', 'waiting', 'error', 'unknown');
    END IF;
    --more types here...
END$$;

CREATE TABLE IF NOT EXISTS users (
  id SERIAL,
  group_id BIGINT NOT NULL,
  email VARCHAR(150) NOT NULL,
  password VARCHAR(255) NOT NULL,
  activation BOOL DEFAULT false,
  name VARCHAR(200),
  surname VARCHAR(200),
  token_secret_key TEXT NOT NULL,
  last_activity_at TIMESTAMP(0),
  updated_at TIMESTAMP(0),
  created_at TIMESTAMP(0) NOT NULL,
  confirm_code CHAR(6) DEFAULT NULL::bpchar,
  confirm_action VARCHAR(10) DEFAULT NULL,
  confirmed_at TIMESTAMP(0) DEFAULT NULL::TIMESTAMP without time zone,
  confirm_status users_confirm_status DEFAULT 'unknown',
  CONSTRAINT users_pkey PRIMARY KEY (id),
  CONSTRAINT users_email_key UNIQUE (email),
  CONSTRAINT groups_group_id_fk FOREIGN KEY (group_id) REFERENCES groups(id)
);