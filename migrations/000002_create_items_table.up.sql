CREATE TABLE items
(
    id         SERIAL PRIMARY KEY,
    name       VARCHAR(100) NOT NULL,
    quantity   INTEGER      NOT NULL DEFAULT 0,
    created_at TIMESTAMP             DEFAULT NOW(),
    updated_at TIMESTAMP             DEFAULT NOW()
);