package xplorentities

import (
	"fmt"
	"strings"

	"github.com/angelbarreiros/XPlorGo/util"
)

type XPloreClubs struct {
	Context    interface{} `json:"@context"` // Cambiado a interface{} para permitir string o object
	ID         string      `json:"@id"`
	Type       string      `json:"@type"`
	Clubs      []XPlorClub `json:"hydra:member"`
	Pagination HydraView   `json:"hydra:view"`
}

type XPlorClub struct {
	ID                string          `json:"@id"` // Cambiado de *string a string (no nullable según spec)
	Type              string          `json:"@type"`
	Context           interface{}     `json:"@context"`          // Agregado, interface{} para string or object
	ClubNumberID      int             `json:"id"`                // Cambiado de ClubNumberID a id para coincidir
	Number            *string         `json:"number"`            // Nullable, string [3..7] chars
	Code              string          `json:"code"`              // Required, string [3..5] chars
	Name              string          `json:"name"`              // Required
	Email             *string         `json:"email"`             // Nullable, string <email>
	Phone             *string         `json:"phone"`             // Nullable
	StreetAddress     *string         `json:"streetAddress"`     // Nullable
	PostalCode        string          `json:"postalCode"`        // Required
	AddressLocality   string          `json:"addressLocality"`   // Required
	AddressCountry    string          `json:"addressCountry"`    // Required
	AddressCountryIso string          `json:"addressCountryIso"` // Required
	OpeningDate       *string         `json:"openingDate"`       // Nullable, string <date-time> (cambiado de *util.LocalTime)
	Description       *string         `json:"description"`       // Nullable
	ClubTags          []ClubTag       `json:"clubTags"`          // Required, Array of object (definir ClubTag si no existe)
	PublicMetadata    *ClubMetadata   `json:"publicMetadata"`    // Nullable, object (definir ClubMetadata)
	Locale            *string         `json:"locale"`            // Nullable
	SaleTerms         []SaleTerms     `json:"saleTerms"`         // Required, Array of object (definir SaleTerms)
	CreatedAt         *util.LocalTime `json:"createdAt"`         // Required, string <date-time> (cambiado de util.LocalTime)
	CreatedBy         string          `json:"createdBy"`         // Required
	DeletedAt         *string         `json:"deletedAt"`         // Nullable, string <date-time>

}

// Subestructuras

type TaxRates struct {
	Available     []int `json:"available"`
	Preferred     int   `json:"preferred"`
	RejectionFee  int   `json:"rejectionFee"`
	LateCancelFee int   `json:"lateCancelFee"`
	NoShowFee     int   `json:"noShowFee"`
}

type ResaboxNotification struct {
	Enabled        bool    `json:"enabled"`
	ContactEmails  *string `json:"contactEmails"`
	ContactNumbers *string `json:"contactNumbers"`
}

// Nuevas estructuras según spec

type ClubTag struct {
	Context interface{} `json:"@context"`
	ID      string      `json:"@id"`
	Type    string      `json:"@type"`
	Title   *string     `json:"title"`
	// ... otros campos según ClubTag.jsonld-read_sale_terms_light_norm
}

type ClubMetadata struct {
	Context     interface{} `json:"@context"`
	ID          string      `json:"@id"`
	Type        string      `json:"@type"`
	Title       *string     `json:"title"`
	Description *string     `json:"description"`
	Locale      *string     `json:"locale"`
}

type SaleTerms struct {
	Context interface{} `json:"@context"`
	ID      string      `json:"@id"`
	Type    string      `json:"@type"`
	// ... otros campos según SaleTerms.jsonld-read_sale_terms_light_norm
}

type DebitDays struct {
	// Definir como object, ejemplo básico
	Monday    bool `json:"monday"`
	Tuesday   bool `json:"tuesday"`
	Wednesday bool `json:"wednesday"`
	Thursday  bool `json:"thursday"`
	Friday    bool `json:"friday"`
	Saturday  bool `json:"saturday"`
	Sunday    bool `json:"sunday"`
}

type AutomateDebits struct {
	Sepa *SepaConfig `json:"sepa"`
	Card *CardConfig `json:"card"`
}

type SepaConfig struct {
	Enabled bool `json:"enabled"`
}

type CardConfig struct {
	Enabled bool `json:"enabled"`
}

type OnlineSuspension struct {
	IsAuthorized              bool               `json:"isAuthorized"`
	ImpactEndOfContract       bool               `json:"impactEndOfContract"`
	ImpactBooking             bool               `json:"impactBooking"`
	MonthlySubscription       SubscriptionConfig `json:"monthlySubscription"`
	WeeklySubscription        SubscriptionConfig `json:"weeklySubscription"`
	IncludedContactTags       []string           `json:"includedContactTags"`
	ExcludedContactTags       []string           `json:"excludedContactTags"`
	IncludedSuspensionMotives []string           `json:"includedSuspensionMotives"`
}

type SubscriptionConfig struct {
	Enabled bool `json:"enabled"`
}

type BlockConfig struct {
	Front                  bool      `json:"front"`
	Back                   bool      `json:"back"`
	CustomProductCode      *string   `json:"customProductCode"`
	LastMinute             []*string `json:"lastMinute"`
	UseOriginalProductCode bool      `json:"useOriginalProductCode"`
}

// Devuelve el ID numérico del club extraído de @id (ej: "/enjoy/clubs/1249" → 1249)
func (c XPlorClub) ClubID() (string, error) {
	// Since ID is now string, extract the numeric part
	parts := strings.Split(c.ID, "/")
	if len(parts) > 0 {
		return parts[len(parts)-1], nil
	}
	return "", fmt.Errorf("invalid club ID format")
}

// Devuelve un identificador legible combinando código y nombre
func (c XPlorClub) Identifier() string {
	return "[" + c.Code + "] " + c.Name
}

// Devuelve la dirección postal completa
func (c XPlorClub) FullAddress() string {
	street := ""
	if c.StreetAddress != nil {
		street = *c.StreetAddress
	}
	return strings.TrimSpace(street + ", " + c.PostalCode + " " + c.AddressLocality + " (" + c.AddressCountryIso + ")")
}

// Indica si el club está activo (no tiene fecha de borrado)
func (c XPlorClub) IsActive() bool {
	return c.DeletedAt == nil
}
