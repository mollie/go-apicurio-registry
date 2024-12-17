package apis

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/mollie/go-apicurio-registry/client"
	"github.com/mollie/go-apicurio-registry/models"
	"github.com/pkg/errors"
	"net/http"
)

type VersionsAPI struct {
	Client *client.Client
}

func NewVersionsAPI(client *client.Client) *VersionsAPI {
	return &VersionsAPI{
		Client: client,
	}
}

// DeleteArtifactVersion deletes a single version of the artifact.
// Parameters `groupId`, `artifactId`, and the unique `versionExpression` are needed.
// This feature must be enabled using the `registry.rest.artifact.deletion.enabled` property.
func (api *VersionsAPI) DeleteArtifactVersion(
	ctx context.Context,
	groupID, artifactID, versionExpression string,
) error {
	// Validate inputs
	if err := validateInput(groupID, regexGroupIDArtifactID, "Group ID"); err != nil {
		return err
	}
	if err := validateInput(artifactID, regexGroupIDArtifactID, "Artifact ID"); err != nil {
		return err
	}
	if err := validateInput(versionExpression, regexVersion, "Version Expression"); err != nil {
		return err
	}

	// Construct the URL
	url := fmt.Sprintf("%s/groups/%s/artifacts/%s/versions/%s", api.Client.BaseURL, groupID, artifactID, versionExpression)

	// Execute the request
	resp, err := api.executeRequest(ctx, http.MethodDelete, url, nil)
	if err != nil {
		return err
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	return handleResponse(resp, http.StatusNoContent, nil)
}

// GetArtifactVersionReferences retrieves all references for a single artifact version.
func (api *VersionsAPI) GetArtifactVersionReferences(ctx context.Context,
	groupId, artifactId, versionExpression string,
	params *models.ArtifactVersionReferencesParams,
) (*[]models.ArtifactReference, error) {
	// Validate inputs
	if err := validateInput(groupId, regexGroupIDArtifactID, "Group ID"); err != nil {
		return nil, err
	}
	if err := validateInput(artifactId, regexGroupIDArtifactID, "Artifact ID"); err != nil {
		return nil, err
	}
	if err := validateInput(versionExpression, regexVersion, "Version Expression"); err != nil {
		return nil, err
	}

	query := ""
	if params != nil {
		query = "?" + params.ToQuery().Encode()
	}

	// Start building the URL
	url := fmt.Sprintf(
		"%s/groups/%s/artifacts/%s/versions/%s/references%s",
		api.Client.BaseURL,
		groupId,
		artifactId,
		versionExpression,
		query,
	)

	// Execute the request
	resp, err := api.executeRequest(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	var references []models.ArtifactReference
	if err = handleResponse(resp, http.StatusOK, &references); err != nil {
		return nil, err
	}

	return &references, nil
}

// GetArtifactVersionComments retrieves all comments for a version of an artifact.
func (api *VersionsAPI) GetArtifactVersionComments(
	ctx context.Context,
	groupId, artifactId, versionExpression string,
) (*[]models.ArtifactComment, error) {
	// Validate inputs
	if err := validateInput(groupId, regexGroupIDArtifactID, "Group ID"); err != nil {
		return nil, err
	}
	if err := validateInput(artifactId, regexGroupIDArtifactID, "Artifact ID"); err != nil {
		return nil, err
	}
	if err := validateInput(versionExpression, regexVersion, "Version Expression"); err != nil {
		return nil, err
	}

	// Construct the URL
	url := fmt.Sprintf("%s/groups/%s/artifacts/%s/versions/%s/comments", api.Client.BaseURL, groupId, artifactId, versionExpression)

	// Execute the request
	resp, err := api.executeRequest(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	// Parse the response
	var comments []models.ArtifactComment
	if err = handleResponse(resp, http.StatusOK, &comments); err != nil {
		return nil, err
	}

	return &comments, nil
}

// AddArtifactVersionComment adds a new comment to a specific artifact version.
func (api *VersionsAPI) AddArtifactVersionComment(
	ctx context.Context,
	groupId, artifactId, versionExpression string,
	commentValue string,
) (*models.ArtifactComment, error) {
	// Validate inputs
	if err := validateInput(groupId, regexGroupIDArtifactID, "Group ID"); err != nil {
		return nil, err
	}
	if err := validateInput(artifactId, regexGroupIDArtifactID, "Artifact ID"); err != nil {
		return nil, err
	}
	if err := validateInput(versionExpression, regexVersion, "Version Expression"); err != nil {
		return nil, err
	}

	// Construct the URL
	url := fmt.Sprintf(
		"%s/groups/%s/artifacts/%s/versions/%s/comments",
		api.Client.BaseURL,
		groupId,
		artifactId,
		versionExpression,
	)

	// Create the request body
	requestBody := map[string]string{
		"value": commentValue,
	}

	// Execute the request
	resp, err := api.executeRequest(ctx, http.MethodPost, url, requestBody)
	if err != nil {
		return nil, err
	}

	// Handle the response
	var comment models.ArtifactComment
	if err := handleResponse(resp, http.StatusOK, &comment); err != nil {
		return nil, err
	}

	return &comment, nil
}

// UpdateArtifactVersionComment updates the value of a single comment in an artifact version.
func (api *VersionsAPI) UpdateArtifactVersionComment(
	ctx context.Context,
	groupId, artifactId, versionExpression, commentId string,
	updatedComment string,
) error {
	if err := validateInput(groupId, regexGroupIDArtifactID, "Group ID"); err != nil {
		return err
	}
	if err := validateInput(artifactId, regexGroupIDArtifactID, "Artifact ID"); err != nil {
		return err
	}
	if err := validateInput(versionExpression, regexVersion, "Version Expression"); err != nil {
		return err
	}
	// Build the URL
	url := fmt.Sprintf(
		"%s/groups/%s/artifacts/%s/versions/%s/comments/%s",
		api.Client.BaseURL,
		groupId,
		artifactId,
		versionExpression,
		commentId,
	)

	// Create the request body
	requestBody := map[string]string{
		"value": updatedComment,
	}

	// Execute the request
	resp, err := api.executeRequest(ctx, http.MethodPut, url, requestBody)
	if err != nil {
		return err
	}

	// Handle the response
	if err := handleResponse(resp, http.StatusNoContent, nil); err != nil {
		return err
	}

	return nil
}

// DeleteArtifactVersionComment deletes a single comment from an artifact version.
func (api *VersionsAPI) DeleteArtifactVersionComment(
	ctx context.Context,
	groupId, artifactId, versionExpression, commentId string,
) error {
	if err := validateInput(groupId, regexGroupIDArtifactID, "Group ID"); err != nil {
		return err
	}
	if err := validateInput(artifactId, regexGroupIDArtifactID, "Artifact ID"); err != nil {
		return err
	}
	if err := validateInput(versionExpression, regexVersion, "Version Expression"); err != nil {
		return err
	}

	url := fmt.Sprintf(
		"%s/groups/%s/artifacts/%s/versions/%s/comments/%s",
		api.Client.BaseURL,
		groupId,
		artifactId,
		versionExpression,
		commentId,
	)

	resp, err := api.executeRequest(ctx, http.MethodDelete, url, nil)
	if err != nil {
		return err
	}

	return handleResponse(resp, http.StatusNoContent, nil)

}

// ListArtifactVersions retrieves all versions of an artifact.
func (api *VersionsAPI) ListArtifactVersions(
	ctx context.Context,
	groupId, artifactId string,
	params *models.ListArtifactsVersionsParams,
) (*[]models.ArtifactVersion, error) {
	if err := validateInput(groupId, regexGroupIDArtifactID, "Group ID"); err != nil {
		return nil, err
	}
	if err := validateInput(artifactId, regexGroupIDArtifactID, "Artifact ID"); err != nil {
		return nil, err
	}

	query := ""
	if params != nil {
		if err := params.Validate(); err != nil {
			return nil, errors.Wrap(err, "invalid parameters provided")
		}
		query = "?" + params.ToQuery().Encode()
	}
	url := fmt.Sprintf("%s/groups/%s/artifacts/%s/versions%s", api.Client.BaseURL, groupId, artifactId, query)

	resp, err := api.executeRequest(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	var versionsResponse = models.ArtifactVersionListResponse{}
	if err = handleResponse(resp, http.StatusOK, &versionsResponse); err != nil {
		return nil, err
	}

	return &versionsResponse.Versions, nil

}

// CreateArtifactVersion creates a new version of the artifact.
func (api *VersionsAPI) CreateArtifactVersion(
	ctx context.Context,
	groupId, artifactId string,
	request *models.CreateVersionRequest,
	dryRun bool,
) (*models.ArtifactVersionDetailed, error) {
	if err := validateInput(groupId, regexGroupIDArtifactID, "Group ID"); err != nil {
		return nil, err
	}
	if err := validateInput(artifactId, regexGroupIDArtifactID, "Artifact ID"); err != nil {
		return nil, err
	}

	url := fmt.Sprintf("%s/groups/%s/artifacts/%s/versions", api.Client.BaseURL, groupId, artifactId)
	if dryRun {
		url = fmt.Sprintf("%s?dryRun=true", url)
	}

	resp, err := api.executeRequest(ctx, http.MethodPost, url, request)
	if err != nil {
		return nil, err
	}

	var version models.ArtifactVersionDetailed
	if err = handleResponse(resp, http.StatusOK, &version); err != nil {
		return nil, err
	}

	return &version, nil

}

// GetArtifactVersionContent retrieves a single version of the artifact.
func (api *VersionsAPI) GetArtifactVersionContent(
	ctx context.Context,
	groupId, artifactId, versionExpression string,
	params *models.ArtifactReferenceParams,
) (*models.ArtifactContent, error) {
	if err := validateInput(groupId, regexGroupIDArtifactID, "Group ID"); err != nil {
		return nil, err
	}
	if err := validateInput(artifactId, regexGroupIDArtifactID, "Artifact ID"); err != nil {
		return nil, err
	}
	if err := validateInput(versionExpression, regexVersion, "Version Expression"); err != nil {
		return nil, err
	}

	query := ""
	if params != nil {
		if err := params.Validate(); err != nil {
			return nil, errors.Wrap(err, "invalid parameters provided")
		}
		query = "?" + params.ToQuery().Encode()
	}
	url := fmt.Sprintf("%s/groups/%s/artifacts/%s/versions/%s/content%s", api.Client.BaseURL, groupId, artifactId, versionExpression, query)

	resp, err := api.executeRequest(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	content, err := handleRawResponse(resp, http.StatusOK)
	if err != nil {
		return nil, err
	}

	return &models.ArtifactContent{
		Content: content,
	}, nil
}

// UpdateArtifactVersionContent updates the content of a single version of the artifact.
func (api *VersionsAPI) UpdateArtifactVersionContent(
	ctx context.Context,
	groupId, artifactId, versionExpression string,
	content *models.CreateContentRequest,
) error {
	if err := validateInput(groupId, regexGroupIDArtifactID, "Group ID"); err != nil {
		return err
	}
	if err := validateInput(artifactId, regexGroupIDArtifactID, "Artifact ID"); err != nil {
		return err
	}
	if err := validateInput(versionExpression, regexVersion, "Version Expression"); err != nil {
		return err
	}

	url := fmt.Sprintf("%s/groups/%s/artifacts/%s/versions/%s/content", api.Client.BaseURL, groupId, artifactId, versionExpression)

	resp, err := api.executeRequest(ctx, http.MethodPut, url, content)
	if err != nil {
		return err
	}

	return handleResponse(resp, http.StatusNoContent, nil)
}

// SearchForArtifactVersions searches for versions of an artifact.
func (api *VersionsAPI) SearchForArtifactVersions(
	ctx context.Context,
	params *models.SearchVersionParams,
) (*[]models.ArtifactVersion, error) {

	query := ""
	if params != nil {
		if err := params.Validate(); err != nil {
			return nil, errors.Wrap(err, "invalid parameters provided")
		}
		query = params.ToQuery().Encode()
	}

	url := fmt.Sprintf("%s/search/versions?%s", api.Client.BaseURL, query)

	resp, err := api.executeRequest(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	var searchVersionsResponse = models.ArtifactVersionListResponse{}
	if err = handleResponse(resp, http.StatusOK, &searchVersionsResponse); err != nil {
		return nil, err
	}

	return &searchVersionsResponse.Versions, nil
}

// SearchForArtifactVersionByContent searches for a version of an artifact by content.
func (api *VersionsAPI) SearchForArtifactVersionByContent(
	ctx context.Context,
	content string,
	params *models.SearchVersionByContentParams,
) (*[]models.ArtifactVersion, error) {
	query := ""
	if params != nil {
		if err := params.Validate(); err != nil {
			return nil, errors.Wrap(err, "invalid parameters provided")
		}
		query = params.ToQuery().Encode()
	}

	url := fmt.Sprintf("%s/search/versions?%s", api.Client.BaseURL, query)

	resp, err := api.executeRequest(ctx, http.MethodPost, url, content)
	if err != nil {
		return nil, err
	}

	var searchVersionsResponse = models.ArtifactVersionListResponse{}
	if err = handleResponse(resp, http.StatusOK, &searchVersionsResponse); err != nil {
		return nil, err
	}

	return &searchVersionsResponse.Versions, nil
}

// GetArtifactVersionState retrieves the current state of an artifact version.
func (api *VersionsAPI) GetArtifactVersionState(
	ctx context.Context,
	groupId, artifactId, versionExpression string,
) (*models.State, error) {
	// Validate inputs
	if err := validateInput(groupId, regexGroupIDArtifactID, "Group ID"); err != nil {
		return nil, err
	}
	if err := validateInput(artifactId, regexGroupIDArtifactID, "Artifact ID"); err != nil {
		return nil, err
	}
	if err := validateInput(versionExpression, regexVersion, "Version Expression"); err != nil {
		return nil, err
	}

	// Build the URL
	url := fmt.Sprintf("%s/groups/%s/artifacts/%s/versions/%s/state", api.Client.BaseURL, groupId, artifactId, versionExpression)

	// Execute the request
	resp, err := api.executeRequest(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	// Parse response
	var stateResponse models.StateResponse
	if err = handleResponse(resp, http.StatusOK, &stateResponse); err != nil {
		return nil, err
	}

	return &stateResponse.State, nil
}

// UpdateArtifactVersionState updates the state of an artifact version.
func (api *VersionsAPI) UpdateArtifactVersionState(
	ctx context.Context,
	groupId, artifactId, versionExpression string,
	state models.State,
	dryRun bool,
) error {
	// Validate inputs
	if err := validateInput(groupId, regexGroupIDArtifactID, "Group ID"); err != nil {
		return err
	}
	if err := validateInput(artifactId, regexGroupIDArtifactID, "Artifact ID"); err != nil {
		return err
	}
	if err := validateInput(versionExpression, regexVersion, "Version Expression"); err != nil {
		return err
	}

	// Construct the URL with optional dryRun parameter
	url := fmt.Sprintf("%s/groups/%s/artifacts/%s/versions/%s/state", api.Client.BaseURL, groupId, artifactId, versionExpression)
	if dryRun {
		url += "?dryRun=true"
	}

	// Create request body
	requestBody := models.StateRequest{
		State: state,
	}

	// Execute the request
	resp, err := api.executeRequest(ctx, http.MethodPut, url, requestBody)
	if err != nil {
		return err
	}

	// Handle response
	if err = handleResponse(resp, http.StatusNoContent, nil); err != nil {
		return err
	}

	return nil
}

// executeRequest handles the creation and execution of an HTTP request.
func (api *VersionsAPI) executeRequest(ctx context.Context, method, url string, body interface{}) (*http.Response, error) {
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
