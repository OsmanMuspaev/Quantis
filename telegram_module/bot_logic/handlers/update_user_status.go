package handlers

import (
	"net/http"

	"bot_logic/update"
)

// UpdateUserStatusHandler triggers a user status check.
func UpdateUserStatusHandler(w http.ResponseWriter, r *http.Request) {
	update.UpdateUserStatus()
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Update process completed"))
}
