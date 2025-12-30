CREATE TABLE sessions (
  id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  user_email varchar(255) NOT NULL,
  refresh_token varchar(512) NOT NULL,
  is_revoked bool NOT NULL DEFAULT false,
  created_at datetime DEFAULT (now())
  expires_at datetime
);
