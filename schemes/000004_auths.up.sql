CREATE TABLE IF NOT EXISTS auths (
  id SERIAL,
  user_id BIGINT NOT NULL,
  access_token VARCHAR(255) NOT NULL,
  refresh_token VARCHAR(255) NOT NULL,
  ip VARCHAR(64) DEFAULT NULL,
  device VARCHAR(30) DEFAULT NULL,
  user_agent VARCHAR(200) DEFAULT NULL,
  created_at TIMESTAMP(0) NOT NULL,
  CONSTRAINT auths_pkey PRIMARY KEY (id),
  CONSTRAINT auths_user_id_fk FOREIGN KEY (user_id) REFERENCES users(id)
);

CREATE INDEX auths_user_agent_key ON auths(user_agent);
CREATE INDEX auths_ip_key ON auths(ip);
CREATE INDEX auths_device_key ON auths(device);
CREATE INDEX auths_refresh_token_key ON auths(refresh_token);