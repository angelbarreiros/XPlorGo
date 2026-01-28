package xplorentities

import (
	"path"
	"slices"
	"time"

	"github.com/angelbarreiros/XPlorGo/util"
)

type UserProperties struct {
	NetworkNodeIds []*string `json:"networkNodeIds"`
}

type XPlorUsers struct {
	Context    *string     `json:"@context"`
	ID         *string     `json:"@id"`
	Type       *string     `json:"@type"`
	Users      []XPlorUser `json:"hydra:member"`
	Pagination HydraView   `json:"hydra:view"`
	// Search can be added if needed
}

type XPlorUser struct {
	AtID           *string         `json:"@id"`
	AtType         *string         `json:"@type"`
	AtContext      interface{}     `json:"@context"`
	PictureLink    *string         `json:"pictureLink"`
	Code           *string         `json:"code"`
	ClubIds        []*string       `json:"clubIds"`
	Mobile         *string         `json:"mobile"`
	Properties     interface{}     `json:"properties"`
	NetworkNodeIds []*string       `json:"networkNodeIds"`
	Email          string          `json:"email"`
	FamilyName     string          `json:"familyName"`
	GivenName      string          `json:"givenName"`
	Active         bool            `json:"active"`
	Roles          []string        `json:"roles"`
	CreatedAt      *util.LocalTime `json:"createdAt"`
	CreatedBy      string          `json:"createdBy"`
	DeletedAt      *util.LocalTime `json:"deletedAt"`
	DeletedBy      *string         `json:"deletedBy"`
	Locale         *string         `json:"locale"`
	ArchivedAt     *string         `json:"archivedAt"`
	ArchivedBy     *string         `json:"archivedBy"`
}

// UserID extracts the user ID from the @id field
func (u XPlorUser) UserID() (string, error) {
	return ExtractID(u.AtID, "user @id field is nil")
}

// ClubIDs extracts the base IDs from club IRIs
func (u XPlorUser) ClubIDs() []string {
	var ids []string
	for _, club := range u.ClubIds {
		if club == nil {
			continue
		}
		base := path.Base(*club)
		ids = append(ids, base)
	}
	return ids
}

// NetworkNodeIDs extracts the base IDs from network node IRIs
func (u XPlorUser) NetworkNodeIDs() []string {
	var ids []string
	for _, node := range u.NetworkNodeIds {
		if node == nil {
			continue
		}
		base := path.Base(*node)
		ids = append(ids, base)
	}
	return ids
}

// PropertiesNetworkNodeIDs extracts the base IDs from properties network node IRIs
func (u XPlorUser) PropertiesNetworkNodeIDs() []string {
	var ids []string
	if props, ok := u.Properties.(map[string]interface{}); ok {
		if nnids, ok := props["networkNodeIds"].([]interface{}); ok {
			for _, n := range nnids {
				if s, ok := n.(string); ok {
					base := path.Base(s)
					ids = append(ids, base)
				}
			}
		}
	}
	return ids
}

// IsActive returns true if the user is active and not deleted or archived
func (u XPlorUser) IsActive() bool {
	return u.Active && !u.IsDeleted() && !u.IsArchived()
}

// IsDeleted returns true if the user has been deleted
func (u XPlorUser) IsDeleted() bool {
	return u.DeletedAt != nil
}

// IsArchived returns true if the user has been archived
func (u XPlorUser) IsArchived() bool {
	return u.ArchivedAt != nil && *u.ArchivedAt != ""
}

// IsInactive returns true if the user's Active flag is false
func (u XPlorUser) IsInactive() bool {
	return !u.Active
}

// FullName returns the full name of the user
func (u XPlorUser) FullName() string {
	return u.GivenName + " " + u.FamilyName
}

// HasRole checks if the user has a specific role
func (u XPlorUser) HasRole(role string) bool {
	return slices.Contains(u.Roles, role)
}

// GetCreatedAt returns the creation time
func (u XPlorUser) GetCreatedAt() time.Time {
	if u.CreatedAt == nil {
		return time.Time{}
	}
	return u.CreatedAt.Time
}

// GetDeletedAt returns the deletion time if the user is deleted
func (u XPlorUser) GetDeletedAt() *time.Time {
	if u.DeletedAt == nil {
		return nil
	}
	t := u.DeletedAt.Time
	return &t
}
