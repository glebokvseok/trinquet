CREATE TABLE IF NOT EXISTS auth_service.credential(
    user_id UUID PRIMARY KEY,
    email VARCHAR(128) UNIQUE NOT NULL,
    password_hash VARCHAR(128) NOT NULL,
    created_on TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_on TIMESTAMPTZ NOT NULL DEFAULT now()
)
