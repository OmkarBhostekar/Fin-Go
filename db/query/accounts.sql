-- name: CreateAccount :one
INSERT INTO accounts(
    owner,
    balance,
    currency
) VALUES (
    $1,
    $2,
    $3
) RETURNING *;

-- name: GetAcountById :one
SELECT * FROM accounts WHERE id = $1 LIMIT 1;

-- name: GetAllAccounts :many
SELECT * FROM accounts 
ORDER BY id
LIMIT $1
OFFSET $2;

-- name: UpdateAccountBalance :exec
UPDATE accounts
SET balance = $1
WHERE id = $2
RETURNING *;

-- name: DeleteAccountById :exec
DELETE FROM accounts
WHERE id = $1;