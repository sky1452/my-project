package repository

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

// GetWeekTypeByDate возвращает "числитель" (true)  или "знаменатель" (false) по дате.
func GetWeekTypeByDate(dbpool *pgxpool.Pool, date time.Time) (*bool, error) {
	var isNumerator bool

	query := `SELECT is_numerator FROM weeks WHERE week_start <= $1 AND week_end >= $1 LIMIT 1`

	err := dbpool.QueryRow(context.Background(), query, date).Scan(&isNumerator)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil // Если данных нет, возвращаем nil
		}
		return nil, fmt.Errorf("ошибка при запросе типа недели для %v: %w", date, err)
	}

	return &isNumerator, nil // Возвращаем указатель на результат
}
