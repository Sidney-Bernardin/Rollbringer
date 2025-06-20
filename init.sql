CREATE EXTENSION hstore;



CREATE SCHEMA users;

CREATE TABLE users.users (
    id UUID PRIMARY KEY,
    username text NOT NULL,

    google_id text,
    google_picture text,

    UNIQUE(google_id),
    UNIQUE(username)
);

CREATE TABLE users.sessions (
    id UUID PRIMARY KEY,

    user_id UUID NOT NULL REFERENCES users.users ON DELETE CASCADE,
    csrf_token text NOT NULL,

    UNIQUE(user_id),
    UNIQUE(csrf_token)
);



CREATE SCHEMA games;

CREATE TABLE games.games (
    id UUID PRIMARY KEY,
    created_at timestamp NOT NULL DEFAULT now(),

    host_id UUID NOT NULL REFERENCES users.users ON DELETE CASCADE,
    name text NOT NULL
);

CREATE TABLE games.pdfs (
    id UUID PRIMARY KEY,
    created_at timestamp NOT NULL DEFAULT now(),

    owner_id UUID NOT NULL REFERENCES users.users ON DELETE CASCADE,
    game_id UUID REFERENCES games.games ON DELETE SET NULL,

    name text NOT NULL,
    schema text NOT NULL,

    pages hstore[] NOT NULL
);

CREATE TABLE games.rolls (
    id UUID PRIMARY KEY,
    created_at timestamp NOT NULL DEFAULT now(),

    owner_id UUID NOT NULL REFERENCES users.users ON DELETE CASCADE,
    game_id UUID NOT NULL REFERENCES games.games ON DELETE CASCADE,

    dice_types int[] NOT NULL,
    dice_results int[] NOT NULL,
    modifiers text NOT NULL
);

CREATE TABLE games.chat_messages (
    id UUID PRIMARY KEY,
    created_at timestamp NOT NULL DEFAULT now(),

    owner_id UUID NOT NULL REFERENCES users.users ON DELETE CASCADE,
    game_id UUID NOT NULL REFERENCES games.games ON DELETE CASCADE,

    message text NOT NULL
);



CREATE TABLE game_users (
    game_id UUID NOT NULL REFERENCES games.games,
    user_id UUID NOT NULL REFERENCES users.users
);
