CREATE TABLE IF NOT EXISTS court_service.court_media(
    court_id UUID NOT NULL REFERENCES court_service.court(id),
    media_id UUID NOT NULL,
    is_preview BOOLEAN NOT NULL DEFAULT FALSE,
    mime_type VARCHAR(64) NOT NULL,
    created_on TIMESTAMPTZ NOT NULL DEFAULT now(),
    PRIMARY KEY (court_id, media_id)
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_court_id_is_preview
ON court_service.court_media(court_id, is_preview)
WHERE is_preview = TRUE;
