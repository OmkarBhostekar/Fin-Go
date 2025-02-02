-- name: CreateVerifyEmail :one
INSERT INTO verify_emails (
    username,
    email,
    secret_code
) VALUES (
    $1,
    $2,
    $3
) RETURNING *;

-- name: UpdateVerifyEmail :one
UPDATE verify_emails
SET
    is_used = TRUE
WHERE
    id = @id
    AND is_used = FALSE
    AND secret_code = @secret_code
    AND expired_at > NOW()
RETURNING *;
