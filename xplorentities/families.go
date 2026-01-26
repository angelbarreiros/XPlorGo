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
	ID              *string          `json:"@id"`
	Type            string           `json:"@type"`
	Name            string           `json:"name"`
	CreatedAt       util.LocalDate   `json:"createdAt"`
	CreatedBy       string           `json:"createdBy"`
	Members         []familyMember   `json:"members"`
	SharedResources []sharedResource `json:"sharedResources"`
}

type familyMember struct {
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
	ID         *string         `json:"@id"`
	Type       string          `json:"@type"`
	BirthDate  *util.LocalDate `json:"birthDate"`
	FamilyName string          `json:"familyName"`
	GivenName  string          `json:"givenName"`
	State      string          `json:"state"`
}

type familyLinkResource struct {
	SharedResource sharedResource `json:"sharedResource"`
	Type           string         `json:"type"`
	Owner          bool           `json:"owner"`
	CreatedAt      util.LocalDate `json:"createdAt"`
}

type sharedResource struct {
	Subscription       *string        `json:"subsregistrar"`
	SubscriptionOption *string        `json:"subscriptionOption"`
	Family             string         `json:"family"`
	Properties         resourceLimits `json:"properties"`
}

type XPlorFamiliesParams struct {
	ContactId  string
	ContactIds []string
}

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

// Returns the subscription ID from the Subscription field (assumes it's a path)
func (s sharedResource) FamilySubscriptionID() (string, error) {
	return ExtractID(s.Subscription, "subscription field is empty")
}

type resourceLimits struct {
	AdultsMax   int `json:"adultsMax"`
	ChildrenMax int `json:"childrenMax"`
}

// metodos
func (f XPlorFamily) Parents() []familyMember {
	var parents []familyMember
	for _, m := range f.Members {
		if m.Responsible {
			parents = append(parents, m)
		}
	}
	return parents
}

// 🔹 Devuelve todos los hijos (no responsables)
func (f XPlorFamily) Children() []familyMember {
	var children []familyMember
	for _, m := range f.Members {
		if !m.Responsible {
			children = append(children, m)
		}
	}
	return children
}

// 🔹 Devuelve todos los nombres de los miembros (útil para debug o listados)
func (f XPlorFamily) MemberNames() []string {
	names := []string{}
	for _, m := range f.Members {
		names = append(names, fmt.Sprintf("%s %s", m.Contact.GivenName, m.Contact.FamilyName))
	}
	return names
}

// 🔹 Agrupa los miembros por estado (client, prospect, lost_client...)
func (f XPlorFamily) MembersByState() map[string][]familyMember {
	grouped := make(map[string][]familyMember)
	for _, m := range f.Members {
		state := strings.ToLower(m.Contact.State)
		grouped[state] = append(grouped[state], m)
	}
	return grouped
}

// 🔹 Busca un miembro por nombre o apellido
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

// 🔹 Edad aproximada de un miembro (si tiene birthDate)
func (m familyMember) Age() *int {
	if m.Contact.BirthDate == nil || m.Contact.BirthDate.IsZero() {
		return nil
	}
	years := int(time.Since(m.Contact.BirthDate.Time).Hours() / 24 / 365)
	return &years
}

func (c XPlorFamily) FamilyID() (string, error) {
	return ExtractID(c.ID, "family ID field is nil")
}

func (m familyMember) FamilyMemberID() (string, error) {
	return ExtractID(m.ID, "family member ID field is nil")
}

func (c contact) ContactID() (string, error) {
	return ExtractID(c.ID, "contact ID field is nil")
}

// 🔹 Pretty print de familia
func (f XPlorFamily) String() string {
	return fmt.Sprintf("Familia %s (%d miembros, %d padres, %d hijos)",
		f.Name,
		len(f.Members),
		len(f.Parents()),
		len(f.Children()))
}
