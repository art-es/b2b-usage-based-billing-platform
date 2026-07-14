SELECT 
    id, 
    email, 
    password_hash,
    verified_at IS NOT NULL AS is_verified
FROM users 
WHERE email = $1;

BEGIN;

INSERT INTO sessions (
    user_id, 
    refresh_token_hash, 
    refresh_token_expires_at
) 
VALUES ($1, $2, $3)
RETURNING id;

COMMIT;
