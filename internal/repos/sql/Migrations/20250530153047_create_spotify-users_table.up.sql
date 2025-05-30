CREATE TABLE IF NOT EXISTS rollbringer.spotify_users (
    spotify_id text NOT null UNIQUE,

    created_at timestamp NOT null DEFAULT now(),
    updated_at timestamp NOT null DEFAULT now(),

    display_name text NOT null,
    email text NOT null
);
