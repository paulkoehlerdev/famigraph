package entity

import (
	"fmt"
	"github.com/go-webauthn/webauthn/webauthn"
)

var _ webauthn.User = (*User)(nil)

type User struct {
	Handle      UserHandle
	Credentials []webauthn.Credential
}

func NewUserWithRandomID() (*User, error) {
	handle, err := NewUserHandle()
	if err != nil {
		return nil, fmt.Errorf("creating user handle: %w", err)
	}

	return &User{
		Handle: handle,
	}, nil
}

func NewUser(handle UserHandle, credentials []webauthn.Credential) *User {
	return &User{
		Handle:      handle,
		Credentials: credentials,
	}
}

func (u *User) AddCredential(credential webauthn.Credential) {
	u.Credentials = append(u.Credentials, credential)
}

func (u *User) WebAuthnID() []byte {
	return u.Handle
}

func (u *User) WebAuthnName() string {
	return ""
}

func (u *User) WebAuthnDisplayName() string {
	return ""
}

func (u *User) WebAuthnCredentials() []webauthn.Credential {
	return u.Credentials
}
