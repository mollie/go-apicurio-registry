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

type GroupOrderBy string

const (
	GroupOrderByName       GroupOrderBy = "name"
	GroupOrderByGroupId    GroupOrderBy = "groupId"
	GroupOrderByCreatedOn  GroupOrderBy = "createdOn"
	GroupOrderByModifiedOn GroupOrderBy = "modifiedOn"
)

type ArtifactSortBy string

const (
	ArtifactSortByName       ArtifactSortBy = "name"
	ArtifactSortByType       ArtifactSortBy = "artifactType"
	ArtifactSortByCreatedOn  ArtifactSortBy = "createdOn"
	ArtifactSortByModifiedOn ArtifactSortBy = "modifiedOn"
	ArtifactSortByGroupID    ArtifactSortBy = "groupId"
	ArtifactSortByArtifactID ArtifactSortBy = "artifactId"
)

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

type GetArtifactByGlobalIDParams struct {
	HandleReferencesType HandleReferencesType `validate:"omitempty,oneof=PRESERVE DEREFERENCE REWRITE"`
	ReturnArtifactType   bool                 `validate:"omitempty"`
}

func (p *GetArtifactByGlobalIDParams) Validate() error {
	return structValidator.Struct(p)
}

func (p *GetArtifactByGlobalIDParams) ToQuery() url.Values {
	query := url.Values{}
	if p.HandleReferencesType != "" {
		query.Set("references", string(p.HandleReferencesType))
	}
	if p.ReturnArtifactType {
		query.Set("returnType", "true")
	}
	return query

}

// SearchArtifactsParams represents the optional parameters for searching artifacts.
type SearchArtifactsParams struct {
	Name         string         // Filter by artifact name
	Offset       int            `validate:"omitempty,gte=0"`                // Default: 0
	Limit        int            `validate:"omitempty,gte=0"`                // Default: 20
	Order        Order          `validate:"omitempty,oneof=asc desc"`       // Default: "asc", Enum: "asc", "desc"
	OrderBy      ArtifactSortBy `validate:"omitempty,oneof=name createdOn"` // Field to sort by, e.g., "name", "createdOn"
	Labels       []string       // Filter by one or more name/value labels
	Description  string         // Filter by description
	GroupID      string         `validate:"omitempty,groupid"` // Filter by artifact group
	GlobalID     int64          // Filter by globalId
	ContentID    int64          // Filter by contentId
	ArtifactID   string         `validate:"omitempty,artifactid"`   // Filter by artifactId
	ArtifactType ArtifactType   `validate:"omitempty,artifacttype"` // Filter by artifact type (e.g., AVRO, JSON)
}

// Validate validates the SearchArtifactsParams struct.
func (p *SearchArtifactsParams) Validate() error {
	return structValidator.Struct(p)
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
	Canonical    bool           // Canonicalize the content
	ArtifactType string         `validate:"omitempty,artifacttype"`         // Artifact type (e.g., AVRO, JSON)
	GroupID      string         `validate:"omitempty,groupid"`              // Filter by group ID
	Offset       int            `validate:"omitempty,gte=0"`                // Number of artifacts to skip
	Limit        int            `validate:"omitempty,gte=0"`                // Number of artifacts to return
	Order        Order          `validate:"omitempty,oneof=asc desc"`       // Sort order (asc, desc)
	OrderBy      ArtifactSortBy `validate:"omitempty,oneof=name createdOn"` // Field to sort by
}

// Validate validates the SearchArtifactsByContentParams struct.
func (p *SearchArtifactsByContentParams) Validate() error {
	return structValidator.Struct(p)
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
	IfExists  IfExistsType `validate:"oneof=FAIL CREATE_VERSION FIND_OR_CREATE_VERSION"` // IfExists behavior @See IfExistsType
	Canonical bool         // Indicates whether to canonicalize the artifact content.
	DryRun    bool         // If true, no changes are made, only checks are performed.
}

// Validate validates the CreateArtifactParams struct.
func (p *CreateArtifactParams) Validate() error {
	return structValidator.Struct(p)
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
	RefType RefType `validate:"omitempty,oneof=INBOUND OUTBOUND"`
}

// Validate validates the ListArtifactReferencesByGlobalIDParams struct.
func (p *ListArtifactReferencesByGlobalIDParams) Validate() error {
	return structValidator.Struct(p)
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
	Offset  int            `validate:"omitempty,gte=0"`                // Number of artifacts to skip
	Limit   int            `validate:"omitempty,gte=0"`                // Number of artifacts to return
	Order   Order          `validate:"omitempty,oneof=asc desc"`       // Sort order (asc, desc)
	OrderBy ArtifactSortBy `validate:"omitempty,oneof=name createdOn"` // Field to sort by
}

// Validate validates the ListArtifactsInGroupParams struct.
func (p *ListArtifactsInGroupParams) Validate() error {
	return structValidator.Struct(p)
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
		query.Set("order", string(p.Order))
	}
	if p.OrderBy != "" {
		query.Set("orderby", string(p.OrderBy))
	}
	return query
}

// ArtifactVersionReferencesParams represents the query parameters for GetArtifactVersionReferences.
type ArtifactVersionReferencesParams struct {
	RefType RefType `validate:"omitempty,oneof=INBOUND OUTBOUND"` // "INBOUND" or "OUTBOUND"
}

// Validate validates the ListArtifactsInGroupParams struct.
func (p *ArtifactVersionReferencesParams) Validate() error {
	return structValidator.Struct(p)
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

// Validate validates the ArtifactReferenceParams struct.
func (p *ArtifactReferenceParams) Validate() error {
	return structValidator.Struct(p)
}

// ToQuery converts the ArtifactReferenceParams into URL query parameters.
func (p *ArtifactReferenceParams) ToQuery() url.Values {
	query := url.Values{}
	if p.HandleReferencesType != "" {
		query.Set("references", string(p.HandleReferencesType))
	}
	return query
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

// ListArtifactsVersionsParams represents the query parameters for listing artifacts in a group.
type ListArtifactsVersionsParams struct {
	Limit   int           `validate:"omitempty,gte=0"`                        // Number of artifacts to return (default: 20)
	Offset  int           `validate:"omitempty,gte=0"`                        // Number of artifacts to skip (default: 0)
	Order   Order         `validate:"omitempty,oneof=asc desc"`               // Enum: "asc", "desc"
	OrderBy VersionSortBy `validate:"omitempty,oneof=name version createdOn"` // Enum: only: name version createdOn
}

func (p *ListArtifactsVersionsParams) Validate() error {
	return structValidator.Struct(p)
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

type ListBranchesParams struct {
	Offset int `validate:"omitempty,gte=0"` // Number of branches to skip
	Limit  int `validate:"omitempty,gte=0"` // Number of branches to return
}

func (p *ListBranchesParams) Validate() error {
	return structValidator.Struct(p)
}

func (p *ListBranchesParams) ToQuery() url.Values {
	query := url.Values{}
	if p.Offset != 0 {
		query.Set("offset", strconv.Itoa(p.Offset))
	}
	if p.Limit != 0 {
		query.Set("limit", strconv.Itoa(p.Limit))
	}
	return query

}

// CustomValidationFunctions registers custom validation functions with the validator.
func CustomValidationFunctions(validate *validator.Validate) error {
	// Validation for Version: ^[a-zA-Z0-9._\-+]{1,256}$
	versionRegex := regexp.MustCompile(`^[a-zA-Z0-9._\-+]{1,256}$`)
	err := validate.RegisterValidation("version", func(fl validator.FieldLevel) bool {
		return versionRegex.MatchString(fl.Field().String())
	})
	if err != nil {
		return err
	}

	// Validation for BranchID: ^[a-zA-Z0-9._\-+]{1,256}$
	branchIDRegex := regexp.MustCompile(`^[a-zA-Z0-9._\-+]{1,256}$`)
	err = validate.RegisterValidation("branchid", func(fl validator.FieldLevel) bool {
		return branchIDRegex.MatchString(fl.Field().String())
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
	if err := CustomValidationFunctions(structValidator); err != nil {
		panic(err)
	}
}
