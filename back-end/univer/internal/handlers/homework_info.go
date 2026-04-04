package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type ProgressResponse struct {
	Success   bool `json:"success"`
	UserID    int  `json:"userId"`
	StudentID int  `json:"studentId"`
	Progress  int  `json:"progress"`
}

func (h *Handler) GetStudentProgressHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	ctx := r.Context()

	vars := mux.Vars(r)
	userIDStr := vars["userId"]
	if userIDStr == "" {
		http.Error(w, "userId is required", http.StatusBadRequest)
		return
	}

	studentID, err := strconv.Atoi(userIDStr)
	if err != nil {
		http.Error(w, "Некорректный userId", http.StatusBadRequest)
		return
	}

	var progress int
	err = h.DB.QueryRow(ctx, `
		SELECT COUNT(DISTINCT task_id)
		FROM submissions
		WHERE student_id = $1
	`, studentID).Scan(&progress)
	if err != nil {
		http.Error(w, fmt.Sprintf("Ошибка при подсчёте прогресса: %v", err), http.StatusInternalServerError)
		return
	}

	resp := ProgressResponse{
		Success:   true,
		UserID:    studentID,
		StudentID: studentID,
		Progress:  progress,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}