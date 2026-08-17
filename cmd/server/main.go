package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"todo-api/internal/config"
	"todo-api/internal/handler"
	"todo-api/internal/repository"
	"todo-api/internal/router"
	"todo-api/internal/service"
	"todo-api/migrations"
)

func main() {
	cfg := config.Load()

	db, err := repository.Open(cfg.DatabasePath)
	if err != nil {
		log.Fatalf("open database: %v", err)
	}
	defer db.Close()

	if err := migrations.Run(db); err != nil {
		log.Fatalf("run migrations: %v", err)
	}

	repo := repository.NewSQLiteRepository(db)
	projectService := service.NewProjectService(repo)
	taskService := service.NewTaskService(repo, repo)
	projectHandler := handler.NewProjectHandler(projectService)
	taskHandler := handler.NewTaskHandler(taskService)

	server := &http.Server{
		Addr:    cfg.Addr(),
		Handler: router.New(projectHandler, taskHandler),
	}

	go func() {
		log.Printf("todo-api listening on %s", cfg.Addr())
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("listen: %v", err)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		log.Printf("shutdown: %v", err)
	}
}
