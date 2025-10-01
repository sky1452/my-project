package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"strconv"
	"time"
	"univer/internal/repository"
	"univer/models"
)

// GetParaNum определяет текущую пару или перемену
// @Summary Определить номер текущей пары
// @Description Возвращает номер текущей пары или сообщает, что сейчас перемена
// @Tags Para
// @Accept  json
// @Produce  json
// @Success 200 {object} int "Номер пары (1-6) или 'Перемена'"
// @Failure 500 {object} string "Ошибка при определении номера пары"
// @Router /paraNum [get]

func (h *Handler) GetParaNum(writer http.ResponseWriter, reader *http.Request) {
	date := time.Now()
	paraNum, err := repository.GetParaNumByTime(h.DB, date)
	writer.Header().Set("Content-Type", "application/json")

	if err != nil {
		fmt.Println("Ошибка:", err)
		json.NewEncoder(writer).Encode("Ошибка")
	} else if paraNum == 0 {
		fmt.Printf("Перемена")
		json.NewEncoder(writer).Encode("Перемена")
	} else {
		fmt.Printf("Пара: %d\n", paraNum)
		writer.Header().Set("Content-Type", "application/json")
		json.NewEncoder(writer).Encode(paraNum)
	}
}

// MyCurrentPara определяет текущую пару
// @Summary Определить текущую пару
// @Description Возвращает текущую пару
// @Tags Para
// @Accept  json
// @Produce  json
// @Success 200 {object} models.Para
// @Failure 404 {object} string "Текущей пары не найдено"
// @Router /myPara [get]

func (h *Handler) MyCurrentPara(writer http.ResponseWriter, reader *http.Request) {

	date := time.Now()
	var id rune = h.Config.Ids.My
	currentPara, err := repository.GetCurrentParaById(h.DB, id, date)

	if err != nil {
		// Логируем ошибку
		log.Printf("текущей пары для %d не найдено: "+err.Error(), id)
		http.Error(writer, "текущей пары не найдено: "+err.Error(), http.StatusNotFound)
		return
	}
	json.NewEncoder(writer).Encode(currentPara)
}

// MyTodayParas возвращает список пар на сегодня для пользователя или группы
// @Summary Получение списка пар на сегодня
// @Description Возвращает массив занятий (пар) на текущий день для указанного пользователя (по ID) или группы (по названию)
// @Tags Para
// @Accept json
// @Produce json
// @Param id path string true "ID пользователя или название группы"
// @Success 200 {array} models.Para "Список занятий на сегодня"
// @Failure 400 {object} string "Некорректный запрос"
// @Failure 404 {object} string "Группа не найдена"
// @Example {json} Ошибка сервера
//
//	{
//	  "error": "Ошибка при поиске группы: группа ПРИб-23211 не найдена"
//	}
//
// @Failure 500 {object} string "Ошибка при получении расписания"
// @Example {json} Ошибка сервера
//
//	{
//	  "error": "Ошибка при определении роли",
//	  "error": "Ошибка при получении группы студента"
//	}
//
// @Router /myTodayParas/{id} [get]

func (h *Handler) MyTodayParas(writer http.ResponseWriter, reader *http.Request) {

	vars := mux.Vars(reader)
	param := vars["id"] // В URL может быть либо id пользователя, либо название группы

	var paras []*models.Para
	var group *models.Group
	var err error

	// 1. Проверяем, является ли param числом (ID пользователя)
	if userID, convErr := strconv.Atoi(param); convErr == nil {
		// Если param - число, значит ищем расписание по ID пользователя
		role, err := repository.GetUserRole(h.DB, userID)
		if err != nil {
			log.Printf("Ошибка при получении роли: %v", err)
			http.Error(writer, "Ошибка при определении роли", http.StatusInternalServerError)
			return
		}

		if role == "teacher" {
			paras, err = repository.GetTodayParasById(h.DB, userID, role, time.Now())
		} else if role == "student" {
			group, err := repository.GetGroupByStudentId(h.DB, userID)
			if err != nil {
				log.Printf("Ошибка при получении группы студента: %v", err)
				http.Error(writer, "Ошибка при получении группы студента", http.StatusInternalServerError)
				return
			}
			paras, err = repository.GetTodayParasByGroup(h.DB, group, time.Now())
		}
	} else {
		// 2. Если param - не число, значит это название группы
		group, err = repository.GetGroupByName(h.DB, param)
		if err != nil {
			log.Printf("Ошибка при поске группы: %v", err)
			http.Error(writer, "Ошибка при поиске группы: "+err.Error(), http.StatusNotFound)
			return
		}
		paras, err = repository.GetTodayParasByGroup(h.DB, group, time.Now())
	}

	// 3. Обрабатываем ошибки
	if err != nil {
		log.Printf("Ошибка при получении расписания: %v", err)
		http.Error(writer, "Ошибка при получении расписания", http.StatusInternalServerError)
		return
	}

	// 4. Отправляем JSON-ответ
	writer.Header().Set("Content-Type", "application/json")
	json.NewEncoder(writer).Encode(paras)

}

// MySchedule возвращает список пар для пользователя или группы
// @Summary Получение списка пар
// @Description Возвращает массив занятий (пар) указанного пользователя (по ID) или группы (по названию)
// @Tags Para
// @Accept json
// @Produce json
// @Param id path string true "ID пользователя или название группы"
// @Success 200 {array} models.Para "Список занятий"
// @Failure 400 {object} string "Некорректный запрос"
// @Failure 404 {object} string "Группа не найдена"
// @Failure 500 {object} string "Ошибка при получении расписания"
// @Router /mySchedule/{id} [get]

func (h *Handler) MySchedule(writer http.ResponseWriter, reader *http.Request) {

	vars := mux.Vars(reader)
	param := vars["id"] // В URL может быть либо id пользователя, либо название группы

	var paras []*models.Para
	var group *models.Group
	var err error

	// 1. Проверяем, является ли param числом (ID пользователя)
	if userID, convErr := strconv.Atoi(param); convErr == nil {
		// Если param - число, значит ищем расписание по ID пользователя
		role, err := repository.GetUserRole(h.DB, userID)
		if err != nil {
			log.Printf("Ошибка при получении роли: %v", err)
			http.Error(writer, "Ошибка при определении роли", http.StatusInternalServerError)
			return
		}

		if role == "teacher" {
			paras, err = repository.GetScheduleById(h.DB, userID, role)
		} else if role == "student" {
			group, err := repository.GetGroupByStudentId(h.DB, userID)
			if err != nil {
				log.Printf("Ошибка при получении группы студента: %v", err)
				http.Error(writer, "Ошибка при получении группы студента", http.StatusInternalServerError)
				return
			}
			paras, err = repository.GetScheduleByGroup(h.DB, group)
		}
	} else {
		// 2. Если param - не число, значит это название группы
		group, err = repository.GetGroupByName(h.DB, param)
		if err != nil {
			log.Printf("Ошибка при поске группы: %v", err)
			http.Error(writer, "Ошибка при поиске группы: "+err.Error(), http.StatusInternalServerError)
			return
		}
		paras, err = repository.GetScheduleByGroup(h.DB, group)
	}

	// 4. Отправляем JSON-ответ
	writer.Header().Set("Content-Type", "application/json")
	json.NewEncoder(writer).Encode(paras)
}
