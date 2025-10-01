package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"univer/models"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// mapRowToLab маппит строку из БД в структуру Lab
func mapRowToLab(row pgx.Row) (*models.Lab, error) {
	var l models.Lab
	var propertiesData []byte

	// Здесь важно, чтобы порядок столбцов соответствовал запросу
	err := row.Scan(
		&l.Id,
		&l.Name,
		&l.Count,
		&propertiesData, // JSON из properties
		&l.Year,
		&l.DefaultGrade,
	)
	if err != nil {
		return nil, fmt.Errorf("ошибка сканирования строки: %w", err)
	}

	// Если properties не пустой, пробуем его распарсить
	if len(propertiesData) > 0 {
		var labProperties models.LabProperties
		if err := json.Unmarshal(propertiesData, &labProperties); err != nil {
			return nil, fmt.Errorf("ошибка декодирования properties: %w", err)
		}
		l.Properties = labProperties
	} else {
		// Если properties пустой, создаем пустую структуру
		l.Properties = models.LabProperties{
			LabsGradeMax:           make(map[string]interface{}),
			ReducePointsEachModule: false,
		}
	}

	return &l, nil
}

// GetLabs получает список всех лабораторных работ
func GetLabs(dbpool *pgxpool.Pool) ([]*models.Lab, error) {
	query := `SELECT id, name, count, properties, year, default_grade FROM lab`

	rows, err := dbpool.Query(context.Background(), query)
	if err != nil {
		return nil, fmt.Errorf("ошибка выполнения запроса: %w", err)
	}
	defer rows.Close()

	var labs []*models.Lab

	for rows.Next() {
		lab, err := mapRowToLab(rows)
		if err != nil {
			return nil, err
		}
		labs = append(labs, lab)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("ошибка обработки строк: %w", err)
	}

	return labs, nil
}
