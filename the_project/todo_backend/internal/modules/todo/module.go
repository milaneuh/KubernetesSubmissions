package todo

import (
	"context"
	"errors"

	"todo_backend/internal/db"
)

const maxTodoLength = 140

var ErrInvalidTodo = errors.New("todo must be between 1 and 140 characters")

type Module struct {
	queries *db.Queries
}

func New(queries *db.Queries) *Module {
	return &Module{queries: queries}
}

func (m *Module) CreateTodo(ctx context.Context, todo string) error {
	if todo == "" || len(todo) > maxTodoLength {
		return ErrInvalidTodo
	}

	return m.queries.CreateTodo(ctx, todo)
}

func (m *Module) ListTodos(ctx context.Context) ([]string, error) {
	return m.queries.ListTodos(ctx)
}
