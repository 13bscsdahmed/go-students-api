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
	"syscall"
	"time"
)

func main() {
	cfg := config.MustLoad()
	router := http.NewServeMux()

	router.HandleFunc("POST /api/students", student.New())
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
