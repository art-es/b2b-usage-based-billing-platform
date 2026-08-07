CREATE TABLE email_verifications (
    token uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id uuid NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    sent_at TIMESTAMPTZ DEFAULT NULL,
    created_at TIMESTAMPTZ DEFAULT current_timestamp
);

CREATE INDEX idx_email_verifications_sent_at_is_null 
    ON email_verifications (created_at)
    WHERE sent_at IS NULL;

CREATE INDEX idx_email_verifications_user_id_token 
    ON email_verifications (user_id, token);
