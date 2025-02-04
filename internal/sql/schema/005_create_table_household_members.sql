-- +goose Up
-- Create household_members table
CREATE TABLE household_members (
    id SERIAL PRIMARY KEY,
    household_id INT REFERENCES households(id) ON DELETE CASCADE NOT NULL,
    name VARCHAR(255) NOT NULL,
    age INT NOT NULL,
    relation VARCHAR(50) NOT NULL,
    created_at TIMESTAMP(0) WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP(0) WITH TIME ZONE NOT NULL DEFAULT NOW(),
    CONSTRAINT unique_household_member UNIQUE (household_id, name)
);

-- Index on household_id
CREATE INDEX idx_household_members_household_id ON household_members(household_id);
-- Index on name
CREATE INDEX idx_household_members_name ON household_members(name);
-- Index on relation
CREATE INDEX idx_household_members_relation ON household_members(relation);

-- +goose Down
DROP TABLE IF EXISTS household_members;