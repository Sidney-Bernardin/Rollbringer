CREATE EXTENSION hstore;



CREATE TABLE users (
    id UUID PRIMARY KEY,

    google_id text,
    username text NOT NULL,

    UNIQUE(google_id),
    UNIQUE(username)
);

CREATE TABLE sessions (
    id UUID PRIMARY KEY,

    user_id UUID NOT NULL REFERENCES users,
    csrf_token text NOT NULL,

    UNIQUE(user_id),
    UNIQUE(csrf_token)
);

CREATE TABLE games (
    id UUID PRIMARY KEY,

    host_id UUID NOT NULL REFERENCES users,
    name text NOT NULL
);

CREATE TABLE pdfs (
    id UUID PRIMARY KEY,

    owner_id UUID NOT NULL REFERENCES users,
    game_id UUID REFERENCES games,

    name text NOT NULL,
    schema text NOT NULL,

    pages hstore[] NOT NULL
);



CREATE TABLE user_joined_games (
    user_id UUID NOT NULL REFERENCES users,
    game_id UUID NOT NULL REFERENCES games
)
