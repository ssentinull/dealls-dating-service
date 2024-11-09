-- +migrate Up
CREATE TABLE IF NOT EXISTS "preferences" (
    "id" BIGSERIAL PRIMARY KEY,
    "gender" VARCHAR(50),
    "min_age" SMALLINT,
    "max_age" SMALLINT,
    "location" VARCHAR(255),
    "user_id" BIGINT,
    "created_at" TIMESTAMP NOT NULL DEFAULT NOW(),
    "updated_at" TIMESTAMP,
    "deleted_at" TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id)
);

CREATE INDEX preferences_user_id_idx ON preferences(user_id);