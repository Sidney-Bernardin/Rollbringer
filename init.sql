CREATE TABLE users (
    id UUID NOT NULL,
    google_id text,
    username VARCHAR(50),

    UNIQUE(google_id),
    UNIQUE(username)
);

CREATE TABLE sessions (
    id UUID NOT NULL,
    user_id UUID NOT NULL,
    csrf_token UUID NOT NULL,

    UNIQUE(user_id)
);
