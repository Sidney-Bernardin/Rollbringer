CREATE TABLE IF NOT EXISTS users (
    id uuid PRIMARY KEY,

    created_at timestamp NOT null DEFAULT now(),
    updated_at timestamp NOT null DEFAULT now(),

    google_id text REFERENCES google_users (google_id) UNIQUE,

    username text NOT null UNIQUE,
    profile_picture text NOT null,
    password_hash bytea,
    password_salt text
);
