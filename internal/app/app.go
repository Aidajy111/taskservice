package app

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"
	"sync/atomic"
	"time"
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
	Logger *log.Logger // для отладки и мониторинга
}

type Task struct {
	ID          int64     `json:"id"`          // Id для задачи
	Types       string    `json:"types"`       // Типа задачи который передал пользователь или сама задача
	Status      string    `json:"status"`      // Статус задачи который будет менятся
	Created_at  time.Time `json:"created_at"`  // Время создания задачи
	Started_at  time.Time `json:"started_at"`  // Время фактического выполнения задачи
	Finished_at time.Time `json:"finished_at"` // Время завершения задачи
}

var (
	tasks = map[int]Task{
		1: {
			ID:          1,
			Types:       "image_processing",
			Status:      "completed",
			Created_at:  time.Date(2025, 6, 28, 10, 0, 0, 0, time.UTC),
			Started_at:  time.Date(2025, 6, 28, 10, 1, 0, 0, time.UTC),
			Finished_at: time.Date(2025, 6, 28, 10, 4, 30, 0, time.UTC),
		},
		2: {
			ID:         2,
			Types:      "data_export",
			Status:     "running",
			Created_at: time.Date(2025, 6, 28, 10, 5, 0, 0, time.UTC),
			Started_at: time.Date(2025, 6, 28, 10, 6, 0, 0, time.UTC),
			// Finished_at не установлено - задача ещё выполняется
		},
		3: {
			ID:    3,
			Types: "report_generation",

			Status:     "pending",
			Created_at: time.Date(2025, 6, 28, 10, 10, 0, 0, time.UTC),
			// Started_at и Finished_at не установлены - задача в очереди
		},
		4: {
			ID:          4,
			Types:       "database_backup",
			Status:      "failed",
			Created_at:  time.Date(2025, 6, 28, 10, 15, 0, 0, time.UTC),
			Started_at:  time.Date(2025, 6, 28, 10, 16, 0, 0, time.UTC),
			Finished_at: time.Date(2025, 6, 28, 10, 17, 30, 0, time.UTC),
		},
	} // Инициализация сразу с тестовыми данными
	tasksLock sync.RWMutex
)

// processTask - заканчивает процес и меняет Status на completed
func processTask(taskId int) {
	time.Sleep(3 * time.Minute) // Имитация долгой обработки (3-5 минут)

	tasksLock.Lock()
	defer tasksLock.Unlock()

	if task, ok := tasks[taskId]; ok {
		task.Status = "completed"
		task.Finished_at = time.Now().UTC()
		tasks[taskId] = task
	}
}

// startTaskProcessing - начинает процесс меняет значение pending на running
func startTaskProcessing(taskId int) {
	tasksLock.Lock()
	defer tasksLock.Unlock()

	if task, ok := tasks[taskId]; ok {
		task.Status = "running"
		task.Started_at = time.Now().UTC()
		tasks[taskId] = task
		go processTask(taskId) // Запуск фоновой обработки
	}
}

// saveTask - сохраняет новую таску в масив
func saveTask(t Task) {
	tasksLock.Lock()
	defer tasksLock.Unlock()
	tasks[int(t.ID)] = t
	return
}

var taskIDCounter int64

// generateID - генерирует id
func generateID() int64 {
	return atomic.AddInt64(&taskIDCounter, 1)
}

// middleware для создания тасков
func (a *Application) HandleCreateTask(w http.ResponseWriter, r *http.Request) {
	var requests struct {
		Types     string    `json:"types"`
		CreatedAd time.Time `json:"created_at"`
	}

	if err := json.NewDecoder(r.Body).Decode(&requests); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	task := Task{
		ID:         generateID(),
		Types:      requests.Types,
		Status:     "pending",
		Created_at: requests.CreatedAd, //Время начала и создания таски одинаковые
	}

	saveTask(task)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(task)
}

// middleware для проверки сервера
func (a *Application) HandleMain(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("The server is running"))
}

func (a *Application) HandleGetTasks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(tasks)
}
