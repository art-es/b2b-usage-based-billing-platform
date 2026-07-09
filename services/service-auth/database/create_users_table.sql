CREATE TABLE users (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(100) NOT NULL,
    email VARCHAR(100) NOT NULL,
    password VARCHAR NOT NULL,
    created_at TIMESTAMP DEFAULT current_timestamp,
    activation_sent_at TIMESTAMP DEFAULT NULL,
    activated_at TIMESTAMP DEFAULT NULL
);
