package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"database/sql"
)

type UpdateStazhRequest struct {
	Name  string `json:"name"`
	Stazh int    `json:"stazh"`
}

func (h *Handler) UpdateStazhHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
		return
	}

	var req UpdateStazhRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Неверный формат запроса", http.StatusBadRequest)
		return
	}

	// Получаем старое значение stazh безопасно для NULL
	var oldStazh sql.NullInt64
	err := h.DB.QueryRow(r.Context(),
		`SELECT stazh FROM users WHERE name=$1`,
		req.Name,
	).Scan(&oldStazh)
	if err != nil {
		log.Printf("[UpdateStazh] Не удалось получить старое значение стажа для %s: %v\n", req.Name, err)
		http.Error(w, "Пользователь не найден", http.StatusNotFound)
		return
	}

	oldStazhVal := 0
	if oldStazh.Valid {
		oldStazhVal = int(oldStazh.Int64)
	}

	// Обновляем stazh
	_, err = h.DB.Exec(r.Context(),
		`UPDATE users SET stazh=$1 WHERE name=$2`,
		req.Stazh, req.Name,
	)
	if err != nil {
		log.Printf("[UpdateStazh failed] Ошибка обновления стажа для %s: %v\n", req.Name, err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"message": "Ошибка при обновлении стажа",
		})
		return
	}

	log.Printf("[UpdateStazh] Пользователь: %s | Стаж: %d → %d\n", req.Name, oldStazhVal, req.Stazh)

	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
	})
}
