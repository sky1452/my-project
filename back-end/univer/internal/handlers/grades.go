package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

// Таблица для оценок
const gradesTable = "grades"

// GradePayload структура для одной оценки
type GradePayload struct {
	StudentID      int  `json:"student_id"`
	GroupName      string `json:"group_name"`
	DisciplineID   int  `json:"discipline_id"`
	DisciplineType string `json:"discipline_type"`
	ModuleNumber   int  `json:"module_number"`
	Score          *int `json:"score"`
}

// UpdateGradesHandler обновляет оценки и логирует изменения
func (h *Handler) UpdateGradesHandler(w http.ResponseWriter, r *http.Request) {
	var payload []GradePayload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		log.Printf("[ERROR] Неверный формат данных: %v", err)
		http.Error(w, "Неверный формат данных", http.StatusBadRequest)
		return
	}

	

	tx, err := h.DB.Begin(r.Context())
	if err != nil {
		log.Printf("[ERROR] Ошибка транзакции: %v", err)
		http.Error(w, "Ошибка транзакции", http.StatusInternalServerError)
		return
	}
	defer tx.Rollback(r.Context())

	for _, g := range payload {
		var studentName string
		err := tx.QueryRow(r.Context(),
			`SELECT name FROM users WHERE user_id=$1 AND role_id=3`, g.StudentID).Scan(&studentName)
		if err != nil {
			studentName = fmt.Sprintf("ID %d", g.StudentID)
		}

		if g.Score == nil {
			// Удаляем оценку
			cmdTag, err := tx.Exec(r.Context(),
				fmt.Sprintf(`DELETE FROM %s WHERE student_id=$1 AND discipline_id=$2 AND module_number=$3`, gradesTable),
				g.StudentID, g.DisciplineID, g.ModuleNumber)
			if err != nil {
				log.Printf("[ERROR] Ошибка удаления оценки: %v | Payload: %+v", err, g)
				http.Error(w, fmt.Sprintf("Ошибка удаления оценки: %v", err), http.StatusInternalServerError)
				return
			}
			log.Printf("[DELETE] Студент=%s, Дисциплина=%s, Модуль=%d, затронуто строк=%d",
				studentName, g.DisciplineType, g.ModuleNumber, cmdTag.RowsAffected())
		} else {
			// Вставка или обновление
			cmdTag, err := tx.Exec(r.Context(),
				fmt.Sprintf(`INSERT INTO %s
				(student_id, group_name, discipline_id, discipline_type, module_number, score, created_at, updated_at)
				VALUES ($1,$2,$3,$4,$5,$6,NOW(),NOW())
				ON CONFLICT (student_id, discipline_id, module_number)
				DO UPDATE SET score=EXCLUDED.score, updated_at=NOW()`, gradesTable),
				g.StudentID, g.GroupName, g.DisciplineID, g.DisciplineType, g.ModuleNumber, *g.Score)
			if err != nil {
				log.Printf("[ERROR] Ошибка сохранения оценки: %v | Payload: %+v", err, g)
				http.Error(w, fmt.Sprintf("Ошибка сохранения оценки: %v", err), http.StatusInternalServerError)
				return
			}
			log.Printf("[UPSERT] Студент=%s, Дисциплина=%s, Модуль=%d, Новая оценка=%d, затронуто строк=%d",
				studentName, g.DisciplineType, g.ModuleNumber, *g.Score, cmdTag.RowsAffected())
		}
	}

	if err := tx.Commit(r.Context()); err != nil {
		log.Printf("[ERROR] Ошибка коммита: %v", err)
		http.Error(w, "Ошибка коммита", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message":"Оценки успешно сохранены"}`))
}

// GetStudentsByGroupForGrades отдаёт студентов по группе (для страницы оценок)
func (h *Handler) GetStudentsByGroupForGrades(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	groupName := vars["group"]
	if groupName == "" {
		http.Error(w, "Не указана группа", http.StatusBadRequest)
		return
	}

	rows, err := h.DB.Query(r.Context(),
		`SELECT user_id, name 
		 FROM users 
		 WHERE properties->>'group'=$1 AND role_id=3`, groupName)
	if err != nil {
		log.Printf("[ERROR] Ошибка запроса студентов: %v", err)
		http.Error(w, fmt.Sprintf("Ошибка запроса студентов: %v", err), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	type Student struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	}
	students := []Student{}

	for rows.Next() {
		var s Student
		if err := rows.Scan(&s.ID, &s.Name); err != nil {
			log.Printf("[ERROR] Ошибка чтения данных: %v", err)
			http.Error(w, fmt.Sprintf("Ошибка чтения данных: %v", err), http.StatusInternalServerError)
			return
		}
		students = append(students, s)
	}

	if err := rows.Err(); err != nil {
		log.Printf("[ERROR] Ошибка после итерации rows: %v", err)
		http.Error(w, fmt.Sprintf("Ошибка после итерации rows: %v", err), http.StatusInternalServerError)
		return
	}

	log.Printf("[INFO] Получено студентов: %d для группы %s", len(students), groupName)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(students)
}
