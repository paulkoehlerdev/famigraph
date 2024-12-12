package entity

import (
	"crypto/rand"
	"fmt"
)

type UserHandle = []byte

func NewUserHandle() (UserHandle, error) {
	id := make([]byte, 32)
	_, err := rand.Read(id)
	if err != nil {
		return nil, fmt.Errorf("creating random user handle: %w", err)
	}

	return id, nil
}
