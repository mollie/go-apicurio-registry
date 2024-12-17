package models

// ========================================
// SECTION: Requests
// ========================================

// CreateArtifactRequest represents the request to create an artifact.
type CreateArtifactRequest struct {
	ArtifactID   string               `json:"artifactId,omitempty" validate:"required,artifactid"`
	ArtifactType ArtifactType         `json:"artifactType" validate:"omitempty,artifacttype"`
	Name         string               `json:"name,omitempty"`
	Description  string               `json:"description,omitempty"`
	Labels       map[string]string    `json:"labels,omitempty"`
	FirstVersion CreateVersionRequest `json:"firstVersion,omitempty"`
}

func (r *CreateArtifactRequest) Validate() error {
	return structValidator.Struct(r)
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

func (r *CreateVersionRequest) Validate() error {
	return structValidator.Struct(r)
}

// CreateContentRequest represents the content of an artifact.
type CreateContentRequest struct {
	Content     string              `json:"content" validate:"required"`
	References  []ArtifactReference `json:"references,omitempty"`
	ContentType string              `json:"contentType" validate:"required"`
}

func (r *CreateContentRequest) Validate() error {
	return structValidator.Struct(r)
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

type CreateUpdateRuleRequest struct {
	RuleType Rule      `json:"ruleType"`
	Config   RuleLevel `json:"config"`
}

type CreateGroupRequest struct {
	GroupID     string            `json:"groupId"`
	Description string            `json:"description"`
	Labels      map[string]string `json:"labels"`
}

type UpdateGroupRequest struct {
	Description string            `json:"description"`
	Labels      map[string]string `json:"labels"`
}

type CreateBranchRequest struct {
	BranchID    string `json:"branchId" validate:"required,branchid"`
	Description string `json:"description,omitempty"`
}

func (r *CreateBranchRequest) Validate() error {
	return structValidator.Struct(r)
}

type UpdateBranchMetaDataRequest struct {
	Description string `json:"description,omitempty"`
}
