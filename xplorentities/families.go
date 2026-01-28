package xplorentities

import (
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/angelbarreiros/XPlorGo/util"
)

type XPlorFamilies struct {
	Context    string        `json:"@context"`
	ID         string        `json:"@id"`
	Type       string        `json:"@type"`
	Families   []XPlorFamily `json:"hydra:member"`
	Pagination *HydraView    `json:"hydra:view,omitempty"`
}

type HydraIriTemplateMapping struct {
	Type     string `json:"@type"`
	Variable string `json:"variable"`
	Property string `json:"property"`
	Required bool   `json:"required"`
}

type XPlorFamily struct {
	Context         string           `json:"@context"`
	ID              *string          `json:"@id"`
	Type            string           `json:"@type"`
	Name            string           `json:"name"`
	CreatedAt       util.LocalDate   `json:"createdAt"`
	CreatedBy       string           `json:"createdBy"`
	Members         []familyMember   `json:"members"`
	SharedResources []sharedResource `json:"sharedResources"`
}

type familyMember struct {
	Context             string               `json:"@context"`
	ID                  *string              `json:"@id"`
	Type                string               `json:"@type"`
	Contact             contact              `json:"contact"`
	FamilyLinkResources []familyLinkResource `json:"familyLinkResources"`
	Role                string               `json:"role"`
	Responsible         bool                 `json:"responsible"`
	CreatedAt           util.LocalDate       `json:"createdAt"`
	CreatedBy           string               `json:"createdBy"`
	CanOverride         any                  `json:"canOverride"`
	OverriddenBy        any                  `json:"overriddenBy"`
}
type contact struct {
	Context    string          `json:"@context"`
	ID         *string         `json:"@id"`
	Type       string          `json:"@type"`
	BirthDate  *util.LocalDate `json:"birthDate"`
	FamilyName string          `json:"familyName"`
	GivenName  string          `json:"givenName"`
	State      string          `json:"state"`
}

type familyLinkResource struct {
	Context        string         `json:"@context"`
	ID             *string        `json:"@id"`
	Type           string         `json:"@type"`
	SharedResource sharedResource `json:"sharedResource"`
	RelationType   string         `json:"type"`
	Owner          bool           `json:"owner"`
	CreatedAt      util.LocalDate `json:"createdAt"`
}

type sharedResource struct {
	Context            string  `json:"@context"`
	ID                 *string `json:"@id"`
	Type               string  `json:"@type"`
	Subscription       *string `json:"subscription"`
	SubscriptionOption *string `json:"subscriptionOption"`
	Family             string  `json:"family"`
	Properties         any     `json:"properties"`
}

type XPlorFamiliesParams struct {
	ContactId  string
	ContactIds []string
}

// ToValues converts the params to url.Values for query parameters
func (p XPlorFamiliesParams) ToValues(values *url.Values) {
	contactId := strings.TrimSpace(p.ContactId)
	if contactId != "" {
		values.Set("members.contact", contactId)
	}

	for _, id := range p.ContactIds {
		if strings.TrimSpace(id) != "" {
			values.Add("members.contact[]", id)
		}
	}
}

// FamilySubscriptionID extracts the subscription ID from the subscription field
func (s sharedResource) FamilySubscriptionID() (string, error) {
	return ExtractID(s.Subscription, "subscription field is empty")
}

// Parents returns all primary responsible members of the family
func (f XPlorFamily) Parents() []familyMember {
	var parents []familyMember
	for _, m := range f.Members {
		if m.Responsible {
			parents = append(parents, m)
		}
	}
	return parents
}

// ObtainParentsSecundary returns secondary parents based on age (>= 18) or owner status
// Includes non-responsible members who are either adults or have owner=true in their FamilyLinkResources
func (f XPlorFamily) ObtainParentsSecundary() []familyMember {
	var secondaryParents []familyMember

	for _, m := range f.Members {
		// No incluir responsables primarios
		if m.Responsible {
			continue
		}

		isAdult := false
		hasOwnerLink := false

		// Verificar si es adulto por edad
		if age := m.Age(); age != nil && *age >= 18 {
			isAdult = true
		}

		// Verificar si tiene familyLinkResource con owner=true
		for _, link := range m.FamilyLinkResources {
			if link.Owner {
				hasOwnerLink = true
				break
			}
		}

		// Incluir si es adulto O tiene owner link
		if isAdult || hasOwnerLink {
			secondaryParents = append(secondaryParents, m)
		}
	}

	return secondaryParents
}

// Children returns all non-responsible members of the family (minors)
func (f XPlorFamily) Children() []familyMember {
	var children []familyMember
	for _, m := range f.Members {
		if !m.Responsible {
			children = append(children, m)
		}
	}
	return children
}

// MemberNames returns all member full names (useful for debugging or listing)
func (f XPlorFamily) MemberNames() []string {
	names := []string{}
	for _, m := range f.Members {
		names = append(names, fmt.Sprintf("%s %s", m.Contact.GivenName, m.Contact.FamilyName))
	}
	return names
}

// MembersByState groups the family members by their contact state
func (f XPlorFamily) MembersByState() map[string][]familyMember {
	grouped := make(map[string][]familyMember)
	for _, m := range f.Members {
		state := strings.ToLower(m.Contact.State)
		grouped[state] = append(grouped[state], m)
	}
	return grouped
}

// FindMember searches for a family member by name (given or family name)
func (f XPlorFamily) FindMember(name string) *familyMember {
	nameLower := strings.ToLower(name)
	for _, m := range f.Members {
		fullName := strings.ToLower(fmt.Sprintf("%s %s", m.Contact.GivenName, m.Contact.FamilyName))
		if strings.Contains(fullName, nameLower) {
			return &m
		}
	}
	return nil
}

// Age calculates the approximate age of a family member (if birth date is available)
func (m familyMember) Age() *int {
	if m.Contact.BirthDate == nil || m.Contact.BirthDate.IsZero() {
		return nil
	}
	years := int(time.Since(m.Contact.BirthDate.Time).Hours() / 24 / 365)
	return &years
}

// FamilyID extracts the family ID from the @id field
func (c XPlorFamily) FamilyID() (string, error) {
	return ExtractID(c.ID, "family ID field is nil")
}

// FamilyMemberID extracts the family member ID from the @id field
func (m familyMember) FamilyMemberID() (string, error) {
	return ExtractID(m.ID, "family member ID field is nil")
}

// ContactID extracts the contact ID from the @id field
func (c contact) ContactID() (string, error) {
	return ExtractID(c.ID, "contact ID field is nil")
}

// String returns a formatted string representation of the family
func (f XPlorFamily) String() string {
	return fmt.Sprintf("Familia %s (%d miembros, %d padres, %d hijos)",
		f.Name,
		len(f.Members),
		len(f.Parents()),
		len(f.Children()))
}
