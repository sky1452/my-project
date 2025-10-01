package repository

import (
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5"
	"log"
	"univer/models"

	"github.com/jackc/pgx/v5/pgxpool"
)

// GetUniqueStudentGroups получает список всех уникальных групп студентов
func GetUniqueStudentGroups(dbpool *pgxpool.Pool) ([]string, error) {
	query := `SELECT DISTINCT properties->>'group' FROM users WHERE role_id = 3 ORDER BY properties->>'group'`

	rows, err := dbpool.Query(context.Background(), query)
	if err != nil {
		return nil, fmt.Errorf("ошибка выполнения запроса: %w", err)
	}
	defer rows.Close()

	var groups []string

	for rows.Next() {
		var groupName string
		if err := rows.Scan(&groupName); err != nil {
			log.Printf("ошибка сканирования данных: %v", err)
			continue
		}
		groups = append(groups, groupName)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("ошибка обработки строк: %w", err)
	}

	return groups, nil
}

func GetGroupByStudentId(dbpool *pgxpool.Pool, id int) (*models.Group, error) {
	var group models.Group
	query := `SELECT 
    u.properties->>'group' AS group_name,
    g.id AS group_id
	FROM public.users u
	LEFT JOIN public.group g ON (u.properties->>'group') = g.name
	WHERE u.role_id = 3 AND u.user_id = $1
	GROUP BY g.id, group_name
	ORDER BY group_name`

	row, err := dbpool.Query(context.Background(), query, id)
	if err != nil {
		return nil, fmt.Errorf("ошибка выполнения запроса: %w", err)
	}
	defer row.Close()

	for row.Next() {

		if err = row.Scan(&group.Name, &group.Id); err != nil {
			return nil, fmt.Errorf("ошибка сканирования: %w", err)
		}
	}
	//err = row.Scan(&group.Name, &group.Id)
	//if err != nil {
	//	return nil, fmt.Errorf("ошибка сканирования строки: %w", err)
	//}

	return &group, nil
}

func GetGroupByName(dbpool *pgxpool.Pool, name string) (*models.Group, error) {
	var group models.Group
	log.Printf("Ищу группу с именем: %s", name)

	query := `SELECT id, name FROM public."group"
          WHERE name ILIKE $1
          ORDER BY id ASC
          LIMIT 1`

	row := dbpool.QueryRow(context.Background(), query, name)
	err := row.Scan(&group.Id, &group.Name)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errors.New("группа " + name + " не найдена")
		}
		return nil, fmt.Errorf("ошибка выполнения запроса: %w", err)
	}

	return &group, nil
}
