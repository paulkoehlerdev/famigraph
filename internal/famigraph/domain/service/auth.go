package service

import "github.com/paulkoehlerdev/famigraph/internal/famigraph/domain/value"

type AuthService interface {
	GetRegistrationChallenge() (value.WebauthnRegistrationChallengeData, value.WebauthnRegistrationSessionData, error)
	Register(response value.WebauthnRegistrationChallengeResponseData, session value.WebauthnRegistrationSessionData)

	GetLoginChallenge() (value.WebauthnLoginChallengeData, value.WebauthnLoginSessionData, error)
	Login(response value.WebauthnLoginChallengeResponseData, session value.WebauthnLoginSessionData)
}
