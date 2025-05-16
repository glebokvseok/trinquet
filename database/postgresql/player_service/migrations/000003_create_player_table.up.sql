CREATE TABLE IF NOT EXISTS player_service.player(
    id UUID PRIMARY KEY,
    username VARCHAR(64) UNIQUE NOT NULL,
    name VARCHAR(128),
    surname VARCHAR(128),
    birth_date DATE,
    gender INT NOT NULL DEFAULT 0,
    height INT,
    country VARCHAR(128),
    city VARCHAR(128),
    created_on TIMESTAMPTZ NOT NULL DEFAULT now(),
    modified_on TIMESTAMPTZ NOT NULL DEFAULT now()
);
