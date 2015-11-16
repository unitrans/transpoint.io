// Copyright 2015 Home24 AG. All rights reserved.
// Proprietary license.

// Package domain business entities
package domain
import (
	"github.com/OneOfOne/xxhash"
	"strconv"
)

// NewUser creates new user
func NewUser() *User {
	return &User{}
}

// User user struct
type User struct {
	Id string         `json:"id"`
	Pass  string      `json:"pass"`
	Keys  []string    `json:"keys"`
}

// IsLogin is user logged in
func (u *User) IsLogin() bool {
	return u.Id != ""
}

func (u *User) AddKey(key string) {
	u.Keys = append(u.Keys, key)
}

func (u *User) Clone() *User {
	clone := &User{}
	*clone = *u
	return clone
}

func HashPassword(password string) string {
	pHash := xxhash.Checksum64([]byte(password))
	var passBytes []byte
	passBytes = strconv.AppendUint(passBytes, pHash, 10)
	return string(passBytes)
}