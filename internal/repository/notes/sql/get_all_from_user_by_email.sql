SELECT id, title, content
FROM notes n
JOIN users u ON n.user_id = u.id
WHERE u.email = $1;