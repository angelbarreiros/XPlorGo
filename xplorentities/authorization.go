package xplorentities

import "time"

type XPlorTokenResponse struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
	TokenType   string `json:"token_type"`
	Scope       string `json:"scope"`
}

// XPlorTokenWithTimestamp envuelve el token con información de cuándo se obtuvo
type XPlorTokenWithTimestamp struct {
	Token      *XPlorTokenResponse
	ObtainedAt time.Time
}

// IsExpired verifica si el token ha expirado
func (t *XPlorTokenWithTimestamp) IsExpired() bool {
	expirationTime := t.ObtainedAt.Add(time.Duration(t.Token.ExpiresIn) * time.Second)
	return time.Now().After(expirationTime)
}

// IsValid verifica si el token es válido (no nil, no vacío, no expirado)
func (t *XPlorTokenWithTimestamp) IsValid() bool {
	if t.Token == nil {
		return false
	}

	if t.Token.AccessToken == "" {
		return false
	}

	return !t.IsExpired()
}
