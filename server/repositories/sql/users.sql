-- name: CreateUser :exec
INSERT INTO users (id, google_id, username, profile_picture)
VALUES (@id, @google_id, @username, @profile_picture);

-- name: GetUser :one
SELECT * FROM users
WHERE id = @user_id;

-- name: GetUserID :one
SELECT id FROM users
WHERE google_id = @google_id;
