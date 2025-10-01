package models

import "time"

type semester struct {
	id           rune
	year         time.Time
	semesterPart bool // 0 - autumn  1 - spring
}
