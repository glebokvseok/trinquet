CREATE TABLE IF NOT EXISTS player_service.achievement(
    id UUID PRIMARY KEY,
    code VARCHAR(256) UNIQUE NOT NULL,
    name VARCHAR(256) NOT NULL,
    description VARCHAR(1024) NOT NULL,
    created_on TIMESTAMPTZ NOT NULL DEFAULT now(),
    modified_on TIMESTAMPTZ NOT NULL DEFAULT now()
)
