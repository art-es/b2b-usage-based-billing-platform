SELECT 
    u.verified_at IS NOT NULL AS is_verified,
    EXISTS (
        SELECT 1 
        FROM email_verifications AS v
        WHERE 
            v.user_id = u.id
            AND v.sent_at IS NULL
    ) AS has_unsent
FROM users AS u
WHERE u.email = $1;


INSERT INTO email_verifications (user_id)
SELECT id FROM users WHERE email = $1;
