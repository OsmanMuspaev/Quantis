package handlers

import (
	"bot_logic/main_module"
	"bot_logic/storage"
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
)

func getToken(r *http.Request) (string, string) {
	chatID := r.URL.Query().Get("chat_id")
	if chatID == "" {
		return "", ""
	}
	userData, _ := storage.AuthRedis.HGetAll(context.Background(), chatID).Result()
	return userData["access_token"], chatID
}

// GetCoursesHandler returns all courses.
func GetCoursesHandler(w http.ResponseWriter, r *http.Request) {
	token, _ := getToken(r)

	courses, err := main_module.GetCourses(token)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(courses)
}

// GetCourseHandler returns a single course by ID.
func GetCourseHandler(w http.ResponseWriter, r *http.Request) {
	token, _ := getToken(r)
	courseID, _ := strconv.Atoi(r.URL.Query().Get("course_id"))

	course, err := main_module.GetCourseByID(token, courseID)
	if err != nil {
		http.Error(w, "Course not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(course)
}

// CreateCourseHandler creates a new course.
func CreateCourseHandler(w http.ResponseWriter, r *http.Request) {
	token, _ := getToken(r)

	var body struct {
		Title       string `json:"title"`
		Description string `json:"description"`
	}
	json.NewDecoder(r.Body).Decode(&body)

	courseID, err := main_module.CreateCourse(token, body.Title, body.Description)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	json.NewEncoder(w).Encode(map[string]int{"course_id": courseID})
}

// UpdateCourseHandler updates an existing course.
func UpdateCourseHandler(w http.ResponseWriter, r *http.Request) {
	token, _ := getToken(r)
	courseID, _ := strconv.Atoi(r.URL.Query().Get("course_id"))

	var body struct {
		Title       string `json:"title"`
		Description string `json:"description"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "Invalid JSON body", http.StatusBadRequest)
		return
	}

	if err := main_module.UpdateCourse(token, courseID, body.Title, body.Description); err != nil {
		log.Printf("Update error: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// DeleteCourseHandler deletes a course.
func DeleteCourseHandler(w http.ResponseWriter, r *http.Request) {
	token, _ := getToken(r)
	courseID, _ := strconv.Atoi(r.URL.Query().Get("course_id"))

	if err := main_module.DeleteCourse(token, courseID); err != nil {
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// JoinCourseHandler enrolls a user in a course.
func JoinCourseHandler(w http.ResponseWriter, r *http.Request) {
	token, _ := getToken(r)

	var body struct {
		CourseID int    `json:"course_id"`
		UserID   string `json:"user_id"`
	}
	json.NewDecoder(r.Body).Decode(&body)

	if err := main_module.JoinCourse(token, body.CourseID, body.UserID); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// LeaveCourseHandler unenrolls a user from a course.
func LeaveCourseHandler(w http.ResponseWriter, r *http.Request) {
	token, _ := getToken(r)

	var body struct {
		CourseID int    `json:"course_id"`
		UserID   string `json:"user_id"`
	}
	json.NewDecoder(r.Body).Decode(&body)

	if err := main_module.LeaveCourse(token, body.CourseID, body.UserID); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
