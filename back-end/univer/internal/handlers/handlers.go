package handlers

import (
	"fmt"
	"net/http"
	"univer/internal/config"

	"github.com/jackc/pgx/v5/pgxpool"
)

// Handler — структура для обработчиков
type Handler struct {
	DB     *pgxpool.Pool
	Config *config.Config
}

// NewHandler создаёт новый объект Handler
func NewHandler(db *pgxpool.Pool, cfg *config.Config) *Handler {
	return &Handler{DB: db, Config: cfg}
}

// HomeHandler — пример обработчика с доступом к конфигу
func (h *Handler) HomeHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Сервер работает на порту: %d\n", h.Config.Server.Port)
}

// GetAllTeachersHandler — обработчик с доступом к конфигу и БД
//func (h *Handler) GetAllTeachersHandler(w http.ResponseWriter, r *http.Request) {
//	// Доступ к конфигу
//	fmt.Fprintf(w, "Подключено к базе данных: %s\n", h.Config.Database.DSN)
//}
