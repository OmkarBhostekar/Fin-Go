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


-- name: GetAcountForUpdate :one
SELECT * FROM accounts WHERE id = $1 LIMIT 1 FOR NO KEY UPDATE;

-- name: GetAllAccounts :many
SELECT * FROM accounts 
WHERE owner = $1
ORDER BY id
LIMIT $2
OFFSET $3;

-- name: UpdateAccountBalance :one
UPDATE accounts
SET balance = balance + sqlc.arg(amount)
WHERE id = sqlc.arg(id)
RETURNING *;

-- name: DeleteAccountById :exec
DELETE FROM accounts
WHERE id = $1;

-- name: GetAccountsCount :one
SELECT COUNT(*) FROM accounts;