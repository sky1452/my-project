package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"univer/internal/repository"
)

func (h *Handler) GetAllLabs(writer http.ResponseWriter, reader *http.Request) {
	// Получаем список групп
	labs, err := repository.GetLabs(h.DB)
	if err != nil {
		// Логируем ошибку
		log.Printf("ошибка при получении списка лабораторных: %v", err)
		http.Error(writer, "Ошибка при получении списка лабораторных", http.StatusInternalServerError)
		return
	}

	// Устанавливаем заголовки, чтобы вернуть результат в формате JSON
	writer.Header().Set("Content-Type", "application/json")

	// Преобразуем данные в JSON и отправляем клиенту
	if err := json.NewEncoder(writer).Encode(labs); err != nil {
		log.Printf("ошибка при кодировании ответа в JSON: %v", err)
		http.Error(writer, "Ошибка при обработке данных", http.StatusInternalServerError)
		return
	}
}
