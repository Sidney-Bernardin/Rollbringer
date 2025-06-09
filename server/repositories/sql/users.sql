-- name: InsertUser :execrows
INSERT INTO users (id, google_id, username, profile_picture, password_hash, password_salt)
VALUES (@id, @google_id, @username, @profile_picture, @password_hash, @password_salt)
ON CONFLICT (username) DO NOTHING;

-- name: SelectUser :one
SELECT id, google_id, username, profile_picture
FROM users
WHERE id = @user_id;

-- name: SelectUserID :one
SELECT id
FROM users
WHERE google_id = @google_id;

-- name: SelectRoomUsers :many
SELECT
    users.id, users.google_id, users.username, users.profile_picture,
    user_rooms.permisions
FROM users
INNER JOIN user_rooms ON user_rooms.room_id = @room_id;

-- name: SelectUserWithPassword :one
SELECT id, google_id, username, profile_picture, password_hash, password_salt
FROM users
WHERE username = @username;
