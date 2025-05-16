CREATE TABLE IF NOT EXISTS chat_service.message_event(
    id SERIAL PRIMARY KEY,
    message_id UUID NOT NULL REFERENCES chat_service.message(id),
    type INT NOT NULL,
    state INT NOT NULL,
    state_description VARCHAR(1024),
    data JSONB,
    created_on TIMESTAMPTZ NOT NULL DEFAULT now(),
    modified_on TIMESTAMPTZ NOT NULL DEFAULT now()
);
