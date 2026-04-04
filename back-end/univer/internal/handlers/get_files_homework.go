package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
)

type SubmissionFileResponse struct {
	FileURL  string `json:"file_url"`
	FileName string `json:"file_name"`
}

type StudentSubmissionFilesResponse struct {
	Success   bool                     `json:"success"`
	TaskID    int                      `json:"taskId"`
	UserID    int                      `json:"userId"`
	StudentID int                      `json:"studentId"`
	Comment   string                   `json:"comment"`
	Files     []SubmissionFileResponse `json:"files"`
}

func extractFileNameFromURL(fileURL string) string {
	parts := strings.Split(fileURL, "/")
	if len(parts) == 0 {
		return ""
	}

	lastPart := parts[len(parts)-1]
	if lastPart == "" {
		return ""
	}

	underscoreIndex := strings.Index(lastPart, "_")
	if underscoreIndex == -1 || underscoreIndex+1 >= len(lastPart) {
		return lastPart
	}

	return lastPart[underscoreIndex+1:]
}

func (h *Handler) GetStudentSubmissionFilesHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	ctx := r.Context()
	vars := mux.Vars(r)

	taskIDStr := vars["taskId"]
	userIDStr := vars["userId"]

	if taskIDStr == "" || userIDStr == "" {
		http.Error(w, "taskId and userId are required", http.StatusBadRequest)
		return
	}

	taskID, err := strconv.Atoi(taskIDStr)
	if err != nil {
		http.Error(w, "Некорректный taskId", http.StatusBadRequest)
		return
	}

	studentID, err := strconv.Atoi(userIDStr)
	if err != nil {
		http.Error(w, "Некорректный userId", http.StatusBadRequest)
		return
	}

	resp := StudentSubmissionFilesResponse{
		Success:   true,
		TaskID:    taskID,
		UserID:    studentID,
		StudentID: studentID,
		Comment:   "",
		Files:     []SubmissionFileResponse{},
	}

	rows, err := h.DB.Query(ctx, `
		SELECT s.id, COALESCE(s.comment, '')
		FROM submissions s
		WHERE s.task_id = $1 AND s.student_id = $2
		ORDER BY s.id DESC
	`, taskID, studentID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Ошибка при получении submissions: %v", err), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	firstCommentSet := false

	for rows.Next() {
		var submissionID int
		var comment string

		if err := rows.Scan(&submissionID, &comment); err != nil {
			http.Error(w, fmt.Sprintf("Ошибка при чтении submissions: %v", err), http.StatusInternalServerError)
			return
		}

		if !firstCommentSet && strings.TrimSpace(comment) != "" {
			resp.Comment = comment
			firstCommentSet = true
		}

		fileRows, err := h.DB.Query(ctx, `
			SELECT file_url
			FROM submission_files
			WHERE submission_id = $1
			ORDER BY id
		`, submissionID)
		if err != nil {
			http.Error(w, fmt.Sprintf("Ошибка при получении файлов: %v", err), http.StatusInternalServerError)
			return
		}

		for fileRows.Next() {
			var fileURL string

			if err := fileRows.Scan(&fileURL); err != nil {
				fileRows.Close()
				http.Error(w, fmt.Sprintf("Ошибка при чтении файла: %v", err), http.StatusInternalServerError)
				return
			}

			resp.Files = append(resp.Files, SubmissionFileResponse{
				FileURL:  fileURL,
				FileName: extractFileNameFromURL(fileURL),
			})
		}

		if err := fileRows.Err(); err != nil {
			fileRows.Close()
			http.Error(w, fmt.Sprintf("Ошибка при обработке файлов: %v", err), http.StatusInternalServerError)
			return
		}

		fileRows.Close()
	}

	if err := rows.Err(); err != nil {
		http.Error(w, fmt.Sprintf("Ошибка при обработке submissions: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}