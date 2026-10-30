package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"
	"todo_app/views"
)

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	todos, err := app.fetchTodos(r.Context())
	if err != nil {
		app.logger.Error(err.Error())
	}

	err = views.Home(app.cfg.imagePath, todos).Render(r.Context(), w)
	if err != nil {
		app.logger.Error(err.Error())
		http.Error(w, http.StatusText(500), 500)
	}
}

func (app *application) image(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Cache-Control", "no-cache")
	http.ServeFile(w, r, app.cfg.imagePath)
}

var backendClient = &http.Client{Timeout: 2 * time.Second}

func (app *application) fetchTodos(ctx context.Context) ([]string, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, app.cfg.todoBackendURL, nil)
	if err != nil {
		return nil, err
	}

	resp, err := backendClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("todo backend could not list the todos: %s", resp.Status)
	}

	var todos []string
	if err := json.NewDecoder(resp.Body).Decode(&todos); err != nil {
		return nil, err
	}
	return todos, nil
}

func (app *application) newTodo(w http.ResponseWriter, r *http.Request) {
	todo := r.FormValue("todo-input")

	if todo == "" || len(todo) > 140 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	resp, err := backendClient.PostForm(app.cfg.todoBackendURL, url.Values{"todo-input": {todo}})
	if err != nil {
		app.logger.Error(err.Error())
		http.Error(w, "could not reach the todo backend", http.StatusBadGateway)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		app.logger.Error("todo backend refused the todo", "status", resp.Status)
		http.Error(w, "could not save the todo", http.StatusBadGateway)
		return
	}

	todos, err := app.fetchTodos(r.Context())
	if err != nil {
		app.logger.Error(err.Error())
		http.Error(w, "could not load the todos", http.StatusBadGateway)
		return
	}

	if err := views.Todo(todos).Render(r.Context(), w); err != nil {
		app.logger.Error(err.Error())
	}
}
