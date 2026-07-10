BEGIN;

SELECT v.token, u.email
FROM email_verifications AS v
JOIN users AS u ON u.id = v.user_id
WHERE v.sent_at is NULL
ORDER BY v.created_at
LIMIT $1
FOR UPDATE OF v SKIP LOCKED;

UPDATE email_verifications
SET sent_at = current_timestamp
WHERE token = ANY($1);

COMMIT;
