package models

import "univer/types"

// @Description описывает входные данные для добавления преподавателя
type AddTeacherRequest struct {
	Name       string           `json:"name" example:"Иван Иванов"`
	Email      types.NullString `json:"email,omitempty" swaggertype:"string" example:"ivan.ivanov@example.com"`
	Properties interface{}      `json:"properties" swaggertype:"array,object"`
}
