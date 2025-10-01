package types

import (
	"database/sql"
	"encoding/json"
	"errors"
)

// NullString - расширение sql.NullString для работы с JSON
// @typedef NullString
// @type object
// @description Тип для работы с sql.NullString в JSON.
// @properties
//
//	value:
//	  type: string
//	  description: Значение строки (если оно присутствует)
//	valid:
//	  type: boolean
//	  description: Флаг, который указывает, является ли значение действительным
type NullString struct {
	sql.NullString
}

// Реализация интерфейса json.Unmarshaler для корректной работы с JSON
func (ns *NullString) UnmarshalJSON(data []byte) error {
	var s *string
	if err := json.Unmarshal(data, &s); err != nil {
		return errors.New("ошибка декодирования: значение должно быть строкой или null")
	}
	if s != nil {
		ns.String = *s
		ns.Valid = true
	} else {
		ns.Valid = false
	}
	return nil
}

// Реализация интерфейса json.Marshaler для корректной сериализации в JSON
func (ns NullString) MarshalJSON() ([]byte, error) {
	if ns.Valid {
		return json.Marshal(ns.String)
	}
	return json.Marshal(nil)
}
