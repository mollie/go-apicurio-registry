package models

// ========================================
// SECTION: Requests
// ========================================

// CreateArtifactRequest represents the request to create an artifact.
type CreateArtifactRequest struct {
	ArtifactID   string               `json:"artifactId,omitempty"`
	ArtifactType ArtifactType         `json:"artifactType"`
	Name         string               `json:"name,omitempty"`
	Description  string               `json:"description,omitempty"`
	Labels       map[string]string    `json:"labels,omitempty"`
	FirstVersion CreateVersionRequest `json:"firstVersion,omitempty"`
}

// CreateVersionRequest represents the request to create a version for an artifact.
type CreateVersionRequest struct {
	Version     string               `json:"version"`
	Content     CreateContentRequest `json:"content" validate:"required"`
	Name        string               `json:"name,omitempty"`
	Description string               `json:"description,omitempty"`
	Labels      map[string]string    `json:"labels,omitempty"`
	Branches    []string             `json:"branches,omitempty"`
	IsDraft     bool                 `json:"isDraft"`
}

// CreateContentRequest represents the content of an artifact.
type CreateContentRequest struct {
	Content     string              `json:"content"`
	References  []ArtifactReference `json:"references,omitempty"`
	ContentType string              `json:"contentType"`
}

// UpdateArtifactMetadataRequest represents the metadata update request.
type UpdateArtifactMetadataRequest struct {
	Name        string            `json:"name,omitempty"`        // Editable name
	Description string            `json:"description,omitempty"` // Editable description
	Labels      map[string]string `json:"labels,omitempty"`      // Editable labels
	Owner       string            `json:"owner,omitempty"`       // Editable owner
}

type StateRequest struct {
	State State `json:"state"`
}

type CreateUpdateGlobalRuleRequest struct {
	RuleType Rule      `json:"ruleType"`
	Config   RuleLevel `json:"config"`
}
