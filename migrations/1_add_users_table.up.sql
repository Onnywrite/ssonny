CREATE TABLE users (
    user_id             UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_nickname       VARCHAR(32) UNIQUE,
    user_email          VARCHAR(345) NOT NULL UNIQUE,
    user_verified       BOOLEAN NOT NULL,
    user_gender         VARCHAR(16),
    user_password_hash  CHAR(60),
    user_birthday       DATE,
    user_created_at     TIMESTAMP NOT NULL DEFAULT NOW(),
    user_updated_at     TIMESTAMP NOT NULL DEFAULT NOW(),
    user_deleted_at     TIMESTAMP
);