package handlers

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"strconv"
	"univer/internal/repository"
)

func (h *Handler) GetAllStudentsHandler(writer http.ResponseWriter, reader *http.Request) {
	// Получаем список студентов из базы данных
	students, err := repository.FetchAllTeachers(h.DB)
	if err != nil {
		// Логируем ошибку
		log.Printf("ошибка при получении списка студентов: %v", err)
		http.Error(writer, "Ошибка при получении списка студентов", http.StatusInternalServerError)
		return
	}

	// Устанавливаем заголовки, чтобы вернуть результат в формате JSON
	writer.Header().Set("Content-Type", "application/json")

	// Преобразуем данные в JSON и отправляем клиенту
	if err := json.NewEncoder(writer).Encode(students); err != nil {
		log.Printf("ошибка при кодировании ответа в JSON: %v", err)
		http.Error(writer, "Ошибка при обработке данных", http.StatusInternalServerError)
		return
	}
}

func (h *Handler) GetStudentsByGroupHandler(writer http.ResponseWriter, reader *http.Request) {
	// Получаем переменную из URL
	vars := mux.Vars(reader)
	var group = vars["group"]
	//if err != nil {
	//	http.Error(writer, "Invalid group", http.StatusBadRequest)
	//	return
	//}
	//fmt.Printf("group: %s\n", group)
	students, err := repository.FetchStudentsByGroupFromDB(h.DB, group)
	if err != nil {
		// Логируем ошибку
		log.Printf("ошибка при получении списка студентов: %v", err)
		http.Error(writer, "Ошибка при получении списка студентов", http.StatusInternalServerError)
		return
	}

	// Устанавливаем заголовки, чтобы вернуть результат в формате JSON
	writer.Header().Set("Content-Type", "application/json")

	// Преобразуем данные в JSON и отправляем клиенту
	if err := json.NewEncoder(writer).Encode(students); err != nil {
		log.Printf("ошибка при кодировании ответа в JSON: %v", err)
		http.Error(writer, "Ошибка при обработке данных", http.StatusInternalServerError)
		return
	}
}
func (h *Handler) GetStudentByIDHandler(writer http.ResponseWriter, reader *http.Request) {
	// Получаем переменную из URL
	vars := mux.Vars(reader)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(writer, "Invalid student ID", http.StatusBadRequest)
		return
	}

	// Получаем студента из базы данных
	student, err := repository.FetchStudentFromDB(h.DB, id)
	if err != nil {
		http.Error(writer, "Ошибка получения студента", http.StatusInternalServerError)
		return
	}
	// Возвращаем преподавателя в JSON формате
	json.NewEncoder(writer).Encode(student)
}
