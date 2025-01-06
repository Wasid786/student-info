package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Wasid786/student-info/internal/config"
	"github.com/Wasid786/student-info/internal/http/handler/student"
	"github.com/Wasid786/student-info/internal/storage/sqlite"
)

func main() {
	cfg := config.MustLoad()
	// for db

	storage, err := sqlite.New(*cfg)
	if err != nil {
		log.Fatal(err)
	}
	slog.Info("Storage initialized", slog.String("env", cfg.Env), slog.String("version", "1.0.0"))

	router := http.NewServeMux()
	router.HandleFunc("POST /api/students", student.New(storage))
	router.HandleFunc("GET /api/students/{id}", student.GetById(storage))

	router.HandleFunc("GET /api/students/list", student.GetList(storage))

	server := http.Server{
		Addr:    cfg.HTTPSERVER.Addr,
		Handler: router,
	}

	slog.Info("Server started", slog.String("address", cfg.HTTPSERVER.Addr))
	fmt.Printf("Server Started %s", cfg.HTTPSERVER.Addr)

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {

		err := server.ListenAndServe()
		if err != nil {
			log.Fatal("Failed to start server")
		}

	}()
	<-done

	slog.Info("Shutting down the Server")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err = server.Shutdown(ctx)

	if err != nil {
		slog.Error("Failed to shutdown Server", slog.String("error", err.Error()))
	}
	slog.Info("Server Shutdown Successfully")
}
