CREATE TABLE users (
    id uuid PRIMARY KEY,

    created_at timestamp NOT null DEFAULT now(),
    updated_at timestamp NOT null DEFAULT now(),

    google_id text REFERENCES google_users (google_id) UNIQUE,

    username text NOT null,
    profile_picture text NOT null
);
