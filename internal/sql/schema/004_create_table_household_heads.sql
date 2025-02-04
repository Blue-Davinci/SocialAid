-- +goose Up
-- Create household_heads table
CREATE TABLE household_heads (
    id SERIAL PRIMARY KEY,
    household_id INT UNIQUE REFERENCES households(id) ON DELETE CASCADE NOT NULL, -- we will use UNIQUE as we will have only one head per household
    name VARCHAR(255) NOT NULL,
    national_id VARCHAR(50) NOT NULL, -- we will not use UNIQUE as we may have multiple heads with the same national id
    phone_number TEXT NOT NULL, -- we will store as a base64 encoded string
    age INT NOT NULL,
    created_at TIMESTAMP(0) WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP(0) WITH TIME ZONE NOT NULL DEFAULT NOW()
);

-- Index on household_id
CREATE INDEX idx_household_heads_household_id ON household_heads(household_id);
-- Index on national_id
CREATE INDEX idx_household_heads_national_id ON household_heads(national_id);

-- +goose Down
DROP TABLE IF EXISTS household_heads;