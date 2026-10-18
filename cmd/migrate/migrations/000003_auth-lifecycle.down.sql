DROP TABLE auth_sessions;
DROP TABLE auth_action_tokens;
DROP INDEX users_username_lower_uidx;
DROP INDEX users_email_lower_uidx;
ALTER TABLE users ADD CONSTRAINT users_username_key UNIQUE (username);
ALTER TABLE users ADD CONSTRAINT users_email_key UNIQUE (email);
