SELECT id, title, content
FROM notes
WHERE user_id = $1;