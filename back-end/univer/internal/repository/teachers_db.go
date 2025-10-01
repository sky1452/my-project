package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"univer/models"
)

func mapRowToTeacher(row pgx.Row) (*models.Teacher, error) {
	var t models.Teacher
	var propertiesData []byte // Поле для хранения JSON-данных

	// Сканируем поля из строки
	err := row.Scan(&t.Id, &t.Name, &t.Email, &propertiesData)
	if err != nil {
		return nil, fmt.Errorf("ошибка сканирования строки: %w", err)
	}

	// Декодируем поле properties в структуру TeacherProperties
	var teacherProperties models.TeacherProperties
	if len(propertiesData) > 0 {
		err = json.Unmarshal(propertiesData, &teacherProperties)
		if err != nil {
			return nil, fmt.Errorf("ошибка декодирования properties: %w", err)
		}
	}
	t.Properties = teacherProperties

	return &t, nil
}

func FetchAllTeachers(dbpool *pgxpool.Pool) ([]*models.Teacher, error) {
	// Выполняем SQL-запрос
	rows, err := dbpool.Query(context.Background(),
		"SELECT user_id, name, email, properties FROM public.users WHERE role_id=2")
	if err != nil {
		return nil, fmt.Errorf("ошибка выполнения запроса: %w", err)
	}
	defer rows.Close()

	var teachers []*models.Teacher

	// Итерируем по результатам и мапим их в структуры Teacher
	for rows.Next() {
		teacher, err := mapRowToTeacher(rows)
		if err != nil {
			return nil, err
		}
		teachers = append(teachers, teacher)
	}

	// Проверяем на наличие ошибок при итерации по строкам
	if rows.Err() != nil {
		return nil, fmt.Errorf("ошибка при обработке строк: %w", rows.Err())
	}

	return teachers, nil
}

// Функция для запроса преподавателя по ID
func FetchTeacherFromDB(dbpool *pgxpool.Pool, id int) (*models.Teacher, error) {
	// Выполняем SQL-запрос с параметром
	row := dbpool.QueryRow(context.Background(), ""+
		"SELECT user_id, name, email, properties FROM public.users WHERE user_id=$1 AND role_id<=2", id)

	// Используем ранее созданную функцию для маппинга данных
	teacher, err := mapRowToTeacher(row)
	if err != nil {
		return nil, fmt.Errorf("ошибка получения преподавателя по ID: %w", err)
	}
	return teacher, nil
}

// Функция для добавления нового преподавателя в БД
func AddTeacher(dbpool *pgxpool.Pool, teacher *models.AddTeacherRequest) error {
	// Преобразуем поле Properties в JSON
	propertiesData, err := json.Marshal(teacher.Properties)
	if err != nil {
		return fmt.Errorf("ошибка кодирования properties в JSON: %w", err)
	}

	// Выполняем SQL-запрос на добавление
	_, err = dbpool.Exec(
		context.Background(),
		`INSERT INTO public.users (name, email, properties, role_id) 
		 VALUES ($1, $2, $3, 2)`, // Мы предполагаем, что role_id=2 — это преподаватель
		teacher.Name,
		teacher.Email.String, // Используем String, так как это sql.NullString
		propertiesData,
	)

	if err != nil {
		return fmt.Errorf("ошибка выполнения запроса на добавление преподавателя: %w", err)
	}

	return nil
}
