CREATE TABLE IF NOT EXISTS chat_service.chat(
    id UUID PRIMARY KEY,
    key VARCHAR(256) UNIQUE NOT NULL,
    type INT NOT NULL,
    created_on TIMESTAMPTZ NOT NULL DEFAULT now()
);
