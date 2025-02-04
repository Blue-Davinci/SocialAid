-- +goose Up
CREATE TABLE households (
    id SERIAL PRIMARY KEY,
    program_id INT REFERENCES programs(id) ON DELETE CASCADE NOT NULL,
    geolocation_id INT REFERENCES geolocations(id) ON DELETE CASCADE NOT NULL,
    name VARCHAR(255) NOT NULL,
    created_at TIMESTAMP(0) WITH TIME ZONE NOT NULL DEFAULT NOW()
);

-- Index on program_id  
CREATE INDEX idx_households_program_id ON households(program_id);
-- Index on geolocation_id
CREATE INDEX idx_households_geolocation_id ON households(geolocation_id);

-- +goose Down
DROP TABLE IF EXISTS households;