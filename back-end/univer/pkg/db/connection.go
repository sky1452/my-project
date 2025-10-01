package db

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Структура для хранения подключения к базе данных
type Connection struct {
	Pool *pgxpool.Pool
}

// Функция для подключения к базе данных
func ConnectDB(connectionString string) (*Connection, error) {
	// Настройка пула подключений
	dbpool, err := pgxpool.New(context.Background(), connectionString)
	if err != nil {
		return nil, fmt.Errorf("не удалось подключиться к базе данных: %w", err)
	}
	err = dbpool.Ping(context.Background())
	if err != nil {
		return nil, fmt.Errorf("ошибка при проверке подключения к базе данных: %w", err)
	}
	return &Connection{Pool: dbpool}, nil
}

// Функция закрытия подключения
func (conn *Connection) Close() {
	conn.Pool.Close()
}
