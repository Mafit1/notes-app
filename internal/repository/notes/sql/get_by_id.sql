SELECT id, title, content, user_id
FROM notes
WHERE id = $1;