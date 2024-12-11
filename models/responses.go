package models

// ========================================
// SECTION: Responses
// ========================================

// SearchArtifactsAPIResponse represents the response from the search artifacts API.
type SearchArtifactsAPIResponse struct {
	Artifacts []SearchedArtifact `json:"artifacts"`
	Count     int                `json:"count"`
}

// ListArtifactsResponse represents the response from the list artifacts API.
type ListArtifactsResponse struct {
	Artifacts []SearchedArtifact `json:"artifacts"`
	Count     int                `json:"count"`
}

// CreateArtifactResponse represents the response from the create artifact API.
type CreateArtifactResponse struct {
	Artifact ArtifactDetail `json:"artifact"`
}

// ArtifactVersionListResponse represents the response of GetArtifactVersions.
type ArtifactVersionListResponse struct {
	Count    int               `json:"count"`
	Versions []ArtifactVersion `json:"versions"`
}

type StateResponse struct {
	State State `json:"state"`
}

type GlobalRuleResponse struct {
	RuleType Rule      `json:"ruleType"`
	Config   RuleLevel `json:"config"`
}
