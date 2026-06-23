package main_module

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

// GetUserData returns aggregated user data (courses, tests, grades).
func GetUserData(token, targetID string, courses, tests, grades bool) (*UserProfileData, error) {
	url := fmt.Sprintf("%s/users/%s/data?c=%t&t=%t&g=%t", baseURL, targetID, courses, tests, grades)
	req, _ := http.NewRequest("GET", url, nil)
	setAuthHeader(req, token)

	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var data UserProfileData
	json.NewDecoder(resp.Body).Decode(&data)
	return &data, nil
}

// GetUserList returns all users.
func GetUserList(token string) ([]UserInfo, error) {
	req, _ := http.NewRequest("GET", baseURL+"/users", nil)
	setAuthHeader(req, token)

	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var list []UserInfo
	json.NewDecoder(resp.Body).Decode(&list)
	return list, nil
}

// GetUserInfo returns user info by ID.
func GetUserInfo(token, targetID string) (map[string]any, error) {
	url := fmt.Sprintf("%s/users/%s", baseURL, targetID)
	req, _ := http.NewRequest("GET", url, nil)
	setAuthHeader(req, token)

	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get user info: status %d", resp.StatusCode)
	}

	var res map[string]any
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return nil, err
	}
	return res, nil
}

// UpdateUserFullName updates a user's full name.
func UpdateUserFullName(token, targetID, newName string) error {
	payload := map[string]string{"full_name": newName}
	body, _ := json.Marshal(payload)

	url := fmt.Sprintf("%s/users/%s/name", baseURL, targetID)
	req, _ := http.NewRequest("PATCH", url, bytes.NewBuffer(body))
	setAuthHeader(req, token)

	resp, err := httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return nil
}

// GetUserRoles returns the roles assigned to a user.
func GetUserRoles(token, targetID string) ([]string, error) {
	url := fmt.Sprintf("%s/users/%s/roles", baseURL, targetID)
	req, _ := http.NewRequest("GET", url, nil)
	setAuthHeader(req, token)

	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var res UserRolesResponse
	json.NewDecoder(resp.Body).Decode(&res)
	return res.Roles, nil
}

// UpdateUserRoles updates the roles assigned to a user.
func UpdateUserRoles(token, targetID string, roles []string) error {
	payload := map[string][]string{"roles": roles}
	body, _ := json.Marshal(payload)

	url := fmt.Sprintf("%s/users/%s/roles", baseURL, targetID)
	req, _ := http.NewRequest("PATCH", url, bytes.NewBuffer(body))
	setAuthHeader(req, token)

	resp, err := httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to update roles: status %d", resp.StatusCode)
	}
	return nil
}

// GetUserBlockStatus returns whether a user is blocked.
func GetUserBlockStatus(token, targetID string) (bool, error) {
	url := fmt.Sprintf("%s/users/%s/block-status", baseURL, targetID)
	req, _ := http.NewRequest("GET", url, nil)
	setAuthHeader(req, token)

	resp, err := httpClient.Do(req)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return false, fmt.Errorf("failed to get block status: status %d", resp.StatusCode)
	}

	var res struct {
		IsBlocked bool `json:"is_blocked"`
	}
	json.NewDecoder(resp.Body).Decode(&res)
	return res.IsBlocked, nil
}

// SetUserBlockStatus blocks or unblocks a user.
func SetUserBlockStatus(token, targetID string, isBlocked bool) error {
	payload := map[string]bool{"is_blocked": isBlocked}
	body, _ := json.Marshal(payload)

	url := fmt.Sprintf("%s/users/%s/block", baseURL, targetID)
	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(body))
	setAuthHeader(req, token)

	resp, err := httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to set block status: status %d", resp.StatusCode)
	}
	return nil
}
