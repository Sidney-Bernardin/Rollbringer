CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    google_id text,

    username text NOT NULL,

    UNIQUE(google_id),
    UNIQUE(username)
);

CREATE TABLE sessions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL,
    csrf_token text NOT NULL,

    UNIQUE(user_id),
    UNIQUE(csrf_token)
);

CREATE TABLE games (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    host_id UUID NOT NULL,

    title text NOT NULL
);

CREATE TABLE pdfs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    owner_id UUID NOT NULL,
    game_id UUID NOT NULL,

    name text NOT NULL,
    schema text NOT NULL,

    fields jsonb NOT NULL
);
