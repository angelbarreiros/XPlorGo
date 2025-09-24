package xplorentities

import (
	"errors"
	"path"
	"strconv"
	"strings"

	"github.com/angelbarreiros/XPlorGo/util"
)

// Estructuras principales
type XPlorArticles struct {
	Context     string         `json:"@context"`
	ID          string         `json:"@id"`
	Type        string         `json:"@type"`
	HydraMember []XPlorArticle `json:"hydra:member"`
	Pagination  *HydraView     `json:"hydra:view"`
}

type XPlorArticle struct {
	ID                        *string           `json:"@id"`
	Type                      string            `json:"@type"`
	Sale                      *string           `json:"sale"`
	ProductName               string            `json:"productName"`
	ProductDescription        string            `json:"productDescription"`
	ProductCode               string            `json:"productCode"`
	ProductType               string            `json:"productType"`
	OfferName                 string            `json:"offerName"`
	RegistrationFeeCode       interface{}       `json:"registrationFeeCode"`
	RegistrationFeeName       interface{}       `json:"registrationFeeName"`
	PriceTE                   float64           `json:"priceTE"`
	PriceTI                   float64           `json:"priceTI"`
	ProratedPriceTI           float64           `json:"proratedPriceTI"`
	ProratedPriceTE           float64           `json:"proratedPriceTE"`
	RegistrationFeeTI         float64           `json:"registrationFeeTI"`
	RegistrationFeeTE         float64           `json:"registrationFeeTE"`
	Tax                       float64           `json:"tax"`
	PriceCurrency             string            `json:"priceCurrency"`
	TaxRate                   float64           `json:"taxRate"`
	ArticleBehaviors          []ArticleBehavior `json:"articleBehaviors"`
	Parent                    interface{}       `json:"parent"`
	CreatedAt                 util.LocalTime    `json:"createdAt"`
	CreatedBy                 string            `json:"createdBy"`
	DeletedAt                 *util.LocalTime   `json:"deletedAt"`
	DeletedBy                 interface{}       `json:"deletedBy"`
	RepaymentSchedule         RepaymentSchedule `json:"repaymentSchedule"`
	ProductID                 *string           `json:"productId"`
	OfferID                   *string           `json:"offerId"`
	ClubID                    *string           `json:"clubId"`
	ContractModelID           *string           `json:"contractModelId"`
	Mandatory                 bool              `json:"mandatory"`
	RegistrationFeeDiscount   bool              `json:"registrationFeeDiscount"`
	ContractID                *string           `json:"contractId"`
	PackageID                 interface{}       `json:"packageId"`
	PackageName               interface{}       `json:"packageName"`
	PriceDiscountTI           float64           `json:"priceDiscountTI"`
	PriceDiscountTE           float64           `json:"priceDiscountTE"`
	RegistrationFeeDiscountTI float64           `json:"registrationFeeDiscountTI"`
	ProrataDiscountTI         float64           `json:"prorataDiscountTI"`
	RegistrationFeeDiscountTE float64           `json:"registrationFeeDiscountTE"`
	TotalTE                   float64           `json:"totalTE"`
	TotalTI                   float64           `json:"totalTI"`
	TotalTaxes                map[string]string `json:"totalTaxes"`
	InvoiceReference          string            `json:"invoiceReference"`
	ContactFamilyName         string            `json:"contactFamilyName"`
	ContactGivenName          string            `json:"contactGivenName"`
	ContactNumber             string            `json:"contactNumber"`
	HasImplementationErrors   bool              `json:"hasImplementationErrors"`
	ImplementationErrors      []interface{}     `json:"implementationErrors"`
	ContactID                 *string           `json:"contactId"`
	RenewalType               string            `json:"renewalType"`
}

type ArticleBehavior struct {
	ID                  *string     `json:"@id"`
	Type                string      `json:"@type"`
	BehaviorID          *string     `json:"behaviorId"`
	Configuration       interface{} `json:"configuration"`
	Implementation      interface{} `json:"implementation"`
	ImplementationError interface{} `json:"implementationError"`
	Result              interface{} `json:"result"`
	PackageElementID    interface{} `json:"packageElementId"`
}

type RepaymentSchedule struct {
	Occurrences []Occurrence   `json:"occurrences"`
	Recurrences []Recurrence   `json:"recurrences,omitempty"`
	SpecificDay bool           `json:"specificDay"`
	DebitDay    int            `json:"debitDay"`
	StartDate   util.LocalTime `json:"startDate"`
}

type Occurrence struct {
	Offset        string  `json:"offset"`
	Interval      string  `json:"interval"`
	Loop          int     `json:"loop"`
	TaxRate       float64 `json:"taxRate"`
	PriceTI       float64 `json:"priceTI"`
	PriceTE       float64 `json:"priceTE"`
	Tax           float64 `json:"tax"`
	PriceCurrency string  `json:"priceCurrency"`
}

type Recurrence struct {
	Offset        string  `json:"offset"`
	Interval      string  `json:"interval"`
	Loop          int     `json:"loop,omitempty"`
	TaxRate       float64 `json:"taxRate"`
	PriceTI       float64 `json:"priceTI"`
	PriceTE       float64 `json:"priceTE"`
	Tax           float64 `json:"tax"`
	PriceCurrency string  `json:"priceCurrency"`
}

// Métodos para extraer IDs
func (a *XPlorArticle) ArticleID() (string, error) {
	return ExtractID(a.ID, "article ID field is nil")
}

func (a *XPlorArticle) SaleID() (string, error) {
	return ExtractID(a.Sale, "sale ID field is nil")
}

func (a *XPlorArticle) ProductIDValue() (string, error) {
	return ExtractID(a.ProductID, "product ID field is nil")
}

func (a *XPlorArticle) OfferIDValue() (string, error) {
	return ExtractID(a.OfferID, "offer ID field is nil")
}

func (a *XPlorArticle) ClubIDValue() (string, error) {
	return ExtractID(a.ClubID, "club ID field is nil")
}

func (a *XPlorArticle) ContractModelIDValue() (string, error) {
	return ExtractID(a.ContractModelID, "contract model ID field is nil")
}

func (a *XPlorArticle) ContractIDValue() (string, error) {
	return ExtractID(a.ContractID, "contract ID field is nil")
}

func (a *XPlorArticle) ContactIDValue() (string, error) {
	return ExtractID(a.ContactID, "contact ID field is nil")
}

func (ab *ArticleBehavior) ArticleBehaviorID() (string, error) {
	return ExtractID(ab.ID, "article behavior ID field is nil")
}

func (ab *ArticleBehavior) BehaviorIDValue() (string, error) {
	return ExtractID(ab.BehaviorID, "behavior ID field is nil")
}

// Método para obtener el ID base de cualquier campo de tipo *string
func ExtractID(field *string, errMsg string) (string, error) {
	if field == nil {
		return "", errors.New(errMsg)
	}

	// Eliminar cualquier parámetro de consulta si existe
	cleanPath := strings.Split(*field, "?")[0]
	base := path.Base(cleanPath)
	return base, nil
}

// Método para obtener el ID numérico de cualquier campo de tipo *string
func ExtractIDInt(field *string, errMsg string) (int, error) {
	if field == nil {
		return 0, errors.New(errMsg)
	}

	// Eliminar cualquier parámetro de consulta si existe
	cleanPath := strings.Split(*field, "?")[0]
	base := path.Base(cleanPath)
	return strconv.Atoi(base)
}
