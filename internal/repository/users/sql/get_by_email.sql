SELECT id, name, email, password, role
FROM users
WHERE email = $1;