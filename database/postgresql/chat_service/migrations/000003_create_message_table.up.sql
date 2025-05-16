CREATE TABLE IF NOT EXISTS chat_service.message(
    id UUID PRIMARY KEY,
    chat_id UUID NOT NULL REFERENCES chat_service.chat(id),
    user_id UUID NOT NULL,
    type INT NOT NULL,
    status INT NOT NULL,
    status_description VARCHAR(1024),
    content JSONB NOT NULL,
    client_created_on TIMESTAMPTZ NOT NULL,
    client_modified_on TIMESTAMPTZ NOT NULL,
    created_on TIMESTAMPTZ NOT NULL DEFAULT now(),
    modified_on TIMESTAMPTZ NOT NULL DEFAULT now()
);
