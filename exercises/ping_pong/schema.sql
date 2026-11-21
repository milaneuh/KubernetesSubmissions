-- Run by the app at startup, so it must be safe to run more than once.
CREATE TABLE IF NOT EXISTS ping_counter (
  id   int PRIMARY KEY,
  hits bigint NOT NULL DEFAULT 0
);
