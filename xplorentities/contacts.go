package xplorentities

import (
	"errors"
	"net/url"
	"strings"
	"time"

	"github.com/angelbarreiros/XPlorGo/util"
)

// Colección
type XPlorContacts struct {
	Context    string         `json:"@context"`
	ID         string         `json:"@id"`
	Type       string         `json:"@type"`
	Contacts   []XPlorContact `json:"hydra:member"`
	Pagination *HydraView     `json:"hydra:view,omitempty"`
	Search     *HydraSearch   `json:"hydra:search,omitempty"`
}

// Entidad Contact
type XPlorContact struct {
	ID                      *string         `json:"@id"`
	Type                    string          `json:"@type"`
	Number                  string          `json:"number"`
	Address                 XPlorAddress    `json:"address"`
	BirthDate               *util.LocalDate `json:"birthDate"`
	Email                   string          `json:"email"`
	FamilyName              string          `json:"familyName"`
	Gender                  string          `json:"gender"`
	GivenName               string          `json:"givenName"`
	PictureID               *string         `json:"pictureId"`
	ClubID                  *string         `json:"clubId"`
	Mobile                  *string         `json:"mobile"`
	SourceID                *string         `json:"sourceId"`
	PrescriberID            *string         `json:"prescriberId"`
	OccupationID            *string         `json:"occupationId"`
	GoalID                  *string         `json:"goalId"`
	GoalIDs                 []*string       `json:"goalIds"`
	MotivationID            *string         `json:"motivationId"`
	MotivationIDs           []*string       `json:"motivationIds"`
	CompanyID               *string         `json:"companyId"`
	IdentificationValidated bool            `json:"identificationValidated"`
	State                   string          `json:"state"`
	InitialSalepersonID     *string         `json:"initialSalepersonId"`
	CurrentSalepersonID     *string         `json:"currentSalepersonId"`
	CurrentSalepersonGiven  string          `json:"currentSalepersonGivenName"`
	CurrentSalepersonFamily string          `json:"currentSalepersonFamilyName"`
	CreatedAt               *util.LocalTime `json:"createdAt"`
	UpdatedAt               *util.LocalTime `json:"updatedAt"`
	PictureAllowed          *bool           `json:"pictureAllowed"`
	SponsorshipCode         string          `json:"sponsorshipCode"`
	ExternalID              *string         `json:"externalId"`
	NationalID              string          `json:"nationalId"`
	NationalIdDocumentID    *string         `json:"nationalIdDocumentId"`
	ProspectingState        *string         `json:"prospectingState"`
	Channel                 string          `json:"channel"`
}

// Substruct para Address
type XPlorAddress struct {
	ID             *string `json:"@id"`
	Type           string  `json:"@type"`
	AddressCountry string  `json:"addressCountry"`
	CountryIso     string  `json:"addressCountryIso"`
	Locality       string  `json:"addressLocality"`
	PostalCode     string  `json:"postalCode"`
	StreetAddress  string  `json:"streetAddress"`
}

// Métodos helpers para extraer IDs -----------------

// ContactID extracts the contact ID from the @id field
func (c XPlorContact) ContactID() (string, error) {
	return ExtractID(c.ID, "contact ID field is nil")
}

// ClubIDValue extracts the club ID from the clubId field
func (c XPlorContact) ClubIDValue() (string, error) {
	return ExtractID(c.ClubID, "club ID field is nil")
}

// PictureIDValue extracts the picture ID from the pictureId field
func (c XPlorContact) PictureIDValue() (string, error) {
	return ExtractID(c.PictureID, "picture ID field is nil")
}

// GoalIDValue extracts the goal ID from the goalId field
func (c XPlorContact) GoalIDValue() (string, error) {
	return ExtractID(c.GoalID, "goal ID field is nil")
}

// InitialSalepersonIDValue extracts the initial salesperson ID from the initialSalepersonId field
func (c XPlorContact) InitialSalepersonIDValue() (string, error) {
	return ExtractID(c.InitialSalepersonID, "initial salesperson ID field is nil")
}

// CurrentSalepersonIDValue extracts the current salesperson ID from the currentSalepersonId field
func (c XPlorContact) CurrentSalepersonIDValue() (string, error) {
	return ExtractID(c.CurrentSalepersonID, "current salesperson ID field is nil")
}

// Para Address

// AddressID extracts the address ID from the @id field
func (a XPlorAddress) AddressID() (string, error) {
	return ExtractID(a.ID, "address ID field is nil")
}

// FullAddress returns the complete concatenated address string
func (a XPlorAddress) FullAddress() string {
	parts := []string{}
	if a.StreetAddress != "" {
		parts = append(parts, a.StreetAddress)
	}
	if a.PostalCode != "" {
		parts = append(parts, a.PostalCode)
	}
	if a.Locality != "" {
		parts = append(parts, a.Locality)
	}
	if a.AddressCountry != "" {
		parts = append(parts, a.AddressCountry)
	}
	return strings.Join(parts, ", ")
}

// Age calculates the contact's age based on their birth date
func (c XPlorContact) Age() (*int, error) {
	if c.BirthDate.IsZero() {
		return nil, errors.New("birth date is empty")
	}

	now := time.Now()
	age := now.Year() - c.BirthDate.Year()

	// Ajustar si el cumpleaños no ha ocurrido este año
	if now.YearDay() < c.BirthDate.YearDay() {
		age--
	}

	return &age, nil
}

// XPlorContactsParams represents the search parameters for contacts
type XPlorContactsParams struct {
	ContactID  string
	ContactIDs []string
	ClubID     string
	ClubIDs    []string
	State      string
	States     []string
	Email      string
	Emails     []string
	Mobile     string
	Number     string
	FamilyName string
	GivenName  string
}

// ToValues converts the params to url.Values for query parameters
func (p XPlorContactsParams) ToValues(values *url.Values) {
	// Contact ID filters
	contactID := strings.TrimSpace(p.ContactID)
	if contactID != "" {
		values.Set("id", contactID)
	}
	for _, id := range p.ContactIDs {
		if strings.TrimSpace(id) != "" {
			values.Add("id[]", id)
		}
	}

	// Club filters
	clubID := strings.TrimSpace(p.ClubID)
	if clubID != "" {
		values.Set("clubId", clubID)
	}
	for _, id := range p.ClubIDs {
		if strings.TrimSpace(id) != "" {
			values.Add("clubId[]", id)
		}
	}

	// State filters
	state := strings.TrimSpace(p.State)
	if state != "" {
		values.Set("state", state)
	}
	for _, s := range p.States {
		if strings.TrimSpace(s) != "" {
			values.Add("state[]", s)
		}
	}

	// Email filters
	if p.Email != "" {
		values.Set("email", p.Email)
	}
	for _, email := range p.Emails {
		if strings.TrimSpace(email) != "" {
			values.Add("email[]", email)
		}
	}

	// Other contact info filters
	if p.Mobile != "" {
		values.Set("mobile", p.Mobile)
	}
	if p.Number != "" {
		values.Set("number", p.Number)
	}
	if p.FamilyName != "" {
		values.Set("familyName", p.FamilyName)
	}
	if p.GivenName != "" {
		values.Set("givenName", p.GivenName)
	}
}
