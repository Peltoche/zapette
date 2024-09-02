CREATE TABLE IF NOT EXISTS timeseries_config (
  "ts_id" TEXT NOT NULL,
  "graph_span_s" INT NOT NULL,
  "tick_span_s" INT NOT NULL
) STRICT;

CREATE UNIQUE INDEX IF NOT EXISTS idx_timeseries_config_tsid ON timeseries_config(ts_id);
