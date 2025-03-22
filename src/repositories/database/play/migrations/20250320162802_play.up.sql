CREATE SCHEMA IF NOT EXISTS play;

CREATE TABLE play.rooms (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    created_at timestamp NOT null DEFAULT now(),

    owner_id uuid NOT null,

    name text NOT null
);
