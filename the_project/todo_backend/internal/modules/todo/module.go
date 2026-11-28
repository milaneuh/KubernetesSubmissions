package todo

import (
	"context"
	"errors"

	"todo_backend/internal/db"
)

const maxTodoLength = 140

var (
	ErrInvalidTodo  = errors.New("todo must be between 1 and 140 characters")
	ErrTodoNotFound = errors.New("todo not found")
)

type Todo struct {
	ID      int64  `json:"id"`
	Content string `json:"content"`
	Done    bool   `json:"done"`
}

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

func (m *Module) ListTodos(ctx context.Context) ([]Todo, error) {
	rows, err := m.queries.ListTodos(ctx)
	if err != nil {
		return nil, err
	}

	todos := make([]Todo, 0, len(rows))
	for _, row := range rows {
		todos = append(todos, Todo{ID: row.ID, Content: row.Content, Done: row.Done})
	}
	return todos, nil
}

func (m *Module) MarkTodoDone(ctx context.Context, id int64) error {
	n, err := m.queries.MarkTodoDone(ctx, id)
	if err != nil {
		return err
	}
	if n == 0 {
		return ErrTodoNotFound
	}
	return nil
}
