package xplorentities

import (
	"net/url"
	"strings"

	"github.com/angelbarreiros/XPlorGo/util"
)

// XPlorContactImages represents the Hydra collection for contact image files.
type XPlorContactImages struct {
	Context       string              `json:"@context"`
	ID            string              `json:"@id"`
	Type          string              `json:"@type"`
	ContactImages []XPlorContactImage `json:"hydra:member"`
	Pagination    *HydraView          `json:"hydra:view,omitempty"`
	Search        *HydraSearch        `json:"hydra:search,omitempty"`
}

// XPlorContactImage represents a contact image file resource.
type XPlorContactImage struct {
	ID           *string         `json:"@id"`
	Type         string          `json:"@type"`
	Contact      *string         `json:"contact,omitempty"`
	ContactID    *string         `json:"contactId,omitempty"`
	ContentURL   *string         `json:"contentUrl,omitempty"`
	FilePath     *string         `json:"filePath,omitempty"`
	FileName     *string         `json:"fileName,omitempty"`
	OriginalName *string         `json:"originalName,omitempty"`
	MimeType     *string         `json:"mimeType,omitempty"`
	Size         *int            `json:"size,omitempty"`
	CreatedAt    *util.LocalTime `json:"createdAt,omitempty"`
	UpdatedAt    *util.LocalTime `json:"updatedAt,omitempty"`
}

// ContactImageID extracts the contact image ID from the @id field.
func (ci XPlorContactImage) ContactImageID() (string, error) {
	return ExtractID(ci.ID, "contact image ID field is nil")
}

// ContactIDValue extracts the contact ID from contactId or contact.
func (ci XPlorContactImage) ContactIDValue() (string, error) {
	if ci.ContactID != nil {
		return ExtractID(ci.ContactID, "contact ID field is nil")
	}
	return ExtractID(ci.Contact, "contact field is nil")
}

// XPlorContactImagesParams represents search parameters for contact image files.
type XPlorContactImagesParams struct {
	ContactImageID  string
	ContactImageIDs []string
	ContactID       string
	ContactIDs      []string
}

// ToValues converts the params to url.Values for query parameters.
func (p XPlorContactImagesParams) ToValues(values *url.Values) {
	contactImageID := strings.TrimSpace(p.ContactImageID)
	if contactImageID != "" {
		values.Set("id", contactImageID)
	}
	for _, id := range p.ContactImageIDs {
		if strings.TrimSpace(id) != "" {
			values.Add("id[]", id)
		}
	}

	contactID := strings.TrimSpace(p.ContactID)
	if contactID != "" {
		values.Set("contact", contactID)
	}
	for _, id := range p.ContactIDs {
		if strings.TrimSpace(id) != "" {
			values.Add("contact[]", id)
		}
	}
}
