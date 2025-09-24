package xplorentities

import (
	"errors"
	"time"

	"github.com/angelbarreiros/XPlorGo/util"
)

// XPlorCounterLine representa una línea de contador individual
type XPlorCounterLine struct {
	ContactID         *string                     `json:"contactId"`
	ContactFamilyName string                      `json:"contactFamilyName"`
	ContactFirstName  string                      `json:"contactFirstName"`
	ContactNumber     string                      `json:"contactNumber"`
	Unit              *string                     `json:"unit"`
	TotalUnities      int                         `json:"totalUnities"`
	RemainingUnities  int                         `json:"remainingUnities"`
	ValidFrom         util.LocalTime              `json:"validFrom"`
	ValidThrough      util.LocalTime              `json:"validThrough"`
	CounterMovements  []string                    `json:"counterMovements"`
	ArticleID         *string                     `json:"articleId"`
	ServiceProperty   *CounterLineServiceProperty `json:"serviceProperty,omitempty"`
	CreatedAt         util.LocalTime              `json:"createdAt"`
	CreatedBy         string                      `json:"createdBy"`
	UpdatedAt         util.LocalTime              `json:"updatedAt"`
	DeletedAt         *util.LocalTime             `json:"deletedAt,omitempty"`
	DeletedBy         interface{}                 `json:"deletedBy,omitempty"`
}

// CounterLineServiceProperty representa las propiedades de servicio de una línea de contador
type CounterLineServiceProperty struct {
	Service    *string                `json:"service"`
	Properties map[string]interface{} `json:"properties"`
}

// XPlorCounterLines representa una colección de líneas de contador
type XPlorCounterLines struct {
	Context      string             `json:"@context"`
	ID           string             `json:"@id"`
	Type         string             `json:"@type"`
	CounterLines []XPlorCounterLine `json:"hydra:member"`
	Pagination   *HydraView         `json:"hydra:view,omitempty"`
}

// Métodos para XPlorCounterLine
func (cl *XPlorCounterLine) ContactIDValue() (string, error) {
	return ExtractID(cl.ContactID, " nil contact ID field")
}

func (cl *XPlorCounterLine) UnitID() (string, error) {
	return ExtractID(cl.Unit, "unit ID field is nil")
}

func (cl *XPlorCounterLine) ArticleIDValue() (string, error) {
	return ExtractID(cl.ArticleID, "article ID field is nil")
}

func (cl *XPlorCounterLine) CounterMovementIDs() ([]string, error) {
	if len(cl.CounterMovements) == 0 {
		return nil, errors.New("no counter movements available")
	}

	ids := make([]string, len(cl.CounterMovements))
	for i, movement := range cl.CounterMovements {
		id, err := ExtractIDFromString(movement, "counter movement ID field is empty")
		if err != nil {
			return nil, err
		}
		ids[i] = id
	}
	return ids, nil
}

func (cl *XPlorCounterLine) ServiceID() (string, error) {
	if cl.ServiceProperty == nil || cl.ServiceProperty.Service == nil {
		return "", errors.New("service ID field is nil")
	}
	return ExtractID(cl.ServiceProperty.Service, "service ID field is nil")
}

// Métodos para verificar estados
func (cl *XPlorCounterLine) IsActive() bool {
	now := util.LocalTime{Time: time.Now()}
	return cl.ValidFrom.Time.Before(now.Time) && cl.ValidThrough.Time.After(now.Time)
}

func (cl *XPlorCounterLine) IsExpired() bool {
	now := util.LocalTime{Time: time.Now()}
	return cl.ValidThrough.Time.Before(now.Time)
}

func (cl *XPlorCounterLine) IsNotStarted() bool {
	now := util.LocalTime{Time: time.Now()}
	return cl.ValidFrom.Time.After(now.Time)
}

func (cl *XPlorCounterLine) IsDeleted() bool {
	return cl.DeletedAt != nil
}

// Métodos para XPlorCounterLines (colección)
func (c *XPlorCounterLines) CollectionID() (string, error) {
	return ExtractIDFromString(c.ID, "collection ID field is empty")
}

func (c *XPlorCounterLines) ContextID() (string, error) {
	return ExtractIDFromString(c.Context, "context ID field is empty")
}

// Método para obtener todos los IDs de la colección
func (c *XPlorCounterLines) AllContactIDs() ([]string, error) {
	if len(c.CounterLines) == 0 {
		return nil, errors.New("no counter lines available")
	}

	contactIDs := make([]string, 0)
	for _, cl := range c.CounterLines {
		contactID, err := cl.ContactIDValue()
		if err == nil { // Solo agregar si no hay error
			contactIDs = append(contactIDs, contactID)
		}
	}

	if len(contactIDs) == 0 {
		return nil, errors.New("no contact IDs found")
	}

	return contactIDs, nil
}

// Método para obtener todas las líneas activas
func (c *XPlorCounterLines) ActiveCounterLines() []XPlorCounterLine {
	activeLines := make([]XPlorCounterLine, 0)
	for _, cl := range c.CounterLines {
		if cl.IsActive() && !cl.IsDeleted() {
			activeLines = append(activeLines, cl)
		}
	}
	return activeLines
}

// Método para obtener todas las líneas expiradas
func (c *XPlorCounterLines) ExpiredCounterLines() []XPlorCounterLine {
	expiredLines := make([]XPlorCounterLine, 0)
	for _, cl := range c.CounterLines {
		if cl.IsExpired() && !cl.IsDeleted() {
			expiredLines = append(expiredLines, cl)
		}
	}
	return expiredLines
}

// Método para obtener todas las líneas no iniciadas
func (c *XPlorCounterLines) NotStartedCounterLines() []XPlorCounterLine {
	notStartedLines := make([]XPlorCounterLine, 0)
	for _, cl := range c.CounterLines {
		if cl.IsNotStarted() && !cl.IsDeleted() {
			notStartedLines = append(notStartedLines, cl)
		}
	}
	return notStartedLines
}

// Método para obtener todas las líneas eliminadas
func (c *XPlorCounterLines) DeletedCounterLines() []XPlorCounterLine {
	deletedLines := make([]XPlorCounterLine, 0)
	for _, cl := range c.CounterLines {
		if cl.IsDeleted() {
			deletedLines = append(deletedLines, cl)
		}
	}
	return deletedLines
}
