CREATE TABLE users (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(100) NOT NULL,
    email VARCHAR(100) NOT NULL,
    password_hash VARCHAR NOT NULL,
    created_at TIMESTAMPTZ DEFAULT current_timestamp,
    verified_at TIMESTAMPTZ DEFAULT NULL
);

CREATE UNIQUE INDEX idx_users_email_unique 
    ON users (email);

CREATE TABLE verifications (
    token uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id uuid NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    email_sent_at TIMESTAMPTZ DEFAULT NULL,
    created_at TIMESTAMPTZ DEFAULT current_timestamp
);

CREATE INDEX idx_verifications_email_sent_at_is_null 
    ON verifications (created_at)
    WHERE email_sent_at IS NULL;

CREATE INDEX idx_verifications_user_id_token 
    ON verifications (user_id, token);
