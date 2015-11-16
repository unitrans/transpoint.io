// Copyright 2015 Home24 AG. All rights reserved.
// Proprietary license.
package domain
import (
	"github.com/satori/go.uuid"
	"github.com/OneOfOne/xxhash/native"
	"strconv"
)

func GenerateKeyPair(u *User) (string, string) {
	uid := uuid.NewV4().String()
	uHash := xxhash.Checksum64([]byte(u.Id))
	var keyBytes []byte
	keyBytes = strconv.AppendUint(keyBytes, uHash, 10)
	keyBytes = append(keyBytes, '.')
	keyBytes = append(keyBytes, []byte(uid)[:8]...)

	var secretBytes []byte
	sHash := xxhash.Checksum64([]byte(uid))
	secretBytes = append(secretBytes, []byte(uid)[9:13]...)
	secretBytes = append(secretBytes, '.')
	secretBytes = strconv.AppendUint(secretBytes, sHash, 10)
	secretBytes = append(secretBytes, '.')
	secretBytes = append(secretBytes, []byte(uid)[24:]...)

	return string(keyBytes), string(secretBytes)
}