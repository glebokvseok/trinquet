CREATE TABLE IF NOT EXISTS court_service.court(
    id UUID PRIMARY KEY,
    name VARCHAR(256) NOT NULL,
    sport_type INT NULL DEFAULT 0,
    setting_type INT NULL DEFAULT 0,
    surface_type INT NULL DEFAULT 0,
    price INT,
    rating FLOAT,
    country VARCHAR(128),
    city VARCHAR(128),
    address VARCHAR(256),
    map_link VARCHAR(256),
    location GEOMETRY(Point, 4326),
    created_on TIMESTAMPTZ NOT NULL DEFAULT now(),
    modified_on TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_court_name_trgm
ON court_service.court USING gin (name gin_trgm_ops);
