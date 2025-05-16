CREATE TABLE IF NOT EXISTS player_service.racquet_match(
    id UUID PRIMARY KEY,
    owner_id UUID NOT NULL,
    sport_type INT NOT NULL,
    match_type INT NOT NULL,
    is_competitive BOOLEAN NOT NULL,
    court_id UUID,
    state INT NOT NULL,
    started_on TIMESTAMPTZ,
    finished_on TIMESTAMPTZ,
    scheduled_on TIMESTAMPTZ NOT NULL,
    created_on TIMESTAMPTZ NOT NULL DEFAULT now(),
    modified_on TIMESTAMPTZ NOT NULL DEFAULT now(),
    FOREIGN KEY(owner_id) REFERENCES player_service.racquet_profile(id)
);
