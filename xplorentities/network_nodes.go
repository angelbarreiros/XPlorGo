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
	NotworkNodeID *string            `json:"@id"`
	Type          string             `json:"@type"`
	ID            int                `json:"id"`
	Name          string             `json:"name"`
	Alias         *string            `json:"alias"`
	NodeType      string             `json:"type"`
	ClubID        *string            `json:"clubId"`
	Children      []XPlorNetworkNode `json:"children"`
}

// ---------------- Métodos helpers ----------------

// Devuelve el ID numérico a partir de @id (ej: "/enjoy/network_nodes/1359" → 1359)
func (n XPlorNetworkNode) NetworkNodeID() (string, error) {
	return ExtractID(n.NotworkNodeID, "network node hydra ID field is nil")
}

// Devuelve el ClubID numérico si existe (ej: "/enjoy/clubs/1249" → 1249)
func (n XPlorNetworkNode) ClubIDValue() (string, error) {
	return ExtractID(n.ClubID, "club ID field is nil")
}

// Saber si es un "club"
func (n XPlorNetworkNode) IsClub() bool {
	return n.NodeType == "club"
}

// Saber si es un "grupo" o "franquicia"
func (n XPlorNetworkNode) IsGroup() bool {
	return n.NodeType == "group"
}

func (n XPlorNetworkNode) IsFranchise() bool {
	return n.NodeType == "franchise"
}
