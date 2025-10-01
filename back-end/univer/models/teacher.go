package models

import "univer/types"

// Teacher представляет модель преподавателя
// @Description Модель, описывающая преподавателя
type Teacher struct {
	Id         rune             `json:"id" example:"1"`                                               // Уникальный идентификатор преподавателя
	Name       string           `json:"name" example:"Иван Иванов"`                                   // Имя преподавателя
	Email      types.NullString `json:"email" swaggertype:"string" example:"ivan.ivanov@example.com"` // Email преподавателя (может отсутствовать)
	Properties interface{}      `json:"properties" swaggertype:"array,object"`                        // Дополнительные свойства преподавателя
}

// TeacherProperties описывает свойства преподавателя
// @Description Модель, содержащая дополнительные характеристики преподавателя
type TeacherProperties struct {
	Position    string `json:"position" example:"Профессор"`              // Должность преподавателя
	Departament string `json:"departament" example:"Кафедра информатики"` // Кафедра, к которой относится преподаватель
}
