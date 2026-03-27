package handlers

import (
    "context"
    "encoding/json"
    "log"
    "net/http"
    "strconv"
    "time"

    "github.com/gorilla/mux"
)

type Homework1 struct {
    ID           int       `json:"id"`
    Title        string    `json:"title"`
    Description  string    `json:"description"`
    MaxScore     int       `json:"max_score"`
    Deadline     time.Time `json:"deadline"`
    DisciplineID int       `json:"discipline_id"`
}

func (h *Handler) GetStudentHomeworks(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)

    disciplineId, err := strconv.Atoi(vars["disciplineId"])
    if err != nil {
        http.Error(w, "invalid disciplineId", http.StatusBadRequest)
        return
    }

    userId, err := strconv.Atoi(vars["userId"])
    if err != nil {
        http.Error(w, "invalid userId", http.StatusBadRequest)
        return
    }

    // Получаем имя группы студента из JSON-поля properties
    var groupName string
    err = h.DB.QueryRow(context.Background(),
        `SELECT properties->>'group' FROM users WHERE user_id = $1`,
        userId).Scan(&groupName)
    if err != nil {
        log.Println("Ошибка получения группы студента:", err)
        w.Header().Set("Content-Type", "application/json")
        w.WriteHeader(http.StatusInternalServerError)
        json.NewEncoder(w).Encode(map[string]string{
            "error": "Ошибка получения группы студента",
        })
        return
    }

    // Находим id группы по названию
    var groupId int
    err = h.DB.QueryRow(context.Background(),
        `SELECT id FROM "group" WHERE name = $1`, groupName).Scan(&groupId)
    if err != nil {
        log.Println("Ошибка получения id группы:", err)
        w.Header().Set("Content-Type", "application/json")
        w.WriteHeader(http.StatusInternalServerError)
        json.NewEncoder(w).Encode(map[string]string{
            "error": "Ошибка получения id группы",
        })
        return
    }

    // Получаем домашние задания по группе и дисциплине
    rows, err := h.DB.Query(context.Background(), `
        SELECT 
            id,
            title,
            description,
            max_score,
            deadline,
            discipline_id
        FROM homeworks
        WHERE group_id = $1 AND discipline_id = $2
        ORDER BY deadline ASC
    `, groupId, disciplineId)
    if err != nil {
        log.Println("Ошибка запроса к БД:", err)
        w.Header().Set("Content-Type", "application/json")
        w.WriteHeader(http.StatusInternalServerError)
        json.NewEncoder(w).Encode(map[string]string{
            "error": "Ошибка получения заданий",
        })
        return
    }
    defer rows.Close()

    var tasks []Homework1
    for rows.Next() {
        var hw Homework1
        if err := rows.Scan(
            &hw.ID,
            &hw.Title,
            &hw.Description,
            &hw.MaxScore,
            &hw.Deadline,
            &hw.DisciplineID,
        ); err != nil {
            log.Println("Ошибка сканирования строки:", err)
            w.Header().Set("Content-Type", "application/json")
            w.WriteHeader(http.StatusInternalServerError)
            json.NewEncoder(w).Encode(map[string]string{
                "error": "Ошибка обработки заданий",
            })
            return
        }
        tasks = append(tasks, hw)
    }

    if err := rows.Err(); err != nil {
        log.Println("Ошибка после rows.Next():", err)
        w.Header().Set("Content-Type", "application/json")
        w.WriteHeader(http.StatusInternalServerError)
        json.NewEncoder(w).Encode(map[string]string{
            "error": "Ошибка обработки заданий",
        })
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(map[string]interface{}{
        "tasks": tasks,
    })
}