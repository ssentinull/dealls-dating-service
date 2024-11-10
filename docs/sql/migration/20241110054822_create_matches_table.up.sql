-- +migrate Up
CREATE TABLE IF NOT EXISTS "matches" (
    "id" BIGSERIAL PRIMARY KEY,
    "my_user_id" BIGINT,
    "matched_user_id" BIGINT,
    "created_at" TIMESTAMP NOT NULL DEFAULT NOW(),
    FOREIGN KEY (my_user_id) REFERENCES users(id),
    FOREIGN KEY (matched_user_id) REFERENCES users(id)
);

CREATE INDEX IF NOT EXISTS matches_my_user_id_idx ON matches(my_user_id);
CREATE INDEX IF NOT EXISTS matches_matched_user_id_idx ON matches(matched_user_id);