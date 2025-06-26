package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"taskservice/internal/app"

	"github.com/go-chi/chi/v5"
)

func main() {
	var cfg app.Config

	// Записываем значения с конфига port и env
	flag.IntVar(&cfg.Port, "port", 8080, "API server port")
	flag.StringVar(&cfg.Env, "env", "development", "Environment (development|staging|production)")
	flag.Parse()

	// Инициализируем новый логгер, с датой и временем
	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)

	// Объявляем экземпляр структуры приложения, которая содержит структуру
	application := &app.Application{
		Config: cfg,
		Logger: logger,
	}

	r := chi.NewRouter()
	r.Post("/tasks", application.HandleCreateTask)
	r.Get("/", application.HandleMain)

	srv := http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.Port),
		Handler:      r,
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	logger.Printf("starting %s server on %s", cfg.Env, srv.Addr)
	err := srv.ListenAndServe()
	logger.Fatal(err)
}
