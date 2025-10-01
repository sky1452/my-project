package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"database/sql"
)

type UpdateDopRequest struct {
	Name string `json:"name"`
	Dop  string `json:"dop"`
}

func (h *Handler) UpdateDopHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
		return
	}

	var req UpdateDopRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Неверный формат запроса", http.StatusBadRequest)
		return
	}

	// Получаем старое значение dop безопасно для NULL
	var oldDop sql.NullString
	err := h.DB.QueryRow(r.Context(),
		`SELECT dop FROM users WHERE name=$1`,
		req.Name,
	).Scan(&oldDop)
	if err != nil {
		log.Printf("[UpdateDop] Не удалось получить старое значение информации для %s: %v\n", req.Name, err)
		http.Error(w, "Пользователь не найден", http.StatusNotFound)
		return
	}

	oldDopStr := ""
	if oldDop.Valid {
		oldDopStr = oldDop.String
	}

	// Обновляем dop
	_, err = h.DB.Exec(r.Context(),
		`UPDATE users SET dop=$1 WHERE name=$2`,
		req.Dop, req.Name,
	)
	if err != nil {
		log.Printf("[UpdateDop failed] Ошибка обновления информации о себе для %s: %v\n", req.Name, err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"message": "Ошибка при обновлении информации о себе",
		})
		return
	}

	log.Printf("[UpdateDop] Пользователь: %s | Доп: %s → %s\n", req.Name, oldDopStr, req.Dop)

	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
	})
}
