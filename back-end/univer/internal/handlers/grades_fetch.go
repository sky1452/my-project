package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

// Grade структура для отдачи на фронт
type Grade struct {
	StudentID    int     `json:"student_id"`
	GroupName    string  `json:"group_name"`
	DisciplineID int     `json:"discipline_id"`
	Discipline   string  `json:"discipline"`
	ModuleNumber int     `json:"module_number"`
	Score        *int    `json:"score"`
}

// GetGradesHandler отдаёт оценки по группе и дисциплине
func (h *Handler) GetGradesHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	groupName := vars["group"]
	disciplineIDStr := vars["discipline"]

	if groupName == "" || disciplineIDStr == "" {
		http.Error(w, "Не указана группа или дисциплина", http.StatusBadRequest)
		return
	}

	disciplineID, err := strconv.Atoi(disciplineIDStr)
	if err != nil {
		http.Error(w, "Некорректный ID дисциплины", http.StatusBadRequest)
		return
	}

	rows, err := h.DB.Query(r.Context(),
		`SELECT student_id, group_name, discipline_id, discipline_type, module_number, score
		 FROM grades
		 WHERE group_name=$1 AND discipline_id=$2`, groupName, disciplineID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Ошибка запроса оценок: %v", err), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var grades []Grade
	for rows.Next() {
		var g Grade
		if err := rows.Scan(&g.StudentID, &g.GroupName, &g.DisciplineID, &g.Discipline, &g.ModuleNumber, &g.Score); err != nil {
			http.Error(w, fmt.Sprintf("Ошибка чтения данных: %v", err), http.StatusInternalServerError)
			return
		}
		grades = append(grades, g)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(grades)
}
