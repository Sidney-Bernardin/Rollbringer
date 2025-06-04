-- name: CreateUser :execrows
INSERT INTO users (id, google_id, username, profile_picture, password_hash, password_salt)
VALUES (@id, @google_id, @username, @profile_picture, @password_hash, @password_salt)
ON CONFLICT (username) DO NOTHING;

-- name: GetUser :one
SELECT id, google_id, username, profile_picture
FROM users
WHERE id = @user_id;

-- name: GetUserID :one
SELECT id
FROM users
WHERE google_id = @google_id;

-- name: GetUserWithPassword :one
SELECT id, google_id, username, profile_picture, password_hash, password_salt
FROM users
WHERE username = @username;
