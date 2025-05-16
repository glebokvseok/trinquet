CREATE INDEX IF NOT EXISTS idx_player_username_trgm
ON player_service.player USING gin (username gin_trgm_ops);

CREATE INDEX IF NOT EXISTS idx_player_name_trgm
ON player_service.player USING gin (name gin_trgm_ops);

CREATE INDEX IF NOT EXISTS idx_player_surname_trgm
ON player_service.player USING gin (surname gin_trgm_ops);
