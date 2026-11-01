-- name: CreateTodo :exec
INSERT INTO todos (content) VALUES ($1);

-- name: ListTodos :many
SELECT content FROM todos ORDER BY id;
