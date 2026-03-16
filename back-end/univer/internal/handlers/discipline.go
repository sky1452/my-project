package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type AcademicSubject struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

func (h *Handler) GetDisciplineById(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["disciplineId"]

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "invalid discipline id", http.StatusBadRequest)
		return
	}

	query := `
		SELECT id, name
		FROM academic_subject
		WHERE id = $1
	`

	var subject AcademicSubject

	err = h.DB.QueryRow(context.Background(), query, id).Scan(
		&subject.Id,
		&subject.Name,
	)

	if err != nil {
		http.Error(w, "discipline not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(subject)
}