package xplorentities

import (
	"fmt"
	"strings"
)

type XPloreClubs struct {
	Context    interface{} `json:"@context"` // Cambiado a interface{} para permitir string o object
	ID         string      `json:"@id"`
	Type       string      `json:"@type"`
	Clubs      []XPlorClub `json:"hydra:member"`
	Pagination HydraView   `json:"hydra:view"`
}

type XPlorClub struct {
	ID                                                     string            `json:"@id"` // Cambiado de *string a string (no nullable según spec)
	Type                                                   string            `json:"@type"`
	Context                                                interface{}       `json:"@context"`          // Agregado, interface{} para string or object
	ClubNumberID                                           int               `json:"id"`                // Cambiado de ClubNumberID a id para coincidir
	Number                                                 *string           `json:"number"`            // Nullable, string [3..7] chars
	Code                                                   string            `json:"code"`              // Required, string [3..5] chars
	Name                                                   string            `json:"name"`              // Required
	Email                                                  *string           `json:"email"`             // Nullable, string <email>
	Phone                                                  *string           `json:"phone"`             // Nullable
	StreetAddress                                          *string           `json:"streetAddress"`     // Nullable
	PostalCode                                             string            `json:"postalCode"`        // Required
	AddressLocality                                        string            `json:"addressLocality"`   // Required
	AddressCountry                                         string            `json:"addressCountry"`    // Required
	AddressCountryIso                                      string            `json:"addressCountryIso"` // Required
	OpeningDate                                            *string           `json:"openingDate"`       // Nullable, string <date-time> (cambiado de *util.LocalTime)
	Description                                            *string           `json:"description"`       // Nullable
	ClubTags                                               []ClubTag         `json:"clubTags"`          // Required, Array of object (definir ClubTag si no existe)
	PublicMetadata                                         *ClubMetadata     `json:"publicMetadata"`    // Nullable, object (definir ClubMetadata)
	Locale                                                 *string           `json:"locale"`            // Nullable
	SaleTerms                                              []SaleTerms       `json:"saleTerms"`         // Required, Array of object (definir SaleTerms)
	CreatedAt                                              string            `json:"createdAt"`         // Required, string <date-time> (cambiado de util.LocalTime)
	CreatedBy                                              string            `json:"createdBy"`         // Required
	DeletedAt                                              *string           `json:"deletedAt"`         // Nullable, string <date-time>
	DebitDays                                              DebitDays         `json:"debitDays"`         // Required, object (definir DebitDays)
	DisabledProcess                                        []string          `json:"disabledProcess"`   // Array of string Nullable (cambiado de struct)
	AutomateDebits                                         *AutomateDebits   `json:"automateDebits"`    // Required Nullable, object (definir AutomateDebits)
	DebitFileType                                          *string           `json:"debitFileType"`     // Nullable, enum
	CanLaunchInsolventContactExportCommand                 bool              `json:"canLaunchInsolventContactExportCommand"`
	AllowedSignatureMethods                                []string          `json:"allowedSignatureMethods"` // Required
	SepaSftp                                               []string          `json:"sepaSftp"`                // Required
	GrantingBarcodeAccess                                  bool              `json:"grantingBarcodeAccess"`
	EmployeeSafeDeposit                                    bool              `json:"employeeSafeDeposit"`              // Required
	DefaultRefundPolicyInitialPeriod                       string            `json:"defaultRefundPolicyInitialPeriod"` // Required
	DefaultRefundPolicy                                    string            `json:"defaultRefundPolicy"`              // Required
	BarcodeModel                                           string            `json:"barcodeModel"`
	HasPrintNode                                           bool              `json:"hasPrintNode"`
	QRCodeValidity                                         *int              `json:"qrCodeValidity"` // Required Nullable, integer
	AllowingPinAccess                                      bool              `json:"allowingPinAccess"`
	MandateBicEnabled                                      bool              `json:"mandateBicEnabled"`
	TaxRates                                               TaxRates          `json:"taxRates"`                  // Object (mantener struct existente si coincide)
	ResaboxNotification                                    []string          `json:"resaboxNotification"`       // Array of string
	VoucherConfig                                          []string          `json:"voucherConfig"`             // Required, Array of string
	InitialCheckoutAmount                                  *int              `json:"initialCheckoutAmount"`     // Required Nullable, integer
	GtmID                                                  *string           `json:"gtmId"`                     // Required Nullable
	OnlineCancellation                                     []*string         `json:"onlineCancellation"`        // Required, Array of string Nullable
	PicturePermissionRequired                              bool              `json:"picturePermissionRequired"` // Required
	OnlineSuspension                                       *OnlineSuspension `json:"onlineSuspension"`          // Required Nullable, object (definir OnlineSuspension)
	Unclosing                                              []string          `json:"unclosing"`                 // Array of string
	SendExportNotification                                 bool              `json:"sendExportNotification"`
	DebtSettlement                                         []*string         `json:"debtSettlement"`                       // Required, Array of string Nullable
	AutomaticContactClosing                                *bool             `json:"automaticContactClosing"`              // Required Nullable
	MarkContactClosing                                     *bool             `json:"markContactClosing"`                   // Required Nullable
	CheckNationalIdDocument                                bool              `json:"checkNationalIdDocument"`              // Required
	BookingUnitPositions                                   []string          `json:"bookingUnitPositions"`                 // Required
	AuthorizeContactClosingOnlinePayment                   *bool             `json:"authorizeContactClosingOnlinePayment"` // Required Nullable
	RecurrentPaymentType                                   []string          `json:"recurrentPaymentType"`                 // Required, Array of string
	SharingBookingClubs                                    []string          `json:"sharingBookingClubs"`                  // Required
	Prospecting                                            bool              `json:"prospecting"`                          // Required
	DefaultLegalAge                                        *int              `json:"defaultLegalAge"`                      // Nullable
	ContactTransferabilityPercent                          *int              `json:"contactTransferabilityPercent"`        // Nullable
	AvoidGeneratingEmptyInvoiceForFutureServiceSales       bool              `json:"avoidGeneratingEmptyInvoiceForFutureServiceSales"`
	RemoveEmptyAnticipatedInvoiceDuringServiceCancellation bool              `json:"removeEmptyAnticipatedInvoiceDuringServiceCancellation"`
	ActivityClubRuleEnabled                                bool              `json:"activityClubRuleEnabled"` // Required
	ContactAutomaticOptin                                  bool              `json:"contactAutomaticOptin"`
	GenerateCancellationCertificate                        bool              `json:"generateCancellationCertificate"` // Required
	FingerprintRequired                                    bool              `json:"fingerprintRequired"`             // Required
	AllowMandateReport                                     bool              `json:"allowMandateReport"`
	ActivateCustomerRisk                                   bool              `json:"activateCustomerRisk"`
	EnableProspectOptinAgreement                           bool              `json:"enableProspectOptinAgreement"` // Required
	ManageActiveSalePerson                                 []string          `json:"manageActiveSalePerson"`       // Required
	ExternalQuota                                          *int              `json:"externalQuota"`                // Nullable
	BlockContactInDebt                                     *BlockConfig      `json:"blockContactInDebt"`           // Nullable, object (definir BlockConfig)
	BlockContactMissingMandate                             *BlockConfig      `json:"blockContactMissingMandate"`   // Nullable, object
	ExternalQuotaCriteria                                  *BlockConfig      `json:"externalQuotaCriteria"`        // Nullable, object
	MandateType                                            string            `json:"mandateType"`
	IsSixPenaltyLevels                                     bool              `json:"isSixPenaltyLevels"`
	DebitSuspendedTagName                                  string            `json:"debitSuspendedTagName"`
	Timezone                                               string            `json:"timezone"`
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

// Indica si el club tiene habilitado el acceso por código de barras
func (c XPlorClub) AllowsBarcodeAccess() bool {
	return c.GrantingBarcodeAccess
}

// Devuelve la tasa de impuesto preferida
func (c XPlorClub) PreferredTaxRate() int {
	return c.TaxRates.Preferred
}

// Indica si el club permite pagos recurrentes SEPA
func (c XPlorClub) SupportsRecurrentSepa() bool {
	for _, paymentType := range c.RecurrentPaymentType {
		if paymentType == "sepa" {
			return true
		}
	}
	return false
}

// Devuelve la edad legal mínima para inscribirse
func (c XPlorClub) MinLegalAge() int {
	if c.DefaultLegalAge != nil {
		return *c.DefaultLegalAge
	}
	return 0 // or some default value
}

// Devuelve si el club tiene notificaciones Resabox activas
func (c XPlorClub) HasResaboxNotification() bool {
	return len(c.ResaboxNotification) > 0
}
