package repository

import "time"

type Signer interface {
	Sign(data []byte, expiry time.Time) (string, error)
	Validate(string) ([]byte, time.Time, error)
}
