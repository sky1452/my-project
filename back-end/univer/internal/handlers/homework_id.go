package handlers

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5"
)

type Teacher struct {
	ID     int    `json:"id"`
	Name   string `json:"name"`
	Avatar []byte `json:"avatar"`
}

type HomeworkResponse struct {
	ID           int      `json:"id"`
	Title        string   `json:"title"`
	Description  string   `json:"description"`
	MaxScore     int      `json:"max_score"`
	Deadline     string   `json:"deadline"`
	DisciplineID int      `json:"discipline_id"`
	TeacherID    int      `json:"teacher_id"`
	CreatedAt    string   `json:"created_at"`
	UpdatedAt    string   `json:"updated_at"`
	Teacher      Teacher  `json:"teacher"`
}

func (h *Handler) GetHomeworkByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	homeworkID, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "invalid homework id", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var (
		resp HomeworkResponse

		teacherID    sql.NullInt64
		teacherName  sql.NullString
		teacherAvatar []byte

		deadline   time.Time
		createdAt  time.Time
		updatedAt  time.Time
	)

	err = h.DB.QueryRow(ctx, `
		SELECT 
			h.id,
			h.title,
			h.description,
			h.max_score,
			h.deadline,
			h.discipline_id,
			h.teacher_id,
			h.created_at,
			h.updated_at,
			u.user_id,
			u.name,
			u.avatar
		FROM homeworks h
		LEFT JOIN users u ON u.user_id = h.teacher_id
		WHERE h.id = $1
	`, homeworkID).Scan(
		&resp.ID,
		&resp.Title,
		&resp.Description,
		&resp.MaxScore,
		&deadline,
		&resp.DisciplineID,
		&resp.TeacherID,
		&createdAt,
		&updatedAt,
		&teacherID,
		&teacherName,
		&teacherAvatar,
	)

	if err != nil {
		log.Println("Ошибка получения задания:", err)

		w.Header().Set("Content-Type", "application/json")

		if errors.Is(err, pgx.ErrNoRows) {
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(map[string]string{
				"error": "Задание не найдено",
			})
			return
		}

		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Ошибка получения задания",
		})
		return
	}

	// форматирование дат
	resp.Deadline = deadline.Format("2006-01-02 15:04:05")
	resp.CreatedAt = createdAt.Format("2006-01-02 15:04:05")
	resp.UpdatedAt = updatedAt.Format("2006-01-02 15:04:05")

	// teacher
	if teacherID.Valid {
		resp.Teacher.ID = int(teacherID.Int64)
	}

	if teacherName.Valid {
		resp.Teacher.Name = teacherName.String
	}

	resp.Teacher.Avatar = teacherAvatar

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}