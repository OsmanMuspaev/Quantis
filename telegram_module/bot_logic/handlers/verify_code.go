package handlers

import (
	"bot_logic/auth"
	"bot_logic/storage"
	"encoding/json"
	"net/http"
	"strconv"
)

type VerifyCodeRequest struct {
	ChatId string `json:"chat_id"`
	Code   string `json:"code"`
}

func VerifyCode(w http.ResponseWriter, r *http.Request) {
	var req VerifyCodeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}
	id, err := strconv.Atoi(req.ChatId)
	if err != nil {
		http.Error(w, "invalid chat_id", http.StatusBadRequest)
		return
	}
	status := storage.GetUserStatus(id)

	if status != "authorized" {
		http.Error(w, "not authorized", http.StatusForbidden)
		return
	}

	if err := auth.VerifyCode(id, req.Code); err != nil {
		http.Error(w, "code verification failed", http.StatusBadRequest)
		return
	}

	w.Write([]byte("Successfully logged in. You can continue on your new device!"))
}
