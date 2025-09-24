package xplorentities

import (
	"path"
)

type UserProperties struct {
	NetworkNodeIds []*string `json:"networkNodeIds"`
}

type XPlorUsers struct {
	Context    *string      `json:"@context"`
	ID         *string      `json:"@id"`
	Type       *string      `json:"@type"`
	Users      []XPlorUser  `json:"hydra:member"`
	Pagination *HydraSearch `json:"hydra:view"`
	// Search can be added if needed
}

type XPlorUser struct {
	AtID           *string     `json:"@id"`
	AtType         *string     `json:"@type"`
	AtContext      interface{} `json:"@context"`
	PictureLink    *string     `json:"pictureLink"`
	Code           *string     `json:"code"`
	ClubIds        []*string   `json:"clubIds"`
	Mobile         *string     `json:"mobile"`
	Properties     interface{} `json:"properties"`
	NetworkNodeIds []*string   `json:"networkNodeIds"`
	Email          string      `json:"email"`
	FamilyName     string      `json:"familyName"`
	GivenName      string      `json:"givenName"`
	Active         bool        `json:"active"`
	Roles          []string    `json:"roles"`
	CreatedAt      string      `json:"createdAt"`
	CreatedBy      string      `json:"createdBy"`
	DeletedAt      *string     `json:"deletedAt"`
	DeletedBy      *string     `json:"deletedBy"`
	Locale         *string     `json:"locale"`
	ArchivedAt     *string     `json:"archivedAt"`
	ArchivedBy     *string     `json:"archivedBy"`
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
