CREATE TABLE IF NOT EXISTS rollbringer.google_users (
    google_id text NOT null UNIQUE,

    created_at timestamp NOT null DEFAULT now(),
    updated_at timestamp NOT null DEFAULT now(),

    given_name text NOT null,
    email text NOT null
);
