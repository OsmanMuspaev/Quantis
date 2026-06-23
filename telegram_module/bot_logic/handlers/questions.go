package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"bot_logic/main_module"
)

// CreateQuestionHandler creates a new question.
func CreateQuestionHandler(w http.ResponseWriter, r *http.Request) {
	token, _ := getToken(r)
	var q main_module.Question
	if err := json.NewDecoder(r.Body).Decode(&q); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	id, err := main_module.CreateQuestion(token, q)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(map[string]int{"id": id})
}

// GetMyQuestionsHandler returns questions owned by the current user.
func GetMyQuestionsHandler(w http.ResponseWriter, r *http.Request) {
	token, _ := getToken(r)
	questions, err := main_module.GetMyQuestions(token)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(questions)
}

// UpdateQuestionHandler updates a question.
func UpdateQuestionHandler(w http.ResponseWriter, r *http.Request) {
	token, _ := getToken(r)
	id, _ := strconv.Atoi(r.URL.Query().Get("id"))

	var q main_module.Question
	json.NewDecoder(r.Body).Decode(&q)

	newVer, err := main_module.UpdateQuestion(token, id, q)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(map[string]int{"new_version": newVer})
}

// DeleteQuestionHandler deletes a question.
func DeleteQuestionHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	token, _ := getToken(r)
	questionID, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil {
		http.Error(w, "Invalid question ID", http.StatusBadRequest)
		return
	}

	if err := main_module.DeleteQuestion(token, questionID); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// GetQuestionDetailHandler returns a specific version of a question.
func GetQuestionDetailHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	token, _ := getToken(r)

	questionID, err1 := strconv.Atoi(r.URL.Query().Get("id"))
	version, err2 := strconv.Atoi(r.URL.Query().Get("version"))

	if err1 != nil || err2 != nil {
		http.Error(w, "Invalid id or version", http.StatusBadRequest)
		return
	}

	question, err := main_module.GetQuestionDetail(token, questionID, version)
	if err != nil {
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(question)
}
