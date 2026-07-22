CREATE TABLE organizations (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR NOT NULL,
    created_at TIMESTAMPTZ DEFAULT current_timestamp,
    updated_at TIMESTAMPTZ
);
