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
