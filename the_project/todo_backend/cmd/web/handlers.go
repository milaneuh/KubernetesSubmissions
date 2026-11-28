package main

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"strconv"

	"todo_backend/internal/modules/todo"
)

func (app *application) newTodo(w http.ResponseWriter, r *http.Request) {
	todoInput := r.FormValue("todo-input")

	err := app.todos.CreateTodo(r.Context(), todoInput)
	if errors.Is(err, todo.ErrInvalidTodo) {
		slog.Error("could not create the todo", "err", err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err != nil {
		slog.Error("could not create the todo", "err", err.Error())
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	slog.Info("todo created", "todo", todoInput)
	w.WriteHeader(http.StatusCreated)
}

func (app *application) listTodos(w http.ResponseWriter, r *http.Request) {
	todos, err := app.todos.ListTodos(r.Context())
	if err != nil {
		app.logger.Error("could not list the todos", "err", err.Error())
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(todos); err != nil {
		app.logger.Error("could not encode the todos", "err", err.Error())
	}
}

func (app *application) markTodoDone(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		http.Error(w, "invalid todo id", http.StatusBadRequest)
		return
	}

	err = app.todos.MarkTodoDone(r.Context(), id)
	if errors.Is(err, todo.ErrTodoNotFound) {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	if err != nil {
		app.logger.Error("could not mark the todo as done", "id", id, "err", err.Error())
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	slog.Info("todo marked as done", "id", id)
	w.WriteHeader(http.StatusNoContent)
}

func (app *application) healthz(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK\n"))
}
