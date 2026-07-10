BEGIN;

SELECT v.token, u.email
FROM verifications AS v
JOIN users AS u ON u.id = v.user_id
WHERE v.email_sent_at is NULL
ORDER BY v.created_at
LIMIT $1
FOR UPDATE OF v SKIP LOCKED;

UPDATE verifications
SET email_sent_at = current_timestamp
WHERE token = ANY($1);

COMMIT;
