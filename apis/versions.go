package apis

import (
	"context"
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
// See https://www.apicur.io/registry/docs/apicurio-registry/3.0.x/assets-attachments/registry-rest-api.htm#tag/Versions/operation/deleteArtifactVersion
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

// GetArtifactVersionReferences Retrieves all references for a single version of an artifact.
// Both the artifactId and the unique version number must be provided.
// Using the refType query parameter, it is possible to retrieve an array of either the inbound or outbound references.
// See https://www.apicur.io/registry/docs/apicurio-registry/3.0.x/assets-attachments/registry-rest-api.htm#tag/Versions/operation/getArtifactVersionReferences
func (api *VersionsAPI) GetArtifactVersionReferences(ctx context.Context,
	groupId, artifactId, versionExpression string,
	params *models.ArtifactVersionReferencesParams,
) ([]models.ArtifactReference, error) {
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
		if err := params.Validate(); err != nil {
			return nil, errors.Wrap(err, "invalid parameters provided")
		}
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

	return references, nil
}

// GetArtifactVersionComments Retrieves all comments for a version of an artifact.
// Both the artifactId and the unique version number must be provided.
// See https://www.apicur.io/registry/docs/apicurio-registry/3.0.x/assets-attachments/registry-rest-api.htm#tag/Versions/operation/getArtifactVersionComments
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

// AddArtifactVersionComment Adds a new comment to the artifact version.
// Both the artifactId and the unique version number must be provided.
// See https://www.apicur.io/registry/docs/apicurio-registry/3.0.x/assets-attachments/registry-rest-api.htm#tag/Versions/operation/addArtifactVersionComment
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

// UpdateArtifactVersionComment Updates the value of a single comment in an artifact version.
// Only the owner of the comment can modify it.
// The artifactId, unique version number, and commentId must be provided.
// See https://www.apicur.io/registry/docs/apicurio-registry/3.0.x/assets-attachments/registry-rest-api.htm#tag/Versions/operation/updateArtifactVersionComment
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

// DeleteArtifactVersionComment Deletes a single comment in an artifact version.
// Only the owner of the comment can delete it.
// The artifactId, unique version number, and commentId must be provided.
// See https://www.apicur.io/registry/docs/apicurio-registry/3.0.x/assets-attachments/registry-rest-api.htm#tag/Versions/operation/deleteArtifactVersionComment
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

	if commentId == "" {
		return errors.New("Comment ID cannot be empty")
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

// ListArtifactVersions Returns a list of all versions of the artifact.
// The result set is paged.
// See https://www.apicur.io/registry/docs/apicurio-registry/3.0.x/assets-attachments/registry-rest-api.htm#tag/Versions/operation/listArtifactVersions
func (api *VersionsAPI) ListArtifactVersions(
	ctx context.Context,
	groupId, artifactId string,
	params *models.ListArtifactsVersionsParams,
) ([]models.ArtifactVersion, error) {
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

	return versionsResponse.Versions, nil

}

// CreateArtifactVersion Creates a new version of the artifact by uploading new content.
// The configured rules for the artifact are applied, and if they all pass, the new content is added as the most recent version of the artifact.
// If any of the rules fail, an error is returned.
// See https://www.apicur.io/registry/docs/apicurio-registry/3.0.x/assets-attachments/registry-rest-api.htm#tag/Versions/operation/createArtifactVersion
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

// GetArtifactVersionContent Retrieves a single version of the artifact content.
// Both the artifactId and the unique version number must be provided.
// The Content-Type of the response depends on the artifact type.
// In most cases, this is application/json, but for some types it may be different (for example, PROTOBUF).
// See https://www.apicur.io/registry/docs/apicurio-registry/3.0.x/assets-attachments/registry-rest-api.htm#tag/Versions/operation/getArtifactVersionContent
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

// UpdateArtifactVersionContent Updates the content of a single version of an artifact.
// See https://www.apicur.io/registry/docs/apicurio-registry/3.0.x/assets-attachments/registry-rest-api.htm#tag/Versions/operation/updateArtifactVersionContent
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

	if err := content.Validate(); err != nil {
		return errors.Wrap(err, "invalid content provided")
	}

	url := fmt.Sprintf("%s/groups/%s/artifacts/%s/versions/%s/content", api.Client.BaseURL, groupId, artifactId, versionExpression)

	resp, err := api.executeRequest(ctx, http.MethodPut, url, content)
	if err != nil {
		return err
	}

	return handleResponse(resp, http.StatusNoContent, nil)
}

// SearchForArtifactVersions Returns a paginated list of all versions that match the provided filter criteria.
// See https://www.apicur.io/registry/docs/apicurio-registry/3.0.x/assets-attachments/registry-rest-api.htm#tag/Versions/operation/searchVersions
func (api *VersionsAPI) SearchForArtifactVersions(
	ctx context.Context,
	params *models.SearchVersionParams,
) ([]models.ArtifactVersion, error) {

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

	return searchVersionsResponse.Versions, nil
}

// SearchForArtifactVersionByContent Returns a paginated list of all versions that match the posted content.
// See https://www.apicur.io/registry/docs/apicurio-registry/3.0.x/assets-attachments/registry-rest-api.htm#tag/Versions/operation/searchVersionsByContent
func (api *VersionsAPI) SearchForArtifactVersionByContent(
	ctx context.Context,
	content string,
	params *models.SearchVersionByContentParams,
) ([]models.ArtifactVersion, error) {
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

	return searchVersionsResponse.Versions, nil
}

// GetArtifactVersionState Gets the current state of an artifact version.
// See https://www.apicur.io/registry/docs/apicurio-registry/3.0.x/assets-attachments/registry-rest-api.htm#tag/Versions/operation/getArtifactVersionState
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

// UpdateArtifactVersionState Updates the state of an artifact version.
// NOTE: There are some restrictions on state transitions.
// Notably a version cannot be transitioned to the DRAFT state from any other state.
// The DRAFT state can only be entered (optionally) when creating a new artifact/version.
// A version in DRAFT state can only be transitioned to ENABLED.
// When this happens, any configured content rules will be applied.
// This may result in a failure to change the state.
// See https://www.apicur.io/registry/docs/apicurio-registry/3.0.x/assets-attachments/registry-rest-api.htm#tag/Versions/operation/updateArtifactVersionState
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
	return executeRequest(ctx, api.Client, method, url, body)
}
