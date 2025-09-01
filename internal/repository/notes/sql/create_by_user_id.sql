INSERT INTO notes (user_id, title, content)
VALUES ($1, $2, $3)
RETURNING id;