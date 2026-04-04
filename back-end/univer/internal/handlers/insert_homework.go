package handlers

import (

	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"path/filepath"
	"strconv"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
)

var bucketName = "cloud-sky-pirson"

func (h *Handler) UploadSubmission(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		log.Println("ParseMultipartForm error:", err)
		http.Error(w, "invalid form", http.StatusBadRequest)
		return
	}

	studentID, err := strconv.Atoi(r.FormValue("student_id"))
	if err != nil {
		log.Println("invalid student_id:", err)
		http.Error(w, "invalid student_id", http.StatusBadRequest)
		return
	}

	taskID, err := strconv.Atoi(r.FormValue("task_id"))
	if err != nil {
		log.Println("invalid task_id:", err)
		http.Error(w, "invalid task_id", http.StatusBadRequest)
		return
	}

	comment := r.FormValue("comment")

	disciplineID, err := strconv.Atoi(r.FormValue("discipline_id"))
	if err != nil {
		log.Println("invalid discipline_id:", err)
		http.Error(w, "invalid discipline_id", http.StatusBadRequest)
		return
	}

	// 🔥 1. Получаем properties пользователя
	var propertiesJSON []byte

err = h.DB.QueryRow(r.Context(), `
	SELECT properties
	FROM users
	WHERE user_id = $1
`, studentID).Scan(&propertiesJSON)

	if err != nil {
		log.Println("get user properties error:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// 🔥 2. Достаём group из JSON
	var props struct {
		Group string `json:"group"`
	}

	err = json.Unmarshal(propertiesJSON, &props)
	if err != nil {
		log.Println("json unmarshal error:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if props.Group == "" {
		http.Error(w, "group not found in user properties", http.StatusBadRequest)
		return
	}

	// 🔥 3. Ищем group_id по имени группы
	var groupID int

	err = h.DB.QueryRow(r.Context(), `
		SELECT id
		FROM "group"
		WHERE name = $1
	`, props.Group).Scan(&groupID)

	if err != nil {
		log.Println("group not found:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// 🔥 4. Вставляем submission с group_id
	var submissionID int

	err = h.DB.QueryRow(r.Context(), `
		INSERT INTO submissions (student_id, task_id, discipline_id, group_id, comment, created_at)
		VALUES ($1, $2, $3, $4, $5, NOW())
		RETURNING id
	`, studentID, taskID, disciplineID, groupID, comment).Scan(&submissionID)

	if err != nil {
		log.Println("DB insert submission error:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// 🔥 5. Файлы
	files := r.MultipartForm.File["files"]

	var urls []string

	for _, file := range files {
		f, err := file.Open()
		if err != nil {
			log.Println("file open error:", err)
			continue
		}

		key := fmt.Sprintf("submissions/%d_%s", time.Now().UnixNano(), filepath.Base(file.Filename))

		_, err = h.S3.PutObject(&s3.PutObjectInput{
			Bucket:      aws.String(bucketName),
			Key:         aws.String(key),
			Body:        f,
			ContentType: aws.String(file.Header.Get("Content-Type")),
		})

		f.Close()

		if err != nil {
			log.Println("S3 upload error:", err)
			continue
		}

		url := fmt.Sprintf("https://%s.storage.yandexcloud.net/%s", bucketName, key)
		urls = append(urls, url)

		_, err = h.DB.Exec(r.Context(), `
			INSERT INTO submission_files (submission_id, file_url)
			VALUES ($1, $2)
		`, submissionID, url)

		if err != nil {
			log.Println("DB insert file error:", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, `{"submission_id": %d, "files_count": %d}`, submissionID, len(urls))
}