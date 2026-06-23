package handlers

import (
	"encoding/json"
	"net/http"

	"bot_logic/auth"
	"bot_logic/storage"
)

type LoginRequest struct {
	ChatId int `json:"chat_id"`
}

func Login(w http.ResponseWriter, r *http.Request) {
	loginType := r.URL.Query().Get("type")

	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	status := storage.GetUserStatus(req.ChatId)

	switch status {
	case "anonymous":
		if err := storage.UpdateUser(req.ChatId); err != nil {
			http.Error(w, "failed to update user", http.StatusBadRequest)
			return
		}
	case "unknown":
		if err := storage.CreateUser(req.ChatId); err != nil {
			http.Error(w, "failed to create user", http.StatusBadRequest)
			return
		}
	}

	entryToken, err := storage.GetEntryToken(req.ChatId)
	if err != nil {
		http.Error(w, "failed to get entry token", http.StatusBadRequest)
		return
	}

	var (
		result interface{}
		resultErr error
	)

	switch loginType {
	case "github":
		result, resultErr = auth.GithubLogin(req.ChatId, entryToken)
	case "yandex":
		result, resultErr = auth.YandexLogin(req.ChatId, entryToken)
	case "code":
		code, err := auth.CodeAuth(req.ChatId, entryToken)
		if err != nil {
			http.Error(w, "failed to get code", http.StatusNotFound)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"code": code})
		return
	default:
		http.Error(w, "invalid login type", http.StatusBadRequest)
		return
	}

	if resultErr != nil {
		http.Error(w, "login request failed", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}
