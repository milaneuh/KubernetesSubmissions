-- name: CreateTodo :exec
INSERT INTO todos (content) VALUES ($1);

-- name: ListTodos :many
SELECT id, content, done FROM todos ORDER BY id;

-- name: MarkTodoDone :execrows
UPDATE todos SET done = true WHERE id = $1;
