package main_module

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

// GetCourses returns all courses.
func GetCourses(token string) ([]Course, error) {
	req, _ := http.NewRequest("GET", baseURL+"/courses", nil)
	setAuthHeader(req, token)

	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get courses: status %d", resp.StatusCode)
	}

	var data CoursesResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, err
	}
	return data.Courses, nil
}

// GetCourseByID returns a single course by ID.
func GetCourseByID(token string, id int) (*Course, error) {
	url := fmt.Sprintf("%s/courses/%d", baseURL, id)
	req, _ := http.NewRequest("GET", url, nil)
	setAuthHeader(req, token)

	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("course %d not found or error: %d", id, resp.StatusCode)
	}

	var course Course
	json.NewDecoder(resp.Body).Decode(&course)
	return &course, nil
}

// CreateCourse creates a new course.
func CreateCourse(token string, title, description string) (int, error) {
	payload := CoursePayload{Title: title, Description: description}
	body, _ := json.Marshal(payload)

	req, _ := http.NewRequest("POST", baseURL+"/courses", bytes.NewBuffer(body))
	setAuthHeader(req, token)

	resp, err := httpClient.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return 0, fmt.Errorf("failed to create course: status %d", resp.StatusCode)
	}

	var res struct{ ID int `json:"id"` }
	json.NewDecoder(resp.Body).Decode(&res)
	return res.ID, nil
}

// UpdateCourse updates an existing course.
func UpdateCourse(token string, id int, title, description string) error {
	payload := CoursePayload{Title: title, Description: description}
	body, _ := json.Marshal(payload)

	url := fmt.Sprintf("%s/courses/%d", baseURL, id)
	req, _ := http.NewRequest("PATCH", url, bytes.NewBuffer(body))
	setAuthHeader(req, token)

	resp, err := httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("failed to update course: status %d", resp.StatusCode)
	}
	return nil
}

// DeleteCourse removes a course by ID.
func DeleteCourse(token string, id int) error {
	url := fmt.Sprintf("%s/courses/%d", baseURL, id)
	req, _ := http.NewRequest("DELETE", url, nil)
	setAuthHeader(req, token)

	resp, err := httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("failed to delete course: status %d", resp.StatusCode)
	}
	return nil
}

// JoinCourse enrolls a user in a course.
func JoinCourse(token string, id int, targetUserID string) error {
	payload := UserPayload{UserID: targetUserID}
	body, _ := json.Marshal(payload)

	url := fmt.Sprintf("%s/courses/%d/join", baseURL, id)
	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(body))
	setAuthHeader(req, token)

	resp, err := httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("failed to join course: status %d", resp.StatusCode)
	}
	return nil
}

// LeaveCourse unenrolls a user from a course.
func LeaveCourse(token string, id int, targetUserID string) error {
	payload := UserPayload{UserID: targetUserID}
	body, _ := json.Marshal(payload)

	url := fmt.Sprintf("%s/courses/%d/leave", baseURL, id)
	req, _ := http.NewRequest("DELETE", url, bytes.NewBuffer(body))
	setAuthHeader(req, token)

	resp, err := httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("failed to leave course: status %d", resp.StatusCode)
	}
	return nil
}

// GetCourseStudents returns the list of students enrolled in a course.
func GetCourseStudents(token string, id int) ([]string, error) {
	url := fmt.Sprintf("%s/courses/%d/students", baseURL, id)
	req, _ := http.NewRequest("GET", url, nil)
	setAuthHeader(req, token)

	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get students: status %d", resp.StatusCode)
	}

	var data StudentsListResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, err
	}
	return data.StudentIDs, nil
}
