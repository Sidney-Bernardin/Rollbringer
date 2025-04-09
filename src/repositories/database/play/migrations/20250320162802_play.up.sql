CREATE SCHEMA IF NOT EXISTS play;

CREATE TABLE IF NOT EXISTS play.rooms (
    id uuid PRIMARY KEY,
    created_at timestamp NOT null DEFAULT now(),

    name text NOT null
);

CREATE TYPE room_user_permision AS ENUM ('OWNER', 'GAME_MASTER', 'PLAYER');
CREATE TABLE IF NOT EXISTS play.rooms_users (
    room_id uuid REFERENCES play.rooms (id) ON DELETE CASCADE,
    user_id uuid NOT NULL,
    permisions room_user_permision[] NOT NULL,
);

CREATE TABLE IF NOT EXISTS play.boards (
    id uuid PRIMARY KEY,
    created_at timestamp NOT null DEFAULT now(),

    name text NOT null,
    canvas jsonb NOT null
);

CREATE TYPE board_user_permision AS ENUM ('OWNER', 'EDIT');
CREATE TABLE IF NOT EXISTS play.board_users (
    baord_id uuid REFERENCES play.boards (id) ON DELETE CASCADE,
    user_id uuid NOT NULL,
    permisions board_user_permision[] NOT NULL,
);
