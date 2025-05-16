CREATE TABLE IF NOT EXISTS player_service.avatar(
    player_id UUID NOT NULL REFERENCES player_service.player(id),
    avatar_id UUID NOT NULL,
    mime_type VARCHAR(64) NOT NULL,
    is_used BOOLEAN NOT NULL,
    created_on TIMESTAMPTZ NOT NULL DEFAULT now(),
    modified_on TIMESTAMPTZ NOT NULL DEFAULT now(),
    PRIMARY KEY (player_id, avatar_id)
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_player_id_is_used
ON player_service.avatar(player_id, is_used)
WHERE is_used = TRUE;
