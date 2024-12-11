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

// MetadataAPI handles metadata-related operations for artifacts.
type MetadataAPI struct {
	Client *client.Client
}

// NewMetadataAPI creates a new MetadataAPI instance.
func NewMetadataAPI(client *client.Client) *MetadataAPI {
	return &MetadataAPI{
		Client: client,
	}
}

// GetArtifactVersionMetadata retrieves metadata for a single artifact version.
func (api *MetadataAPI) GetArtifactVersionMetadata(ctx context.Context, groupId, artifactId, versionExpression string) (*models.ArtifactVersionMetadata, error) {
	if err := validateInput(groupId, regexGroupIDArtifactID, "Group ID"); err != nil {
		return nil, err
	}
	if err := validateInput(artifactId, regexGroupIDArtifactID, "Artifact ID"); err != nil {
		return nil, err
	}
	if err := validateInput(versionExpression, regexVersion, "Version Expression"); err != nil {
		return nil, err
	}

	url := fmt.Sprintf("%s/groups/%s/artifacts/%s/versions/%s", api.Client.BaseURL, groupId, artifactId, versionExpression)

	resp, err := api.executeRequest(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	var metadata models.ArtifactVersionMetadata
	if err := handleResponse(resp, http.StatusOK, &metadata); err != nil {
		return nil, err
	}

	return &metadata, nil
}

// UpdateArtifactVersionMetadata updates the user-editable metadata of an artifact version.
func (api *MetadataAPI) UpdateArtifactVersionMetadata(ctx context.Context, groupId, artifactId, versionExpression string, metadata models.UpdateArtifactMetadataRequest) error {
	if err := validateInput(groupId, regexGroupIDArtifactID, "Group ID"); err != nil {
		return err
	}
	if err := validateInput(artifactId, regexGroupIDArtifactID, "Artifact ID"); err != nil {
		return err
	}
	if err := validateInput(versionExpression, regexVersion, "Version Expression"); err != nil {
		return err
	}

	url := fmt.Sprintf("%s/groups/%s/artifacts/%s/versions/%s", api.Client.BaseURL, groupId, artifactId, versionExpression)

	resp, err := api.executeRequest(ctx, http.MethodPut, url, metadata)
	if err != nil {
		return err
	}

	return handleResponse(resp, http.StatusNoContent, nil)
}

// GetArtifactMetadata retrieves metadata for an artifact based on the latest version or the next available non-disabled version.
func (api *MetadataAPI) GetArtifactMetadata(ctx context.Context, groupId, artifactId string) (*models.ArtifactMetadata, error) {
	if err := validateInput(groupId, regexGroupIDArtifactID, "Group ID"); err != nil {
		return nil, err
	}
	if err := validateInput(artifactId, regexGroupIDArtifactID, "Artifact ID"); err != nil {
		return nil, err
	}

	url := fmt.Sprintf("%s/groups/%s/artifacts/%s", api.Client.BaseURL, groupId, artifactId)

	resp, err := api.executeRequest(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	var metadata models.ArtifactMetadata
	if err := handleResponse(resp, http.StatusOK, &metadata); err != nil {
		return nil, err
	}

	return &metadata, nil
}

// UpdateArtifactMetadata updates the editable parts of an artifact's metadata.
func (api *MetadataAPI) UpdateArtifactMetadata(ctx context.Context, groupId, artifactId string, metadata models.UpdateArtifactMetadataRequest) error {
	if err := validateInput(groupId, regexGroupIDArtifactID, "Group ID"); err != nil {
		return err
	}
	if err := validateInput(artifactId, regexGroupIDArtifactID, "Artifact ID"); err != nil {
		return err
	}

	// Construct the URL
	url := fmt.Sprintf("%s/groups/%s/artifacts/%s", api.Client.BaseURL, groupId, artifactId)

	resp, err := api.executeRequest(ctx, http.MethodPut, url, metadata)
	if err != nil {
		return err
	}

	return handleResponse(resp, http.StatusNoContent, nil)
}

// executeRequest executes an HTTP request with the given method, URL, and body.
func (api *MetadataAPI) executeRequest(ctx context.Context, method, url string, body interface{}) (*http.Response, error) {
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
