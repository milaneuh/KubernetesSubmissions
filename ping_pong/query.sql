-- name: InitCounter :exec
INSERT INTO ping_counter (id, hits) VALUES (1, 0)
ON CONFLICT (id) DO NOTHING;

-- name: GetCounter :one
SELECT hits FROM ping_counter WHERE id = 1;

-- name: IncrementCounter :one
UPDATE ping_counter SET hits = hits + 1 WHERE id = 1
RETURNING hits;
