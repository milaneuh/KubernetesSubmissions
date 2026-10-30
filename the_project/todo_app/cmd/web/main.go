package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
	"time"
	"todo_app/assets"
)

type application struct {
	logger *slog.Logger
	cfg    config
}

func main() {
	cfg := loadConfig()

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	app := &application{
		logger: logger,
		cfg:    cfg,
	}

	go app.ImageLookupLoop()

	fileserver := http.FileServerFS(assets.StaticFiles)

	mux := http.NewServeMux()

	// Hypermedia endpoints
	mux.Handle("GET /static/", http.StripPrefix("/static", fileserver))
	mux.HandleFunc("GET /{$}", app.home)
	mux.HandleFunc("GET /image", app.image)
	mux.HandleFunc("POST /new-todo", app.newTodo)

	log.Println("Server started in port " + cfg.port)
	err := http.ListenAndServe(":"+cfg.port, mux)
	if err != nil {
		log.Fatal("Could not start the server\n", err)
	}

}

func (app *application) ImageLookupLoop() {
	for {
		if err := app.ImageStore(); err != nil {
			slog.Error("image refresh failed, keeping the current image", "err", err)
		}
		time.Sleep(app.cfg.imageRefreshInterval)
	}
}

func (app *application) ImageStore() error {
	path := app.cfg.imagePath
	if !IsFilePresent(path) {
		slog.Info("no image stored yet, pulling one", "path", path)
		return PullImage(app.cfg.imageURL, path)
	}
	if IsFileOlderThan(path, app.cfg.imageMaxAge) {
		slog.Info("image is too old, pulling a new one", "path", path, "max_age", app.cfg.imageMaxAge)
		return PullImage(app.cfg.imageURL, path)
	}
	slog.Debug("image is still fresh", "path", path)
	return nil
}

func IsFilePresent(path string) bool {
	_, err := os.Stat(path)
	return !errors.Is(err, os.ErrNotExist)
}

func IsFileOlderThan(path string, d time.Duration) bool {
	stat, err := os.Stat(path)
	if err != nil {
		return true
	}
	return time.Since(stat.ModTime()) > d
}

func PullImage(url, path string) error {
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return err
	}

	start := time.Now()
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("image source %s returned %s", url, resp.Status)
	}

	out, err := os.CreateTemp(filepath.Dir(path), "image-*.tmp")
	if err != nil {
		return err
	}
	defer os.Remove(out.Name())

	n, err := io.Copy(out, resp.Body)
	if err != nil {
		out.Close()
		return err
	}
	if err := out.Close(); err != nil {
		return err
	}

	if err := os.Rename(out.Name(), path); err != nil {
		return err
	}

	slog.Info("image saved", "path", path, "bytes", n, "duration", time.Since(start).Round(time.Millisecond))
	return nil
}
