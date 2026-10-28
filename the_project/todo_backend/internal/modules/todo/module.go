package todo

import (
	"context"
	"errors"
	"slices"
	"sync"
)

const maxTodoLength = 140

var ErrInvalidTodo = errors.New("todo must be between 1 and 140 characters")

type Module struct {
	mu    sync.Mutex
	todos []string
}

func New() *Module {
	return &Module{}
}

func (m *Module) CreateTodo(ctx context.Context, todo string) error {
	if todo == "" || len(todo) > maxTodoLength {
		return ErrInvalidTodo
	}

	m.mu.Lock()
	defer m.mu.Unlock()
	m.todos = append(m.todos, todo)
	return nil
}

func (m *Module) ListTodos(ctx context.Context) ([]string, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	return slices.Clone(m.todos), nil
}
