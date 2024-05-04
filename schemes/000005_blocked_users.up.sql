CREATE TABLE IF NOT EXISTS blocked_users (
  id SERIAL,
  user_id BIGINT NOT NULL,
  reason_id BIGINT NOT NULL,
  is_expire BOOLEAN DEFAULT false,
  expired_at TIMESTAMP(0) DEFAULT NULL,
  created_at TIMESTAMP(0) NOT NULL,
  CONSTRAINT blocked_users_pkey PRIMARY KEY (id),
  CONSTRAINT blocked_users_user_id_fk FOREIGN KEY (user_id) REFERENCES users(id),
  CONSTRAINT blocked_users_reason_id_fk FOREIGN KEY (reason_id) REFERENCES reasons_bans(id)
);