package main

import (
	"context"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Harshith121212/crud-project/internal/config"
	"github.com/Harshith121212/crud-project/internal/http/handlers/student"
)

func main() {
	//configuration loading.
	cfg := config.MustLoad()
	//setting up database
	//router setup
	router := http.NewServeMux()

	router.HandleFunc("POST /api/students", student.New())
	//server setup
	server := http.Server{
		Addr:    cfg.Addr,
		Handler: router,
	}
	slog.Info("Server is running on port", slog.String("address", cfg.Addr))
	//fmt.Printf("Server is running on port %s", cfg.HTTPServer.Addr)

	done := make(chan os.Signal, 1)

	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		err := server.ListenAndServe()
		if err != nil {
			log.Fatal("failed to start server")
		}
	}()

	<-done

	slog.Info("shutting down server")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := server.Shutdown(ctx)
	if err != nil {
		slog.Error("Failed to shut down server", slog.String("error", err.Error()))
	}

	slog.Info("server shut down successfully")

}
