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
func (api *BranchAPI) ListBranches(
	ctx context.Context,
	groupId, artifactId string,
	params *models.ListBranchesParams,
) ([]models.BranchInfo, error) {
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

	urlPath := fmt.Sprintf(
		"%s/groups/%s/artifacts/%s/branches%s",
		api.Client.BaseURL,
		url.PathEscape(groupId),
		url.PathEscape(artifactId),
		query,
	)
	resp, err := api.executeRequest(ctx, http.MethodGet, urlPath, nil)
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
func (api *BranchAPI) CreateBranch(
	ctx context.Context,
	groupId, artifactId string,
	branch *models.CreateBranchRequest,
) (*models.BranchInfo, error) {
	if err := validateInput(artifactId, regexGroupIDArtifactID, "Artifact ID"); err != nil {
		return nil, err
	}
	if err := validateInput(groupId, regexGroupIDArtifactID, "Group ID"); err != nil {
		return nil, err
	}

	if err := branch.Validate(); err != nil {
		return nil, errors.Wrap(err, "invalid branch provided")
	}

	urlPath := fmt.Sprintf(
		"%s/groups/%s/artifacts/%s/branches",
		api.Client.BaseURL,
		url.PathEscape(groupId),
		url.PathEscape(artifactId),
	)
	resp, err := api.executeRequest(ctx, http.MethodPost, urlPath, branch)
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
func (api *BranchAPI) GetBranchMetaData(
	ctx context.Context,
	groupId, artifactId, branchId string,
) (*models.BranchInfo, error) {
	if err := validateInput(artifactId, regexGroupIDArtifactID, "Artifact ID"); err != nil {
		return nil, err
	}
	if err := validateInput(groupId, regexGroupIDArtifactID, "Group ID"); err != nil {
		return nil, err
	}
	if err := validateInput(branchId, regexBranchID, "Branch ID"); err != nil {
		return nil, err
	}

	urlPath := fmt.Sprintf(
		"%s/groups/%s/artifacts/%s/branches/%s",
		api.Client.BaseURL,
		url.PathEscape(groupId),
		url.PathEscape(artifactId),
		url.PathEscape(branchId),
	)
	resp, err := api.executeRequest(ctx, http.MethodGet, urlPath, nil)
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
func (api *BranchAPI) UpdateBranchMetaData(
	ctx context.Context,
	groupId, artifactId, branchId, description string,
) error {
	if err := validateInput(artifactId, regexGroupIDArtifactID, "Artifact ID"); err != nil {
		return err
	}
	if err := validateInput(groupId, regexGroupIDArtifactID, "Group ID"); err != nil {
		return err
	}
	if err := validateInput(branchId, regexBranchID, "Branch ID"); err != nil {
		return err
	}

	urlPath := fmt.Sprintf(
		"%s/groups/%s/artifacts/%s/branches/%s",
		api.Client.BaseURL,
		url.PathEscape(groupId),
		url.PathEscape(artifactId),
		url.PathEscape(branchId),
	)

	branchMetaData := models.UpdateBranchMetaDataRequest{
		Description: description,
	}
	resp, err := api.executeRequest(ctx, http.MethodPut, urlPath, branchMetaData)
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
func (api *BranchAPI) DeleteBranch(
	ctx context.Context,
	groupId, artifactId, branchId string,
) error {
	if err := validateInput(artifactId, regexGroupIDArtifactID, "Artifact ID"); err != nil {
		return err
	}
	if err := validateInput(groupId, regexGroupIDArtifactID, "Group ID"); err != nil {
		return err
	}
	if err := validateInput(branchId, regexBranchID, "Branch ID"); err != nil {
		return err
	}

	urlPath := fmt.Sprintf(
		"%s/groups/%s/artifacts/%s/branches/%s",
		api.Client.BaseURL,
		url.PathEscape(groupId),
		url.PathEscape(artifactId),
		url.PathEscape(branchId),
	)
	resp, err := api.executeRequest(ctx, http.MethodDelete, urlPath, nil)
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
func (api *BranchAPI) GetVersionsInBranch(
	ctx context.Context,
	groupId, artifactId, branchId string,
	params *models.ListBranchesParams,
) ([]models.ArtifactVersion, error) {
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

	urlPath := fmt.Sprintf(
		"%s/groups/%s/artifacts/%s/branches/%s/versions%s",
		api.Client.BaseURL,
		url.PathEscape(groupId),
		url.PathEscape(artifactId),
		url.PathEscape(branchId),
		query,
	)
	resp, err := api.executeRequest(ctx, http.MethodGet, urlPath, nil)
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
func (api *BranchAPI) ReplaceVersionsInBranch(
	ctx context.Context,
	groupId, artifactId, branchId string,
	versions []string,
) error {
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

	urlPath := fmt.Sprintf(
		"%s/groups/%s/artifacts/%s/branches/%s/versions",
		api.Client.BaseURL,
		url.PathEscape(groupId),
		url.PathEscape(artifactId),
		url.PathEscape(branchId),
	)

	requestBody := map[string]interface{}{
		"versions": versions,
	}

	resp, err := api.executeRequest(ctx, http.MethodPut, urlPath, requestBody)
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
func (api *BranchAPI) AddVersionToBranch(
	ctx context.Context,
	groupId, artifactId, branchId, version string,
) error {
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

	urlPath := fmt.Sprintf(
		"%s/groups/%s/artifacts/%s/branches/%s/versions",
		api.Client.BaseURL,
		url.PathEscape(groupId),
		url.PathEscape(artifactId),
		url.PathEscape(branchId),
	)

	requestBody := map[string]interface{}{
		"version": version,
	}
	resp, err := api.executeRequest(ctx, http.MethodPost, urlPath, requestBody)
	if err != nil {
		return err
	}

	if err := handleResponse(resp, http.StatusNoContent, nil); err != nil {
		return err
	}

	return nil
}

// executeRequest handles the creation and execution of an HTTP request.
func (api *BranchAPI) executeRequest(
	ctx context.Context,
	method, url string,
	body interface{},
) (*http.Response, error) {
	return executeRequest(ctx, api.Client, method, url, body)
}
