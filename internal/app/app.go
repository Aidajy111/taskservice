package app

import (
	"encoding/json"
	"log"
	"net/http"
)

// Структура конфигурации
type Config struct {
	Port int
	Env  string
}

// Определим структуру приложения, которая будет содержать зависимости для
// обработчиков HTTP, вспомогательных функций и middleware.
type Application struct {
	Config Config
	Logger *log.Logger
}

// middleware для создания тасков
func (a *Application) HandleCreateTask(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "Application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "path - /tasks"})
}

// middleware для проверки сервера
func (a *Application) HandleMain(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("The server is running"))
}
