package apis

import (
	"context"
	"fmt"
	"net/http"
	"net/url"

	"github.com/mollie/go-apicurio-registry/client"
	"github.com/mollie/go-apicurio-registry/models"
	"github.com/pkg/errors"
)

type ArtifactsAPI struct {
	Client *client.Client
}

func NewArtifactsAPI(client *client.Client) *ArtifactsAPI {
	return &ArtifactsAPI{
		Client: client,
	}
}

// GetArtifactByGlobalID Gets the content for an artifact version in the registry using its globally unique identifier.
// See https://www.apicur.io/registry/docs/apicurio-registry/3.0.x/assets-attachments/registry-rest-api.htm#tag/Artifacts/operation/getContentByGlobalId
func (api *ArtifactsAPI) GetArtifactByGlobalID(
	ctx context.Context,
	globalID int64,
	params *models.GetArtifactByGlobalIDParams,
) (*models.ArtifactContent, error) {
	returnArtifactType := false
	query := ""
	if params != nil {
		if err := params.Validate(); err != nil {
			return nil, errors.Wrap(err, "invalid parameters provided")
		}
		returnArtifactType = params.ReturnArtifactType
		query = "?" + params.ToQuery().Encode()
	}

	urlPath := fmt.Sprintf("%s/ids/globalIds/%d%s", api.Client.BaseURL, globalID, query)
	resp, err := api.executeRequest(ctx, http.MethodGet, urlPath, nil)
	if err != nil {
		return nil, err
	}

	content, err := handleRawResponse(resp, http.StatusOK)
	if err != nil {
		return nil, err
	}

	var artifactType models.ArtifactType
	if returnArtifactType {
		// Parse artifact type header
		aType, err := parseArtifactTypeHeader(resp)
		if err != nil {
			return nil, err
		}
		artifactType = aType
	}

	return &models.ArtifactContent{
		Content:      content,
		ArtifactType: artifactType,
	}, nil
}

// SearchArtifacts - Search for artifacts using the given filter parameters.
// Search for artifacts using the given filter parameters.
// See https://www.apicur.io/registry/docs/apicurio-registry/3.0.x/assets-attachments/registry-rest-api.htm#tag/Artifacts/operation/searchArtifacts
func (api *ArtifactsAPI) SearchArtifacts(
	ctx context.Context,
	params *models.SearchArtifactsParams,
) ([]models.SearchedArtifact, error) {
	query := ""
	if params != nil {
		if err := params.Validate(); err != nil {
			return nil, errors.Wrap(err, "invalid parameters provided")
		}
		query = "?" + params.ToQuery().Encode()
	}

	urlPath := fmt.Sprintf("%s/search/artifacts%s", api.Client.BaseURL, query)
	resp, err := api.executeRequest(ctx, http.MethodGet, urlPath, nil)
	if err != nil {
		return nil, err
	}

	var result models.SearchArtifactsAPIResponse
	if err := handleResponse(resp, http.StatusOK, &result); err != nil {
		return nil, err
	}

	return result.Artifacts, nil
}

// SearchArtifactsByContent searches for artifacts that match the provided content.
// Returns a paginated list of all artifacts with at least one version that matches the posted content.
// See https://www.apicur.io/registry/docs/apicurio-registry/3.0.x/assets-attachments/registry-rest-api.htm#tag/Artifacts/operation/searchArtifactsByContent
func (api *ArtifactsAPI) SearchArtifactsByContent(
	ctx context.Context,
	content []byte,
	params *models.SearchArtifactsByContentParams,
) ([]models.SearchedArtifact, error) {
	// Convert params to query string
	query := ""
	if params != nil {
		if err := params.Validate(); err != nil {
			return nil, errors.Wrap(err, "invalid parameters provided")
		}
		query = "?" + params.ToQuery().Encode()
	}

	url := fmt.Sprintf("%s/search/artifacts%s", api.Client.BaseURL, query)
	resp, err := api.executeRequest(ctx, http.MethodPost, url, content)
	if err != nil {
		return nil, err
	}

	var result models.SearchArtifactsAPIResponse
	if err := handleResponse(resp, http.StatusOK, &result); err != nil {
		return nil, err
	}

	return result.Artifacts, nil
}

// ListArtifactReferences Returns a list containing all the artifact references using the artifact content ID.
// See https://www.apicur.io/registry/docs/apicurio-registry/3.0.x/assets-attachments/registry-rest-api.htm#tag/Artifacts/operation/referencesByContentId
func (api *ArtifactsAPI) ListArtifactReferences(
	ctx context.Context,
	contentID int64,
) (*[]models.ArtifactReference, error) {
	urlPath := fmt.Sprintf("%s/ids/contentId/%d/references", api.Client.BaseURL, contentID)
	resp, err := api.executeRequest(ctx, http.MethodGet, urlPath, nil)
	if err != nil {
		return nil, err
	}

	var references []models.ArtifactReference
	if err := handleResponse(resp, http.StatusOK, &references); err != nil {
		return nil, err
	}

	return &references, nil
}

// ListArtifactReferencesByGlobalID Returns a list containing all the artifact references using the artifact global ID.
// See https://www.apicur.io/registry/docs/apicurio-registry/3.0.x/assets-attachments/registry-rest-api.htm#tag/Artifacts/operation/referencesByContentHash
func (api *ArtifactsAPI) ListArtifactReferencesByGlobalID(
	ctx context.Context,
	globalID int64,
	params *models.ListArtifactReferencesByGlobalIDParams,
) (*[]models.ArtifactReference, error) {
	query := ""
	if params != nil {
		if err := params.Validate(); err != nil {
			return nil, errors.Wrap(err, "invalid parameters provided")
		}
		query = "?" + params.ToQuery().Encode()
	}

	urlPath := fmt.Sprintf("%s/ids/globalIds/%d/references%s", api.Client.BaseURL, globalID, query)
	resp, err := api.executeRequest(ctx, http.MethodGet, urlPath, nil)
	if err != nil {
		return nil, err
	}

	var references []models.ArtifactReference
	if err := handleResponse(resp, http.StatusOK, &references); err != nil {
		return nil, err
	}

	return &references, nil
}

// ListArtifactReferencesByHash Returns a list containing all the artifact references using the artifact content hash.
// See https://www.apicur.io/registry/docs/apicurio-registry/3.0.x/assets-attachments/registry-rest-api.htm#tag/Artifacts/operation/referencesByContentHash
func (api *ArtifactsAPI) ListArtifactReferencesByHash(
	ctx context.Context,
	contentHash string,
) ([]models.ArtifactReference, error) {
	urlPath := fmt.Sprintf(
		"%s/ids/contentHashes/%s/references",
		api.Client.BaseURL,
		url.PathEscape(contentHash),
	)
	resp, err := api.executeRequest(ctx, http.MethodGet, urlPath, nil)
	if err != nil {
		return nil, err
	}

	var references []models.ArtifactReference
	if err := handleResponse(resp, http.StatusOK, &references); err != nil {
		return nil, err
	}

	return references, nil
}

// ListArtifactsInGroup lists all artifacts in a specified group.
// See https://www.apicur.io/registry/docs/apicurio-registry/3.0.x/assets-attachments/registry-rest-api.htm#tag/Artifacts/operation/referencesByContentHash
func (api *ArtifactsAPI) ListArtifactsInGroup(
	ctx context.Context,
	groupID string,
	params *models.ListArtifactsInGroupParams,
) (*models.ListArtifactsResponse, error) {
	if err := validateInput(groupID, regexGroupIDArtifactID, "Group ID"); err != nil {
		return nil, err
	}

	query := ""
	if params != nil {
		if err := params.Validate(); err != nil {
			return nil, errors.Wrap(err, "invalid parameters provided")
		}
		query = "?" + params.ToQuery().Encode()
	}

	urlPath := fmt.Sprintf(
		"%s/groups/%s/artifacts%s",
		api.Client.BaseURL,
		url.PathEscape(groupID),
		query,
	)
	resp, err := api.executeRequest(ctx, http.MethodGet, urlPath, nil)
	if err != nil {
		return nil, err
	}

	var result models.ListArtifactsResponse
	if err := handleResponse(resp, http.StatusOK, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// GetArtifactContentByHash Gets the content for an artifact version in the registry using the SHA-256 hash of the content
// This content hash may be shared by multiple artifact versions in the case where the artifact versions have identical content.
// See https://www.apicur.io/registry/docs/apicurio-registry/3.0.x/assets-attachments/registry-rest-api.htm#tag/Artifacts/operation/getContentByHash
func (api *ArtifactsAPI) GetArtifactContentByHash(
	ctx context.Context,
	contentHash string,
) (*models.ArtifactContent, error) {
	urlPath := fmt.Sprintf(
		"%s/ids/contentHashes/%s",
		api.Client.BaseURL,
		url.PathEscape(contentHash),
	)
	resp, err := api.executeRequest(ctx, http.MethodGet, urlPath, nil)
	if err != nil {
		return nil, err
	}

	content, err := handleRawResponse(resp, http.StatusOK)
	if err != nil {
		return nil, err
	}

	// Parse artifact type header
	artifactType, err := parseArtifactTypeHeader(resp)
	if err != nil {
		return nil, err
	}

	return &models.ArtifactContent{
		Content:      content,
		ArtifactType: artifactType,
	}, nil
}

// GetArtifactContentByID Gets the content for an artifact version in the registry using the unique content identifier for that content
// This content ID may be shared by multiple artifact versions in the case where the artifact versions are identical.
// See https://www.apicur.io/registry/docs/apicurio-registry/3.0.x/assets-attachments/registry-rest-api.htm#tag/Artifacts/operation/getContentById
func (api *ArtifactsAPI) GetArtifactContentByID(
	ctx context.Context,
	contentID int64,
) (*models.ArtifactContent, error) {
	urlPath := fmt.Sprintf("%s/ids/contentIds/%d", api.Client.BaseURL, contentID)
	resp, err := api.executeRequest(ctx, http.MethodGet, urlPath, nil)
	if err != nil {
		return nil, err
	}

	content, err := handleRawResponse(resp, http.StatusOK)
	if err != nil {
		return nil, err
	}

	// Parse artifact type header
	artifactType, err := parseArtifactTypeHeader(resp)
	if err != nil {
		return nil, err
	}

	return &models.ArtifactContent{
		Content:      content,
		ArtifactType: artifactType,
	}, nil
}

// DeleteArtifactsInGroup deletes all artifacts in a given group.
// Deletes all the artifacts that exist in a given group.
// See https://www.apicur.io/registry/docs/apicurio-registry/3.0.x/assets-attachments/registry-rest-api.htm#tag/Artifacts/operation/deleteArtifactsInGroup
func (api *ArtifactsAPI) DeleteArtifactsInGroup(ctx context.Context, groupID string) error {
	if err := validateInput(groupID, regexGroupIDArtifactID, "Group ID"); err != nil {
		return err
	}

	urlPath := fmt.Sprintf("%s/groups/%s/artifacts", api.Client.BaseURL, url.PathEscape(groupID))
	resp, err := api.executeRequest(ctx, http.MethodDelete, urlPath, nil)
	if err != nil {
		return err
	}

	return handleResponse(resp, http.StatusNoContent, nil)
}

// DeleteArtifact deletes a specific artifact identified by groupId and artifactId.
// Deletes an artifact completely, resulting in all versions of the artifact also being deleted. This may fail for one of the following reasons:
// See https://www.apicur.io/registry/docs/apicurio-registry/3.0.x/assets-attachments/registry-rest-api.htm#tag/Artifacts/operation/deleteArtifact
func (api *ArtifactsAPI) DeleteArtifact(ctx context.Context, groupID, artifactId string) error {
	if err := validateInput(groupID, regexGroupIDArtifactID, "Group ID"); err != nil {
		return err
	}
	if err := validateInput(artifactId, regexGroupIDArtifactID, "Artifact ID"); err != nil {
		return err
	}

	urlPath := fmt.Sprintf(
		"%s/groups/%s/artifacts/%s",
		api.Client.BaseURL,
		url.PathEscape(groupID),
		url.PathEscape(artifactId),
	)
	resp, err := api.executeRequest(ctx, http.MethodDelete, urlPath, nil)
	if err != nil {
		return err
	}

	return handleResponse(resp, http.StatusNoContent, nil)
}

// CreateArtifact Creates a new artifact.
// See https://www.apicur.io/registry/docs/apicurio-registry/3.0.x/assets-attachments/registry-rest-api.htm#tag/Artifacts/operation/createArtifact
func (api *ArtifactsAPI) CreateArtifact(
	ctx context.Context,
	groupId string,
	artifact models.CreateArtifactRequest,
	params *models.CreateArtifactParams,
) (*models.ArtifactDetail, error) {
	if err := validateInput(groupId, regexGroupIDArtifactID, "Group ID"); err != nil {
		return nil, err
	}

	if err := artifact.Validate(); err != nil {
		return nil, errors.Wrap(err, "invalid artifact provided")
	}

	query := ""
	if params != nil {
		if err := params.Validate(); err != nil {
			return nil, errors.Wrap(err, "invalid parameters provided")
		}
		query = "?" + params.ToQuery().Encode()
	}
	urlPath := fmt.Sprintf(
		"%s/groups/%s/artifacts%s",
		api.Client.BaseURL,
		url.PathEscape(groupId),
		query,
	)

	resp, err := api.executeRequest(ctx, http.MethodPost, urlPath, artifact)
	if err != nil {
		return nil, err
	}

	var response models.CreateArtifactResponse
	if err := handleResponse(resp, http.StatusOK, &response); err != nil {
		return nil, err
	}

	return &response.Artifact, nil
}

// ListArtifactRules lists all artifact rules for a given artifact.
// See https://www.apicur.io/registry/docs/apicurio-registry/3.0.x/assets-attachments/registry-rest-api.htm#tag/Artifact-rules/operation/createArtifactRule
func (api *ArtifactsAPI) ListArtifactRules(
	ctx context.Context,
	groupID, artifactId string,
) ([]models.Rule, error) {
	urlPath := fmt.Sprintf(
		"%s/groups/%s/artifacts/%s/rules",
		api.Client.BaseURL,
		url.PathEscape(groupID),
		url.PathEscape(artifactId),
	)
	resp, err := api.executeRequest(ctx, http.MethodGet, urlPath, nil)
	if err != nil {
		return nil, err
	}

	var rules []models.Rule
	if err := handleResponse(resp, http.StatusOK, &rules); err != nil {
		return nil, err
	}

	return rules, nil
}

// CreateArtifactRule creates a new artifact rule for a given artifact.
// See https://www.apicur.io/registry/docs/apicurio-registry/3.0.x/assets-attachments/registry-rest-api.htm#tag/Artifact-rules/operation/createArtifactRule
func (api *ArtifactsAPI) CreateArtifactRule(
	ctx context.Context,
	groupID, artifactId string,
	rule models.Rule,
	level models.RuleLevel,
) error {
	urlPath := fmt.Sprintf(
		"%s/groups/%s/artifacts/%s/rules",
		api.Client.BaseURL,
		url.PathEscape(groupID),
		url.PathEscape(artifactId),
	)

	// Prepare the request body
	body := models.CreateUpdateRuleRequest{
		RuleType: rule,
		Config:   level,
	}
	resp, err := api.executeRequest(ctx, http.MethodPost, urlPath, body)
	if err != nil {
		return err
	}

	return handleResponse(resp, http.StatusNoContent, nil)
}

// DeleteAllArtifactRule deletes all artifact rules for a given artifact.
// See https://www.apicur.io/registry/docs/apicurio-registry/3.0.x/assets-attachments/registry-rest-api.htm#tag/Artifact-rules/operation/deleteArtifactRules
func (api *ArtifactsAPI) DeleteAllArtifactRule(
	ctx context.Context,
	groupID, artifactId string,
) error {
	urlPath := fmt.Sprintf(
		"%s/groups/%s/artifacts/%s/rules",
		api.Client.BaseURL,
		url.PathEscape(groupID),
		url.PathEscape(artifactId),
	)
	resp, err := api.executeRequest(ctx, http.MethodDelete, urlPath, nil)
	if err != nil {
		return err
	}

	return handleResponse(resp, http.StatusNoContent, nil)
}

// GetArtifactRule gets the rule level for a given artifact rule.
// See https://www.apicur.io/registry/docs/apicurio-registry/3.0.x/assets-attachments/registry-rest-api.htm#tag/Artifact-rules/operation/getArtifactRuleConfig
func (api *ArtifactsAPI) GetArtifactRule(
	ctx context.Context,
	groupID, artifactId string,
	rule models.Rule,
) (models.RuleLevel, error) {
	urlPath := fmt.Sprintf(
		"%s/groups/%s/artifacts/%s/rules/%s",
		api.Client.BaseURL,
		url.PathEscape(groupID),
		url.PathEscape(artifactId),
		string(rule),
	)
	resp, err := api.executeRequest(ctx, http.MethodGet, urlPath, nil)
	if err != nil {
		return "", err
	}

	var globalRule models.RuleResponse
	if err := handleResponse(resp, http.StatusOK, &globalRule); err != nil {
		return "", err
	}

	return globalRule.Config, nil
}

// UpdateArtifactRule updates the rule level for a given artifact rule.
// See https://www.apicur.io/registry/docs/apicurio-registry/3.0.x/assets-attachments/registry-rest-api.htm#tag/Artifact-rules/operation/updateArtifactRuleConfig
func (api *ArtifactsAPI) UpdateArtifactRule(
	ctx context.Context,
	groupID, artifactId string,
	rule models.Rule,
	level models.RuleLevel,
) error {
	urlPath := fmt.Sprintf(
		"%s/groups/%s/artifacts/%s/rules/%s",
		api.Client.BaseURL,
		url.PathEscape(groupID),
		url.PathEscape(artifactId),
		string(rule),
	)

	// Prepare the request body
	body := models.CreateUpdateRuleRequest{
		RuleType: rule,
		Config:   level,
	}
	resp, err := api.executeRequest(ctx, http.MethodPut, urlPath, body)
	if err != nil {
		return err
	}

	var globalRule models.RuleResponse
	if err := handleResponse(resp, http.StatusOK, &globalRule); err != nil {
		return err
	}

	return nil
}

// DeleteArtifactRule deletes a specific artifact rule for a given artifact.
// See https://www.apicur.io/registry/docs/apicurio-registry/3.0.x/assets-attachments/registry-rest-api.htm#tag/Artifact-rules/operation/deleteArtifactRule
func (api *ArtifactsAPI) DeleteArtifactRule(
	ctx context.Context,
	groupID, artifactId string,
	rule models.Rule,
) error {
	urlPath := fmt.Sprintf(
		"%s/groups/%s/artifacts/%s/rules/%s",
		api.Client.BaseURL,
		url.PathEscape(groupID),
		url.PathEscape(artifactId),
		string(rule),
	)
	resp, err := api.executeRequest(ctx, http.MethodDelete, urlPath, nil)
	if err != nil {
		return err
	}

	return handleResponse(resp, http.StatusNoContent, nil)
}

// executeRequest handles the creation and execution of an HTTP request.
func (api *ArtifactsAPI) executeRequest(
	ctx context.Context,
	method, url string,
	body interface{},
) (*http.Response, error) {
	return executeRequest(ctx, api.Client, method, url, body)
}
