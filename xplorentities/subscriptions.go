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
	Context                 string               `json:"@context,omitempty"`
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
	Tags                    [][]string           `json:"tags"`
	SuspensionQuota         int                  `json:"suspensionQuota"`
	Warranties              []Warranty           `json:"warranties"`
	WarrantyInfo            *WarrantyInfo        `json:"warrantyInfo"`
	WarrantyState           string               `json:"warrantyState"`
	ClubId                  string               `json:"clubId"`
	NoticePeriod            *string              `json:"noticePeriod"`
	FixedPeriod             bool                 `json:"fixedPeriod"`
	EngagementRenewal       EngagementRenewal    `json:"engagementRenewal"`
	StatisticsDisabled      bool                 `json:"statisticsDisabled"`
	ServiceProperty         string               `json:"serviceProperty"`
	InitialInfo             InitialInfo          `json:"initialInfo"`
	SharedResource          string               `json:"sharedResource"`
	AutoRenewal             bool                 `json:"autoRenewal"`
	RenewalExpiredAt        string               `json:"renewalExpiredAt"`
	RenewalInfo             []RenewalInfo        `json:"renewalInfo"`
	NextRenewalDate         string               `json:"nextRenewalDate"`
	RenewalContact          Contact              `json:"renewalContact"`
	CounterLineAutoRenewal  map[string]any       `json:"counterLineAutoRenewal"`
	InclusiveValidThrough   string               `json:"inclusiveValidThrough"`
	InclusiveEngagedThrough string               `json:"inclusiveEngagedThrough"`
	InclusiveEndDate        string               `json:"inclusiveEndDate"`
	RenewalType             string               `json:"renewalType"`
	PaymentInfo             []PaymentInfoItem    `json:"paymentInfo"`
	RegularDebitDay         int                  `json:"regularDebitDay"`
	Family                  string               `json:"family"`
	Properties              map[string]any       `json:"properties"`
}

// Contact representa la información de contacto
type Contact struct {
	Id         *string `json:"@id,omitempty"`
	Type       string  `json:"@type,omitempty"`
	Number     string  `json:"number"`
	FamilyName string  `json:"familyName"`
	GivenName  string  `json:"givenName"`
	ClubId     *string `json:"clubId,omitempty"`
}

// InitialInfo representa la información inicial de la suscripción
type InitialInfo struct {
	OfferName             string             `json:"offerName"`
	ProductionCode        string             `json:"productionCode"`
	ProductName           string             `json:"productName"`
	BillingRhythm         string             `json:"billingRhythm"`
	RhythmBilling         string             `json:"rhythmBilling"`
	ProductDescription    string             `json:"productDescription"`
	IsOldComputedSchedule bool               `json:"isOldComputedSchedule"`
	Payments              []PaymentReference `json:"payments"`
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
	Amount             *Amount        `json:"amount"`
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
	Tags                    [][]string        `json:"tags"`
	SuspensionQuota         int               `json:"suspensionQuota"`
	Warranties              []Warranty        `json:"warranties"`
	WarrantyInfo            *WarrantyInfo     `json:"warrantyInfo"`
	WarrantyState           string            `json:"warrantyState"`
	ClubId                  string            `json:"clubId"`
	NoticePeriod            *string           `json:"noticePeriod"`
	FixedPeriod             bool              `json:"fixedPeriod"`
	EngagementRenewal       EngagementRenewal `json:"engagementRenewal"`
	StatisticsDisabled      bool              `json:"statisticsDisabled"`
	ServiceProperty         string            `json:"serviceProperty"`
	InitialInfo             InitialInfo       `json:"initialInfo"`
	SharedResource          string            `json:"sharedResource"`
	AutoRenewal             bool              `json:"autoRenewal"`
	RenewalExpiredAt        string            `json:"renewalExpiredAt"`
	RenewalInfo             []RenewalInfo     `json:"renewalInfo"`
	NextRenewalDate         string            `json:"nextRenewalDate"`
	RenewalContact          Contact           `json:"renewalContact"`
	CounterLineAutoRenewal  map[string]any    `json:"counterLineAutoRenewal"`
	InclusiveValidThrough   string            `json:"inclusiveValidThrough"`
	InclusiveEngagedThrough string            `json:"inclusiveEngagedThrough"`
	InclusiveEndDate        string            `json:"inclusiveEndDate"`
	RenewalType             string            `json:"renewalType"`
	PaymentInfo             []PaymentInfoItem `json:"paymentInfo"`
	RegularDebitDay         int               `json:"regularDebitDay"`
	Family                  string            `json:"family"`
	Properties              map[string]any    `json:"properties"`
}

// Warranty representa una garantía
type Warranty struct{}

// WarrantyInfo representa la información de garantía
type WarrantyInfo struct {
	Amount       int    `json:"amount"`
	Instructions string `json:"instructions"`
}

// EngagementRenewal representa la renovación de compromiso
type EngagementRenewal struct {
	MonthBeforeEnd int    `json:"monthBeforeEnd"`
	Period         string `json:"period"`
}

// PaymentReference representa una referencia a un pago
type PaymentReference struct {
	Id *string `json:"@id,omitempty"`
}

// PaymentInfoItem representa un elemento de información de pago (actualmente vacío según la API)
type PaymentInfoItem struct{}

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
	Page                   int
	ItemsPerPage           int
	ContactId              string
	ContactIds             []string
	ContactGivenName       string
	ContactFamilyName      string
	ContactClubId          string
	ContactClubIds         []string
	Name                   string
	TagName                string
	TagNames               []string
	ArticleId              string
	ArticleIds             []string
	WarrantyState          string
	WarrantyStates         []string
	ClubId                 string
	ClubIds                []string
	ValidFromBefore        string // Y-m-d H:i:s
	ValidFromAfter         string // Y-m-d H:i:s
	ValidThroughBefore     string // Y-m-d H:i:s
	ValidThroughAfter      string // Y-m-d H:i:s
	EngagedThroughBefore   string // Y-m-d H:i:s
	EngagedThroughAfter    string // Y-m-d H:i:s
	CreatedAtBefore        string // Y-m-d H:i:s
	CreatedAtAfter         string // Y-m-d H:i:s
	TerminatedAtBefore     string // Y-m-d H:i:s
	TerminatedAtAfter      string // Y-m-d H:i:s
	UpdatedAtBefore        string // Y-m-d H:i:s
	UpdatedAtAfter         string // Y-m-d H:i:s
	InclusiveEndDateBefore string // Y-m-d H:i:s
	InclusiveEndDateAfter  string // Y-m-d H:i:s
	InclusiveEngagedBefore string // Y-m-d H:i:s
	InclusiveEngagedAfter  string // Y-m-d H:i:s
	Consumed               *bool
	FixedPeriod            *bool
	Unlimited              *bool
	Current                *bool
	IsTerminated           *bool
	IsActive               *bool
	OrderBy                string // "validFrom", "createdAt", "validThrough"
	OrderDirection         string // "asc", "desc"
	Query                  string // Search on contact number, firstname or lastname
	ExistsTerminatedAt     *bool
}

// ToValues devuelve los parámetros como url.Values

// ToValues converts the subscription parameters to url.Values for query parameters
func (p XPlorSubscriptionsParams) ToValues(values *url.Values) {
	if p.Page > 0 {
		values.Set("page", strconv.Itoa(p.Page))
	}
	if p.ItemsPerPage > 0 {
		values.Set("itemsPerPage", strconv.Itoa(p.ItemsPerPage))
	}
	if strings.TrimSpace(p.ContactId) != "" {
		values.Set("contact", p.ContactId)
	}
	for _, id := range p.ContactIds {
		if strings.TrimSpace(id) != "" {
			values.Add("contact[]", id)
		}
	}
	if strings.TrimSpace(p.ContactGivenName) != "" {
		values.Set("contact.givenName", p.ContactGivenName)
	}
	if strings.TrimSpace(p.ContactFamilyName) != "" {
		values.Set("contact.familyName", p.ContactFamilyName)
	}
	if strings.TrimSpace(p.ContactClubId) != "" {
		values.Set("contact.clubId", p.ContactClubId)
	}
	for _, id := range p.ContactClubIds {
		if strings.TrimSpace(id) != "" {
			values.Add("contact.clubId[]", id)
		}
	}
	if strings.TrimSpace(p.Name) != "" {
		values.Set("name", p.Name)
	}
	if strings.TrimSpace(p.TagName) != "" {
		values.Set("tagName", p.TagName)
	}
	for _, tag := range p.TagNames {
		if strings.TrimSpace(tag) != "" {
			values.Add("tagName[]", tag)
		}
	}
	if strings.TrimSpace(p.ArticleId) != "" {
		values.Set("articleId", p.ArticleId)
	}
	for _, id := range p.ArticleIds {
		if strings.TrimSpace(id) != "" {
			values.Add("articleId[]", id)
		}
	}
	if strings.TrimSpace(p.WarrantyState) != "" {
		values.Set("warrantyState", p.WarrantyState)
	}
	for _, state := range p.WarrantyStates {
		if strings.TrimSpace(state) != "" {
			values.Add("warrantyState[]", state)
		}
	}
	if strings.TrimSpace(p.ClubId) != "" {
		values.Set("clubId", p.ClubId)
	}
	for _, id := range p.ClubIds {
		if strings.TrimSpace(id) != "" {
			values.Add("clubId[]", id)
		}
	}
	if strings.TrimSpace(p.ValidFromBefore) != "" {
		values.Set("validFrom[before]", p.ValidFromBefore)
	}
	if strings.TrimSpace(p.ValidFromAfter) != "" {
		values.Set("validFrom[after]", p.ValidFromAfter)
	}
	if strings.TrimSpace(p.ValidThroughBefore) != "" {
		values.Set("validThrough[before]", p.ValidThroughBefore)
	}
	if strings.TrimSpace(p.ValidThroughAfter) != "" {
		values.Set("validThrough[after]", p.ValidThroughAfter)
	}
	if strings.TrimSpace(p.EngagedThroughBefore) != "" {
		values.Set("engagedThrough[before]", p.EngagedThroughBefore)
	}
	if strings.TrimSpace(p.EngagedThroughAfter) != "" {
		values.Set("engagedThrough[after]", p.EngagedThroughAfter)
	}
	if strings.TrimSpace(p.CreatedAtBefore) != "" {
		values.Set("createdAt[before]", p.CreatedAtBefore)
	}
	if strings.TrimSpace(p.CreatedAtAfter) != "" {
		values.Set("createdAt[after]", p.CreatedAtAfter)
	}
	if strings.TrimSpace(p.TerminatedAtBefore) != "" {
		values.Set("terminatedAt[before]", p.TerminatedAtBefore)
	}
	if strings.TrimSpace(p.TerminatedAtAfter) != "" {
		values.Set("terminatedAt[after]", p.TerminatedAtAfter)
	}
	if strings.TrimSpace(p.UpdatedAtBefore) != "" {
		values.Set("updatedAt[before]", p.UpdatedAtBefore)
	}
	if strings.TrimSpace(p.UpdatedAtAfter) != "" {
		values.Set("updatedAt[after]", p.UpdatedAtAfter)
	}
	if strings.TrimSpace(p.InclusiveEndDateBefore) != "" {
		values.Set("inclusiveEndDate[before]", p.InclusiveEndDateBefore)
	}
	if strings.TrimSpace(p.InclusiveEndDateAfter) != "" {
		values.Set("inclusiveEndDate[after]", p.InclusiveEndDateAfter)
	}
	if strings.TrimSpace(p.InclusiveEngagedBefore) != "" {
		values.Set("inclusiveEngagedThrough[before]", p.InclusiveEngagedBefore)
	}
	if strings.TrimSpace(p.InclusiveEngagedAfter) != "" {
		values.Set("inclusiveEngagedThrough[after]", p.InclusiveEngagedAfter)
	}
	if p.Consumed != nil {
		values.Set("consumed", strconv.FormatBool(*p.Consumed))
	}
	if p.FixedPeriod != nil {
		values.Set("fixedPeriod", strconv.FormatBool(*p.FixedPeriod))
	}
	if p.Unlimited != nil {
		values.Set("unlimited", strconv.FormatBool(*p.Unlimited))
	}
	if p.Current != nil {
		values.Set("current", strconv.FormatBool(*p.Current))
	}
	if p.IsTerminated != nil {
		values.Set("isTerminated", strconv.FormatBool(*p.IsTerminated))
	}
	if p.IsActive != nil {
		values.Set("active", strconv.FormatBool(*p.IsActive))
	}
	if strings.TrimSpace(p.Query) != "" {
		values.Set("q", p.Query)
	}
	if strings.TrimSpace(p.OrderBy) != "" {
		if strings.TrimSpace(p.OrderDirection) != "" {
			values.Set("order["+p.OrderBy+"]", p.OrderDirection)
		}
	}
	if p.ExistsTerminatedAt != nil {
		values.Set("exists[terminatedAt]", strconv.FormatBool(*p.ExistsTerminatedAt))
	}
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
