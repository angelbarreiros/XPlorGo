package xplorentities

import (
	"net/url"
	"strconv"
	"strings"

	"github.com/angelbarreiros/XPlorGo/util"
)

// XPlorSubscriptions representa la colección de suscripciones
type XPlorSubscriptions struct {
	Context       string              `json:"@context"`
	ID            string              `json:"@id"`
	Type          string              `json:"@type"`
	Subscriptions []XPlorSubscription `json:"hydra:member"`
	Pagination    HydraView           `json:"hydra:view"`
}

// XPlorSubscription representa una suscripción individual
type XPlorSubscription struct {
	Id                      *string              `json:"@id"`
	Type                    string               `json:"@type"`
	Contact                 Contact              `json:"contact"`
	Name                    string               `json:"name"`
	TagName                 string               `json:"tagName"`
	ValidFrom               string               `json:"validFrom"`
	ValidThrough            string               `json:"validThrough"`
	EngagedThrough          string               `json:"engagedThrough"`
	TerminatedAt            string               `json:"terminatedAt"`
	Unlimited               bool                 `json:"unlimited"`
	CreatedAt               string               `json:"createdAt"`
	CreatedBy               string               `json:"createdBy"`
	UpdatedAt               string               `json:"updatedAt"`
	Consumed                bool                 `json:"consumed"`
	SubscriptionOptions     []SubscriptionOption `json:"subscriptionOptions"`
	ArticleId               string               `json:"articleId"`
	Tags                    []string             `json:"tags"`
	SuspensionQuota         int                  `json:"suspensionQuota"`
	Warranties              []Warranty           `json:"warranties"`
	WarrantyInfo            any                  `json:"warrantyInfo"`
	WarrantyState           string               `json:"warrantyState"`
	ClubId                  string               `json:"clubId"`
	NoticePeriod            *string              `json:"noticePeriod"`
	FixedPeriod             bool                 `json:"fixedPeriod"`
	EngagementRenewal       EngagementRenewal    `json:"engagementRenewal"`
	StatisticsDisabled      bool                 `json:"statisticsDisabled"`
	ServiceProperty         ServiceProperty      `json:"serviceProperty"`
	InitialInfo             InitialInfo          `json:"initialInfo"`
	Payments                []Payment            `json:"payments"`
	OfferName               string               `json:"offerName"`
	ProductionCode          string               `json:"productionCode"`
	ProductName             *string              `json:"productName"`
	BillingRhythm           string               `json:"billingRhythm"`
	RhythmBilling           string               `json:"rhythmBilling"`
	ProductDescription      string               `json:"productDescription"`
	IsOldComputedSchedule   bool                 `json:"isOldComputedSchedule"`
	SharedResource          SharedResource       `json:"sharedResource"`
	AutoRenewal             bool                 `json:"autoRenewal"`
	RenewalExpiredAt        string               `json:"renewalExpiredAt"`
	RenewalInfo             []RenewalInfo        `json:"renewalInfo"`
	NextRenewalDate         string               `json:"nextRenewalDate"`
	RenewalContact          Contact              `json:"renewalContact"`
	InclusiveValidThrough   string               `json:"inclusiveValidThrough"`
	InclusiveEngagedThrough string               `json:"inclusiveEngagedThrough"`
	InclusiveEndDate        string               `json:"inclusiveEndDate"`
	RenewalType             string               `json:"renewalType"`
	PaymentInfo             PaymentInfoItem      `json:"paymentInfo"`
	RegularDebitDay         int                  `json:"regularDebitDay"`
	Family                  string               `json:"family"`
	Properties              map[string]any       `json:"properties"`
	CounterLineAutoRenewal  any                  `json:"counterLineAutoRenewal"`
	CurrentDate             int                  `json:"currentDate"`
	Counters                []Counter            `json:"counters"`
}

// Contact representa la información de contacto
type Contact struct {
	Id         *string `json:"@id"`
	Type       string  `json:"@type"`
	Number     string  `json:"number"`
	FamilyName string  `json:"familyName"`
	GivenName  string  `json:"givenName"`
	ClubId     *string `json:"clubId"`
}

// InitialInfo representa la información inicial de la suscripción
type InitialInfo struct {
	OfferName             string `json:"offerName"`
	ProductCode           string `json:"productCode"`
	ProductDescription    string `json:"productDescription"`
	ProductName           string `json:"productName"`
	Payments              any    `json:"payments"`
	BillingRhythm         string `json:"billingRhythm"`
	RhythmBilling         string `json:"rhythmBilling"`
	IsOldComputedSchedule bool   `json:"isOldComputedSchedule"`
}

// RenewalInfo representa la información de renovación
type RenewalInfo struct {
	ActivatedAt        util.LocalDate `json:"activatedAt"`
	RenewalDay         any            `json:"renewalDay"`
	RenewalWeekDay     any            `json:"renewalWeekDay"`
	RenewalPeriod      string         `json:"renewalPeriod"`
	ProductName        string         `json:"productName"`
	ProductCode        string         `json:"productCode"`
	ProductDescription string         `json:"productDescription"`
	Amount             Amount         `json:"amount"`
	TaxRate            int            `json:"taxRate"`
	PriceCurrency      string         `json:"priceCurrency"`
}

// Amount representa el monto de pago
type Amount struct {
	Type    string `json:"type"`
	PriceTE int    `json:"priceTE"`
	PriceTI int    `json:"priceTI"`
	Tax     int    `json:"tax"`
	Month   any    `json:"month"`
}

// PaymentInfo representa la información de pago
type PaymentInfo struct {
	Current  *CurrentPayment `json:"current"`
	Payments []Payment       `json:"payments"`
}

// CurrentPayment representa el pago actual
type CurrentPayment struct {
	ActivatedAt   util.LocalDate `json:"activatedAt"`
	Day           any            `json:"day"`
	PriceTE       int            `json:"priceTE"`
	PriceTI       int            `json:"priceTI"`
	Tax           int            `json:"tax"`
	TaxRate       int            `json:"taxRate"`
	Period        string         `json:"period"`
	PriceCurrency string         `json:"priceCurrency"`
	Type          string         `json:"type"`
	Month         int            `json:"month"`
	Week          int            `json:"week"`
}

// Payment representa un pago individual
type Payment struct {
	ActivatedAt   util.LocalDate `json:"activatedAt"`
	Day           any            `json:"day"`
	PriceTE       int            `json:"priceTE"`
	PriceTI       int            `json:"priceTI"`
	Tax           int            `json:"tax"`
	TaxRate       int            `json:"taxRate"`
	Period        string         `json:"period"`
	PriceCurrency string         `json:"priceCurrency"`
	Type          string         `json:"type"`
	Month         int            `json:"month"`
	Week          int            `json:"week"`
}

// SubscriptionOption representa una opción de suscripción
type SubscriptionOption struct {
	ArticleId               string            `json:"articleId"`
	Tags                    []string          `json:"tags"`
	SuspensionQuota         int               `json:"suspensionQuota"`
	Warranties              []Warranty        `json:"warranties"`
	WarrantyInfo            any               `json:"warrantyInfo"`
	WarrantyState           string            `json:"warrantyState"`
	ClubId                  string            `json:"clubId"`
	NoticePeriod            *string           `json:"noticePeriod"`
	FixedPeriod             bool              `json:"fixedPeriod"`
	EngagementRenewal       EngagementRenewal `json:"engagementRenewal"`
	StatisticsDisabled      bool              `json:"statisticsDisabled"`
	ServiceProperty         ServiceProperty   `json:"serviceProperty"`
	InitialInfo             InitialInfo       `json:"initialInfo"`
	Payments                []Payment         `json:"payments"`
	OfferName               string            `json:"offerName"`
	ProductionCode          string            `json:"productionCode"`
	ProductName             *string           `json:"productName"`
	BillingRhythm           string            `json:"billingRhythm"`
	RhythmBilling           string            `json:"rhythmBilling"`
	ProductDescription      string            `json:"productDescription"`
	IsOldComputedSchedule   bool              `json:"isOldComputedSchedule"`
	SharedResource          SharedResource    `json:"sharedResource"`
	AutoRenewal             bool              `json:"autoRenewal"`
	RenewalExpiredAt        string            `json:"renewalExpiredAt"`
	RenewalInfo             []RenewalInfo     `json:"renewalInfo"`
	NextRenewalDate         string            `json:"nextRenewalDate"`
	RenewalContact          Contact           `json:"renewalContact"`
	InclusiveValidThrough   string            `json:"inclusiveValidThrough"`
	InclusiveEngagedThrough string            `json:"inclusiveEngagedThrough"`
	InclusiveEndDate        string            `json:"inclusiveEndDate"`
	RenewalType             string            `json:"renewalType"`
	PaymentInfo             []PaymentInfoItem `json:"paymentInfo"`
	RegularDebitDay         int               `json:"regularDebitDay"`
	Family                  string            `json:"family"`
	Properties              map[string]any    `json:"properties"`
	CounterLineAutoRenewal  any               `json:"counterLineAutoRenewal"`
	CurrentDate             int               `json:"currentDate"`
	Counters                []Counter         `json:"counters"`
}

// Warranty representa una garantía
type Warranty struct {
	// Define fields if known, otherwise use any
}

// WarrantyInfo representa la información de garantía
type WarrantyInfo struct {
	// Define fields if known
}

// EngagementRenewal representa la renovación de compromiso
type EngagementRenewal struct {
	MonthBeforeEnd int    `json:"monthBeforeEnd"`
	Period         string `json:"period"`
}

// ServiceProperty representa la propiedad del servicio
type ServiceProperty struct {
	Service    string         `json:"service"`
	Properties map[string]any `json:"properties"`
}

// SharedResource representa el recurso compartido
type SharedResource struct {
	Subscription       *XPlorSubscription  `json:"subscription"`
	SubscriptionOption *SubscriptionOption `json:"subscriptionOption"`
	ValidFrom          string              `json:"validFrom"`
	TerminatedAt       string              `json:"terminatedAt"`
	EngagedThrough     string              `json:"engagedThrough"`
	Name               string              `json:"name"`
	AutoRenewal        bool                `json:"autoRenewal"`
	RenewalExpiredAt   string              `json:"renewalExpiredAt"`
	RenewalInfo        []RenewalInfo       `json:"renewalInfo"`
	NextRenewalDate    string              `json:"nextRenewalDate"`
	RenewalContact     Contact             `json:"renewalContact"`
	InitialInfo        InitialInfo         `json:"initialInfo"`
	ValidThrough       string              `json:"validThrough"`
	InclusiveEndDate   string              `json:"inclusiveEndDate"`
	SpecialOption      bool                `json:"specialOption"`
	PaymentInfo        []PaymentInfoItem   `json:"paymentInfo"`
	RegularDebitDay    int                 `json:"regularDebitDay"`
	Family             string              `json:"family"`
	Properties         map[string]any      `json:"properties"`
}

// PaymentInfoItem representa un elemento de información de pago
type PaymentInfoItem struct {
	// Define fields if known
}

// CounterLineAutoRenewal representa la renovación automática de línea de contador
type CounterLineAutoRenewal struct {
	Subscription *XPlorSubscription `json:"subscription"`
	ValidFrom    string             `json:"validFrom"`
	CurrentDate  int                `json:"currentDate"`
	Counters     []Counter          `json:"counters"`
}

// Counter representa un contador
type Counter struct {
	ArticleId               string            `json:"articleId"`
	InclusiveValidThrough   string            `json:"inclusiveValidThrough"`
	InclusiveEngagedThrough string            `json:"inclusiveEngagedThrough"`
	InclusiveEndDate        string            `json:"inclusiveEndDate"`
	RenewalType             string            `json:"renewalType"`
	PaymentInfo             []PaymentInfoItem `json:"paymentInfo"`
	RegularDebitDay         int               `json:"regularDebitDay"`
}

type XPlorSubscriptionsParams struct {
	ContactId    string
	ContactIds   []string
	Name         string
	IsTerminated bool
	IsActive     bool
}

// ToValues devuelve los parámetros como url.Values

// ToValues converts the subscription parameters to url.Values for query parameters
func (p XPlorSubscriptionsParams) ToValues(values *url.Values) {
	contactId := p.ContactId
	if strings.TrimSpace(contactId) != "" {
		values.Set("contact", contactId)
	}
	if strings.TrimSpace(p.Name) != "" {
		values.Set("name", p.Name)
	}
	for _, id := range p.ContactIds {
		if strings.TrimSpace(id) != "" {
			values.Add("contact[]", id)
		}
	}
	values.Set("isTerminated", strconv.FormatBool(p.IsTerminated))
	values.Set("active", strconv.FormatBool(p.IsActive))
}

// Métodos para obtener IDs

// SubscriptionID extracts the subscription ID from the @id field
func (s XPlorSubscription) SubscriptionID() (string, error) {
	return ExtractID(s.Id, "subscription ID field is nil")
}

// ArticleID extracts the article ID from the articleId field
func (s XPlorSubscription) ArticleID() (string, error) {
	return ExtractID(&s.ArticleId, "article ID field is empty")
}

// ClubID extracts the club ID from the clubId field
func (s XPlorSubscription) ClubID() (string, error) {
	return ExtractID(&s.ClubId, "club ID field is empty")
}

// ContactID extracts the contact ID from the @id field
func (c Contact) ContactID() (string, error) {
	return ExtractID(c.Id, "contact ID field is nil")
}

// ContactClubID extracts the contact club ID from the clubId field
func (c Contact) ContactClubID() (string, error) {
	return ExtractID(c.ClubId, "contact club ID field is nil")
}
