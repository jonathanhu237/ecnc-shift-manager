CREATE TABLE IF NOT EXISTS roles (
    id BIGSERIAL PRIMARY KEY,
    name TEXT NOT NULL UNIQUE,
    level INTEGER NOT NULL
);

INSERT INTO roles (id, name, level) VALUES 
    (1, '普通助理', 1),
    (2, '资深助理', 2),
    (3, '黑心', 3);