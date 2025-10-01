package repository

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"time"
	"univer/models"
)

func mapRowToPara(row pgx.Row) (*models.Para, error) {
	var para models.Para

	// Сканируем поля из строки
	err := row.Scan(
		&para.Id, &para.DayOfWeek, &para.ParaNum, &para.DivisionIntoWeek, &para.WeekType,
		&para.DivisionIntoGroups, &para.CabinetIds, &para.TeacherIds,
		&para.LabId, &para.LabName, &para.LectureId, &para.LectureName,
		&para.PracticId, &para.PracticName, &para.TeacherNames, &para.CabinetNames, &para.GroupName)
	if err != nil {
		return nil, fmt.Errorf("Нет строк для сканирования: %w", err)
	}

	return &para, nil
}

func GetParaNumByTime(dbpool *pgxpool.Pool, date time.Time) (int, error) {
	var paraNum int

	query := `SELECT para_time_id FROM para_time WHERE para_time_start <= $1::TIME AND para_time_end >= $1::TIME LIMIT 1`

	// Форматируем дату в HH:MM:SS для корректного сравнения с TIME
	timeString := date.Format("15:04:05")

	err := dbpool.QueryRow(context.Background(), query, timeString).Scan(&paraNum)

	//если вернулся id >  500 - значит это перемена
	if paraNum > 500 {
		paraNum = 0
	}
	if err != nil {
		if err == pgx.ErrNoRows {
			return 0, nil // 0 означает, что пара не найдена
		}
		return 0, fmt.Errorf("ошибка при запросе номера пары: %w", err)
	}

	return paraNum, nil
}

func GetCurrentParaById(dbpool *pgxpool.Pool, id rune, date time.Time) (*models.Para, error) {

	weekday := int(date.Weekday())
	//timeString := date.Format("15:04:05")

	currentPara, err := GetParaNumByTime(dbpool, date)

	if err != nil {
		return nil, fmt.Errorf("ошибка получения номера текущей пары", err)
	}
	if currentPara == 0 {
		return nil, fmt.Errorf("сейчас перемена")
	}
	//query := `SELECT * FROM para WHERE day_of_week = $1 AND teachers_id @> ARRAY[$2]::BIGINT[] AND para_num = (
	//SELECT para_time_id FROM para_time WHERE para_time_start <= $3::TIME AND para_time_end >= $3::TIME LIMIT 1)`

	query := `SELECT 
    para.id,
    para.day_of_week, para.para_num, para.division_into_weeks, 
    para.para_week_type, para.division_into_groups,
    para.para_cabinet_ids, para.teachers_id,
    
    para.lab_id, lab.name as lab_name, 
    para.lecture_id, lecture.name as lecture_name, 
    para.practic_id, practic.name as practic_name,
    
    -- DISTINCT убирает дубли преподавателей
    COALESCE(array_agg(DISTINCT USER1.name), '{}') AS teachers_name,  
    -- DISTINCT убирает дубли кабинетов
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
		

	GROUP BY para.id, 
	         para.group_id, 
	         para.day_of_week, 
	         para.para_num, 
	         group1.name, 
	         lab.name, 
	         lecture.name, 
	         practic.name

ORDER BY para.group_id, para.day_of_week, para.para_num ASC;
`
	row := dbpool.QueryRow(context.Background(), query, weekday, id, currentPara)

	// Используем ранее созданную функцию для маппинга данных
	para, err := mapRowToPara(row)
	if err != nil {
		return nil, fmt.Errorf("Не надено пары для преподавателя с id:%d  %w", id, err)
	}
	return para, nil
}

func GetTodayParasById(dbpool *pgxpool.Pool, id int, role string, date time.Time) ([]*models.Para, error) {
	weekday := int(date.Weekday())
	var paras []*models.Para

	query := `SELECT 
        para.id,
        para.day_of_week, para.para_num, para.division_into_weeks, 
        para.para_week_type, para.division_into_groups,
        para.para_cabinet_ids, para.teachers_id,
        
        para.lab_id, lab.name as lab_name, 
        para.lecture_id, lecture.name as lecture_name, 
        para.practic_id, practic.name as practic_name,

        -- DISTINCT убирает дубли преподавателей
        COALESCE(array_agg(DISTINCT USER1.name), '{}') AS teachers_name,  
        -- DISTINCT убирает дубли кабинетов
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
      AND (
          CASE 
              WHEN $3 = 'teacher' THEN para.teachers_id @> ARRAY[$2]::BIGINT[]
              WHEN $3 = 'student' THEN group1.name = (
                  SELECT u.properties->>'group' FROM public.users u WHERE u.user_id = $2 LIMIT 1
              )
              ELSE false
          END
      )

    GROUP BY para.id, 
             para.group_id, 
             para.day_of_week, 
             para.para_num, 
             group1.name, 
             lab.name, 
             lecture.name, 
             practic.name

        ORDER BY para.day_of_week, para.para_num ASC;
    `

	rows, err := dbpool.Query(context.Background(), query, weekday, id, role)
	if err != nil {
		return nil, fmt.Errorf("ошибка выполнения запроса: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		para, err := mapRowToPara(rows)
		if err != nil {
			return nil, err
		}
		paras = append(paras, para)
	}

	if err != nil {
		return nil, fmt.Errorf("Не найдено пар на сегодня для %s с id:%d  %w", role, id, err)
	}
	return paras, nil
}
func GetTodayParasByGroup(dbpool *pgxpool.Pool, group *models.Group, date time.Time) ([]*models.Para, error) {
	weekday := int(date.Weekday())
	var paras []*models.Para

	query := `SELECT 
        para.id,
        para.day_of_week, para.para_num, para.division_into_weeks, 
        para.para_week_type, para.division_into_groups,
        para.para_cabinet_ids, para.teachers_id,
        
        para.lab_id, lab.name as lab_name, 
        para.lecture_id, lecture.name as lecture_name, 
        para.practic_id, practic.name as practic_name,

        -- DISTINCT убирает дубли преподавателей
        COALESCE(array_agg(DISTINCT USER1.name), '{}') AS teachers_name,  
        -- DISTINCT убирает дубли кабинетов
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
      AND group_id = $2

    GROUP BY para.id, 
             para.group_id, 
             para.day_of_week, 
             para.para_num, 
             group1.name, 
             lab.name, 
             lecture.name, 
             practic.name

    ORDER BY para.group_id, para.day_of_week, para.para_num ASC;
    `

	rows, err := dbpool.Query(context.Background(), query, weekday, group.Id)
	if err != nil {
		return nil, fmt.Errorf("ошибка выполнения запроса: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		para, err := mapRowToPara(rows)
		if err != nil {
			return nil, err
		}
		paras = append(paras, para)
	}

	if err != nil {
		return nil, fmt.Errorf("Не найдено пар на сегодня для %s  %w", group.Name, err)
	}
	return paras, nil
}

func GetScheduleById(dbpool *pgxpool.Pool, id int, role string) ([]*models.Para, error) {
	var paras []*models.Para

	query := `SELECT 
        para.id,
        para.day_of_week, para.para_num, para.division_into_weeks, 
        para.para_week_type, para.division_into_groups,
        para.para_cabinet_ids, para.teachers_id,
        
        para.lab_id, lab.name as lab_name, 
        para.lecture_id, lecture.name as lecture_name, 
        para.practic_id, practic.name as practic_name,

        -- DISTINCT убирает дубли преподавателей
        COALESCE(array_agg(DISTINCT USER1.name), '{}') AS teachers_name,  
        -- DISTINCT убирает дубли кабинетов
        COALESCE(array_agg(DISTINCT CABINET.name), '{}') AS cabinets_name,

        group1.name AS group_name  
    FROM public.para para
    LEFT JOIN public.users USER1 ON USER1.user_id = ANY(para.teachers_id)
    LEFT JOIN public.lab lab ON para.lab_id = lab.id
    LEFT JOIN public.lecture lecture ON para.lecture_id = lecture.id
    LEFT JOIN public.practic practic ON para.practic_id = practic.id
    LEFT JOIN public.cabinet CABINET ON CABINET.id = ANY(para.para_cabinet_ids)
    INNER JOIN public.group group1 ON para.group_id = group1.id

    WHERE CASE 
              WHEN $2 = 'teacher' THEN para.teachers_id @> ARRAY[$1]::BIGINT[]
              WHEN $2 = 'student' THEN group1.name = (
                  SELECT u.properties->>'group' FROM public.users u WHERE u.user_id = $1 LIMIT 1
              )
              ELSE false
          END

    GROUP BY para.id, 
             para.group_id, 
             para.day_of_week, 
             para.para_num, 
             group1.name, 
             lab.name, 
             lecture.name, 
             practic.name

    ORDER BY para.group_id, para.day_of_week, para.para_num ASC;
    `

	rows, err := dbpool.Query(context.Background(), query, id, role)
	if err != nil {
		return nil, fmt.Errorf("ошибка выполнения запроса: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		para, err := mapRowToPara(rows)
		if err != nil {
			return nil, err
		}
		paras = append(paras, para)
	}

	if err != nil {
		return nil, fmt.Errorf("Не найдено пар для %d  %w", id, err)
	}
	return paras, nil
}
func GetScheduleByGroup(dbpool *pgxpool.Pool, group *models.Group) ([]*models.Para, error) {
	var paras []*models.Para

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
    WHERE group_id = $1
    GROUP BY para.id, para.group_id, para.day_of_week, para.para_num, 
             group1.name, lab.name, lecture.name, practic.name
    ORDER BY para.group_id, para.day_of_week, para.para_num ASC;`

	rows, err := dbpool.Query(context.Background(), query, group.Id)
	if err != nil {
		return nil, fmt.Errorf("ошибка выполнения запроса: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		para, err := mapRowToPara(rows)
		if err != nil {
			fmt.Printf("ошибка: %w", err)
			return nil, fmt.Errorf("ошибка маппинга строки: %w", err)
		}
		
		paras = append(paras, para)
	}

	//if err := rows.Err(); err != nil {
	//	return nil, fmt.Errorf("ошибка при обработке строк: %w", err)
	//}
	//
	//// Если ни одной строки не найдено, можно вернуть пустой срез, а не nil.
	//if paras == nil {
	//	paras = []*models.Para{}
	//}

	return paras, nil
}
