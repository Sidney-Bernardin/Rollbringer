CREATE EXTENSION hstore;



CREATE SCHEMA users;

CREATE TABLE users.users (
    id UUID PRIMARY KEY,

    google_id text,
    username text NOT NULL,

    UNIQUE(google_id),
    UNIQUE(username)
);

CREATE TABLE users.sessions (
    id UUID PRIMARY KEY,

    user_id UUID NOT NULL REFERENCES users ON DELETE CASCADE,
    csrf_token text NOT NULL,

    UNIQUE(user_id),
    UNIQUE(csrf_token)
);



CREATE SCHEMA games;

CREATE TABLE games.games (
    id UUID PRIMARY KEY,

    host_id UUID NOT NULL REFERENCES users ON DELETE CASCADE,
    name text NOT NULL
);

CREATE TABLE games.pdfs (
    id UUID PRIMARY KEY,

    owner_id UUID NOT NULL REFERENCES users ON DELETE CASCADE,
    game_id UUID REFERENCES games ON DELETE SET NULL,

    name text NOT NULL,
    schema text NOT NULL,

    fields hstore[] NOT NULL
);

CREATE TABLE games.rolls (
    id UUID PRIMARY KEY,

    owner_id UUID NOT NULL REFERENCES users ON DELETE CASCADE,
    game_id UUID NOT NULL REFERENCES games ON DELETE CASCADE,

    dice_names int32[] NOT NULL,
    dice_results int32[] NOT NULL
);



CREATE TABLE GameGuests (
    guest_id UUID NOT NULL REFERENCES users,
    game_id UUID NOT NULL REFERENCES games,
)
