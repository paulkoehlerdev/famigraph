package entity

import "github.com/go-webauthn/webauthn/webauthn"

type UserCredential struct {
	UserHandle UserHandle
	Credential *webauthn.Credential
}
