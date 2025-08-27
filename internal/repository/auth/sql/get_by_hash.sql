SELECT token_id, user_id, token_hash, expires_at, revoked, revoked_at
FROM refresh_tokens
WHERE token_hash = $1