CREATE TABLE IF NOT EXISTS sysstats (
  "time" INTEGER NOT NULL,
  "content" BLOB NOT NULL
) STRICT;

CREATE UNIQUE INDEX IF NOT EXISTS idx_sysstats_time ON sysstats(time);
