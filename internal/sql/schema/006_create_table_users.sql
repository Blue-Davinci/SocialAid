-- +goose Up
CREATE EXTENSION IF NOT EXISTS citext;

CREATE TABLE users(
    id SERIAL PRIMARY KEY,
    email CITEXT UNIQUE NOT NULL,
    name VARCHAR(255) NOT NULL,
    api_key BYTEA NOT NULL,
    created_at TIMESTAMP(0) WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP(0) WITH TIME ZONE NOT NULL DEFAULT NOW()
);
-- Index on email
CREATE INDEX idx_users_email ON users(email);

-- create a new registration with the following plaintext api key = D5HDKGF5I5BK3J2FAAHBD37M7M
insert into users (email, name, api_key) VALUES('test@gmail.com', 'test name', decode('gglNdEaLzPEZfyUuCTPrb
MXlsMb7yhThZZi6t5CaHv0=','base64'));


-- +goose Down
DROP TABLE IF EXISTS users; 