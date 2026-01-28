package xplorentities

// Colección de NetworkNodes
type XPlorNetworkNodes struct {
	Context      string             `json:"@context"`
	ID           string             `json:"@id"`
	Type         string             `json:"@type"`
	NetworkNodes []XPlorNetworkNode `json:"hydra:member"`
	Pagination   HydraView          `json:"hydra:view"`
}

// Entidad NetworkNode
type XPlorNetworkNode struct {
	NotworkNodeID *string `json:"@id"`
	Type          string  `json:"@type"`
	ID            int     `json:"id"`
	Name          string  `json:"name"`
	Alias         *string `json:"alias"`
	NodeType      string  `json:"type"`
	ClubID        *string `json:"clubId"`
	Children      any     `json:"children"`
}

// ---------------- Métodos helpers ----------------

// NetworkNodeID extracts the network node ID from the @id field
func (n XPlorNetworkNode) NetworkNodeID() (string, error) {
	return ExtractID(n.NotworkNodeID, "network node hydra ID field is nil")
}

// ClubIDValue extracts the club ID from the clubId field if it exists
func (n XPlorNetworkNode) ClubIDValue() (string, error) {
	return ExtractID(n.ClubID, "club ID field is nil")
}

// IsClub checks if the network node type is "club"
func (n XPlorNetworkNode) IsClub() bool {
	return n.NodeType == "club"
}

// IsGroup checks if the network node type is "group" or "franchise"
func (n XPlorNetworkNode) IsGroup() bool {
	return n.NodeType == "group"
}

// IsFranchise checks if the network node type is "franchise"
func (n XPlorNetworkNode) IsFranchise() bool {
	return n.NodeType == "franchise"
}
