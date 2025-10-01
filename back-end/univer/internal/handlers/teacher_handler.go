package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"strconv"
	"univer/internal/repository"
	"univer/models"
)

func (h *Handler) GetAllTeachersHandler(writer http.ResponseWriter, r *http.Request) {
	// Получаем список преподавателей из базы данных
	teachers, err := repository.FetchAllTeachers(h.DB)
	if err != nil {
		log.Printf("Ошибка при получении списка преподавателей: %v", err)
		http.Error(writer, "Ошибка при получении списка преподавателей", http.StatusInternalServerError)
		return
	}

	// Устанавливаем заголовки, чтобы вернуть результат в формате JSON
	writer.Header().Set("Content-Type", "application/json")

	// Логируем информацию (пример использования конфига)
	log.Printf("Отправляем список преподавателей. Сервер работает на порту %d", h.Config.Server.Port)

	// Преобразуем данные в JSON и отправляем клиенту
	if err := json.NewEncoder(writer).Encode(teachers); err != nil {
		log.Printf("Ошибка при кодировании ответа в JSON: %v", err)
		http.Error(writer, "Ошибка при обработке данных", http.StatusInternalServerError)
		return
	}
}

func (h *Handler) GetTeacherByIDHandler(writer http.ResponseWriter, r *http.Request) {
	// Получаем переменную из URL
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(writer, "Некорректный ID преподавателя", http.StatusBadRequest)
		return
	}

	// Получаем преподавателя из базы данных
	teacher, err := repository.FetchTeacherFromDB(h.DB, id)
	if err != nil {
		log.Printf("Ошибка получения преподавателя ID %d: %v", id, err)
		http.Error(writer, "Ошибка получения преподавателя", http.StatusInternalServerError)
		return
	}

	// Логируем информацию (пример использования конфига)
	log.Printf("Отправляем данные преподавателя ID %d. Сервер работает на порту %d", id, h.Config.Server.Port)

	// Устанавливаем заголовки, чтобы вернуть результат в формате JSON
	writer.Header().Set("Content-Type", "application/json")

	// Возвращаем преподавателя в JSON формате
	if err := json.NewEncoder(writer).Encode(teacher); err != nil {
		log.Printf("Ошибка кодирования ответа в JSON: %v", err)
		http.Error(writer, "Ошибка обработки данных", http.StatusInternalServerError)
		return
	}
}

// AddTeacherHandler добавляет нового преподавателя в базу данных
// @Summary Добавление преподавателя
// @Description Принимает JSON-данные о преподавателе и добавляет его в базу данных
// @Tags Преподаватели
// @Accept json
// @Produce json
// @Param teacher body models.AddTeacherRequest true "Данные преподавателя"
// @Success 201 {string} string "Преподаватель успешно добавлен"
// @Example {json} Статус 201
//
//	{
//	  "error": "Преподаватель успешно добавлен"
//	}
//
// @Failure 400 {object} string "Ошибка чтения тела запроса или декодирования данных"
// @Example {json} Ошибка 400
//
//	{
//	  "error": "Ошибка декодирования данных"
//	}
//
// @Failure 500 {object} string "Ошибка при добавлении преподавателя"
// @Example {json} Ошибка 500
//
//	{
//	  "error": "Ошибка при добавлении преподавателя в БД"
//	}
//
// @Router /teachers [post]
func (h *Handler) AddTeacherHandler(writer http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(writer, "Метод не поддерживается", http.StatusMethodNotAllowed)
		return
	}

	log.Printf("Метод запроса: %s", r.Method)
	log.Printf("Заголовки запроса: %+v", r.Header)

	// Декодируем тело запроса в структуру Teacher
	var newTeacher models.AddTeacherRequest
	if err := json.NewDecoder(r.Body).Decode(&newTeacher); err != nil {
		http.Error(writer, "ошибка декодирования данных", http.StatusBadRequest)
		log.Printf("Ошибка декодирования JSON: %v", err)
		return
	}

	log.Printf("Полученный объект: %+v", newTeacher)

	// Добавляем преподавателя в БД
	err := repository.AddTeacher(h.DB, &newTeacher)
	if err != nil {
		http.Error(writer, fmt.Sprintf("ошибка при добавлении преподавателя: %v", err), http.StatusInternalServerError)
		log.Printf("Ошибка при добавлении преподавателя в БД: %v", err)
		return
	}

	// Отправляем успешный ответ
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusCreated)

	response := map[string]interface{}{
		"message": "Преподаватель успешно добавлен",
		"teacher": newTeacher,
	}
	json.NewEncoder(writer).Encode(response)
}
