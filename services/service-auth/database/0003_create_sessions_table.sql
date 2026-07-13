CREATE TABLE sessions (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id uuid NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    refresh_token_hash VARCHAR NOT NULL,
    updated_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ DEFAULT current_timestamp
);
