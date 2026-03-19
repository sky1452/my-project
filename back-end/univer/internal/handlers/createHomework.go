package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type CreateHomeworkRequest struct {
	TeacherID    int    `json:"teacher_id"`
	GroupName    string `json:"groupName"`
	DisciplineID int    `json:"disciplineId"`
	Title        string `json:"title"`
	Description  string `json:"description"`
	MaxScore     int    `json:"max_score"`
	Deadline     string `json:"deadline"`
}

func (h *Handler) CreateHomeworkHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	ctx := r.Context()

	var req CreateHomeworkRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Неверный JSON", http.StatusBadRequest)
		return
	}

	// Проверка обязательных полей
	if req.TeacherID == 0 || req.GroupName == "" || req.Title == "" ||
		req.Description == "" || req.MaxScore == 0 || req.Deadline == "" || req.DisciplineID == 0 {
		http.Error(w, "Заполните все обязательные поля", http.StatusBadRequest)
		return
	}

	var groupID int
	err := h.DB.QueryRow(ctx, "SELECT id FROM \"group\" WHERE name=$1", req.GroupName).Scan(&groupID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Группа не найдена: %v", err), http.StatusBadRequest)
		return
	}

	// (опционально, но правильно) проверка дисциплины
	var disciplineExists bool
	err = h.DB.QueryRow(ctx, "SELECT EXISTS(SELECT 1 FROM academic_subject WHERE id=$1)", req.DisciplineID).Scan(&disciplineExists)
	if err != nil || !disciplineExists {
		http.Error(w, "Дисциплина не найдена", http.StatusBadRequest)
		return
	}

	deadlineTime, err := time.Parse(time.RFC3339, req.Deadline)
	if err != nil {
		http.Error(w, "Неверный формат deadline. Используйте ISO 8601", http.StatusBadRequest)
		return
	}

	var homeworkID int
	query := `
		INSERT INTO homeworks 
		(teacher_id, group_id, discipline_id, title, description, max_score, deadline, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, NOW(), NOW())
		RETURNING id
	`

	err = h.DB.QueryRow(ctx, query,
		req.TeacherID,
		groupID,
		req.DisciplineID,
		req.Title,
		req.Description,
		req.MaxScore,
		deadlineTime,
	).Scan(&homeworkID)

	if err != nil {
		http.Error(w, fmt.Sprintf("Ошибка при создании задания: %v", err), http.StatusInternalServerError)
		return
	}

	resp := map[string]interface{}{
		"success":     true,
		"homework_id": homeworkID,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}
