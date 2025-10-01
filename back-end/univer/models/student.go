package models

import "univer/types"

type Student struct {
	Id         rune             `json:"id"`
	Name       string           `json:"username"`
	Email      types.NullString `json:"email"`
	Properties interface{}      `json:"properties"`
}

type StudentProperties struct {
	Group string `json:"group"`
}
