package main

import (
	"encoding/json"
	"errors"
	"net/http"

	"todo_backend/internal/modules/todo"
)

func (app *application) newTodo(w http.ResponseWriter, r *http.Request) {
	todoInput := r.FormValue("todo-input")

	err := app.todos.CreateTodo(r.Context(), todoInput)
	if errors.Is(err, todo.ErrInvalidTodo) {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err != nil {
		app.logger.Error(err.Error())
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (app *application) listTodos(w http.ResponseWriter, r *http.Request) {
	todos, err := app.todos.ListTodos(r.Context())
	if err != nil {
		app.logger.Error(err.Error())
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	if todos == nil {
		todos = []string{}
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(todos); err != nil {
		app.logger.Error(err.Error())
	}
}

func (app *application) healthz(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK\n"))
}
