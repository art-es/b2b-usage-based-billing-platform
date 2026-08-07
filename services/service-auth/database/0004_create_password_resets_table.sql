CREATE TABLE password_resets (
    token uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id uuid NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    sent_at TIMESTAMPTZ DEFAULT NULL,
    created_at TIMESTAMPTZ DEFAULT current_timestamp
);

CREATE INDEX idx_password_resets_sent_at_is_null 
    ON password_resets (created_at)
    WHERE sent_at IS NULL;

CREATE INDEX idx_password_resets_user_id_token 
    ON password_resets (user_id, token);
