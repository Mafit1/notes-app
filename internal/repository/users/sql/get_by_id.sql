SELECT id, name, email, password, role
FROM users
WHERE id = $1;