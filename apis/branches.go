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

type BranchAPI struct {
	Client *client.Client
}

func NewBranchAPI(client *client.Client) *BranchAPI {
	return &BranchAPI{
		Client: client,
	}
}

// ListBranches Returns a list of all branches in the artifact.
// Each branch is a list of version identifiers, ordered from the latest (tip of the branch) to the oldest.
// See https://www.apicur.io/registry/docs/apicurio-registry/3.0.x/assets-attachments/registry-rest-api.htm#tag/Branches/operation/listBranches
func (api *BranchAPI) ListBranches(ctx context.Context, groupId, artifactId string, params *models.ListBranchesParams) ([]models.BranchInfo, error) {
	if err := validateInput(artifactId, regexGroupIDArtifactID, "Artifact ID"); err != nil {
		return nil, err
	}
	if err := validateInput(groupId, regexGroupIDArtifactID, "Group ID"); err != nil {
		return nil, err
	}

	query := ""
	if params != nil {
		if err := params.Validate(); err != nil {
			return nil, errors.Wrap(err, "invalid parameters provided")
		}
		query = "?" + params.ToQuery().Encode()
	}

	url := fmt.Sprintf("%s/groups/%s/artifacts/%s/branches%s", api.Client.BaseURL, groupId, artifactId, query)
	resp, err := api.executeRequest(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	var result models.BranchesInfoResponse
	if err := handleResponse(resp, http.StatusOK, &result); err != nil {
		return nil, err
	}

	return result.Branches, nil

}

// CreateBranch Creates a new branch for the artifact.
// A new branch consists of metadata and a list of versions.
// See https://www.apicur.io/registry/docs/apicurio-registry/3.0.x/assets-attachments/registry-rest-api.htm#tag/Branches/operation/createBranch
func (api *BranchAPI) CreateBranch(ctx context.Context, groupId, artifactId string, branch *models.CreateBranchRequest) (*models.BranchInfo, error) {
	if err := validateInput(artifactId, regexGroupIDArtifactID, "Artifact ID"); err != nil {
		return nil, err
	}
	if err := validateInput(groupId, regexGroupIDArtifactID, "Group ID"); err != nil {
		return nil, err
	}

	if err := branch.Validate(); err != nil {
		return nil, errors.Wrap(err, "invalid branch provided")
	}

	url := fmt.Sprintf("%s/groups/%s/artifacts/%s/branches", api.Client.BaseURL, groupId, artifactId)
	resp, err := api.executeRequest(ctx, http.MethodPost, url, branch)
	if err != nil {
		return nil, err
	}

	var result models.BranchInfo
	if err := handleResponse(resp, http.StatusOK, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// GetBranchMetaData Get branch metaData
// See https://www.apicur.io/registry/docs/apicurio-registry/3.0.x/assets-attachments/registry-rest-api.htm#tag/Branches/operation/getBranchMetaData
func (api *BranchAPI) GetBranchMetaData(ctx context.Context, groupId, artifactId, branchId string) (*models.BranchInfo, error) {
	if err := validateInput(artifactId, regexGroupIDArtifactID, "Artifact ID"); err != nil {
		return nil, err
	}
	if err := validateInput(groupId, regexGroupIDArtifactID, "Group ID"); err != nil {
		return nil, err
	}
	if err := validateInput(branchId, regexBranchID, "Branch ID"); err != nil {
		return nil, err
	}

	url := fmt.Sprintf("%s/groups/%s/artifacts/%s/branches/%s", api.Client.BaseURL, groupId, artifactId, branchId)
	resp, err := api.executeRequest(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	var result models.BranchInfo
	if err := handleResponse(resp, http.StatusOK, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// UpdateBranchMetaData Update branch metaData
// See https://www.apicur.io/registry/docs/apicurio-registry/3.0.x/assets-attachments/registry-rest-api.htm#tag/Branches/operation/updateBranchMetaData
func (api *BranchAPI) UpdateBranchMetaData(ctx context.Context, groupId, artifactId, branchId, description string) error {
	if err := validateInput(artifactId, regexGroupIDArtifactID, "Artifact ID"); err != nil {
		return err
	}
	if err := validateInput(groupId, regexGroupIDArtifactID, "Group ID"); err != nil {
		return err
	}
	if err := validateInput(branchId, regexBranchID, "Branch ID"); err != nil {
		return err
	}

	url := fmt.Sprintf("%s/groups/%s/artifacts/%s/branches/%s", api.Client.BaseURL, groupId, artifactId, branchId)

	branchMetaData := models.UpdateBranchMetaDataRequest{
		Description: description,
	}
	resp, err := api.executeRequest(ctx, http.MethodPut, url, branchMetaData)
	if err != nil {
		return err
	}

	if err := handleResponse(resp, http.StatusNoContent, nil); err != nil {
		return err
	}

	return nil

}

// DeleteBranch Deletes a single branch in the artifact.
// See https://www.apicur.io/registry/docs/apicurio-registry/3.0.x/assets-attachments/registry-rest-api.htm#tag/Branches/operation/deleteBranch
func (api *BranchAPI) DeleteBranch(ctx context.Context, groupId, artifactId, branchId string) error {
	if err := validateInput(artifactId, regexGroupIDArtifactID, "Artifact ID"); err != nil {
		return err
	}
	if err := validateInput(groupId, regexGroupIDArtifactID, "Group ID"); err != nil {
		return err
	}
	if err := validateInput(branchId, regexBranchID, "Branch ID"); err != nil {
		return err
	}

	url := fmt.Sprintf("%s/groups/%s/artifacts/%s/branches/%s", api.Client.BaseURL, groupId, artifactId, branchId)
	resp, err := api.executeRequest(ctx, http.MethodDelete, url, nil)
	if err != nil {
		return err
	}

	if err := handleResponse(resp, http.StatusNoContent, nil); err != nil {
		return err
	}

	return nil
}

// GetVersionsInBranch Get a list of all versions in the branch.
// Returns a list of version identifiers in the branch, ordered from the latest (tip of the branch) to the oldest.
// See https://www.apicur.io/registry/docs/apicurio-registry/3.0.x/assets-attachments/registry-rest-api.htm#tag/Branches/operation/listBranchVersions
func (api *BranchAPI) GetVersionsInBranch(ctx context.Context, groupId, artifactId, branchId string, params *models.ListBranchesParams) ([]models.ArtifactVersion, error) {
	if err := validateInput(artifactId, regexGroupIDArtifactID, "Artifact ID"); err != nil {
		return nil, err
	}
	if err := validateInput(groupId, regexGroupIDArtifactID, "Group ID"); err != nil {
		return nil, err
	}
	if err := validateInput(branchId, regexBranchID, "Branch ID"); err != nil {
		return nil, err
	}

	query := ""
	if params != nil {
		if err := params.Validate(); err != nil {
			return nil, errors.Wrap(err, "invalid parameters provided")
		}
		query = "?" + params.ToQuery().Encode()
	}

	url := fmt.Sprintf("%s/groups/%s/artifacts/%s/branches/%s/versions%s", api.Client.BaseURL, groupId, artifactId, branchId, query)
	resp, err := api.executeRequest(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	var result models.ArtifactVersionListResponse
	if err := handleResponse(resp, http.StatusOK, &result); err != nil {
		return nil, err
	}

	return result.Versions, nil
}

// ReplaceVersionsInBranch Add a new version to an artifact branch. Branch is created if it does not exist.
// Returns a list of version identifiers in the artifact branch, ordered from the latest (tip of the branch) to the oldest.
// See https://www.apicur.io/registry/docs/apicurio-registry/3.0.x/assets-attachments/registry-rest-api.htm#tag/Branches/operation/replaceBranchVersions
func (api *BranchAPI) ReplaceVersionsInBranch(ctx context.Context, groupId, artifactId, branchId string, versions []string) error {
	if err := validateInput(artifactId, regexGroupIDArtifactID, "Artifact ID"); err != nil {
		return err
	}
	if err := validateInput(groupId, regexGroupIDArtifactID, "Group ID"); err != nil {
		return err
	}
	if err := validateInput(branchId, regexBranchID, "Branch ID"); err != nil {
		return err
	}

	if len(versions) == 0 {
		return errors.New("versions must not be empty")
	}

	for _, version := range versions {
		err := validateInput(version, regexVersion, "Version")
		if err != nil {
			return err
		}
	}

	url := fmt.Sprintf("%s/groups/%s/artifacts/%s/branches/%s/versions", api.Client.BaseURL, groupId, artifactId, branchId)

	requestBody := map[string]interface{}{
		"versions": versions,
	}

	resp, err := api.executeRequest(ctx, http.MethodPut, url, requestBody)
	if err != nil {
		return err
	}

	if err := handleResponse(resp, http.StatusNoContent, nil); err != nil {
		return err
	}

	return nil

}

// AddVersionToBranch Add a new version to an artifact branch.
// See https://www.apicur.io/registry/docs/apicurio-registry/3.0.x/assets-attachments/registry-rest-api.htm#tag/Branches/operation/addVersionToBranch
func (api *BranchAPI) AddVersionToBranch(ctx context.Context, groupId, artifactId, branchId, version string) error {
	if err := validateInput(artifactId, regexGroupIDArtifactID, "Artifact ID"); err != nil {
		return err
	}
	if err := validateInput(groupId, regexGroupIDArtifactID, "Group ID"); err != nil {
		return err
	}
	if err := validateInput(branchId, regexBranchID, "Branch ID"); err != nil {
		return err
	}
	if err := validateInput(version, regexVersion, "Version"); err != nil {
		return err
	}

	url := fmt.Sprintf("%s/groups/%s/artifacts/%s/branches/%s/versions", api.Client.BaseURL, groupId, artifactId, branchId)

	requestBody := map[string]interface{}{
		"version": version,
	}
	resp, err := api.executeRequest(ctx, http.MethodPost, url, requestBody)
	if err != nil {
		return err
	}

	if err := handleResponse(resp, http.StatusNoContent, nil); err != nil {
		return err
	}

	return nil
}

// executeRequest handles the creation and execution of an HTTP request.
func (api *BranchAPI) executeRequest(ctx context.Context, method, url string, body interface{}) (*http.Response, error) {
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
