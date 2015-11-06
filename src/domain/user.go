// Copyright 2015 Home24 AG. All rights reserved.
// Proprietary license.

// Package domain business entities
package domain

// NewUser creates new user
func NewUser() *User {
	return &User{}
}

// User user struct
type User struct {
	Id string         `json:"id"`
	Email string      `json:"email"`
	Token string      `json:"token"`
	Pass  string      `json:"pass"`
	Keys  []string    `json:"keys"`
}

// IsLogin is user logged in
func (u *User) IsLogin() bool {
	return u.Id != ""
}