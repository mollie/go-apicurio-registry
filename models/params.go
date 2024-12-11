package models

import (
	"net/url"
	"strconv"
	"strings"
)

// ========================================
// SECTION: Params
// ========================================

// SearchArtifactsParams represents the optional parameters for searching artifacts.
type SearchArtifactsParams struct {
	Name         string       // Filter by artifact name
	Offset       int          // Default: 0
	Limit        int          // Default: 20
	Order        Order        // Default: "asc", Enum: "asc", "desc"
	OrderBy      OrderBy      // Field to sort by, e.g., "name", "createdOn"
	Labels       []string     // Filter by one or more name/value labels
	Description  string       // Filter by description
	GroupID      string       // Filter by artifact group
	GlobalID     int64        // Filter by globalId
	ContentID    int64        // Filter by contentId
	ArtifactID   string       // Filter by artifactId
	ArtifactType ArtifactType // Filter by artifact type (e.g., AVRO, JSON)
}

// ToQuery converts the SearchArtifactsParams struct to URL query parameters.
func (p *SearchArtifactsParams) ToQuery() url.Values {
	query := url.Values{}

	if p.Name != "" {
		query.Set("name", p.Name)
	}
	if p.Offset != 0 {
		query.Set("offset", strconv.Itoa(p.Offset))
	}
	if p.Limit != 0 {
		query.Set("limit", strconv.Itoa(p.Limit))
	}
	if p.Order != "" {
		query.Set("order", string(p.Order))
	}
	if p.OrderBy != "" {
		query.Set("orderby", string(p.OrderBy))
	}
	if len(p.Labels) > 0 {
		query.Set("labels", strings.Join(p.Labels, ","))
	}
	if p.Description != "" {
		query.Set("description", p.Description)
	}
	if p.GroupID != "" {
		query.Set("groupId", p.GroupID)
	}
	if p.GlobalID != 0 {
		query.Set("globalId", strconv.FormatInt(p.GlobalID, 10))
	}
	if p.ContentID != 0 {
		query.Set("contentId", strconv.FormatInt(p.ContentID, 10))
	}
	if p.ArtifactID != "" {
		query.Set("artifactId", p.ArtifactID)
	}
	if p.ArtifactType != "" {
		query.Set("artifactType", string(p.ArtifactType))
	}

	return query
}

// SearchArtifactsByContentParams represents the query parameters for the search by content API.
type SearchArtifactsByContentParams struct {
	Canonical    bool    // Canonicalize the content
	ArtifactType string  // Artifact type (e.g., AVRO, JSON)
	GroupID      string  // Filter by group ID
	Offset       int     // Number of artifacts to skip
	Limit        int     // Number of artifacts to return
	Order        Order   // Sort order (asc, desc)
	OrderBy      OrderBy // Field to sort by
}

// ToQuery converts the SearchArtifactsByContentParams struct to query parameters.
func (p *SearchArtifactsByContentParams) ToQuery() url.Values {
	query := url.Values{}

	if p.Canonical {
		query.Set("canonical", "true")
	}
	if p.ArtifactType != "" {
		query.Set("artifactType", p.ArtifactType)
	}
	if p.GroupID != "" {
		query.Set("groupId", p.GroupID)
	}
	if p.Offset != 0 {
		query.Set("offset", strconv.Itoa(p.Offset))
	}
	if p.Limit != 0 {
		query.Set("limit", strconv.Itoa(p.Limit))
	}
	if p.Order != "" {
		query.Set("order", string(p.Order))
	}
	if p.OrderBy != "" {
		query.Set("orderby", string(p.OrderBy))
	}

	return query
}

// CreateArtifactParams represents the parameters for creating an artifact.
type CreateArtifactParams struct {
	IfExists  IfExistsType // IfExists behavior @See IfExistsType
	Canonical bool         // Indicates whether to canonicalize the artifact content.
	DryRun    bool         // If true, no changes are made, only checks are performed.
}

// ToQuery converts the parameters into a query string.
func (p *CreateArtifactParams) ToQuery() url.Values {
	query := url.Values{}
	if p.IfExists != "" {
		query.Set("ifExists", string(p.IfExists))
	}
	if p.Canonical {
		query.Set("canonical", "true")
	}
	if p.DryRun {
		query.Set("dryRun", "true")
	}
	return query
}

// ListArtifactReferencesByGlobalIDParams represents the optional parameters for listing references by global ID.
type ListArtifactReferencesByGlobalIDParams struct {
	RefType RefType
}

// ToQuery converts the params struct to URL query parameters.
func (p *ListArtifactReferencesByGlobalIDParams) ToQuery() url.Values {
	query := url.Values{}
	if p != nil && p.RefType != "" {
		query.Set("refType", string(p.RefType))
	}
	return query
}

// ListArtifactsInGroupParams represents the query parameters for listing artifacts in a group.
type ListArtifactsInGroupParams struct {
	Limit   int    // Number of artifacts to return (default: 20)
	Offset  int    // Number of artifacts to skip (default: 0)
	Order   string // Enum: "asc", "desc"
	OrderBy string // Enum: "groupId", "artifactId", "createdOn", "modifiedOn", "artifactType", "name"
}

// ToQuery converts the ListArtifactsInGroupParams struct to query parameters.
func (p *ListArtifactsInGroupParams) ToQuery() url.Values {
	query := url.Values{}
	if p.Limit != 0 {
		query.Set("limit", strconv.Itoa(p.Limit))
	}
	if p.Offset != 0 {
		query.Set("offset", strconv.Itoa(p.Offset))
	}
	if p.Order != "" {
		query.Set("order", p.Order)
	}
	if p.OrderBy != "" {
		query.Set("orderby", p.OrderBy)
	}
	return query
}

// ArtifactVersionReferencesParams represents the query parameters for GetArtifactVersionReferences.
type ArtifactVersionReferencesParams struct {
	RefType RefType // "INBOUND" or "OUTBOUND"
}

// ToQuery converts the ArtifactVersionReferencesParams struct to URL query parameters.
func (p *ArtifactVersionReferencesParams) ToQuery() url.Values {
	query := url.Values{}
	if p != nil && p.RefType != "" {
		query.Set("refType", string(p.RefType))
	}
	return query
}

// ArtifactReferenceParams represents the query parameters for artifact references.
type ArtifactReferenceParams struct {
	HandleReferencesType HandleReferencesType
}

// ToQuery converts the ArtifactReferenceParams into URL query parameters.
func (p ArtifactReferenceParams) ToQuery() url.Values {
	query := url.Values{}
	if p.HandleReferencesType != "" {
		query.Set("references", string(p.HandleReferencesType))
	}
	return query
}

// SearchVersionParams represents the query parameters for searching artifact versions.
type SearchVersionParams struct {
	Version      string
	Offset       int
	Limit        int
	Order        Order
	OrderBy      OrderBy
	Labels       []string
	Description  string
	GroupID      string
	GlobalID     int64
	ContentID    int64
	ArtifactID   string
	Name         string
	State        State
	ArtifactType ArtifactType
}

// ToQuery converts the SearchVersionParams into URL query parameters.
func (p *SearchVersionParams) ToQuery() url.Values {
	query := url.Values{}
	if p.Version != "" {
		query.Set("version", p.Version)
	}
	if p.Offset > 0 {
		query.Set("offset", strconv.Itoa(p.Offset))
	}
	if p.Limit > 0 {
		query.Set("limit", strconv.Itoa(p.Limit))
	}
	if p.Order != "" {
		query.Set("order", string(p.Order))
	}
	if p.OrderBy != "" {
		query.Set("orderby", string(p.OrderBy))
	}
	if len(p.Labels) > 0 {
		query.Set("labels", strings.Join(p.Labels, ","))
	}
	if p.Description != "" {
		query.Set("description", p.Description)
	}
	if p.GroupID != "" {
		query.Set("groupId", p.GroupID)
	}
	if p.GlobalID > 0 {
		query.Set("globalId", strconv.FormatInt(p.GlobalID, 10))
	}
	if p.ContentID > 0 {
		query.Set("contentId", strconv.FormatInt(p.ContentID, 10))
	}
	if p.ArtifactID != "" {
		query.Set("artifactId", p.ArtifactID)
	}
	if p.Name != "" {
		query.Set("name", p.Name)
	}
	if p.State != "" {
		query.Set("state", string(p.State))
	}
	if p.ArtifactType != "" {
		query.Set("artifactType", string(p.ArtifactType))
	}
	return query
}

// SearchVersionByContentParams defines the query parameters for searching artifact versions by content.
type SearchVersionByContentParams struct {
	Canonical    *bool
	ArtifactType ArtifactType
	Offset       int
	Limit        int
	Order        Order
	OrderBy      OrderBy
	GroupID      string
	ArtifactID   string
}

// ToQuery converts the SearchVersionByContentParams into URL query parameters.
func (p *SearchVersionByContentParams) ToQuery() url.Values {
	query := url.Values{}
	if p.Canonical != nil {
		query.Set("canonical", strconv.FormatBool(*p.Canonical))
	}
	if p.ArtifactType != "" {
		query.Set("artifactType", string(p.ArtifactType))
	}
	if p.Offset > 0 {
		query.Set("offset", strconv.Itoa(p.Offset))
	}
	if p.Limit > 0 {
		query.Set("limit", strconv.Itoa(p.Limit))
	}
	if p.Order != "" {
		query.Set("order", string(p.Order))
	}
	if p.OrderBy != "" {
		query.Set("orderby", string(p.OrderBy))
	}
	if p.GroupID != "" {
		query.Set("groupId", p.GroupID)
	}
	if p.ArtifactID != "" {
		query.Set("artifactId", p.ArtifactID)
	}
	return query
}
