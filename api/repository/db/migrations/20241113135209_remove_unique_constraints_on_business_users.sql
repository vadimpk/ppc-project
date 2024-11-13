-- +goose Up
-- +goose StatementBegin
-- First drop existing constraints
ALTER TABLE users
    DROP CONSTRAINT IF EXISTS users_business_id_email_key,
    DROP CONSTRAINT IF EXISTS users_business_id_phone_key;

-- Add new unique constraint on email
ALTER TABLE users
    ADD CONSTRAINT users_email_key UNIQUE (email);
ALTER TABLE users
    ADD CONSTRAINT users_phone_key UNIQUE (phone);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
-- Remove unique constraint on email
ALTER TABLE users
    DROP CONSTRAINT IF EXISTS users_email_key;
ALTER TABLE users
    DROP CONSTRAINT IF EXISTS users_phone_key;

-- Restore original constraints
ALTER TABLE users
    ADD CONSTRAINT users_business_id_email_key UNIQUE (business_id, email),
    ADD CONSTRAINT users_business_id_phone_key UNIQUE (business_id, phone);
-- +goose StatementEnd