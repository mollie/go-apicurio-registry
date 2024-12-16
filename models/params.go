package models

import (
	"fmt"
	"github.com/go-playground/validator/v10"
	"net/url"
	"regexp"
	"strconv"
	"strings"
)

var structValidator *validator.Validate

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
	HandleReferencesType HandleReferencesType `validate:"omitempty,oneof=PRESERVE DEREFERENCE REWRITE"`
}

// ToQuery converts the ArtifactReferenceParams into URL query parameters.
func (p *ArtifactReferenceParams) ToQuery() url.Values {
	query := url.Values{}
	if p.HandleReferencesType != "" {
		query.Set("references", string(p.HandleReferencesType))
	}
	return query
}

// Validate validates the ArtifactReferenceParams struct.
func (p *ArtifactReferenceParams) Validate() error {
	return structValidator.Struct(p)
}

// SearchVersionParams represents the query parameters for searching artifact versions.
type SearchVersionParams struct {
	Version      string  `validate:"omitempty,version"`
	Offset       int     `validate:"omitempty,gte=0"`
	Limit        int     `validate:"omitempty,gte=0"`
	Order        Order   `validate:"omitempty,oneof=asc desc"`
	OrderBy      OrderBy `validate:"omitempty,oneof=name createdOn"`
	Labels       map[string]string
	Description  string
	GroupID      string `validate:"omitempty,groupid"`
	GlobalID     int64
	ContentID    int64
	ArtifactID   string `validate:"omitempty,artifactid"`
	Name         string
	State        State
	ArtifactType ArtifactType `validate:"omitempty,artifacttype"`
}

// Validate validates the SearchVersionParams struct.
func (p *SearchVersionParams) Validate() error {
	return structValidator.Struct(p)
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
	if p.Labels != nil {
		for k, v := range p.Labels {
			query.Add("labels", fmt.Sprintf("%s:%s", k, v))
		}
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
	ArtifactType ArtifactType `validate:"omitempty,artifacttype"`
	Offset       int          `validate:"omitempty,gte=0"`
	Limit        int          `validate:"omitempty,gte=0"`
	Order        Order        `validate:"omitempty,oneof=asc desc"`
	OrderBy      OrderBy      `validate:"omitempty,oneof=name createdOn"`
	GroupID      string       `validate:"omitempty,groupid"`
	ArtifactID   string       `validate:"omitempty,artifactid"`
}

// Validate validates the SearchVersionByContentParams struct.
func (p *SearchVersionByContentParams) Validate() error {
	return structValidator.Struct(p)
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

type GroupOrderBy string

const (
	GroupOrderByName       GroupOrderBy = "name"
	GroupOrderByGroupId    GroupOrderBy = "groupId"
	GroupOrderByCreatedOn  GroupOrderBy = "createdOn"
	GroupOrderByModifiedOn GroupOrderBy = "modifiedOn"
)

// ListGroupsParams represents the query parameters for listing groups.
type ListGroupsParams struct {
	Limit   int          `validate:"omitempty,gte=0"` // Number of artifacts to return (default: 20)
	Offset  int          `validate:"omitempty,gte=0"` // Number of artifacts to skip (default: 0)
	Order   Order        `validate:"omitempty,oneof=asc desc"`
	OrderBy GroupOrderBy `validate:"omitempty,oneof=name createdOn"`
}

func (p *ListGroupsParams) Validate() error {
	return structValidator.Struct(p)
}

// ToQuery converts the ListGroupsParams struct to query parameters.
func (p *ListGroupsParams) ToQuery() url.Values {
	query := url.Values{}
	if p.Limit != 0 {
		query.Set("limit", strconv.Itoa(p.Limit))
	}
	if p.Offset != 0 {
		query.Set("offset", strconv.Itoa(p.Offset))
	}
	if p.Order != "" {
		query.Set("order", string(p.Order))
	}
	if p.OrderBy != "" {
		query.Set("orderby", string(p.OrderBy))
	}
	return query
}

// SearchGroupsParams represents the query parameters for searching groups.
type SearchGroupsParams struct {
	Offset      int               `validate:"omitempty,gte=0"`
	Limit       int               `validate:"omitempty,gte=0"`
	Order       Order             `validate:"omitempty,oneof=asc desc"`
	OrderBy     GroupOrderBy      `validate:"omitempty,oneof=name createdOn"`
	Labels      map[string]string `validate:"omitempty"`
	Description string            `validate:"omitempty"`
	GroupID     string            `validate:"omitempty,groupid"`
}

// Validate validates the SearchGroupsParams struct.
func (p *SearchGroupsParams) Validate() error {
	return structValidator.Struct(p)
}

// ToQuery converts the SearchGroupsParams struct to URL query parameters.
func (p *SearchGroupsParams) ToQuery() url.Values {
	query := url.Values{}
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
		for k, v := range p.Labels {
			query.Add("labels", fmt.Sprintf("%s:%s", k, v))
		}
	}
	if p.Description != "" {
		query.Set("description", p.Description)
	}
	if p.GroupID != "" {
		query.Set("groupId", p.GroupID)
	}
	return query
}

type VersionSortBy string

const (
	VersionSortByVersion    VersionSortBy = "version"
	VersionSortByGlobalID   VersionSortBy = "globalId"
	VersionSortByCreatedOn  VersionSortBy = "createdOn"
	VersionSortByModifiedOn VersionSortBy = "modifiedOn"
	VersionSortByArtifactID VersionSortBy = "artifactId"
	VersionSortByGroupID    VersionSortBy = "groupId"
	VersionSortByName       VersionSortBy = "name"
)

// ListArtifactsVersionsParams represents the query parameters for listing artifacts in a group.
type ListArtifactsVersionsParams struct {
	Limit   int           `validate:"omitempty,gte=0"`                        // Number of artifacts to return (default: 20)
	Offset  int           `validate:"omitempty,gte=0"`                        // Number of artifacts to skip (default: 0)
	Order   Order         `validate:"omitempty,oneof=asc desc"`               // Enum: "asc", "desc"
	OrderBy VersionSortBy `validate:"omitempty,oneof=name version createdOn"` // Enum: only: name version createdOn
}

// ToQuery converts the ListArtifactsInGroupParams struct to query parameters.
func (p *ListArtifactsVersionsParams) ToQuery() url.Values {
	query := url.Values{}
	if p.Limit != 0 {
		query.Set("limit", strconv.Itoa(p.Limit))
	}
	if p.Offset != 0 {
		query.Set("offset", strconv.Itoa(p.Offset))
	}
	if p.Order != "" {
		query.Set("order", string(p.Order))
	}
	if p.OrderBy != "" {
		query.Set("orderby", string(p.OrderBy))
	}
	return query
}

func (p *ListArtifactsVersionsParams) Validate() error {
	return structValidator.Struct(p)
}

// CustomValidationFuncs registers custom validation rules.
func CustomValidationFuncs(validate *validator.Validate) error {
	// Validation for Version: ^[a-zA-Z0-9._\-+]{1,256}$
	versionRegex := regexp.MustCompile(`^[a-zA-Z0-9._\-+]{1,256}$`)
	err := validate.RegisterValidation("version", func(fl validator.FieldLevel) bool {
		return versionRegex.MatchString(fl.Field().String())
	})
	if err != nil {
		return err
	}

	// Validation for ArtifactID: ^.{1,512}$
	artifactIDRegex := regexp.MustCompile(`^.{1,512}$`)
	err = validate.RegisterValidation("artifactid", func(fl validator.FieldLevel) bool {
		return artifactIDRegex.MatchString(fl.Field().String())
	})
	if err != nil {
		return err
	}

	// Validation for GroupID: ^.{1,512}$
	groupIDRegex := regexp.MustCompile(`^.{1,512}$`)
	err = validate.RegisterValidation("groupid", func(fl validator.FieldLevel) bool {
		return groupIDRegex.MatchString(fl.Field().String())
	})
	if err != nil {
		return err
	}

	// Validation for artifactTypes
	artifactTypes := map[ArtifactType]struct{}{
		Avro:     {},
		Protobuf: {},
		Json:     {},
		KConnect: {},
		OpenAPI:  {},
		AsyncAPI: {},
		GraphQL:  {},
		WSDL:     {},
		XSD:      {},
		XML:      {},
	}

	err = validate.RegisterValidation("artifacttype", func(fl validator.FieldLevel) bool {
		value := fl.Field().String()
		_, valid := artifactTypes[ArtifactType(value)]
		return valid
	})
	if err != nil {
		return err
	}

	return nil
}

func init() {
	structValidator = validator.New()
	if err := CustomValidationFuncs(structValidator); err != nil {
		panic(err)
	}
}
