package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"time"
	"todo_app/views"
)

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	todos, err := app.fetchTodos(r.Context())
	if err != nil {
		app.logger.Error("could not fetch the todos from backend", "err", err.Error())
	}

	err = views.Home(app.cfg.imagePath, todos).Render(r.Context(), w)
	if err != nil {
		app.logger.Error("could not load the view", "err", err.Error())
		http.Error(w, http.StatusText(500), 500)
	}
}

func (app *application) image(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Cache-Control", "no-cache")
	http.ServeFile(w, r, app.cfg.imagePath)
}

var backendClient = &http.Client{Timeout: 2 * time.Second}

func (app *application) fetchTodos(ctx context.Context) ([]views.TodoItem, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, app.cfg.todoBackendURL, nil)
	if err != nil {
		app.logger.Error("could not create the http request", "err", err.Error())
		return nil, err
	}

	resp, err := backendClient.Do(req)
	if err != nil {
		app.logger.Error("could not request the backend", "err", err.Error())
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("todo backend could not list the todos: %s", resp.Status)
	}

	var todos []views.TodoItem
	if err := json.NewDecoder(resp.Body).Decode(&todos); err != nil {
		app.logger.Error("could not decode the json response", "err", err.Error())
		return nil, err
	}
	return todos, nil
}

func (app *application) newTodo(w http.ResponseWriter, r *http.Request) {
	todo := r.FormValue("todo-input")

	if todo == "" || len(todo) > 140 {
		app.logger.Error("could not create the todo", "err", errors.New("malformed todo, the text is nil or longer than 140 characters"))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	resp, err := backendClient.PostForm(app.cfg.todoBackendURL, url.Values{"todo-input": {todo}})
	if err != nil {
		app.logger.Error("could not post the todo", "err", err.Error())
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
		app.logger.Error("could not fecth the todos from backend", "err", err.Error())
		http.Error(w, "could not load the todos", http.StatusBadGateway)
		return
	}

	if err := views.Todo(todos).Render(r.Context(), w); err != nil {
		app.logger.Error("could not render the todos", "err", err.Error())
	}
}

func (app *application) markTodoDone(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		http.Error(w, "invalid todo id", http.StatusBadRequest)
		return
	}

	req, err := http.NewRequestWithContext(r.Context(), http.MethodPut, fmt.Sprintf("%s/%d", app.cfg.todoBackendURL, id), nil)
	if err != nil {
		app.logger.Error("could not create the http request", "err", err.Error())
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	resp, err := backendClient.Do(req)
	if err != nil {
		app.logger.Error("could not mark the todo as done", "err", err.Error())
		http.Error(w, "could not reach the todo backend", http.StatusBadGateway)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		app.logger.Error("todo backend refused to mark the todo as done", "id", id, "status", resp.Status)
		http.Error(w, "could not update the todo", http.StatusBadGateway)
		return
	}

	todos, err := app.fetchTodos(r.Context())
	if err != nil {
		app.logger.Error("could not fetch the todos from backend", "err", err.Error())
		http.Error(w, "could not load the todos", http.StatusBadGateway)
		return
	}

	if err := views.Todo(todos).Render(r.Context(), w); err != nil {
		app.logger.Error("could not render the todos", "err", err.Error())
	}
}

func (app *application) breakApp(w http.ResponseWriter, r *http.Request) {
	app.broken.Store(true)
	app.logger.Warn("the app has been broken on purpose")
	w.Header().Set("HX-Refresh", "true")
	w.WriteHeader(http.StatusOK)
}

func (app *application) healthz(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK\n"))
}

func (app *application) failWhenBroken(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if app.broken.Load() {
			http.Error(w, "the app is broken", http.StatusInternalServerError)
			return
		}
		next.ServeHTTP(w, r)
	})
}
