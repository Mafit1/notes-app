UPDATE refresh_tokens
SET revoked = TRUE,
    revoked_at = now()
WHERE token_id = $1 AND revoked = FALSE
