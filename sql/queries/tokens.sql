-- name: StoreRefreshToken :exec
INSERT INTO refresh_tokens (token, created_at, updated_at, user_id, expires_at)
VALUES($1, now(), now(), $2, $3);

-- name: GetToken :one
SELECT *
FROM refresh_tokens
WHERE token=$1;

-- name: RevokeToken :exec
UPDATE refresh_tokens
SET revoked_at = CURRENT_TIMESTAMP,
    updated_at = CURRENT_TIMESTAMP
WHERE token = $1
AND revoked_at IS NULL AND expires_at > now();

