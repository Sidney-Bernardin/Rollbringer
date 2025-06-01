-- name: InsertUser :exec
INSERT INTO users (id, google_id, spotify_id, username, profile_picture)
VALUES (@id, @google_id, @spotify_id, @username, @profile_picture);

-- name: GetUser :one
SELECT * FROM users
WHERE id = @user_id;

-- name: InsertGoogleUser :execrows
INSERT INTO google_users (google_id, given_name, email)
VALUES (@google_id, @given_name, @email)
ON CONFLICT (google_id) DO NOTHING;

-- name: UpdateGoogleUser :execrows
UPDATE google_users
SET
    google_id = coalese(@google_id, google_id),
    given_name = coalese(@given_name, given_name),
    email = coalese(@email, email);
