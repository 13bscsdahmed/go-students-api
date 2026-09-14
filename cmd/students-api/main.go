package main

import (
	"context"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"students-api/internal/config"
	"students-api/internal/http/handlers/student"
	"students-api/internal/storage/sqlite"
	"syscall"
	"time"
)

func main() {
	cfg := config.MustLoad()

	storage, err := sqlite.New(cfg)
	if err != nil {
		slog.Error("Failed to create storage", slog.String("error", err.Error()))
	}
	slog.Info("Storage Initialized", slog.String("env", cfg.Env), slog.String("version", "1.0.0"))
	router := http.NewServeMux()

	router.HandleFunc("POST /api/students", student.New(storage))
	server := http.Server{
		Addr:    cfg.HTTPServer.Address,
		Handler: router,
	}
	slog.Info("Server Started", slog.String("address", cfg.HTTPServer.Address))
	done := make(chan os.Signal, 1)

	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		err := server.ListenAndServe()
		if err != nil {
			log.Fatal("Failed to start server: ")
		}
	}()
	<-done

	slog.Info("Shutting down the server")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	error := server.Shutdown(ctx)
	if error != nil {
		slog.Error("Failed to shutdown the server", slog.String("error", error.Error()))
	}
	slog.Info("Server shutdown successfully")
}
