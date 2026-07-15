CREATE TABLE sessions (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id uuid NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    organization_id uuid DEFAULT NULL,
    refresh_token_hash VARCHAR NOT NULL,
    refresh_token_expires_at TIMESTAMPTZ NOT NULL,
    updated_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ DEFAULT current_timestamp
);
