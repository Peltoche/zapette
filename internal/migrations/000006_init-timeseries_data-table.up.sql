CREATE TABLE IF NOT EXISTS timeseries_data (
  "ts_id" TEXT NOT NULL,
  "time_unix_sec" INTEGER NOT NULL,
  "content" BLOB NOT NULL,
  FOREIGN KEY(ts_id) REFERENCES timeseries_config(ts_id) ON UPDATE RESTRICT ON DELETE RESTRICT
) STRICT;

CREATE UNIQUE INDEX IF NOT EXISTS idx_timeseries_tsid_time ON timeseries_data(ts_id, time_unix_sec);
