-- +goose Up
CREATE TABLE geolocations (
    id SERIAL PRIMARY KEY,
    county VARCHAR(255) UNIQUE NOT NULL,
    sub_county VARCHAR(255) NOT NULL,
    location VARCHAR(255) NOT NULL,
    sub_location VARCHAR(255) NOT NULL,
    created_at TIMESTAMP(0) WITH TIME ZONE NOT NULL DEFAULT NOW()
);

-- Index on county
CREATE INDEX idx_geolocations_county ON geolocations(county);
-- Index on sub_county
CREATE INDEX idx_geolocations_sub_county ON geolocations(sub_county);


-- +goose Down
DROP TABLE IF EXISTS geolocations;