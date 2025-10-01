package models

import "time"

type labWork struct {
	id                rune
	studentId         rune
	labDate           time.Time
	teacherId         rune
	labGrade          int8
	complited_lab_num int8
	labId             rune
	wasPresent        bool
}
