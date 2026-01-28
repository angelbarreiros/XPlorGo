package xplorentities

import (
	"errors"
	"fmt"
	"net/url"
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
	SubscriptionOption any             `json:"subscriptionOption"`
	Name               string          `json:"name"`
	ValidFrom          util.LocalTime  `json:"validFrom"`
	ValidThrough       *util.LocalTime `json:"validThrough,omitempty"`
	CreatedAt          *util.LocalTime `json:"createdAt"`
	DeletedAt          *util.LocalTime `json:"deletedAt,omitempty"`
	DeletedBy          any             `json:"deletedBy,omitempty"`
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

// ContactTagID extracts the contact tag ID from the @id field
func (ct *XPlorContactTag) ContactTagID() (string, error) {
	return ExtractID(ct.ID, "contact tag ID field is nil")
}

// ContactID extracts the contact ID from the contact field
func (ct *XPlorContactTag) ContactID() (string, error) {
	return ExtractID(ct.Contact, "contact ID field is nil")
}

// SubscriptionID extracts the subscription ID from the subscription field
func (ct *XPlorContactTag) SubscriptionID() (string, error) {
	return ExtractID(ct.Subscription, "subscription ID field is nil")
}

// Métodos para verificar estados

// IsActive checks if the contact tag is currently valid and active
func (ct *XPlorContactTag) IsActive() bool {
	now := util.LocalTime{Time: time.Now()}

	// Si ValidThrough es nil, el tag es permanente (siempre activo)
	if ct.ValidThrough == nil {
		return ct.ValidFrom.Time.Before(now.Time) || ct.ValidFrom.Time.Equal(now.Time)
	}

	// Si tiene fecha de expiración, verificar rango
	return ct.ValidFrom.Time.Before(now.Time) && ct.ValidThrough.Time.After(now.Time)
}

// IsExpired checks if the contact tag has expired
func (ct *XPlorContactTag) IsExpired() bool {
	if ct.ValidThrough == nil {
		return false // Los tags sin fecha de expiración nunca expiran
	}

	now := util.LocalTime{Time: time.Now()}
	return ct.ValidThrough.Time.Before(now.Time)
}

// IsNotStarted checks if the contact tag validity period has not started yet
func (ct *XPlorContactTag) IsNotStarted() bool {
	now := util.LocalTime{Time: time.Now()}
	return ct.ValidFrom.Time.After(now.Time)
}

// IsPermanent checks if the contact tag has no expiration date
func (ct *XPlorContactTag) IsPermanent() bool {
	return ct.ValidThrough == nil
}

// IsDeleted checks if the contact tag has been deleted
func (ct *XPlorContactTag) IsDeleted() bool {
	return ct.DeletedAt != nil
}

// Métodos para XPlorContactTags (coleción)

// CollectionID extracts the collection ID from the @id field
func (c *XPlorContactTags) CollectionID() (string, error) {
	return ExtractIDFromString(c.ID, "collection ID field is invalid")
}

// ContextID extracts the context ID from the @context field
func (c *XPlorContactTags) ContextID() (string, error) {
	return ExtractIDFromString(c.Context, " context ID field is invalid")
}

// Método para obtener todos los IDs de la colección

// AllContactTagIDs returns all contact tag IDs from the collection
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

// ActiveTags returns all active tags that have not been deleted
func (c *XPlorContactTags) ActiveTags() []XPlorContactTag {
	activeTags := make([]XPlorContactTag, 0)
	for _, tag := range c.ContactTags {
		if tag.IsActive() && !tag.IsDeleted() {
			activeTags = append(activeTags, tag)
		}
	}
	return activeTags
}

// ExpiredTags returns all expired tags that have not been deleted
func (c *XPlorContactTags) ExpiredTags() []XPlorContactTag {
	expiredTags := make([]XPlorContactTag, 0)
	for _, tag := range c.ContactTags {
		if tag.IsExpired() && !tag.IsDeleted() {
			expiredTags = append(expiredTags, tag)
		}
	}
	return expiredTags
}

// PermanentTags returns all permanent tags (no expiration date) that have not been deleted
func (c *XPlorContactTags) PermanentTags() []XPlorContactTag {
	permanentTags := make([]XPlorContactTag, 0)
	for _, tag := range c.ContactTags {
		if tag.IsPermanent() && !tag.IsDeleted() {
			permanentTags = append(permanentTags, tag)
		}
	}
	return permanentTags
}

// TagsByName returns all tags with the specified name that have not been deleted
func (c *XPlorContactTags) TagsByName(name string) []XPlorContactTag {
	tags := make([]XPlorContactTag, 0)
	for _, tag := range c.ContactTags {
		if tag.Name == name && !tag.IsDeleted() {
			tags = append(tags, tag)
		}
	}
	return tags
}

// TagsByContact returns all tags for the specified contact ID that have not been deleted
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

// TagsBySubscription returns all tags for the specified subscription ID that have not been deleted
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

// UniqueTagNamesForContact returns unique tag names for the specified contact that have not been deleted
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

// String returns a formatted string representation of the contact tag
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

// XPlorContactTagsParams represents the search parameters for contact tags

// ToValues converts the params to url.Values for query parameters
type XPlorContactTagsParams struct {
	ContactID      string
	ContactIDs     []string
	Email          string
	Emails         []string
	TagName        string
	TagNames       []string
	SubscriptionID string
	Active         *bool
}

// ToValues converts the params to url.Values for query parameters
func (p XPlorContactTagsParams) ToValues(values *url.Values) {
	// Contact ID filters
	if p.ContactID != "" {
		values.Set("contact", p.ContactID)
	}
	for _, id := range p.ContactIDs {
		if id != "" {
			values.Add("contact[]", id)
		}
	}

	// Email filters - assuming the API supports filtering by contact email
	if p.Email != "" {
		values.Set("contact.email", p.Email)
	}
	for _, email := range p.Emails {
		if email != "" {
			values.Add("contact.email[]", email)
		}
	}

	// Tag name filters
	if p.TagName != "" {
		values.Set("name", p.TagName)
	}
	for _, name := range p.TagNames {
		if name != "" {
			values.Add("name[]", name)
		}
	}

	// Subscription ID filter
	if p.SubscriptionID != "" {
		values.Set("subscription", p.SubscriptionID)
	}

	// Active filter
	if p.Active != nil {
		if *p.Active {
			values.Set("active", "true")
		} else {
			values.Set("active", "false")
		}
	}
}
