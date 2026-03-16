package repository

import (
	"context"
	"fmt"
	"strings"
	"time"
	"univer/models"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

func slugify(text string) string {
	replacer := strings.NewReplacer(
		"ё", "e", "й", "i", "ц", "ts", "у", "u",
		"к", "k", "е", "e", "н", "n", "г", "g",
		"ш", "sh", "щ", "sch", "з", "z", "х", "h",
		"ъ", "", "ф", "f", "ы", "y", "в", "v",
		"а", "a", "п", "p", "р", "r", "о", "o",
		"л", "l", "д", "d", "ж", "zh", "э", "e",
		"я", "ya", "ч", "ch", "с", "s", "м", "m",
		"и", "i", "т", "t", "ь", "", "б", "b",
		"ю", "yu",
	)

	text = strings.ToLower(text)
	text = replacer.Replace(text)
	text = strings.ReplaceAll(text, " ", "-")

	return text
}

func mapRowToPara(row pgx.Row, dbpool *pgxpool.Pool) (*models.Para, error) {
	var para models.Para

	err := row.Scan(
		&para.Id, &para.DayOfWeek, &para.ParaNum, &para.DivisionIntoWeek, &para.WeekType,
		&para.DivisionIntoGroups, &para.CabinetIds, &para.TeacherIds,
		&para.LabId, &para.LabName, &para.LectureId, &para.LectureName,
		&para.PracticId, &para.PracticName, &para.TeacherNames, &para.CabinetNames, &para.GroupName)
	if err != nil {
		return nil, fmt.Errorf("Нет строк для сканирования: %w", err)
	}

	var disciplineName string

	if para.LabName.Valid {
		disciplineName = para.LabName.String
	} else if para.LectureName.Valid {
		disciplineName = para.LectureName.String
	} else if para.PracticName.Valid {
		disciplineName = para.PracticName.String
	}

	if disciplineName != "" {
		para.DisciplineSlug = slugify(disciplineName)

		var disciplineId int64
		err := dbpool.QueryRow(
			context.Background(),
			`SELECT id FROM academic_subject WHERE name = $1 LIMIT 1`,
			disciplineName,
		).Scan(&disciplineId)

		if err == nil {
			para.DisciplineId = disciplineId
		}
	}

	return &para, nil
}

func GetParaNumByTime(dbpool *pgxpool.Pool, date time.Time) (int, error) {
	var paraNum int

	query := `SELECT para_time_id FROM para_time WHERE para_time_start <= $1::TIME AND para_time_end >= $1::TIME LIMIT 1`
	timeString := date.Format("15:04:05")

	err := dbpool.QueryRow(context.Background(), query, timeString).Scan(&paraNum)

	if paraNum > 500 {
		paraNum = 0
	}
	if err != nil {
		if err == pgx.ErrNoRows {
			return 0, nil
		}
		return 0, fmt.Errorf("ошибка при запросе номера пары: %w", err)
	}

	return paraNum, nil
}

func GetCurrentParaById(dbpool *pgxpool.Pool, id rune, date time.Time) (*models.Para, error) {
	weekday := int(date.Weekday())

	currentPara, err := GetParaNumByTime(dbpool, date)
	if err != nil {
		return nil, fmt.Errorf("ошибка получения номера текущей пары")
	}
	if currentPara == 0 {
		return nil, fmt.Errorf("сейчас перемена")
	}

	query := `SELECT 
    para.id,
    para.day_of_week, para.para_num, para.division_into_weeks, 
    para.para_week_type, para.division_into_groups,
    para.para_cabinet_ids, para.teachers_id,
    para.lab_id, lab.name as lab_name, 
    para.lecture_id, lecture.name as lecture_name, 
    para.practic_id, practic.name as practic_name,
    COALESCE(array_agg(DISTINCT USER1.name), '{}') AS teachers_name,  
    COALESCE(array_agg(DISTINCT CABINET.name), '{}') AS cabinets_name,
    group1.name AS group_name  
	FROM public.para para
	LEFT JOIN public.users USER1 ON USER1.user_id = ANY(para.teachers_id)
	LEFT JOIN public.lab lab ON para.lab_id = lab.id
	LEFT JOIN public.lecture lecture ON para.lecture_id = lecture.id
	LEFT JOIN public.practic practic ON para.practic_id = practic.id
	LEFT JOIN public.cabinet CABINET ON CABINET.id = ANY(para.para_cabinet_ids)
	INNER JOIN public.group group1 ON para.group_id = group1.id
	WHERE day_of_week = $1 
	  AND teachers_id @> ARRAY[$2]::BIGINT[] 
	  AND para_num = $3
	GROUP BY para.id, para.group_id, para.day_of_week, para.para_num, 
	         group1.name, lab.name, lecture.name, practic.name
	ORDER BY para.group_id, para.day_of_week, para.para_num ASC;`

	row := dbpool.QueryRow(context.Background(), query, weekday, id, currentPara)

	para, err := mapRowToPara(row, dbpool)
	if err != nil {
		return nil, fmt.Errorf("Не найдено пары для преподавателя с id:%d  %w", id, err)
	}
	return para, nil
}
func GetParasByTeacherAndDay(
	dbpool *pgxpool.Pool,
	teacherId rune,
	dayOfWeek int,
) ([]*models.Para, error) {

	query := `SELECT 
    para.id,
    para.day_of_week, para.para_num, para.division_into_weeks, 
    para.para_week_type, para.division_into_groups,
    para.para_cabinet_ids, para.teachers_id,
    para.lab_id, lab.name as lab_name, 
    para.lecture_id, lecture.name as lecture_name, 
    para.practic_id, practic.name as practic_name,
    COALESCE(array_agg(DISTINCT USER1.name), '{}') AS teachers_name,  
    COALESCE(array_agg(DISTINCT CABINET.name), '{}') AS cabinets_name,
    group1.name AS group_name  
	FROM public.para para
	LEFT JOIN public.users USER1 ON USER1.user_id = ANY(para.teachers_id)
	LEFT JOIN public.lab lab ON para.lab_id = lab.id
	LEFT JOIN public.lecture lecture ON para.lecture_id = lecture.id
	LEFT JOIN public.practic practic ON para.practic_id = practic.id
	LEFT JOIN public.cabinet CABINET ON CABINET.id = ANY(para.para_cabinet_ids)
	INNER JOIN public.group group1 ON para.group_id = group1.id
	WHERE day_of_week = $1 
	  AND teachers_id @> ARRAY[$2]::BIGINT[]
	GROUP BY para.id, para.group_id, para.day_of_week, para.para_num, 
	         group1.name, lab.name, lecture.name, practic.name
	ORDER BY para.para_num ASC;`

	rows, err := dbpool.Query(context.Background(), query, dayOfWeek, teacherId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var paras []*models.Para

	for rows.Next() {
		para, err := mapRowToPara(rows, dbpool)
		if err != nil {
			return nil, err
		}
		paras = append(paras, para)
	}

	return paras, nil
}
func GetParasByGroupAndDay(
	dbpool *pgxpool.Pool,
	groupId int64,
	dayOfWeek int,
) ([]*models.Para, error) {

	query := `SELECT 
    para.id,
    para.day_of_week, para.para_num, para.division_into_weeks, 
    para.para_week_type, para.division_into_groups,
    para.para_cabinet_ids, para.teachers_id,
    para.lab_id, lab.name as lab_name, 
    para.lecture_id, lecture.name as lecture_name, 
    para.practic_id, practic.name as practic_name,
    COALESCE(array_agg(DISTINCT USER1.name), '{}') AS teachers_name,  
    COALESCE(array_agg(DISTINCT CABINET.name), '{}') AS cabinets_name,
    group1.name AS group_name  
	FROM public.para para
	LEFT JOIN public.users USER1 ON USER1.user_id = ANY(para.teachers_id)
	LEFT JOIN public.lab lab ON para.lab_id = lab.id
	LEFT JOIN public.lecture lecture ON para.lecture_id = lecture.id
	LEFT JOIN public.practic practic ON para.practic_id = practic.id
	LEFT JOIN public.cabinet CABINET ON CABINET.id = ANY(para.para_cabinet_ids)
	INNER JOIN public.group group1 ON para.group_id = group1.id
	WHERE day_of_week = $1 
	  AND para.group_id = $2
	GROUP BY para.id, para.group_id, para.day_of_week, para.para_num, 
	         group1.name, lab.name, lecture.name, practic.name
	ORDER BY para.para_num ASC;`

	rows, err := dbpool.Query(context.Background(), query, dayOfWeek, groupId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var paras []*models.Para

	for rows.Next() {
		para, err := mapRowToPara(rows, dbpool)
		if err != nil {
			return nil, err
		}
		paras = append(paras, para)
	}

	return paras, nil
}
func GetParaById(
	dbpool *pgxpool.Pool,
	paraId rune,
) (*models.Para, error) {

	query := `SELECT 
    para.id,
    para.day_of_week, para.para_num, para.division_into_weeks, 
    para.para_week_type, para.division_into_groups,
    para.para_cabinet_ids, para.teachers_id,
    para.lab_id, lab.name as lab_name, 
    para.lecture_id, lecture.name as lecture_name, 
    para.practic_id, practic.name as practic_name,
    COALESCE(array_agg(DISTINCT USER1.name), '{}') AS teachers_name,  
    COALESCE(array_agg(DISTINCT CABINET.name), '{}') AS cabinets_name,
    group1.name AS group_name  
	FROM public.para para
	LEFT JOIN public.users USER1 ON USER1.user_id = ANY(para.teachers_id)
	LEFT JOIN public.lab lab ON para.lab_id = lab.id
	LEFT JOIN public.lecture lecture ON para.lecture_id = lecture.id
	LEFT JOIN public.practic practic ON para.practic_id = practic.id
	LEFT JOIN public.cabinet CABINET ON CABINET.id = ANY(para.para_cabinet_ids)
	INNER JOIN public.group group1 ON para.group_id = group1.id
	WHERE para.id = $1
	GROUP BY para.id, para.group_id, para.day_of_week, para.para_num, 
	         group1.name, lab.name, lecture.name, practic.name
	LIMIT 1;`

	row := dbpool.QueryRow(context.Background(), query, paraId)

	return mapRowToPara(row, dbpool)
}
func GetAllParas(dbpool *pgxpool.Pool) ([]*models.Para, error) {

	query := `SELECT 
    para.id,
    para.day_of_week, para.para_num, para.division_into_weeks, 
    para.para_week_type, para.division_into_groups,
    para.para_cabinet_ids, para.teachers_id,
    para.lab_id, lab.name as lab_name, 
    para.lecture_id, lecture.name as lecture_name, 
    para.practic_id, practic.name as practic_name,
    COALESCE(array_agg(DISTINCT USER1.name), '{}') AS teachers_name,  
    COALESCE(array_agg(DISTINCT CABINET.name), '{}') AS cabinets_name,
    group1.name AS group_name  
	FROM public.para para
	LEFT JOIN public.users USER1 ON USER1.user_id = ANY(para.teachers_id)
	LEFT JOIN public.lab lab ON para.lab_id = lab.id
	LEFT JOIN public.lecture lecture ON para.lecture_id = lecture.id
	LEFT JOIN public.practic practic ON para.practic_id = practic.id
	LEFT JOIN public.cabinet CABINET ON CABINET.id = ANY(para.para_cabinet_ids)
	INNER JOIN public.group group1 ON para.group_id = group1.id
	GROUP BY para.id, para.group_id, para.day_of_week, para.para_num, 
	         group1.name, lab.name, lecture.name, practic.name
	ORDER BY para.day_of_week, para.para_num;`

	rows, err := dbpool.Query(context.Background(), query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var paras []*models.Para

	for rows.Next() {
		para, err := mapRowToPara(rows, dbpool)
		if err != nil {
			return nil, err
		}
		paras = append(paras, para)
	}

	return paras, nil
}
func GetTodayParasById(dbpool *pgxpool.Pool, id int, role string, date time.Time) ([]*models.Para, error) {
	weekday := int(date.Weekday())
	return GetParasByTeacherAndDay(dbpool, rune(id), weekday)
}

func GetTodayParasByGroup(dbpool *pgxpool.Pool, group *models.Group, date time.Time) ([]*models.Para, error) {
	weekday := int(date.Weekday())
	return GetParasByGroupAndDay(dbpool, int64(group.Id), weekday)
}

func GetScheduleById(dbpool *pgxpool.Pool, id int, role string) ([]*models.Para, error) {
	var paras []*models.Para

	for day := 0; day <= 6; day++ {
		dayParas, err := GetParasByTeacherAndDay(dbpool, rune(id), day)
		if err != nil {
			return nil, err
		}
		paras = append(paras, dayParas...)
	}

	return paras, nil
}

func GetScheduleByGroup(dbpool *pgxpool.Pool, group *models.Group) ([]*models.Para, error) {
	var paras []*models.Para

	for day := 0; day <= 6; day++ {
		dayParas, err := GetParasByGroupAndDay(dbpool, int64(group.Id), day)
		if err != nil {
			return nil, err
		}
		paras = append(paras, dayParas...)
	}

	return paras, nil
}