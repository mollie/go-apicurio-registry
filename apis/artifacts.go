package apis

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/mollie/go-apicurio-registry/client"
	"github.com/mollie/go-apicurio-registry/models"
	"github.com/pkg/errors"
	"io"
	"net/http"
)

type ArtifactsAPI struct {
	Client *client.Client
}

func NewArtifactsAPI(client *client.Client) *ArtifactsAPI {
	return &ArtifactsAPI{
		Client: client,
	}
}

var (
	ErrArtifactNotFound = errors.New("artifact not found")
	ErrMethodNotAllowed = errors.New("method not allowed or disabled on the server")
	ErrInvalidInput     = errors.New("input must be between 1 and 512 characters")
)

// SearchArtifacts - Search for artifacts using the given filter parameters.
// Search for artifacts using the given filter parameters.
// See:
func (api *ArtifactsAPI) SearchArtifacts(ctx context.Context, params *models.SearchArtifactsParams) (*[]models.SearchedArtifact, error) {
	query := ""
	if params != nil {
		query = "?" + params.ToQuery().Encode()
	}

	url := fmt.Sprintf("%s/search/artifacts%s", api.Client.BaseURL, query)
	resp, err := api.executeRequest(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	var result models.SearchArtifactsAPIResponse
	if err := handleResponse(resp, http.StatusOK, &result); err != nil {
		return nil, err
	}

	return &result.Artifacts, nil
}

// SearchArtifactsByContent searches for artifacts that match the provided content.
// Returns a paginated list of all artifacts with at least one version that matches the posted content.
// See: https://www.apicur.io/registry/docs/apicurio-registry/3.0.x/assets-attachments/registry-rest-api.htm#tag/Artifacts/operation/searchArtifactsByContent
func (api *ArtifactsAPI) SearchArtifactsByContent(ctx context.Context, content []byte, params *models.SearchArtifactsByContentParams) (*[]models.SearchedArtifact, error) {
	// Convert params to query string
	query := ""
	if params != nil {
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

	return &result.Artifacts, nil
}

// ListArtifactReferences Returns a list containing all the artifact references using the artifact content ID.
// See: https://www.apicur.io/registry/docs/apicurio-registry/3.0.x/assets-attachments/registry-rest-api.htm#tag/Artifacts/operation/referencesByContentId
func (api *ArtifactsAPI) ListArtifactReferences(ctx context.Context, contentID int64) (*[]models.ArtifactReference, error) {
	url := fmt.Sprintf("%s/ids/contentId/%d/references", api.Client.BaseURL, contentID)
	resp, err := api.executeRequest(ctx, http.MethodGet, url, nil)
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
// See: https://www.apicur.io/registry/docs/apicurio-registry/3.0.x/assets-attachments/registry-rest-api.htm#tag/Artifacts/operation/referencesByContentHash
func (api *ArtifactsAPI) ListArtifactReferencesByGlobalID(ctx context.Context, globalID int64, params *models.ListArtifactReferencesByGlobalIDParams) (*[]models.ArtifactReference, error) {
	query := ""
	if params != nil {
		query = "?" + params.ToQuery().Encode()
	}

	url := fmt.Sprintf("%s/ids/globalIds/%d/references%s", api.Client.BaseURL, globalID, query)
	resp, err := api.executeRequest(ctx, http.MethodGet, url, nil)
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
// See: https://www.apicur.io/registry/docs/apicurio-registry/3.0.x/assets-attachments/registry-rest-api.htm#tag/Artifacts/operation/referencesByContentHash
func (api *ArtifactsAPI) ListArtifactReferencesByHash(ctx context.Context, contentHash string) (*[]models.ArtifactReference, error) {
	url := fmt.Sprintf("%s/ids/contentHashes/%s/references", api.Client.BaseURL, contentHash)
	resp, err := api.executeRequest(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	var references []models.ArtifactReference
	if err := handleResponse(resp, http.StatusOK, &references); err != nil {
		return nil, err
	}

	return &references, nil
}

// ListArtifactsInGroup lists all artifacts in a specified group.
// See: https://www.apicur.io/registry/docs/apicurio-registry/3.0.x/assets-attachments/registry-rest-api.htm#tag/Artifacts/operation/referencesByContentHash
func (api *ArtifactsAPI) ListArtifactsInGroup(ctx context.Context, groupID string, params *models.ListArtifactsInGroupParams) (*models.ListArtifactsResponse, error) {
	if err := validateInput(groupID, regexGroupIDArtifactID, "Group ID"); err != nil {
		return nil, err
	}

	query := ""
	if params != nil {
		query = "?" + params.ToQuery().Encode()
	}

	url := fmt.Sprintf("%s/groups/%s/artifacts%s", api.Client.BaseURL, groupID, query)
	resp, err := api.executeRequest(ctx, http.MethodGet, url, nil)
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
// See: https://www.apicur.io/registry/docs/apicurio-registry/3.0.x/assets-attachments/registry-rest-api.htm#tag/Artifacts/operation/getContentByHash
func (api *ArtifactsAPI) GetArtifactContentByHash(ctx context.Context, contentHash string) (*models.ArtifactContent, error) {
	url := fmt.Sprintf("%s/ids/contentHashes/%s", api.Client.BaseURL, contentHash)
	resp, err := api.executeRequest(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode == http.StatusNotFound {
		return nil, errors.Wrapf(ErrArtifactNotFound, "content hash: %s", contentHash)
	}

	if resp.StatusCode != http.StatusOK {
		apiError, parseErr := parseAPIError(resp)
		if parseErr != nil {
			return nil, errors.Wrap(parseErr, "unexpected error")
		}
		return nil, apiError
	}

	// Parse artifact type header
	artifactType, err := parseArtifactTypeHeader(resp)
	if err != nil {
		return nil, err
	}

	// Parse the response body
	content, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.Wrap(err, "failed to read response body")
	}

	return &models.ArtifactContent{
		Content:      string(content),
		ArtifactType: artifactType,
	}, nil
}

// GetArtifactContentByID Gets the content for an artifact version in the registry using the unique content identifier for that content
// This content ID may be shared by multiple artifact versions in the case where the artifact versions are identical.
// See: https://www.apicur.io/registry/docs/apicurio-registry/3.0.x/assets-attachments/registry-rest-api.htm#tag/Artifacts/operation/getContentById
func (api *ArtifactsAPI) GetArtifactContentByID(ctx context.Context, contentID int64) (*models.ArtifactContent, error) {
	url := fmt.Sprintf("%s/ids/contentIds/%d", api.Client.BaseURL, contentID)
	resp, err := api.executeRequest(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode == http.StatusNotFound {
		return nil, errors.Wrapf(ErrArtifactNotFound, "content ID: %d", contentID)
	}

	if resp.StatusCode != http.StatusOK {
		apiError, parseErr := parseAPIError(resp)
		if parseErr != nil {
			return nil, errors.Wrap(parseErr, "unexpected error")
		}
		return nil, apiError
	}

	// Parse artifact type header
	artifactType, err := parseArtifactTypeHeader(resp)
	if err != nil {
		return nil, err
	}

	// Parse the response body
	content, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.Wrap(err, "failed to read response body")
	}

	return &models.ArtifactContent{
		Content:      string(content),
		ArtifactType: artifactType,
	}, nil
}

// DeleteArtifactsInGroup deletes all artifacts in a given group.
// Deletes all the artifacts that exist in a given group.
// See: https://www.apicur.io/registry/docs/apicurio-registry/3.0.x/assets-attachments/registry-rest-api.htm#tag/Artifacts/operation/deleteArtifactsInGroup
func (api *ArtifactsAPI) DeleteArtifactsInGroup(ctx context.Context, groupID string) error {
	if err := validateInput(groupID, regexGroupIDArtifactID, "Group ID"); err != nil {
		return err
	}

	url := fmt.Sprintf("%s/groups/%s/artifacts", api.Client.BaseURL, groupID)
	resp, err := api.executeRequest(ctx, http.MethodDelete, url, nil)
	if err != nil {
		return err
	}

	return handleResponse(resp, http.StatusNoContent, nil)
}

// DeleteArtifact deletes a specific artifact identified by groupId and artifactId.
// Deletes an artifact completely, resulting in all versions of the artifact also being deleted. This may fail for one of the following reasons:
// See: https://www.apicur.io/registry/docs/apicurio-registry/3.0.x/assets-attachments/registry-rest-api.htm#tag/Artifacts/operation/deleteArtifact
func (api *ArtifactsAPI) DeleteArtifact(ctx context.Context, groupID, artifactId string) error {
	if err := validateInput(groupID, regexGroupIDArtifactID, "Group ID"); err != nil {
		return err
	}
	if err := validateInput(artifactId, regexGroupIDArtifactID, "Artifact ID"); err != nil {
		return err
	}

	url := fmt.Sprintf("%s/groups/%s/artifacts/%s", api.Client.BaseURL, groupID, artifactId)
	resp, err := api.executeRequest(ctx, http.MethodDelete, url, nil)
	if err != nil {
		return err
	}

	if resp.StatusCode == http.StatusMethodNotAllowed {
		return ErrMethodNotAllowed
	}

	return handleResponse(resp, http.StatusNoContent, nil)
}

// CreateArtifact Creates a new artifact.
// See: https://www.apicur.io/registry/docs/apicurio-registry/3.0.x/assets-attachments/registry-rest-api.htm#tag/Artifacts/operation/createArtifact
func (api *ArtifactsAPI) CreateArtifact(ctx context.Context, groupId string, artifact models.CreateArtifactRequest, params *models.CreateArtifactParams) (*models.ArtifactDetail, error) {
	if err := validateInput(groupId, regexGroupIDArtifactID, "Group ID"); err != nil {
		return nil, err
	}

	query := ""
	if params != nil {
		query = "?" + params.ToQuery().Encode()
	}
	url := fmt.Sprintf("%s/groups/%s/artifacts%s", api.Client.BaseURL, groupId, query)

	resp, err := api.executeRequest(ctx, http.MethodPost, url, artifact)
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
// See: https://www.apicur.io/registry/docs/apicurio-registry/3.0.x/assets-attachments/registry-rest-api.htm#tag/Artifact-rules/operation/createArtifactRule
func (api *ArtifactsAPI) ListArtifactRules(ctx context.Context, groupID, artifactId string) ([]models.Rule, error) {
	url := fmt.Sprintf("%s/groups/%s/artifacts/%s/rules", api.Client.BaseURL, groupID, artifactId)
	resp, err := api.executeRequest(ctx, http.MethodGet, url, nil)
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
// See: https://www.apicur.io/registry/docs/apicurio-registry/3.0.x/assets-attachments/registry-rest-api.htm#tag/Artifact-rules/operation/createArtifactRule
func (api *ArtifactsAPI) CreateArtifactRule(ctx context.Context, groupID, artifactId string, rule models.Rule, level models.RuleLevel) error {
	url := fmt.Sprintf("%s/groups/%s/artifacts/%s/rules", api.Client.BaseURL, groupID, artifactId)

	// Prepare the request body
	body := models.CreateUpdateGlobalRuleRequest{
		RuleType: rule,
		Config:   level,
	}
	resp, err := api.executeRequest(ctx, http.MethodPost, url, body)
	if err != nil {
		return err
	}

	return handleResponse(resp, http.StatusNoContent, nil)
}

// DeleteAllArtifactRule deletes all artifact rules for a given artifact.
// See: https://www.apicur.io/registry/docs/apicurio-registry/3.0.x/assets-attachments/registry-rest-api.htm#tag/Artifact-rules/operation/deleteArtifactRules
func (api *ArtifactsAPI) DeleteAllArtifactRule(ctx context.Context, groupID, artifactId string) error {
	url := fmt.Sprintf("%s/groups/%s/artifacts/%s/rules", api.Client.BaseURL, groupID, artifactId)
	resp, err := api.executeRequest(ctx, http.MethodDelete, url, nil)
	if err != nil {
		return err
	}

	return handleResponse(resp, http.StatusNoContent, nil)
}

// GetArtifactRule gets the rule level for a given artifact rule.
// See: https://www.apicur.io/registry/docs/apicurio-registry/3.0.x/assets-attachments/registry-rest-api.htm#tag/Artifact-rules/operation/getArtifactRuleConfig
func (api *ArtifactsAPI) GetArtifactRule(ctx context.Context, groupID, artifactId string, rule models.Rule) (models.RuleLevel, error) {
	url := fmt.Sprintf("%s/groups/%s/artifacts/%s/rules/%s", api.Client.BaseURL, groupID, artifactId, rule)
	resp, err := api.executeRequest(ctx, http.MethodGet, url, nil)
	if err != nil {
		return "", err
	}

	var globalRule models.GlobalRuleResponse
	if err := handleResponse(resp, http.StatusOK, &globalRule); err != nil {
		return "", err
	}

	return globalRule.Config, nil
}

// UpdateArtifactRule updates the rule level for a given artifact rule.
// See: https://www.apicur.io/registry/docs/apicurio-registry/3.0.x/assets-attachments/registry-rest-api.htm#tag/Artifact-rules/operation/updateArtifactRuleConfig
func (api *ArtifactsAPI) UpdateArtifactRule(ctx context.Context, groupID, artifactId string, rule models.Rule, level models.RuleLevel) error {
	url := fmt.Sprintf("%s/groups/%s/artifacts/%s/rules/%s", api.Client.BaseURL, groupID, artifactId, rule)

	// Prepare the request body
	body := models.CreateUpdateGlobalRuleRequest{
		RuleType: rule,
		Config:   level,
	}
	resp, err := api.executeRequest(ctx, http.MethodPut, url, body)
	if err != nil {
		return err
	}

	var globalRule models.GlobalRuleResponse
	if err := handleResponse(resp, http.StatusOK, &globalRule); err != nil {
		return err
	}

	return nil
}

// DeleteArtifactRule deletes a specific artifact rule for a given artifact.
// See: https://www.apicur.io/registry/docs/apicurio-registry/3.0.x/assets-attachments/registry-rest-api.htm#tag/Artifact-rules/operation/deleteArtifactRule
func (api *ArtifactsAPI) DeleteArtifactRule(ctx context.Context, groupID, artifactId string, rule models.Rule) error {
	url := fmt.Sprintf("%s/groups/%s/artifacts/%s/rules/%s", api.Client.BaseURL, groupID, artifactId, rule)
	resp, err := api.executeRequest(ctx, http.MethodDelete, url, nil)
	if err != nil {
		return err
	}

	return handleResponse(resp, http.StatusNoContent, nil)
}

// executeRequest handles the creation and execution of an HTTP request.
func (api *ArtifactsAPI) executeRequest(ctx context.Context, method, url string, body interface{}) (*http.Response, error) {
	var reqBody []byte
	var err error
	contentType := "*/*"

	switch v := body.(type) {
	case string:
		reqBody = []byte(v)
		contentType = "*/*"
	case []byte:
		reqBody = v
		contentType = "*/*"
	default:
		contentType = "application/json"
		reqBody, err = json.Marshal(body)
		if err != nil {
			return nil, errors.Wrap(err, "failed to marshal request body as JSON")
		}
	}

	// Create the HTTP request
	req, err := http.NewRequestWithContext(ctx, method, url, bytes.NewReader(reqBody))
	if err != nil {
		return nil, errors.Wrap(err, "failed to create HTTP request")
	}

	// Set appropriate Content-Type header
	if body != nil {
		req.Header.Set("Content-Type", contentType)
	}

	// Execute the request
	resp, err := api.Client.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "failed to execute HTTP request")
	}

	return resp, nil
}
