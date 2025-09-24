package xplorentities

import (
	"errors"
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
}

// Entidad Contact
type XPlorContact struct {
	ID                      *string        `json:"@id"`
	Type                    string         `json:"@type"`
	Number                  string         `json:"number"`
	Address                 XPlorAddress   `json:"address"`
	BirthDate               util.LocalDate `json:"birthDate"`
	Email                   string         `json:"email"`
	FamilyName              string         `json:"familyName"`
	Gender                  string         `json:"gender"`
	GivenName               string         `json:"givenName"`
	PictureID               *string        `json:"pictureId"`
	ClubID                  *string        `json:"clubId"`
	Mobile                  *string        `json:"mobile"`
	SourceID                *string        `json:"sourceId"`
	PrescriberID            *string        `json:"prescriberId"`
	OccupationID            *string        `json:"occupationId"`
	GoalID                  *string        `json:"goalId"`
	GoalIDs                 []*string      `json:"goalIds"`
	MotivationID            *string        `json:"motivationId"`
	MotivationIDs           []*string      `json:"motivationIds"`
	CompanyID               *string        `json:"companyId"`
	IdentificationValidated bool           `json:"identificationValidated"`
	State                   string         `json:"state"`
	InitialSalepersonID     *string        `json:"initialSalepersonId"`
	CurrentSalepersonID     *string        `json:"currentSalepersonId"`
	CurrentSalepersonGiven  string         `json:"currentSalepersonGivenName"`
	CurrentSalepersonFamily string         `json:"currentSalepersonFamilyName"`
	CreatedAt               util.LocalTime `json:"createdAt"`
	UpdatedAt               util.LocalTime `json:"updatedAt"`
	PictureAllowed          *bool          `json:"pictureAllowed"`
	SponsorshipCode         string         `json:"sponsorshipCode"`
	ExternalID              *string        `json:"externalId"`
	NationalID              string         `json:"nationalId"`
	NationalIdDocumentID    *string        `json:"nationalIdDocumentId"`
	ProspectingState        *string        `json:"prospectingState"`
	Channel                 string         `json:"channel"`
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

func (c XPlorContact) ContactID() (int, error) {
	return ExtractIDInt(c.ID, "contact ID field is nil")
}

func (c XPlorContact) ClubIDValue() (int, error) {
	return ExtractIDInt(c.ClubID, "club ID field is nil")
}

func (c XPlorContact) PictureIDValue() (int, error) {
	return ExtractIDInt(c.PictureID, "picture ID field is nil")
}

func (c XPlorContact) GoalIDValue() (int, error) {
	return ExtractIDInt(c.GoalID, "goal ID field is nil")
}

func (c XPlorContact) InitialSalepersonIDValue() (int, error) {
	return ExtractIDInt(c.InitialSalepersonID, "initial salesperson ID field is nil")
}

func (c XPlorContact) CurrentSalepersonIDValue() (int, error) {
	return ExtractIDInt(c.CurrentSalepersonID, "current salesperson ID field is nil")
}

// Para Address
func (a XPlorAddress) AddressID() (int, error) {
	return ExtractIDInt(a.ID, "address ID field is nil")
}

// Helper para concatenar dirección completa
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
