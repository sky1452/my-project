package handlers

import (
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Role     int    `json:"role"`
}

type LoginResponse struct {
	Success     bool   `json:"success"`
	Message     string `json:"message,omitempty"`
	Role        int    `json:"role"`
	FullName    string `json:"fullName,omitempty"`
	Email       string `json:"email,omitempty"`
	Position    string `json:"position,omitempty"`
	Departament string `json:"departament,omitempty"`
	Stazh       int    `json:"stazh,omitempty"`
	Dop         string `json:"dop,omitempty"`
	Avatar      string `json:"avatar,omitempty"`
	UserID      int    `json:"userId,omitempty"`
}

func (h *Handler) LoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
		return
	}

	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Неверный формат запроса", http.StatusBadRequest)
		return
	}

	log.Printf("DEBUG: login attempt username='%s' role=%d", req.Username, req.Role)

	var userID int
	var hashedPassword, fullName, email, position, departament string
	var roleFromDB int
	var stazh sql.NullInt64
	var dop sql.NullString
	var avatarBytes []byte

	err := h.DB.QueryRow(r.Context(),
		`SELECT user_id, password, 
				role_id, 
				name, 
				email, 
				COALESCE(properties->>'position', '') AS position,     -- FIX: COALESCE
				COALESCE(properties->>'departament', '') AS departament, -- FIX: COALESCE
				stazh,
				dop,
				avatar
		 FROM users 
		 WHERE name = $1`, req.Username,
	).Scan(&userID, &hashedPassword, &roleFromDB, &fullName, &email,
		&position, &departament, &stazh, &dop, &avatarBytes)

	if err != nil {
		// DEBUG: выводим реальную ошибку SQL
		log.Printf("[Login failed] Query error for username='%s': %v", req.Username, err)
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(LoginResponse{
			Success: false,
			Message: "Пользователь не найден",
		})
		return
	}



	if err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(req.Password)); err != nil {
		log.Printf("[Login failed] Неверный пароль для пользователя: %s (%s)", req.Username, fullName)
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(LoginResponse{
			Success: false,
			Message: "Неверный пароль",
		})
		return
	}

	if req.Role != roleFromDB {
		log.Printf("[Login failed] Роль не совпадает для пользователя: %s (%s). Ожидалось %d, пришло %d",
			req.Username, fullName, roleFromDB, req.Role)
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(LoginResponse{
			Success: false,
			Message: fmt.Sprintf("Роль не совпадает: ожидалось %d, пришло %d", roleFromDB, req.Role),
		})
		return
	}

	realStazh := 0
	if stazh.Valid {
		realStazh = int(stazh.Int64)
	}
	realDop := ""
	if dop.Valid {
		realDop = dop.String
	}

	avatarData := ""
	if len(avatarBytes) > 0 {
		avatarData = "data:image/jpeg;base64," + base64.StdEncoding.EncodeToString(avatarBytes)
	}

	log.Printf("[Login success] Пользователь вошёл: %s (%s)", req.Username, fullName)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(LoginResponse{
		Success:     true,
		Role:        roleFromDB,
		FullName:    fullName,
		Email:       email,
		Position:    position,
		Departament: departament,
		Stazh:       realStazh,
		Dop:         realDop,
		Avatar:      avatarData,
		UserID:      userID,
	})
}
