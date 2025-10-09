package xplorentities

import "path"

type XPloreCoaches struct {
	Context    *string       `json:"@context"`
	ID         *string       `json:"@id"`
	Type       *string       `json:"@type"`
	Coaches    []XPloreCoach `json:"hydra:member"`
	Pagination *HydraView    `json:"hydra:view"`
}

type XPloreCoach struct {
	Id            *string   `json:"@id"`
	Type          *string   `json:"@type"`
	GivenName     *string   `json:"givenName"`
	FamilyName    *string   `json:"familyName"`
	AlternateName *string   `json:"alternateName"`
	Email         *string   `json:"email"`
	Mobile        *string   `json:"mobile"`
	Activities    []*string `json:"activities"`
	CreatedAt     *string   `json:"createdAt"`
	CreatedBy     *string   `json:"createdBy"`
	ArchivedAt    *string   `json:"archivedAt"`
	ArchivedBy    *string   `json:"archivedBy"`
}

func (c XPloreCoach) ActivityIDs() []string {
	var ids []string
	for _, act := range c.Activities {
		if act == nil {
			continue
		}
		base := path.Base(*act)
		ids = append(ids, base)
	}
	return ids
}

func (c XPloreCoach) CoachID() (string, error) {
	return ExtractID(c.Id, "coach ID field is nil")
}
