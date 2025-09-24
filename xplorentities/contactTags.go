package xplorentities

import (
	"errors"
	"fmt"
	"time"

	"github.com/angelbarreiros/XPlorGo/util"
	// Ajusta esta importación según tu estructura de proyecto
)

// XPlorContactTag representa un tag de contacto individual
type XPlorContactTag struct {
	ID                 *string         `json:"@id"`
	Type               string          `json:"@type"`
	Contact            *string         `json:"contact"`
	Subscription       *string         `json:"subscription"`
	SubscriptionOption interface{}     `json:"subscriptionOption"`
	Name               string          `json:"name"`
	ValidFrom          util.LocalTime  `json:"validFrom"`
	ValidThrough       *util.LocalTime `json:"validThrough,omitempty"`
	CreatedAt          util.LocalTime  `json:"createdAt"`
	DeletedAt          *util.LocalTime `json:"deletedAt,omitempty"`
	DeletedBy          interface{}     `json:"deletedBy,omitempty"`
}

// XPlorContactTags representa una colección de tags de contacto
type XPlorContactTags struct {
	Context     string            `json:"@context"`
	ID          string            `json:"@id"`
	Type        string            `json:"@type"`
	ContactTags []XPlorContactTag `json:"hydra:member"`
	Pagination  *HydraView        `json:"hydra:view,omitempty"`
}

// Métodos para XPlorContactTag
func (ct *XPlorContactTag) ContactTagID() (string, error) {
	return ExtractID(ct.ID, "contact tag ID field is nil")
}

func (ct *XPlorContactTag) ContactID() (string, error) {
	return ExtractID(ct.Contact, "contact ID field is nil")
}

func (ct *XPlorContactTag) SubscriptionID() (string, error) {
	return ExtractID(ct.Subscription, "subscription ID field is nil")
}

// Métodos para verificar estados
func (ct *XPlorContactTag) IsActive() bool {
	now := util.LocalTime{Time: time.Now()}

	// Si ValidThrough es nil, el tag es permanente (siempre activo)
	if ct.ValidThrough == nil {
		return ct.ValidFrom.Time.Before(now.Time) || ct.ValidFrom.Time.Equal(now.Time)
	}

	// Si tiene fecha de expiración, verificar rango
	return ct.ValidFrom.Time.Before(now.Time) && ct.ValidThrough.Time.After(now.Time)
}

func (ct *XPlorContactTag) IsExpired() bool {
	if ct.ValidThrough == nil {
		return false // Los tags sin fecha de expiración nunca expiran
	}

	now := util.LocalTime{Time: time.Now()}
	return ct.ValidThrough.Time.Before(now.Time)
}

func (ct *XPlorContactTag) IsNotStarted() bool {
	now := util.LocalTime{Time: time.Now()}
	return ct.ValidFrom.Time.After(now.Time)
}

func (ct *XPlorContactTag) IsPermanent() bool {
	return ct.ValidThrough == nil
}

func (ct *XPlorContactTag) IsDeleted() bool {
	return ct.DeletedAt != nil
}

// Métodos para XPlorContactTags (colección)
func (c *XPlorContactTags) CollectionID() (string, error) {
	return ExtractIDFromString(c.ID, "collection ID field is invalid")
}

func (c *XPlorContactTags) ContextID() (string, error) {
	return ExtractIDFromString(c.Context, " context ID field is invalid")
}

// Método para obtener todos los IDs de la colección
func (c *XPlorContactTags) AllContactTagIDs() ([]string, error) {
	if len(c.ContactTags) == 0 {
		return nil, errors.New("no contact tags available")
	}

	ids := make([]string, len(c.ContactTags))
	for i, tag := range c.ContactTags {
		id, err := tag.ContactTagID()
		if err != nil {
			return nil, err
		}
		ids[i] = id
	}
	return ids, nil
}

// Método para obtener todos los contact IDs únicos
func (c *XPlorContactTags) AllContactIDs() ([]string, error) {
	if len(c.ContactTags) == 0 {
		return nil, errors.New("no contact tags available")
	}

	contactIDs := make(map[string]bool)
	for _, tag := range c.ContactTags {
		contactID, err := tag.ContactID()
		if err == nil {
			contactIDs[contactID] = true
		}
	}

	result := make([]string, 0, len(contactIDs))
	for id := range contactIDs {
		result = append(result, id)
	}

	if len(result) == 0 {
		return nil, errors.New("no contact IDs found")
	}

	return result, nil
}

// Método para obtener todos los subscription IDs únicos
func (c *XPlorContactTags) AllSubscriptionIDs() ([]string, error) {
	if len(c.ContactTags) == 0 {
		return nil, errors.New("no contact tags available")
	}

	subscriptionIDs := make(map[string]bool)
	for _, tag := range c.ContactTags {
		subscriptionID, err := tag.SubscriptionID()
		if err == nil {
			subscriptionIDs[subscriptionID] = true
		}
	}

	result := make([]string, 0, len(subscriptionIDs))
	for id := range subscriptionIDs {
		result = append(result, id)
	}

	if len(result) == 0 {
		return nil, errors.New("no subscription IDs found")
	}

	return result, nil
}

// Métodos para filtrar tags
func (c *XPlorContactTags) ActiveTags() []XPlorContactTag {
	activeTags := make([]XPlorContactTag, 0)
	for _, tag := range c.ContactTags {
		if tag.IsActive() && !tag.IsDeleted() {
			activeTags = append(activeTags, tag)
		}
	}
	return activeTags
}

func (c *XPlorContactTags) ExpiredTags() []XPlorContactTag {
	expiredTags := make([]XPlorContactTag, 0)
	for _, tag := range c.ContactTags {
		if tag.IsExpired() && !tag.IsDeleted() {
			expiredTags = append(expiredTags, tag)
		}
	}
	return expiredTags
}

func (c *XPlorContactTags) PermanentTags() []XPlorContactTag {
	permanentTags := make([]XPlorContactTag, 0)
	for _, tag := range c.ContactTags {
		if tag.IsPermanent() && !tag.IsDeleted() {
			permanentTags = append(permanentTags, tag)
		}
	}
	return permanentTags
}

func (c *XPlorContactTags) TagsByName(name string) []XPlorContactTag {
	tags := make([]XPlorContactTag, 0)
	for _, tag := range c.ContactTags {
		if tag.Name == name && !tag.IsDeleted() {
			tags = append(tags, tag)
		}
	}
	return tags
}

func (c *XPlorContactTags) TagsByContact(contactID string) []XPlorContactTag {
	tags := make([]XPlorContactTag, 0)
	for _, tag := range c.ContactTags {
		id, err := tag.ContactID()
		if err == nil && id == contactID && !tag.IsDeleted() {
			tags = append(tags, tag)
		}
	}
	return tags
}

func (c *XPlorContactTags) TagsBySubscription(subscriptionID string) []XPlorContactTag {
	tags := make([]XPlorContactTag, 0)
	for _, tag := range c.ContactTags {
		id, err := tag.SubscriptionID()
		if err == nil && id == subscriptionID && !tag.IsDeleted() {
			tags = append(tags, tag)
		}
	}
	return tags
}

// Método para obtener tags únicos por nombre para un contacto
func (c *XPlorContactTags) UniqueTagNamesForContact(contactID string) []string {
	tagNames := make(map[string]bool)
	for _, tag := range c.ContactTags {
		id, err := tag.ContactID()
		if err == nil && id == contactID && !tag.IsDeleted() {
			tagNames[tag.Name] = true
		}
	}

	result := make([]string, 0, len(tagNames))
	for name := range tagNames {
		result = append(result, name)
	}
	return result
}
func (c *XPlorContactTag) String() string {
	return fmt.Sprintf("ContactTag %s (Name: %s, ContactID: %s, SubscriptionID: %s, Active: %t)", *c.ID, c.Name, func() string {
		if c.Contact != nil {
			return *c.Contact
		}
		return "nil"
	}(), func() string {
		if c.Subscription != nil {
			return *c.Subscription
		}
		return "nil"
	}(), c.IsActive())

}
