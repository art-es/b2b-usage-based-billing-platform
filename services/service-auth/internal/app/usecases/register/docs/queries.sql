BEGIN;

INSERT INTO users (name, email, password)
VALUES ($1, $2, $3)
RETURNING id;

INSERT INTO email_verifications (user_id)
VALUES ($1);

COMMIT;
