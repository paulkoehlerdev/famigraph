package entity

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
)

type UserHandle []byte

func NewUserHandle() (UserHandle, error) {
	id := make([]byte, 32)
	_, err := rand.Read(id)
	if err != nil {
		return nil, fmt.Errorf("creating random user handle: %w", err)
	}

	return id, nil
}

func HandleFromString(s string) (UserHandle, error) {
	return base64.URLEncoding.DecodeString(s)
}

func (u UserHandle) String() string {
	return base64.URLEncoding.EncodeToString(u)
}
