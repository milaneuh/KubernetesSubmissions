package main

import (
	"log"
	"log/slog"
	"net/http"
	"os"
	"todo_app/assets"
)

type application struct {
	logger *slog.Logger
	html   *htmlRenderer
}

func main() {
	port := os.Getenv("PORT")

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	htmlRenderer, err := newHTMLRenderer(assets.HTMLFiles, "base.tmpl")
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	app := &application{
		logger: logger,
		html:   htmlRenderer,
	}

	fileserver := http.FileServerFS(assets.StaticFiles)

	mux := http.NewServeMux()

	// Hypermedia endpoints
	mux.Handle("GET /static/", http.StripPrefix("/static", fileserver))
	mux.HandleFunc("GET /{$}", app.home)

	// JSON endpoints
	mux.HandleFunc("GET /api/healthz", app.healthz)

	log.Println("Server started in port " + port)
	err = http.ListenAndServe(":"+port, mux)
	if err != nil {
		log.Fatal("Could not start the server\n", err)
	}

}
