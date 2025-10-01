package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
	"univer/internal/repository"
)

func (h *Handler) GetWeekType(writer http.ResponseWriter, reader *http.Request) {

	date := time.Now()
	weekType, err := repository.GetWeekTypeByDate(h.DB, date)
	writer.Header().Set("Content-Type", "application/json")

	if err != nil {
		fmt.Println("Ошибка:", err)
		json.NewEncoder(writer).Encode("ошибка")
	} else if weekType == nil {
		fmt.Println("Нет данных о неделе для этой даты")
		json.NewEncoder(writer).Encode("Нет данных о неделе для этой даты")
	} else if *weekType {
		fmt.Println("Числитель")
		json.NewEncoder(writer).Encode("Числитель")
	} else {
		fmt.Println("Знаменатель")
		json.NewEncoder(writer).Encode("Знаменатель")
	}

}
