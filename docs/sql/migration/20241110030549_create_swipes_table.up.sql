-- +migrate Up
CREATE TABLE IF NOT EXISTS "swipes" (
    "id" BIGSERIAL PRIMARY KEY,
    "from_user_id" BIGINT,
    "to_user_id" BIGINT,
    "swipe_type" VARCHAR(50),
    "created_at" TIMESTAMP NOT NULL DEFAULT NOW(),
    FOREIGN KEY (from_user_id) REFERENCES users(id),
    FOREIGN KEY (to_user_id) REFERENCES users(id)
);

CREATE INDEX IF NOT EXISTS swipes_from_user_id_idx ON swipes(from_user_id);
CREATE INDEX IF NOT EXISTS swipes_to_user_id_idx ON swipes(to_user_id);
CREATE INDEX IF NOT EXISTS swipes_created_at_idx ON swipes(created_at);