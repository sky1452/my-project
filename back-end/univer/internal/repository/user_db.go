package repository

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
)

func GetUserRole(dbpool *pgxpool.Pool, userID int) (string, error) {
	var role string

	query := `SELECT 
        CASE 
            WHEN role_id = 3 THEN 'student'
--             WHEN role_id = 2 THEN 'teacher'
			WHEN role_id <= 2 THEN 'teacher'
            ELSE 'unknown' 
        END 
    FROM public.users 
    WHERE user_id = $1`

	err := dbpool.QueryRow(context.Background(), query, userID).Scan(&role)
	if err != nil {
		return "", fmt.Errorf("ошибка получения роли пользователя: %w", err)
	}
	return role, nil
}
