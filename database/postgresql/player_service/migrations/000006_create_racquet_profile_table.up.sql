CREATE TABLE IF NOT EXISTS player_service.racquet_profile(
    id UUID PRIMARY KEY,
    player_id UUID NOT NULL,
    sport_type INT NOT NULL,
    best_hand INT NOT NULL DEFAULT 0,
    court_side INT NOT NULL DEFAULT 0,
    rating INT NOT NULL DEFAULT 0,
    match_count INT NOT NULL DEFAULT 0,
    win_count INT NOT NULL DEFAULT 0,
    loss_count INT NOT NULL DEFAULT 0,
    created_on TIMESTAMPTZ NOT NULL DEFAULT now(),
    modified_on TIMESTAMPTZ NOT NULL DEFAULT now(),
    FOREIGN KEY (player_id) REFERENCES player_service.player(id)
);

ALTER TABLE player_service.racquet_profile
ADD CONSTRAINT uq_player_id_sport_type UNIQUE (player_id, sport_type);
