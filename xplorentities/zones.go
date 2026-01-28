package xplorentities

import (
	"net/url"
	"strings"

	"github.com/angelbarreiros/XPlorGo/util"
)

// XPlorZones represents the collection of zones
type XPlorZones struct {
	Context    string      `json:"@context"`
	ID         string      `json:"@id"`
	Type       string      `json:"@type"`
	Zones      []XPlorZone `json:"hydra:member"`
	Pagination *HydraView  `json:"hydra:view,omitempty"`
}

// XPlorZone represents an individual zone
type XPlorZone struct {
	Context       interface{}     `json:"@context,omitempty"`
	ID            *string         `json:"@id"`
	Type          string          `json:"@type"`
	Name          string          `json:"name"`
	Description   *string         `json:"description"`
	ClubID        string          `json:"clubId"`
	OpeningTime   *util.LocalTime `json:"openingTime"`
	ClosingTime   *util.LocalTime `json:"closingTime"`
	ClosedFrom    *util.LocalTime `json:"closedFrom"`
	ClosedTo      *util.LocalTime `json:"closedTo"`
	Parent        *ZoneParent     `json:"parent"`
	Configuration *ZoneConfig     `json:"configuration"`
	Entrance      bool            `json:"entrance"`
}

// ZoneParent represents a parent zone reference
type ZoneParent struct {
	ID   *string `json:"@id"`
	Type string  `json:"@type"`
	Name string  `json:"name"`
}

// ZoneConfig represents zone configuration
type ZoneConfig struct {
	// Add specific configuration fields as needed based on API response
	Settings map[string]interface{} `json:"settings,omitempty"`
}

// XPlorZonesParams represents the search parameters for zones
type XPlorZonesParams struct {
	ClubID  string
	ClubIDs []string
	Name    string
}

// ToValues converts the zone search parameters to url.Values for query parameters
func (p XPlorZonesParams) ToValues(values *url.Values) {
	clubID := strings.TrimSpace(p.ClubID)
	if clubID != "" {
		values.Set("clubId", clubID)
	}

	for _, id := range p.ClubIDs {
		if strings.TrimSpace(id) != "" {
			values.Add("clubId[]", id)
		}
	}

	if p.Name != "" {
		values.Set("name", p.Name)
	}
}

// ZoneID extracts the zone ID from the @id field
func (z XPlorZone) ZoneID() (string, error) {
	return ExtractID(z.ID, "zone ID field is nil")
}

// ClubIDInt returns the club ID as an integer
func (z XPlorZone) ClubIDInt() (int, error) {
	clubIDPtr := &z.ClubID
	return ExtractIDInt(clubIDPtr, "club ID field is nil")
}

// ParentZoneID returns the parent zone ID if it exists
func (z XPlorZone) ParentZoneID() (string, error) {
	if z.Parent == nil {
		return "", nil
	}
	return ExtractID(z.Parent.ID, "parent zone ID field is nil")
}

// IsEntrance checks if the zone is an entrance zone
func (z XPlorZone) IsEntrance() bool {
	return z.Entrance
}

// HasParent checks if the zone has a parent
func (z XPlorZone) HasParent() bool {
	return z.Parent != nil
}

// IsOpen checks if the zone is currently open based on opening/closing times
func (z XPlorZone) IsOpen() bool {
	// If no opening/closing times are set, assume it's always open
	if z.OpeningTime == nil && z.ClosingTime == nil {
		return true
	}

	// Add logic here if needed to check current time against opening/closing times
	return true
}

// IsClosed checks if the zone is closed based on closedFrom/closedTo dates
func (z XPlorZone) IsClosed() bool {
	if z.ClosedFrom == nil || z.ClosedTo == nil {
		return false
	}

	// Add logic here if needed to check current date against closed dates
	return false
}
