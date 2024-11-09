-- +migrate Up
CREATE TABLE IF NOT EXISTS "books" (
  "id" BIGINT PRIMARY KEY,
  "title" TEXT NOT NULL,
  "author" TEXT NOT NULL DEFAULT ''
);
