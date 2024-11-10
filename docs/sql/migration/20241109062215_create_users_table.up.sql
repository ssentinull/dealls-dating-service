-- +migrate Up
CREATE TABLE IF NOT EXISTS "users" (
    "id" BIGSERIAL PRIMARY KEY,
    "email" VARCHAR(255) NOT NULL UNIQUE,
    "password" TEXT NOT NULL,
    "name" VARCHAR(255) NOT NULL,
    "gender" VARCHAR(50),
    "birth_date" DATE,
    "location" VARCHAR(255),
    "created_at" TIMESTAMP NOT NULL DEFAULT NOW(),
    "updated_at" TIMESTAMP,
    "deleted_at" TIMESTAMP
);

CREATE INDEX IF NOT EXISTS users_email_idx ON users(email);
CREATE INDEX IF NOT EXISTS users_birth_date_idx ON users(birth_date);
CREATE INDEX IF NOT EXISTS users_deleted_at_idx ON users(deleted_at);