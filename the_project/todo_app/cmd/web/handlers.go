package main

import (
	"net/http"
	"todo_app/views"
)

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	err := views.Home(imagePath()).Render(r.Context(), w)
	if err != nil {
		app.logger.Error(err.Error())
		http.Error(w, http.StatusText(500), 500)
	}
}

func (app *application) image(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Cache-Control", "no-cache")
	http.ServeFile(w, r, imagePath())
}

func (app *application) newTodo(w http.ResponseWriter, r *http.Request) {
	todo := r.FormValue("todo-input")

	if todo == "" || len(todo) > 140 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err := views.Todo(todo).Render(r.Context(), w)

	if err != nil {
		app.logger.Error(err.Error())
		http.Error(w, http.StatusText(500), 500)
	}
}

func (app *application) healthz(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK\n"))
}
