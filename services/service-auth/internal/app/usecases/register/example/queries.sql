BEGIN;

INSERT INTO users (name, email, password)
VALUES ($1, $2, $3)
RETURNING id;

INSERT INTO verifications (user_id)
VALUES ($1);

COMMIT;
