package xplorentities

import (
	"strings"

	"github.com/angelbarreiros/XPlorGo/util"
)

// ---------- Colección ----------
type XPlorStudios struct {
	Context    string        `json:"@context"`
	ID         string        `json:"@id"`
	Type       string        `json:"@type"`
	Studios    []XPlorStudio `json:"hydra:member"`
	Pagination HydraView     `json:"hydra:view"`
}

// ---------- Entidad Studio ----------
type XPlorStudio struct {
	ID              *string         `json:"@id"`   // ej: "/enjoy/studios/2552"
	Type            string          `json:"@type"` // "Studio"
	Name            string          `json:"name"`
	Club            *string         `json:"club"`   // ej: "/enjoy/clubs/1249"
	ZoneID          *string         `json:"zoneId"` // puede ser null
	Capacity        *int            `json:"capacity"`
	Overbooking     *int            `json:"overbooking"`
	StreetAddress   string          `json:"streetAddress"`
	PostalCode      string          `json:"postalCode"`
	AddressLocality string          `json:"addressLocality"`
	AddressCountry  string          `json:"addressCountry"`
	CreatedAt       util.LocalTime  `json:"createdAt"`
	CreatedBy       string          `json:"createdBy"`
	Tags            *string         `json:"tags"`
	ArchivedAt      *util.LocalTime `json:"archivedAt"`
	ArchivedBy      *string         `json:"archivedBy"`
}

// ---------- Métodos para extraer IDs ----------

// StudioID devuelve el número final del @id, ej: "/enjoy/studios/2552" → 2552
func (s XPlorStudio) StudioID() (string, error) {
	return ExtractID(s.ID, "studio ID field is nil")
}

// ClubID devuelve el número final del club, ej: "/enjoy/clubs/1249" → "1249"
func (s XPlorStudio) ClubID() (string, error) {
	return ExtractID(s.Club, "club ID field is nil")
}
func (s XPlorStudio) Address() string {
	parts := []string{}

	if s.StreetAddress != "" {
		parts = append(parts, s.StreetAddress)
	}
	if s.PostalCode != "" {
		parts = append(parts, s.PostalCode)
	}
	if s.AddressLocality != "" {
		parts = append(parts, s.AddressLocality)
	}
	if s.AddressCountry != "" {
		parts = append(parts, s.AddressCountry)
	}

	return strings.Join(parts, ", ")
}
