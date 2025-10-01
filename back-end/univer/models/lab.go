package models

import "univer/types"

// @Description Структура для лабораторной работы
type Lab struct {
	Id           rune             `json:"id"`
	Name         string           `json:"lab_name"`
	Count        int              `json:"lab_count"`
	Year         types.NullString `json:"lab_year"`
	DefaultGrade int              `json:"default_grade"`
	Properties   LabProperties    `json:"lab_properties"` // Можно использовать конкретный тип
}

// @Description Структура для свойств лабораторной работы
type LabProperties struct {
	LabsGradeMax           map[string]interface{} `json:"labs_grade_max"`
	ReducePointsEachModule bool                   `json:"reduce_points_each_module"`
}
