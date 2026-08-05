BEGIN;

SELECT 
    id, 
    user_id, 
    organization_id, 
    refresh_token_hash,
    refresh_token_expires_at
FROM sessions
WHERE id = $1 
FOR UPDATE SKIP LOCKED;

UPDATE sessions
SET organization_id = $2
WHERE id = $1

COMMIT;