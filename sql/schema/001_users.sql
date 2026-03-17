-- +goose Up
CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    is_chirpy_red BOOLEAN NOT NULL DEFAULT FALSE,
    email TEXT NOT NULL UNIQUE CHECK ( length(trim(email)) > 0 ),
    password TEXT NOT NULL CHECK ( length(trim(email)) > 3 )
);

-- +goose Down
DROP TABLE IF EXISTS users;
