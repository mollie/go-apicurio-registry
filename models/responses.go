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

type RuleResponse struct {
	RuleType Rule      `json:"ruleType"`
	Config   RuleLevel `json:"config"`
}

type SystemInfoResponse struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Version     string `json:"version"`
	BuiltOn     string `json:"builtOn"`
}

type SystemResourceLimitInfoResponse struct {
	MaxTotalSchemasCount              int `json:"maxTotalSchemasCount"`
	MaxSchemaSizeBytes                int `json:"maxSchemaSizeBytes"`
	MaxArtifactsCount                 int `json:"maxArtifactsCount"`
	MaxVersionsPerArtifactCount       int `json:"maxVersionsPerArtifactCount"`
	MaxArtifactPropertiesCount        int `json:"maxArtifactPropertiesCount"`
	MaxPropertyKeySizeBytes           int `json:"maxPropertyKeySizeBytes"`
	MaxPropertyValueSizeBytes         int `json:"maxPropertyValueSizeBytes"`
	MaxArtifactLabelsCount            int `json:"maxArtifactLabelsCount"`
	MaxLabelSizeBytes                 int `json:"maxLabelSizeBytes"`
	MaxArtifactNameLengthChars        int `json:"maxArtifactNameLengthChars"`
	MaxArtifactDescriptionLengthChars int `json:"maxArtifactDescriptionLengthChars"`
	MaxRequestsPerSecondCount         int `json:"maxRequestsPerSecondCount"`
}

type UIConfig struct {
	ContextPath   string `json:"contextPath"`
	NavPrefixPath string `json:"navPrefixPath"`
	OaiDocsUrl    string `json:"oaiDocsUrl"`
}

// AuthOptions represents the options for authentication configuration.
type AuthOptions struct {
	Url         string `json:"url"`
	RedirectUri string `json:"redirectUri"`
	ClientId    string `json:"clientId"`
}

// AuthConfig represents the authentication-related configuration settings.
type AuthConfig struct {
	Type        string      `json:"type"`
	RbacEnabled bool        `json:"rbacEnabled"`
	ObacEnabled bool        `json:"obacEnabled"`
	Options     AuthOptions `json:"options"`
}

// FeatureFlags represents the available feature flags in the system.
type FeatureFlags struct {
	ReadOnly        bool `json:"readOnly"`
	Breadcrumbs     bool `json:"breadcrumbs"`
	RoleManagement  bool `json:"roleManagement"`
	Settings        bool `json:"settings"`
	DeleteGroup     bool `json:"deleteGroup"`
	DeleteArtifact  bool `json:"deleteArtifact"`
	DeleteVersion   bool `json:"deleteVersion"`
	DraftMutability bool `json:"draftMutability"`
}

// SystemUIConfigResponse represents the overall UI configuration response.
type SystemUIConfigResponse struct {
	Ui       UIConfig     `json:"ui"`
	Auth     AuthConfig   `json:"auth"`
	Features FeatureFlags `json:"features"`
}

type GroupInfoResponse struct {
	Groups []GroupInfo `json:"groups"`
	Count  int         `json:"count"`
}

type ArtifactTypeResponse struct {
	Name ArtifactType `json:"name"`
}
