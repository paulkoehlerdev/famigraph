package entity

import "github.com/go-webauthn/webauthn/webauthn"

var _ webauthn.User = (*User)(nil)

type User struct {
	Handle      UserHandle
	Credentials []webauthn.Credential
}

func (u User) WebAuthnID() []byte {
	return u.Handle
}

func (u User) WebAuthnName() string {
	return ""
}

func (u User) WebAuthnDisplayName() string {
	return ""
}

func (u User) WebAuthnCredentials() []webauthn.Credential {
	return u.Credentials
}
