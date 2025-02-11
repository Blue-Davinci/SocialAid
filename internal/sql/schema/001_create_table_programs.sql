-- +goose Up
CREATE TABLE programs (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) UNIQUE NOT NULL,
    category VARCHAR(255) NOT NULL,
    description TEXT NOT NULL,
    created_at TIMESTAMP(0) WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP(0) WITH TIME ZONE NOT NULL DEFAULT NOW()
);
-- Index on name
CREATE INDEX idx_programs_name ON programs(name);   

-- +goose Down
DROP TABLE IF EXISTS programs;