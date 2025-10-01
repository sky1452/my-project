package handlers

import (
	"encoding/base64"
	"encoding/json"
	"log"
	"net/http"
)

type UpdateAvatarRequest struct {
	Name   string `json:"name"`
	Avatar string `json:"avatar"` // base64 строки
}

func (h *Handler) UpdateAvatarHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
		return
	}

	var req UpdateAvatarRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Неверный формат запроса", http.StatusBadRequest)
		return
	}

	// Декодируем base64
	avatarBytes, err := base64.StdEncoding.DecodeString(req.Avatar)
	if err != nil {
		http.Error(w, "Ошибка декодирования изображения", http.StatusBadRequest)
		return
	}

	_, err = h.DB.Exec(r.Context(),
		`UPDATE users SET avatar=$1 WHERE name=$2`,
		avatarBytes, req.Name,
	)
	if err != nil {
		log.Printf("[UpdateAvatar failed] Ошибка обновления аватарки для %s: %v\n", req.Name, err)
		http.Error(w, "Ошибка обновления аватарки", http.StatusInternalServerError)
		return
	}

	log.Printf("[UpdateAvatar] Пользователь: %s обновил аватарку\n", req.Name)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{"success": true})
	log.Printf("[DEBUG] Аватарка для пользователя %s успешно обновлена в базе данных\n", req.Name)
}
