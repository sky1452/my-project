package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"univer/models"
)

func mapRowToStudent(row pgx.Row) (*models.Student, error) {
	var s models.Student
	var propertiesData []byte // Поле для хранения JSON-данных

	// Сканируем поля из строки
	err := row.Scan(&s.Id, &s.Name, &s.Email, &propertiesData)
	if err != nil {
		return nil, fmt.Errorf("ошибка сканирования строки: %w", err)
	}

	// Декодируем поле properties в структуру TeacherProperties
	var studentProperties models.StudentProperties
	if len(propertiesData) > 0 {
		err = json.Unmarshal(propertiesData, &studentProperties)
		if err != nil {
			return nil, fmt.Errorf("ошибка декодирования properties: %w", err)
		}
	}
	s.Properties = studentProperties

	return &s, nil
}

func FetchAllStudents(dbpool *pgxpool.Pool) ([]*models.Student, error) {
	// Выполняем SQL-запрос
	rows, err := dbpool.Query(context.Background(),
		"SELECT user_id, name, email, properties FROM public.users WHERE role_id=3")
	if err != nil {
		return nil, fmt.Errorf("ошибка выполнения запроса: %w", err)
	}
	defer rows.Close()

	var students []*models.Student

	// Итерируем по результатам и мапим их в структуры Teacher
	for rows.Next() {
		student, err := mapRowToStudent(rows)
		if err != nil {
			return nil, err
		}
		students = append(students, student)
	}

	// Проверяем на наличие ошибок при итерации по строкам
	if rows.Err() != nil {
		return nil, fmt.Errorf("ошибка при обработке строк: %w", rows.Err())
	}

	return students, nil
}

func FetchStudentFromDB(dbpool *pgxpool.Pool, id int) (*models.Student, error) {
	// Выполняем SQL-запрос с параметром
	row := dbpool.QueryRow(context.Background(), ""+
		"SELECT user_id, name, email, properties FROM public.users WHERE user_id=$1 AND role_id=3", id)

	// Используем ранее созданную функцию для маппинга данных
	student, err := mapRowToStudent(row)
	if err != nil {
		return nil, fmt.Errorf("ошибка получения преподавателя по ID: %w", err)
	}
	return student, nil
}

func FetchStudentsByGroupFromDB(dbpool *pgxpool.Pool, groupName string) ([]*models.Student, error) {
	// Выполняем SQL-запрос

	rows, err := dbpool.Query(context.Background(),
		`SELECT user_id, name, email, properties FROM public.users WHERE role_id=3 AND properties->>'group' = $1`, groupName)
	if err != nil {
		return nil, fmt.Errorf("ошибка выполнения запроса: %w", err)
	}
	defer rows.Close()

	var students []*models.Student

	// Итерируем по результатам и мапим их в структуры Teacher
	for rows.Next() {
		student, err := mapRowToStudent(rows)
		if err != nil {
			return nil, err
		}
		students = append(students, student)
	}

	// Проверяем на наличие ошибок при итерации по строкам
	if rows.Err() != nil {
		return nil, fmt.Errorf("ошибка при обработке строк: %w", rows.Err())
	}

	return students, nil
}
