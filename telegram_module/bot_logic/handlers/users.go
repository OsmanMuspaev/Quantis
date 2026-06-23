package handlers

import (
	"encoding/json"
	"net/http"
	"bot_logic/main_module"
)

// GetUserDataProfileHandler returns aggregated user profile data.
func GetUserDataProfileHandler(w http.ResponseWriter, r *http.Request) {
	token, _ := getToken(r)
	targetID := r.URL.Query().Get("target_id")

	c := r.URL.Query().Get("courses") == "true"
	t := r.URL.Query().Get("tests") == "true"
	g := r.URL.Query().Get("grades") == "true"

	data, err := main_module.GetUserData(token, targetID, c, t, g)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(data)
}

// GetUsersListHandler returns the list of all users.
func GetUsersListHandler(w http.ResponseWriter, r *http.Request) {
	token, _ := getToken(r)
	list, err := main_module.GetUserList(token)
	if err != nil {
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	}
	json.NewEncoder(w).Encode(list)
}

// GetUserInfoHandler returns user info by target_id.
func GetUserInfoHandler(w http.ResponseWriter, r *http.Request) {
	token, _ := getToken(r)
	targetID := r.URL.Query().Get("target_id")

	if targetID == "" {
		http.Error(w, "target_id is required", http.StatusBadRequest)
		return
	}

	info, err := main_module.GetUserInfo(token, targetID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(info)
}

// UpdateFullNameHandler updates a user's full name.
func UpdateFullNameHandler(w http.ResponseWriter, r *http.Request) {
	token, _ := getToken(r)
	targetID := r.URL.Query().Get("target_id")

	var body struct {
		NewName string `json:"full_name"`
	}
	json.NewDecoder(r.Body).Decode(&body)

	if err := main_module.UpdateUserFullName(token, targetID, body.NewName); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

// GetUserRolesHandler returns roles for a user.
func GetUserRolesHandler(w http.ResponseWriter, r *http.Request) {
	token, _ := getToken(r)
	targetID := r.URL.Query().Get("target_id")

	roles, err := main_module.GetUserRoles(token, targetID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(map[string][]string{"roles": roles})
}

// UpdateUserRolesHandler updates roles for a user.
func UpdateUserRolesHandler(w http.ResponseWriter, r *http.Request) {
	token, _ := getToken(r)
	targetID := r.URL.Query().Get("target_id")

	var body struct {
		Roles []string `json:"roles"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	if err := main_module.UpdateUserRoles(token, targetID, body.Roles); err != nil {
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	}
	w.WriteHeader(http.StatusOK)
}

// GetUserBlockStatusHandler returns whether a user is blocked.
func GetUserBlockStatusHandler(w http.ResponseWriter, r *http.Request) {
	token, _ := getToken(r)
	targetID := r.URL.Query().Get("target_id")

	isBlocked, err := main_module.GetUserBlockStatus(token, targetID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(map[string]bool{"is_blocked": isBlocked})
}

// SetUserBlockHandler blocks or unblocks a user.
func SetUserBlockHandler(w http.ResponseWriter, r *http.Request) {
	token, _ := getToken(r)
	targetID := r.URL.Query().Get("target_id")

	var body struct {
		IsBlocked bool `json:"is_blocked"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	if err := main_module.SetUserBlockStatus(token, targetID, body.IsBlocked); err != nil {
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	}
	w.WriteHeader(http.StatusOK)
}
