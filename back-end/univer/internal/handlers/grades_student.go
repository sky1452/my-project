package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

// StudentGrade структура для отдачи на фронт
type StudentGrade struct {
	DisciplineID   int             `json:"discipline_id"`
	DisciplineType string          `json:"discipline_type"`
	Modules        map[int]*int    `json:"modules"` // модуль → оценка (NULL поддерживаем через *int)
}

// GetStudentGrades возвращает оценки студента, сгруппированные по дисциплинам и модулям
func (h *Handler) GetStudentGrades(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	studentIDStr := vars["studentId"]

	if studentIDStr == "" {
		http.Error(w, "Не указан studentId", http.StatusBadRequest)
		return
	}

	studentID, err := strconv.Atoi(studentIDStr)
	if err != nil {
		http.Error(w, "Некорректный studentId", http.StatusBadRequest)
		return
	}

	rows, err := h.DB.Query(r.Context(),
		`SELECT discipline_id, discipline_type, module_number, score
		 FROM grades
		 WHERE student_id=$1
		 ORDER BY discipline_type, module_number`, studentID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Ошибка запроса оценок: %v", err), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	result := make(map[int]*StudentGrade)

	for rows.Next() {
		var disciplineID, moduleNumber int
		var disciplineType string
		var score *int

		if err := rows.Scan(&disciplineID, &disciplineType, &moduleNumber, &score); err != nil {
			http.Error(w, fmt.Sprintf("Ошибка чтения данных: %v", err), http.StatusInternalServerError)
			return
		}

		// если дисциплина ещё не в карте — добавляем
		if _, exists := result[disciplineID]; !exists {
			result[disciplineID] = &StudentGrade{
				DisciplineID:   disciplineID,
				DisciplineType: disciplineType,
				Modules:        make(map[int]*int),
			}
		}

		// сохраняем оценку по модулю
		result[disciplineID].Modules[moduleNumber] = score
	}

	// map -> slice
	final := make([]*StudentGrade, 0, len(result))
	for _, v := range result {
		final = append(final, v)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(final)
}
