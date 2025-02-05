CREATE SCHEMA IF NOT EXISTS accounts;

CREATE TABLE IF NOT EXISTS accounts.users (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    created_at timestamp NOT null DEFAULT now(),

    username text NOT null,

    UNIQUE(username)
);

CREATE TABLE IF NOT EXISTS accounts.google_users (
    created_at timestamp NOT null DEFAULT now(),

    user_id uuid REFERENCES accounts.users ON DELETE CASCADE,

    google_id text NOT null,
    given_name text NOT null,
    profile_picture text NOT null,

    UNIQUE(google_id)
);

CREATE TABLE IF NOT EXISTS accounts.spotify_users (
    created_at timestamp NOT null DEFAULT now(),

    user_id uuid REFERENCES accounts.users ON DELETE CASCADE,

    spotify_id text NOT null,
    display_name text NOT null,
    profile_picture text NOT null,

    UNIQUE(spotify_id)
);

CREATE TABLE IF NOT EXISTS accounts.sessions (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    created_at timestamp NOT null DEFAULT now(),

    user_id uuid REFERENCES accounts.users ON DELETE CASCADE,

    csrf_token text NOT null
);
