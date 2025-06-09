CREATE TABLE IF NOT EXISTS rooms (
    id uuid PRIMARY KEY,

    created_at timestamp NOT null DEFAULT now(),
    updated_at timestamp NOT null DEFAULT now(),

    name text NOT null
);
