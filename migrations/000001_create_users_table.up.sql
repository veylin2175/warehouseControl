CREATE TABLE users
(
    id            SERIAL PRIMARY KEY,
    username      VARCHAR(50) UNIQUE NOT NULL,
    password_hash TEXT               NOT NULL,
    role          VARCHAR(20)        NOT NULL CHECK (role IN ('admin', 'manager', 'viewer'))
);

INSERT INTO users (username, password_hash, role)
VALUES ('admin', '$2a$10$AO3opEdtEYoNYXJqxImcceq1f/iIvPGmNg8u13fWUDWCtrS4TYsiq', 'admin');