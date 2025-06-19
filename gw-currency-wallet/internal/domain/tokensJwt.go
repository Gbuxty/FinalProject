package domain

import "time"

type RefreshToken string
type AccessToken string

type TokensPair struct {
	RefreshToken RefreshToken
	AccessToken  AccessToken
	RefreshTTL   time.Time
	AccessTTL    time.Time
}
