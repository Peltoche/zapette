CREATE TABLE IF NOT EXISTS systats (
  "time" TEXT NOT NULL,
  "content" BLOB NOT NULL
) STRICT;

CREATE UNIQUE INDEX IF NOT EXISTS idx_systats_time ON systats(time);
