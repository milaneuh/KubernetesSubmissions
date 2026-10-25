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
}

func main() {
	port := os.Getenv("PORT")

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	go ImageLookupLoop()

	app := &application{
		logger: logger,
	}

	fileserver := http.FileServerFS(assets.StaticFiles)

	mux := http.NewServeMux()

	// Hypermedia endpoints
	mux.Handle("GET /static/", http.StripPrefix("/static", fileserver))
	mux.HandleFunc("GET /{$}", app.home)
	mux.HandleFunc("GET /image", app.image)
	mux.HandleFunc("POST /new-todo", app.newTodo)

	// JSON endpoints
	mux.HandleFunc("GET /api/healthz", app.healthz)

	log.Println("Server started in port " + port)
	err := http.ListenAndServe(":"+port, mux)
	if err != nil {
		log.Fatal("Could not start the server\n", err)
	}

}

func ImageLookupLoop() {
	for {
		if err := ImageStore(); err != nil {
			slog.Error("image refresh failed, keeping the current image", "err", err)
		}
		time.Sleep(30 * time.Second)
	}
}

func imagePath() string {
	if p := os.Getenv("IMAGE_PATH"); p != "" {
		return p
	}
	return "/usr/src/app/files/image.jpg"
}

func ImageStore() error {
	path := imagePath()
	if !IsFilePresent(path) {
		slog.Info("no image stored yet, pulling one", "path", path)
		return PullImage(path)
	}
	if IsFileOlderThan(path, 10*time.Minute) {
		slog.Info("image is older than 10 minutes, pulling a new one", "path", path)
		return PullImage(path)
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

func PullImage(path string) error {
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return err
	}

	start := time.Now()
	resp, err := http.Get("https://picsum.photos/1200")
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("picsum returned %s", resp.Status)
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
