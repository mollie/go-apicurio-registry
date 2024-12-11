package models

// IfExistsType represents the IfExists types for creating an artifact.
type IfExistsType string

const (
	IfExistsFail                IfExistsType = "FAIL"                   // (default) - server rejects the content with a 409 error
	IfExistsCreate              IfExistsType = "CREATE_VERSION"         // server creates a new version of the existing artifact and returns it
	IfExistsFindOrCreateVersion IfExistsType = "FIND_OR_CREATE_VERSION" // server returns an existing version that matches the provided content if such a version exists, otherwise a new version is created
)

// State represents the state of an artifact.
type State string

const (
	StateEnabled    State = "ENABLED"
	StateDisabled   State = "DISABLED"
	StateDeprecated State = "DEPRECATED"
	StateDraft      State = "DRAFT"
)

// Order represents the order of the results.
type Order string

const (
	OrderAsc  Order = "asc"
	OrderDesc Order = "desc"
)

// OrderBy represents the field to sort by.
type OrderBy string

const (
	OrderByGroupId    OrderBy = "groupId"
	OrderByArtifactId OrderBy = "artifactId"
	OrderByVersion    OrderBy = "version"
	OrderByName       OrderBy = "name"
	OrderByCreatedOn  OrderBy = "createdOn"
	OrderByModifiedOn OrderBy = "modifiedOn"
	OrderByGlobalId   OrderBy = "globalId"
)

// HandleReferencesType represents the type of handling references.
type HandleReferencesType string

const (
	HandleReferencesTypePreserve    HandleReferencesType = "PRESERVE"
	HandleReferencesTypeDereference HandleReferencesType = "DEREFERENCE"
	HandleReferencesTypeRewrite     HandleReferencesType = "REWRITE"
)

// RefType represents the type of reference.
type RefType string

const (
	OutBound RefType = "OUTBOUND"
	InBound  RefType = "INBOUND"
)

// ArtifactType represents the type of artifact.
type ArtifactType string

const (
	Avro     ArtifactType = "AVRO"     // Avro artifact type
	Protobuf ArtifactType = "PROTOBUF" // Protobuf artifact type
	Json     ArtifactType = "JSON"     // JSON artifact type
	KConnect ArtifactType = "KCONNECT" // Kafka Connect artifact type
	OpenAPI  ArtifactType = "OPENAPI"  // OpenAPI artifact type
	AsyncAPI ArtifactType = "ASYNCAPI" // AsyncAPI artifact type
	GraphQL  ArtifactType = "GRAPHQL"  // GraphQL artifact type
	WSDL     ArtifactType = "WSDL"     // WSDL artifact type
	XSD      ArtifactType = "XSD"      // XSD artifact type
)

// ParseArtifactType parses a string and returns the corresponding ArtifactType.
func ParseArtifactType(artifactType string) (ArtifactType, error) {
	switch artifactType {
	case string(Avro):
		return Avro, nil
	case string(Protobuf):
		return Protobuf, nil
	case string(Json):
		return Json, nil
	case string(KConnect):
		return KConnect, nil
	case string(OpenAPI):
		return OpenAPI, nil
	case string(AsyncAPI):
		return AsyncAPI, nil
	case string(GraphQL):
		return GraphQL, nil
	case string(WSDL):
		return WSDL, nil
	case string(XSD):
		return XSD, nil
	default:
		return "", ErrUnknownArtifactType
	}
}

type Rule string

const (
	RuleValidity      Rule = "VALIDITY"
	RuleCompatibility Rule = "COMPATIBILITY"
	RuleIntegrity     Rule = "INTEGRITY"
)

// RuleLevel represents the level of different rules for VALIDITY, COMPATIBILITY, and INTEGRITY.
type RuleLevel string

const (
	IntegrityLevelNone          RuleLevel = "NONE"
	IntegrityLevelRefsExist     RuleLevel = "REFS_EXIST"
	IntegrityLevelAllRefsMapped RuleLevel = "ALL_REFS_MAPPED"
	IntegrityLevelNoDuplicates  RuleLevel = "NO_DUPLICATES"
	IntegrityLevelFull          RuleLevel = "FULL"

	CompatibilityLevelBackward           RuleLevel = "BACKWARD"
	CompatibilityLevelBackwardTransitive RuleLevel = "BACKWARD_TRANSITIVE"
	CompatibilityLevelForward            RuleLevel = "FORWARD"
	CompatibilityLevelForwardTransitive  RuleLevel = "FORWARD_TRANSITIVE"
	CompatibilityLevelFull               RuleLevel = "FULL"
	CompatibilityLevelFullTransitive     RuleLevel = "FULL_TRANSITIVE"
	CompatibilityLevelNone               RuleLevel = "NONE"

	ValidityLevelNone       RuleLevel = "NONE"
	ValidityLevelSyntaxOnly RuleLevel = "SYNTAX_ONLY"
	ValidityLevelFull       RuleLevel = "FULL"
)
