package handlers

import (
	"fmt"
	"net/http"
	"univer/internal/config"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/aws/aws-sdk-go/service/s3"
)

// Handler — структура для обработчиков
type Handler struct {
	DB     *pgxpool.Pool
	Config *config.Config
	S3 *s3.S3
}

// NewHandler создаёт новый объект Handler
func NewHandler(db *pgxpool.Pool, cfg *config.Config, s3Client *s3.S3) *Handler {
	return &Handler{
		DB:     db,
		Config: cfg,
		S3:     s3Client,
	}
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
