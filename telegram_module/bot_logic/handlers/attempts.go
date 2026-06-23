package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"bot_logic/main_module"
)

// GetPassedUsersHandler returns users who completed a test.
func GetPassedUsersHandler(w http.ResponseWriter, r *http.Request) {
	token, _ := getToken(r)
	testID, _ := strconv.Atoi(r.URL.Query().Get("test_id"))

	data, err := main_module.GetPassedUsers(token, testID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(data)
}

// GetScoresHandler returns scores for a test.
func GetScoresHandler(w http.ResponseWriter, r *http.Request) {
	token, _ := getToken(r)
	testID, _ := strconv.Atoi(r.URL.Query().Get("test_id"))

	data, err := main_module.GetTestScores(token, testID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	}
	json.NewEncoder(w).Encode(data)
}

// GetAnswersHandler returns answers for a test.
func GetAnswersHandler(w http.ResponseWriter, r *http.Request) {
	token, _ := getToken(r)
	testID, _ := strconv.Atoi(r.URL.Query().Get("test_id"))

	data, err := main_module.GetTestAnswers(token, testID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	}
	json.NewEncoder(w).Encode(data)
}

// StartAttemptHandler starts a new test attempt.
func StartAttemptHandler(w http.ResponseWriter, r *http.Request) {
	token, _ := getToken(r)
	testID, _ := strconv.Atoi(r.URL.Query().Get("test_id"))

	id, err := main_module.StartAttempt(token, testID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	json.NewEncoder(w).Encode(map[string]int{"attempt_id": id})
}

// SubmitAnswerHandler submits an answer for an attempt.
func SubmitAnswerHandler(w http.ResponseWriter, r *http.Request) {
	token, _ := getToken(r)
	attemptID, _ := strconv.Atoi(r.URL.Query().Get("attempt_id"))

	var body struct {
		QuestionID  int `json:"question_id"`
		AnswerIndex int `json:"answer_index"`
	}
	json.NewDecoder(r.Body).Decode(&body)

	if err := main_module.SubmitAnswer(token, attemptID, body.QuestionID, body.AnswerIndex); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// DeleteAnswerHandler deletes an answer for a specific question.
func DeleteAnswerHandler(w http.ResponseWriter, r *http.Request) {
	token, _ := getToken(r)
	attemptID, _ := strconv.Atoi(r.URL.Query().Get("attempt_id"))
	questionID, _ := strconv.Atoi(r.URL.Query().Get("question_id"))

	if err := main_module.DeleteAttemptAnswer(token, attemptID, questionID); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// UpdateAnswerHandler updates an answer for a specific question.
func UpdateAnswerHandler(w http.ResponseWriter, r *http.Request) {
	token, _ := getToken(r)
	attemptID, _ := strconv.Atoi(r.URL.Query().Get("attempt_id"))
	questionID, _ := strconv.Atoi(r.URL.Query().Get("question_id"))

	var body struct {
		AnswerIndex int `json:"answer_index"`
	}
	json.NewDecoder(r.Body).Decode(&body)

	if err := main_module.UpdateAttemptAnswer(token, attemptID, questionID, body.AnswerIndex); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// CompleteAttemptHandler finalizes a test attempt.
func CompleteAttemptHandler(w http.ResponseWriter, r *http.Request) {
	token, _ := getToken(r)
	attemptID, _ := strconv.Atoi(r.URL.Query().Get("attempt_id"))

	if err := main_module.CompleteAttempt(token, attemptID); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
}

// GetUserAttemptHandler returns attempt data for a specific user.
func GetUserAttemptHandler(w http.ResponseWriter, r *http.Request) {
	token, _ := getToken(r)
	testID, _ := strconv.Atoi(r.URL.Query().Get("test_id"))
	targetUserID := r.URL.Query().Get("target_user_id")

	data, err := main_module.GetAttemptData(token, testID, targetUserID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	}
	json.NewEncoder(w).Encode(data)
}

// GetUserAttemptAnswersHandler returns answers for a specific user's attempt.
func GetUserAttemptAnswersHandler(w http.ResponseWriter, r *http.Request) {
	token, _ := getToken(r)
	testID, _ := strconv.Atoi(r.URL.Query().Get("test_id"))
	targetUserID := r.URL.Query().Get("target_user_id")

	data, err := main_module.GetAttemptAnswers(token, testID, targetUserID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	}
	json.NewEncoder(w).Encode(data)
}
