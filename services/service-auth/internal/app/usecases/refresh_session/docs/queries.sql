BEGIN;

SELECT 
    id, 
    user_id, 
    organization_id, 
    refresh_token_hash,
    refresh_token_expires_at
FROM sessions
WHERE refresh_token_hash = $1 
FOR UPDATE SKIP LOCKED;

UPDATE sessions
SET 
    refresh_token_hash = $1,
    refresh_token_expires_at = $2,
    updated_at = current_timestamp
WHERE id = $3

COMMIT;