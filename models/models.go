package models

// ========================================
// SECTION: Models
// ========================================

// ArtifactReference represents a reference to an artifact.
type ArtifactReference struct {
	GroupID    string `json:"groupId"`
	ArtifactID string `json:"artifactId"`
	Version    string `json:"version"`
	Name       string `json:"name"`
}

// SearchedArtifact represents the search result of an artifact.
type SearchedArtifact struct {
	GroupId      string       `json:"groupId"`
	ArtifactId   string       `json:"artifactId"`
	Name         string       `json:"name"`
	Description  string       `json:"description"`
	ArtifactType ArtifactType `json:"artifactType"`
	Owner        string       `json:"owner"`
	CreatedOn    string       `json:"createdOn"`
	ModifiedBy   string       `json:"modifiedBy"`
	ModifiedOn   string       `json:"modifiedOn"`
}

// ArtifactContent represents the content of an artifact + the type of the artifact.
type ArtifactContent struct {
	Content      string       `json:"content"`
	ArtifactType ArtifactType `json:"artifactType"`
}

// ArtifactDetail represents the detailed information about an artifact.
type ArtifactDetail struct {
	GroupID     string            `json:"groupId"`
	ArtifactID  string            `json:"artifactId"`
	Name        string            `json:"name"`
	Description string            `json:"description"`
	Version     string            `json:"version"`
	CreatedOn   string            `json:"createdOn"`
	ModifiedOn  string            `json:"modifiedOn"`
	ContentID   int64             `json:"contentId"`
	Labels      map[string]string `json:"labels"`
}

// BaseMetadata contains common fields shared by both artifact and artifact version metadata.
type BaseMetadata struct {
	GroupID      string            `json:"groupId"`
	ArtifactID   string            `json:"artifactId"`
	Name         string            `json:"name"`
	Description  string            `json:"description"`
	ArtifactType string            `json:"artifactType"`
	Owner        string            `json:"owner"`
	CreatedOn    string            `json:"createdOn"`
	Labels       map[string]string `json:"labels"`
}

// ArtifactVersionMetadata represents metadata for a single artifact version.
type ArtifactVersionMetadata struct {
	BaseMetadata
	Version   string `json:"version"`
	GlobalID  int64  `json:"globalId"`
	ContentID int64  `json:"contentId"`
}

// ArtifactMetadata represents metadata for an artifact.
type ArtifactMetadata struct {
	BaseMetadata
	ModifiedBy string `json:"modifiedBy"`
	ModifiedOn string `json:"modifiedOn"`
}

// ArtifactComment represents a comment on a specific artifact version.
// It's used in the response of GetArtifactVersionComments
type ArtifactComment struct {
	CommentID string `json:"commentId"` // Unique identifier for the comment.
	Value     string `json:"value"`     // The content of the comment.
	Owner     string `json:"owner"`     // The user who created the comment.
	CreatedOn string `json:"createdOn"` // The timestamp when the comment was created.
}

// ArtifactVersion represents a single version of an artifact. it has the minimum information
// required to identify an artifact version. while ArtifactVersionDetailed has more information
type ArtifactVersion struct {
	Version      string       `json:"version" validate:"required"`                                                  // A single version of the artifact
	Owner        string       `json:"owner" validate:"required"`                                                    // Owner of the artifact version
	CreatedOn    string       `json:"createdOn" validate:"required"`                                                // Creation timestamp
	ArtifactType ArtifactType `json:"artifactType" validate:"required"`                                             // Type of the artifact
	GlobalID     int64        `json:"globalId" validate:"required"`                                                 // Global identifier for the artifact version
	State        State        `json:"state,omitempty" validate:"omitempty,oneof=ENABLED DISABLED DEPRECATED DRAFT"` // State of the artifact version
	ContentID    int64        `json:"contentId" validate:"required"`                                                // Content ID of the artifact version
	ArtifactID   string       `json:"artifactId" validate:"required,max=512"`                                       // Artifact ID
	GroupID      string       `json:"groupId,omitempty" validate:"omitempty,max=512"`                               // Artifact group ID
	ModifiedBy   string       `json:"modifiedBy,omitempty"`                                                         // User who last modified the artifact version
	ModifiedOn   string       `json:"modifiedOn,omitempty"`                                                         // Last modification timestamp
}

// ArtifactVersionDetailed represents a single version of an artifact with additional information.
type ArtifactVersionDetailed struct {
	ArtifactVersion                   // Embedding ArtifactVersion
	Name            string            `json:"name,omitempty"`        // Name of the artifact version
	Description     string            `json:"description,omitempty"` // Description of the artifact version
	Labels          map[string]string `json:"labels,omitempty"`      // User-defined name-value pairs
}
