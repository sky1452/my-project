package handlers

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

type updateHomeworkRequest struct {
	DisciplineID int    `json:"disciplineId"`
	Group        string `json:"group"`
	Title        string `json:"title"`
	Description  string `json:"description"`
	Deadline     string `json:"deadline"`
	MaxScore     int    `json:"maxScore"`
}

func (h *Handler) UpdateHomework(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(r)

	homeworkIDStr := vars["homeworkId"]
	teacherIDStr := vars["teacherId"]

	homeworkID, err := strconv.Atoi(homeworkIDStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "invalid homeworkId",
		})
		return
	}

	teacherID, err := strconv.Atoi(teacherIDStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "invalid teacherId",
		})
		return
	}

	var req updateHomeworkRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "invalid request body",
		})
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	var groupID int

	groupQuery := `
		SELECT id
		FROM "group"
		WHERE name = $1
		LIMIT 1
	`

	err = h.DB.QueryRow(ctx, groupQuery, req.Group).Scan(&groupID)
	if err != nil {
		log.Println("UpdateHomework group lookup error:", err)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "group not found",
		})
		return
	}

	updateQuery := `
		UPDATE homeworks
		SET title = $1,
		    description = $2,
		    deadline = $3,
		    max_score = $4,
		    discipline_id = $5,
		    group_id = $6,
		    updated_at = NOW()
		WHERE id = $7 AND teacher_id = $8
	`

	result, err := h.DB.Exec(
		ctx,
		updateQuery,
		req.Title,
		req.Description,
		req.Deadline,
		req.MaxScore,
		req.DisciplineID,
		groupID,
		homeworkID,
		teacherID,
	)

	if err != nil {
		log.Println("UpdateHomework Exec error:", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{
			"error": err.Error(),
		})
		return
	}

	rowsAffected := result.RowsAffected()

	if rowsAffected == 0 {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "homework not found or access denied",
		})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]any{
		"success": true,
		"message": "homework updated successfully",
	})
}