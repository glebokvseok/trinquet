CREATE TABLE IF NOT EXISTS player_service.completed_achievement(
    player_id UUID NOT NULL,
    achievement_id UUID NOT NULL,
    completed_on TIMESTAMPTZ NOT NULL DEFAULT now(),
    PRIMARY KEY(player_id, achievement_id),
    FOREIGN KEY (player_id) REFERENCES player_service.player(id),
    FOREIGN KEY (achievement_id) REFERENCES player_service.achievement(id)
);
