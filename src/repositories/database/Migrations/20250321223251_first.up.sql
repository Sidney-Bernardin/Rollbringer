CREATE SCHEMA IF NOT EXISTS accounts;

CREATE TABLE IF NOT EXISTS accounts.google_users (
    google_id text NOT null,
    created_at timestamp NOT null DEFAULT now(),

    given_name text NOT null,
    email text NOT null,
    profile_picture text NOT null,

    UNIQUE(google_id)
);

CREATE TABLE IF NOT EXISTS accounts.spotify_users (
    spotify_id text NOT null,
    created_at timestamp NOT null DEFAULT now(),

    display_name text NOT null,
    email text NOT null,
    profile_picture text,

    UNIQUE(spotify_id)
);

CREATE TABLE accounts.users (
    id uuid PRIMARY KEY,
    created_at timestamp NOT null DEFAULT now(),

    google_id text REFERENCES accounts.google_users (google_id),
    spotify_id text REFERENCES accounts.spotify_users (spotify_id),

    username text NOT null,
    profile_picture text,

    UNIQUE(google_id),
    UNIQUE(spotify_id)
);

CREATE TABLE IF NOT EXISTS accounts.sessions (
    id uuid PRIMARY KEY,
    created_at timestamp NOT null DEFAULT now(),

    user_id uuid REFERENCES accounts.users (id) ON DELETE CASCADE,

    csrf_token text NOT null,

    UNIQUE(user_id)
);

-- =====

CREATE SCHEMA IF NOT EXISTS play;

CREATE TABLE IF NOT EXISTS play.rooms (
    id uuid PRIMARY KEY,
    created_at timestamp NOT null DEFAULT now(),
    name text NOT null
);

-- =====

CREATE TYPE room_user_permision AS ENUM ('OWNER', 'GAME_MASTER', 'PLAYER');
CREATE TABLE IF NOT EXISTS room_users (
    room_id uuid REFERENCES play.rooms (id) ON DELETE CASCADE,
    user_id uuid REFERENCES accounts.users (id) ON DELETE CASCADE, 
    permisions room_user_permision[] NOT NULL,

    UNIQUE(room_id, user_id)
);
