-- name: InsertGoogleUser :execrows
INSERT INTO google_users (google_id, given_name, email)
VALUES (@google_id, @given_name, @email)
ON CONFLICT (google_id) DO NOTHING;

-- name: UpdateGoogleUser :execrows
UPDATE google_users
SET
    given_name = coalesce(@given_name, given_name),
    email = coalesce(@email, email)
WHERE google_id = @google_id;
