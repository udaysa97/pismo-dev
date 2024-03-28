package types

type Metadata struct {
	TotalCount int `json:"total_count"`
}

type Entities struct {
	Id         string         `json:"id"`
	EntityType string         `json:"entityType"`
	Details    map[string]any `json:"details"`
}

type GetAllEntitiesResponseDTO struct {
	Metadata Metadata   `json:"metadata,omitempty"`
	Entities []Entities `json:"entities,omitempty"`
}
