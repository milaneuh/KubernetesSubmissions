CREATE TABLE IF NOT EXISTS todos (
  id         bigserial PRIMARY KEY,
  content    text NOT NULL,
  created_at timestamptz NOT NULL DEFAULT now()
);
