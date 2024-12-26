package repository

import (
	"net/url"
	"time"
)

type Signer interface {
	Sign(data []byte, expiry time.Time) (string, error)
	Validate(string) ([]byte, time.Time, error)
}

type URLSigner interface {
	Sign(url *url.URL, expiry time.Time) (string, error)
	Validate(url *url.URL) (time.Time, error)
}
