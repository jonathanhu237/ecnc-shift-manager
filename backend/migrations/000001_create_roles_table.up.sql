CREATE TABLE IF NOT EXISTS roles (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL UNIQUE,
    level INTEGER NOT NULL
)