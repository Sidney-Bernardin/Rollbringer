-- name: InsertRoom :one
WITH
inserted_room AS (
    INSERT INTO rooms (id, name)
    VALUES (@room_id, @name)
    RETURNING *
),
inserted_user_room AS (
    INSERT INTO user_rooms (user_id, room_id, permisions)
    VALUES (@user_id::uuid, @room_id, @permisions)
    RETURNING *
)
SELECT
    inserted_room.id, inserted_room.name,
    array_agg(users.id)::uuid[] AS user_ids,
    array_agg(users.username)::text[] AS usernames,
    array_agg(users.profile_picture)::text[] AS profile_pictures,
    array_agg(inserted_user_room.permisions)::user_room_permision[][] AS permisions
FROM inserted_user_room
INNER JOIN inserted_room ON inserted_room.id = inserted_user_room.room_id
INNER JOIN users ON users.id = inserted_user_room.user_id
GROUP BY inserted_room.id, inserted_room.name;

-- name: SelectRoom :one
SELECT id, name
FROM rooms
WHERE id = @id;

-- name: SelectUserRooms :many
SELECT
    rooms.id, rooms.name,
    array_agg(users.id)::uuid[] AS user_ids,
    array_agg(users.username)::text[] AS usernames,
    array_agg(users.profile_picture)::text[] AS profile_pictures,
    array_agg(user_rooms.permisions)::user_room_permision[][] AS permisions
FROM user_rooms
INNER JOIN rooms ON rooms.id = user_rooms.room_id
INNER JOIN users ON users.id = user_rooms.user_id
WHERE user_rooms.user_id = @user_id::uuid
GROUP BY rooms.id;

-- name: DeleteRoom :exec
DELETE FROM rooms WHERE 
    rooms.id = @room_id::uuid AND
    EXISTS (
        SELECT 1 FROM user_rooms
        WHERE
            user_rooms.user_id = @user_id::uuid AND
            user_rooms.permisions @> ARRAY['OWNER']::user_room_permision[]
    );
