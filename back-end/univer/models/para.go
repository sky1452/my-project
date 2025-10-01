package models

import "univer/types"

// @Description Модель для пары
// @docExpansion full
type Para struct {
	Id                 rune               `json:"id"`
	LectureId          types.NullString   `json:"lectureId" swaggertype:"string"`
	LectureName        types.NullString   `json:"lectureName" swaggertype:"string"`
	LabId              types.NullString   `json:"labId" swaggertype:"string"`
	LabName            types.NullString   `json:"labName" swaggertype:"string"`
	PracticId          types.NullString   `json:"practicId" swaggertype:"string"`
	PracticName        types.NullString   `json:"practicName" swaggertype:"string"`
	SemesterId         rune               `json:"semesterId"`
	TeacherIds         []rune             `json:"teacherIds"`
	TeacherNames       []types.NullString `json:"teacherNames"`
	DayOfWeek          int8               `json:"dayOfWeek"`
	ParaNum            int8               `json:"paraNum"`
	DivisionIntoWeek   types.NullString   `json:"divisionIntoWeek" swaggertype:"string"`
	WeekType           types.NullString   `json:"paraWeekType" swaggertype:"string"`
	GroupId            rune               `json:"groupId"`
	GroupName          string             `json:"groupName"`
	DivisionIntoGroups types.NullString   `json:"divisionIntoGroups" swaggertype:"string"`
	CabinetIds         []rune             `json:"cabinetIds"`
	CabinetNames       []types.NullString `json:"cabinetNames"`
}
