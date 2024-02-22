CREATE TABLE users (
    id UUID NOT NULL,
    google_id text,

    username text NOT NULL,

    UNIQUE(id),
    UNIQUE(google_id),
    UNIQUE(username)
);

CREATE TABLE sessions (
    id UUID NOT NULL,
    user_id UUID NOT NULL,
    csrf_token UUID NOT NULL,

    UNIQUE(id),
    UNIQUE(user_id),
    UNIQUE(csrf_token)
);

CREATE TABLE games (
    id UUID NOT NULL,
    host_id UUID NOT NULL,

    title text NOT NULL,

    UNIQUE(id)
);

CREATE TABLE rolls (
    id UUID NOT NULL,
    game_id UUID NOT NULL,

    die_expressions integer[] NOT NULL,
    die_results integer[] NOT NULL,

    modifier_expression text NOT NULL,
    modifier_result decimal NOT NULL,

    unique(id)
);

CREATE TABLE pdfs (
    id UUID NOT NULL,
    owner_id UUID NOT NULL,

    name text NOT NULL,
    schema text NOT NULL,

    pages text[] NOT NULL,

    UNIQUE(id)
);
