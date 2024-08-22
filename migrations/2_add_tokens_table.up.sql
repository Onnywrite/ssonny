CREATE TABLE tokens (
	 token_id BIGSERIAL PRIMARY KEY,
	 token_user_fk UUID NOT NULL REFERENCES users(user_id) ON DELETE CASCADE,
	 token_app_fk BIGINT,
	 token_rotation BIGINT NOT NULL,
	 token_rotated_at TIMESTAMP NOT NULL,
	 token_platform VARCHAR(255) NOT NULL,
	 token_agent VARCHAR(255) NOT NULL
);