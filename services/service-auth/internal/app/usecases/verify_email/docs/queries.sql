BEGIN;

SELECT v.token
FROM email_verifications AS v
JOIN users AS u ON u.id = v.user_id
WHERE v.token = $1 AND u.verified_at IS NULL
FOR UPDATE OF v SKIP LOCKED; 

UPDATE users
SET verified_at = current_timestamp
WHERE id = $1;

DELETE FROM email_verifications
WHERE user_id = $1;

COMMIT;
