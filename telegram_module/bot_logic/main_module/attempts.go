package main_module

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

// GetPassedUsers returns users who have completed a test.
func GetPassedUsers(token string, testID int) ([]string, error) {
	url := fmt.Sprintf("%s/tests/%d/passed-users", baseURL, testID)
	req, _ := http.NewRequest("GET", url, nil)
	setAuthHeader(req, token)

	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var res struct {
		UserIDs []string `json:"user_ids"`
	}
	json.NewDecoder(resp.Body).Decode(&res)
	return res.UserIDs, nil
}

// GetTestScores returns scores for a test.
func GetTestScores(token string, testID int) ([]Score, error) {
	url := fmt.Sprintf("%s/tests/%d/scores", baseURL, testID)
	req, _ := http.NewRequest("GET", url, nil)
	setAuthHeader(req, token)

	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var res struct {
		Scores []Score `json:"scores"`
	}
	json.NewDecoder(resp.Body).Decode(&res)
	return res.Scores, nil
}

// GetTestAnswers returns answers for a test.
func GetTestAnswers(token string, testID int) ([]AttemptDetails, error) {
	url := fmt.Sprintf("%s/tests/%d/answers", baseURL, testID)
	req, _ := http.NewRequest("GET", url, nil)
	setAuthHeader(req, token)

	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var res struct {
		Attempts []AttemptDetails `json:"attempts"`
	}
	json.NewDecoder(resp.Body).Decode(&res)
	return res.Attempts, nil
}

// StartAttempt creates a new test attempt.
func StartAttempt(token string, testID int) (int, error) {
	url := fmt.Sprintf("%s/tests/%d/start", baseURL, testID)
	req, _ := http.NewRequest("POST", url, nil)
	setAuthHeader(req, token)

	resp, err := httpClient.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return 0, fmt.Errorf("failed to start test: status %d", resp.StatusCode)
	}

	var res struct {
		AttemptID int `json:"attempt_id"`
	}
	json.NewDecoder(resp.Body).Decode(&res)
	return res.AttemptID, nil
}

// SubmitAnswer submits an answer for an attempt.
func SubmitAnswer(token string, attemptID, questionID, answerIndex int) error {
	url := fmt.Sprintf("%s/attempts/%d/answers", baseURL, attemptID)
	payload := map[string]int{
		"question_id":  questionID,
		"answer_index": answerIndex,
	}
	body, _ := json.Marshal(payload)

	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(body))
	setAuthHeader(req, token)

	resp, err := httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("failed to submit answer: %d", resp.StatusCode)
	}
	return nil
}

// UpdateAttemptAnswer updates an answer for a specific question.
func UpdateAttemptAnswer(token string, attemptID, questionID, answerIndex int) error {
	url := fmt.Sprintf("%s/attempts/%d/questions/%d/answer", baseURL, attemptID, questionID)
	payload := map[string]int{"answer_index": answerIndex}
	body, _ := json.Marshal(payload)

	req, _ := http.NewRequest("PATCH", url, bytes.NewBuffer(body))
	setAuthHeader(req, token)

	resp, err := httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("failed to update answer: %d", resp.StatusCode)
	}
	return nil
}

// DeleteAttemptAnswer removes an answer for a specific question.
func DeleteAttemptAnswer(token string, attemptID, questionID int) error {
	url := fmt.Sprintf("%s/attempts/%d/questions/%d/answer", baseURL, attemptID, questionID)

	req, _ := http.NewRequest("DELETE", url, nil)
	setAuthHeader(req, token)

	resp, err := httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("failed to delete answer: %d", resp.StatusCode)
	}
	return nil
}

// CompleteAttempt finalizes a test attempt.
func CompleteAttempt(token string, attemptID int) error {
	url := fmt.Sprintf("%s/attempts/%d/complete", baseURL, attemptID)
	req, _ := http.NewRequest("POST", url, nil)
	setAuthHeader(req, token)

	resp, err := httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to complete attempt: %d", resp.StatusCode)
	}
	return nil
}

// GetAttemptData returns attempt data for a specific user.
func GetAttemptData(token string, testID int, targetUserID string) (map[string]any, error) {
	url := fmt.Sprintf("%s/tests/%d/attempts/%s", baseURL, testID, targetUserID)
	req, _ := http.NewRequest("GET", url, nil)
	setAuthHeader(req, token)

	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var data map[string]any
	json.NewDecoder(resp.Body).Decode(&data)
	return data, nil
}

// GetAttemptAnswers returns answers for a specific user's attempt.
func GetAttemptAnswers(token string, testID int, targetUserID string) (map[string]any, error) {
	url := fmt.Sprintf("%s/tests/%d/attempts/%s/answers", baseURL, testID, targetUserID)
	req, _ := http.NewRequest("GET", url, nil)
	setAuthHeader(req, token)

	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var data map[string]any
	json.NewDecoder(resp.Body).Decode(&data)
	return data, nil
}
