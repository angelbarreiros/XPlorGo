package xplorentities

import "time"

type XPlorTokenResponse struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
	TokenType   string `json:"token_type"`
	Scope       string `json:"scope"`
}

// XPlorTokenWithTimestamp wraps a token response with its obtention timestamp
type XPlorTokenWithTimestamp struct {
	Token      *XPlorTokenResponse
	ObtainedAt time.Time
}

// IsExpired checks if the token has expired based on the obtained timestamp and expiration time
func (t *XPlorTokenWithTimestamp) IsExpired() bool {
	expirationTime := t.ObtainedAt.Add(time.Duration(t.Token.ExpiresIn) * time.Second)
	return time.Now().After(expirationTime)
}

// IsValid checks if the token is valid (not nil, not empty, and not expired)
func (t *XPlorTokenWithTimestamp) IsValid() bool {
	if t.Token == nil {
		return false
	}

	if t.Token.AccessToken == "" {
		return false
	}

	return !t.IsExpired()
}
