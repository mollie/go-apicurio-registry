package apis

import (
	"context"
	"fmt"
	"github.com/mollie/go-apicurio-registry/client"
	"github.com/mollie/go-apicurio-registry/models"
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
	return executeRequest(ctx, api.Client, method, url, body)
}
