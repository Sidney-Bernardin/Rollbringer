CREATE TABLE IF NOT EXISTS user_rooms (
    user_id uuid REFERENCES users (id) ON DELETE CASCADE,
    room_id uuid REFERENCES rooms (id) ON DELETE CASCADE,
    permisions user_room_permision[] NOT NULL
);
