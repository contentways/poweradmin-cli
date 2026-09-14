// Copyright (c) 2026 Contentways
// SPDX-License-Identifier: MIT
package schema

import "contentways.dev/contentways/poweradmin-go/v2/poweradmin"

// User is the CLI output schema for a Poweradmin user.
type User struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Fullname string `json:"fullname,omitempty"`
	Active   bool   `json:"active"`
}

// UserFromSDK converts a poweradmin SDK User to the CLI output schema.
func UserFromSDK(u *poweradmin.User) User {
	return User{
		ID:       u.ID,
		Username: u.Username,
		Email:    u.Email,
		Fullname: u.Fullname,
		Active:   u.Active,
	}
}

// UserList wraps a slice of users in a root object for JSON output.
type UserList struct {
	Users []User `json:"users"`
	Count int    `json:"count"`
}

// UserListFromSDK converts a slice of SDK Users to the CLI output schema.
func UserListFromSDK(users []*poweradmin.User) UserList {
	out := make([]User, len(users))
	for i, u := range users {
		out[i] = UserFromSDK(u)
	}
	return UserList{Users: out, Count: len(out)}
}
