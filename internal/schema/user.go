// Copyright (c) 2026 Contentways
// SPDX-License-Identifier: MIT
package schema

import "github.com/contentways/poweradmin-go/v4/poweradmin"

// User is the CLI output schema for a Poweradmin user.
type User struct {
	ID       int    `json:"id" yaml:"id"`
	Username string `json:"username" yaml:"username"`
	Email    string `json:"email" yaml:"email"`
	Fullname string `json:"fullname,omitempty" yaml:"fullname,omitempty"`
	Active   bool   `json:"active" yaml:"active"`
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
	Users []User `json:"users" yaml:"users"`
	Count int    `json:"count" yaml:"count"`
}

// UserListFromSDK converts a slice of SDK Users to the CLI output schema.
func UserListFromSDK(users []*poweradmin.User) UserList {
	out := make([]User, len(users))
	for i, u := range users {
		out[i] = UserFromSDK(u)
	}
	return UserList{Users: out, Count: len(out)}
}
