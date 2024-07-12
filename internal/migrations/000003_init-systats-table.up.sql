CREATE TABLE IF NOT EXISTS systats (
  "time" TEXT NOT NULL,
  "total_mem" INTEGER NOT NULL,
  "available_mem" INTEGER NOT NULL,
  "free_mem" INTEGER NOT NULL,
  "buffers" INTEGER NOT NULL,
  "cached" INTEGER NOT NULL,
  "s_reclaimable" INTEGER NOT NULL,
  "sh_mem" INTEGER NOT NULL,
  "total_swap" INTEGER NOT NULL,
  "free_swap" INTEGER NOT NULL
) STRICT;

CREATE UNIQUE INDEX IF NOT EXISTS idx_systats_time ON systats(time);
