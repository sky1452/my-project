package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type Homework struct {
	ID          int    `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	MaxScore    int    `json:"max_score"`
	CreatedAt   string `json:"created_at"` // строка вместо time.Time
	UpdatedAt   string `json:"updated_at"`
	Deadline    string `json:"deadline"`
}

func (h *Handler) GetHomeworks(w http.ResponseWriter, r *http.Request) {
	disciplineId := r.URL.Query().Get("disciplineId")
	groupName := r.URL.Query().Get("group")
	teacherId := r.URL.Query().Get("teacherId")

	if disciplineId == "" || groupName == "" || teacherId == "" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "missing query params"})
		return
	}

	query := `
		SELECT h.id, h.title, h.description, h.max_score,
		       h.created_at, h.updated_at, h.deadline
		FROM homeworks h
		JOIN "group" g ON h.group_id = g.id
		WHERE g.name = $1
		  AND h.teacher_id = $2
		  AND h.discipline_id = $3
	`

	rows, err := h.DB.Query(r.Context(), query, groupName, teacherId, disciplineId)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": fmt.Sprintf("db error: %v", err)})
		return
	}
	defer rows.Close()

	var homeworks []Homework

	for rows.Next() {
		var id, maxScore int
		var title, description string
		var createdAt, updatedAt, deadline time.Time

		err := rows.Scan(&id, &title, &description, &maxScore, &createdAt, &updatedAt, &deadline)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"error": fmt.Sprintf("scan error: %v", err)})
			return
		}

		// Форматируем время в строку "2006-01-02 15:04"
		hw := Homework{
			ID:          id,
			Title:       title,
			Description: description,
			MaxScore:    maxScore,
			CreatedAt:   createdAt.Format("2006.01.02 15:04"),
			UpdatedAt:   updatedAt.Format("2006.01.02 15:04"),
			Deadline:    deadline.Format("2006.01.02 15:04"),
		}
		homeworks = append(homeworks, hw)
	}

	if homeworks == nil {
		homeworks = []Homework{}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(homeworks)
}