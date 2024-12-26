package value

import "time"

type WebauthnRegistrationChallengeData []byte
type WebauthnRegistrationChallengeResponseData []byte
type WebauthnRegistrationSessionData struct {
	Raw    []byte
	Expiry time.Time
}

type WebauthnLoginChallengeData []byte
type WebauthnLoginChallengeResponseData []byte
type WebauthnLoginSessionData struct {
	Raw    []byte
	Expiry time.Time
}
