CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS roles (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name TEXT NOT NULL UNIQUE,
    level INTEGER NOT NULL
);

INSERT INTO roles (name, level) VALUES 
    ('普通助理', 1),
    ('资深助理', 2),
    ('黑心', 3);