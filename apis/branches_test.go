package apis_test

import (
	"context"
	"encoding/json"
	"github.com/mollie/go-apicurio-registry/apis"
	"github.com/mollie/go-apicurio-registry/client"
	"github.com/mollie/go-apicurio-registry/models"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestBranchAPI_ListBranches(t *testing.T) {
	expectedURL := "/groups/" + stubGroupId + "/artifacts/" + stubArtifactId + "/branches"

	t.Run("Success", func(t *testing.T) {
		mockResponse := models.BranchesInfoResponse{
			Branches: []models.BranchInfo{
				{
					GroupId:       stubGroupId,
					ArtifactId:    stubArtifactId,
					BranchId:      stubBranchID,
					Description:   stubDescription,
					SystemDefined: false,
					CreatedOn:     "2018-02-10T09:30Z",
					ModifiedOn:    "2018-02-10T09:30Z",
					ModifiedBy:    "2018-02-10T09:30Z",
				},
			},
			Count: 1,
		}
		server := setupMockServer(t, http.StatusOK, mockResponse, expectedURL, http.MethodGet)
		defer server.Close()

		mockClient := &client.Client{BaseURL: server.URL, HTTPClient: server.Client()}
		api := apis.NewBranchAPI(mockClient)

		branches, err := api.ListBranches(context.Background(), stubGroupId, stubArtifactId, nil)
		assert.NoError(t, err)
		assert.NotNil(t, branches)
		assert.Len(t, branches, 1)
		assert.Equal(t, stubGroupId, branches[0].GroupId)
		assert.Equal(t, stubArtifactId, branches[0].ArtifactId)
		assert.Equal(t, stubBranchID, branches[0].BranchId)
		assert.Equal(t, stubDescription, branches[0].Description)
	})

	t.Run("Validation Errors", func(t *testing.T) {
		mockClient := &client.Client{BaseURL: "http://mock.server", HTTPClient: http.DefaultClient}
		api := apis.NewBranchAPI(mockClient)

		_, err := api.ListBranches(context.Background(), "", stubArtifactId, nil)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "Group ID")

		_, err = api.ListBranches(context.Background(), stubGroupId, "", nil)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "Artifact ID")
	})

	t.Run("Not Found", func(t *testing.T) {
		mockErrorResponse := models.APIError{Status: http.StatusNotFound, Title: TitleNotFound}

		server := setupMockServer(t, http.StatusNotFound, mockErrorResponse, expectedURL, http.MethodGet)
		defer server.Close()

		mockClient := &client.Client{BaseURL: server.URL, HTTPClient: server.Client()}
		api := apis.NewBranchAPI(mockClient)

		branches, err := api.ListBranches(context.Background(), stubGroupId, stubArtifactId, nil)
		assert.Error(t, err)
		assert.Nil(t, branches)

		assertAPIError(t, err, http.StatusNotFound, TitleNotFound)
	})

	t.Run("Internal Server Error", func(t *testing.T) {
		mockErrorResponse := models.APIError{Status: http.StatusInternalServerError, Title: TitleInternalServerError}

		server := setupMockServer(t, http.StatusInternalServerError, mockErrorResponse, expectedURL, http.MethodGet)
		defer server.Close()

		mockClient := &client.Client{BaseURL: server.URL, HTTPClient: server.Client()}
		api := apis.NewBranchAPI(mockClient)

		branches, err := api.ListBranches(context.Background(), stubGroupId, stubArtifactId, nil)
		assert.Error(t, err)
		assert.Nil(t, branches)

		assertAPIError(t, err, http.StatusInternalServerError, TitleInternalServerError)
	})
}

func TestBranchAPI_CreateBranch(t *testing.T) {
	expectedURL := "/groups/" + stubGroupId + "/artifacts/" + stubArtifactId + "/branches"

	t.Run("Success", func(t *testing.T) {
		mockResponse := models.BranchInfo{
			GroupId:       stubGroupId,
			ArtifactId:    stubArtifactId,
			BranchId:      stubBranchID,
			Description:   stubDescription,
			SystemDefined: false,
			CreatedOn:     "2018-02-10T09:30Z",
			ModifiedOn:    "2018-02-10T09:30Z",
			ModifiedBy:    "2018-02-10T09:30Z",
		}

		server := setupMockServer(t, http.StatusOK, mockResponse, expectedURL, http.MethodPost)
		defer server.Close()

		mockClient := &client.Client{BaseURL: server.URL, HTTPClient: server.Client()}
		api := apis.NewBranchAPI(mockClient)

		branchInfo := models.CreateBranchRequest{
			BranchID:    stubBranchID,
			Description: stubDescription,
		}
		branch, err := api.CreateBranch(context.Background(), stubGroupId, stubArtifactId, &branchInfo)
		assert.NoError(t, err)
		assert.NotNil(t, branch)
		assert.Equal(t, stubGroupId, branch.GroupId)
		assert.Equal(t, stubArtifactId, branch.ArtifactId)
		assert.Equal(t, stubBranchID, branch.BranchId)
		assert.Equal(t, stubDescription, branch.Description)
	})

	t.Run("Validation Errors", func(t *testing.T) {
		mockClient := &client.Client{BaseURL: "http://mock.server", HTTPClient: http.DefaultClient}
		api := apis.NewBranchAPI(mockClient)

		invalidBranch := models.CreateBranchRequest{BranchID: "", Description: ""}

		_, err := api.CreateBranch(context.Background(), "", stubArtifactId, &invalidBranch)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "Group ID")

		_, err = api.CreateBranch(context.Background(), stubGroupId, "", &invalidBranch)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "Artifact ID")

		_, err = api.CreateBranch(context.Background(), stubGroupId, stubArtifactId, &invalidBranch)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid branch provided")
	})

	t.Run("Conflict", func(t *testing.T) {
		mockErrorResponse := models.APIError{Status: http.StatusConflict, Title: TitleConflict}

		server := setupMockServer(t, http.StatusConflict, mockErrorResponse, expectedURL, http.MethodPost)
		defer server.Close()

		mockClient := &client.Client{BaseURL: server.URL, HTTPClient: server.Client()}
		api := apis.NewBranchAPI(mockClient)

		branchInfo := models.CreateBranchRequest{
			BranchID:    stubBranchID,
			Description: stubDescription,
		}
		branch, err := api.CreateBranch(context.Background(), stubGroupId, stubArtifactId, &branchInfo)
		assert.Error(t, err)
		assert.Nil(t, branch)

		assertAPIError(t, err, http.StatusConflict, TitleConflict)
	})

	t.Run("Not Found", func(t *testing.T) {
		mockErrorResponse := models.APIError{Status: http.StatusNotFound, Title: TitleNotFound}

		server := setupMockServer(t, http.StatusNotFound, mockErrorResponse, expectedURL, http.MethodPost)
		defer server.Close()

		mockClient := &client.Client{BaseURL: server.URL, HTTPClient: server.Client()}
		api := apis.NewBranchAPI(mockClient)

		branchInfo := models.CreateBranchRequest{
			BranchID:    stubBranchID,
			Description: stubDescription,
		}
		branch, err := api.CreateBranch(context.Background(), stubGroupId, stubArtifactId, &branchInfo)
		assert.Error(t, err)
		assert.Nil(t, branch)

		assertAPIError(t, err, http.StatusNotFound, TitleNotFound)
	})

	t.Run("Internal Server Error", func(t *testing.T) {
		mockErrorResponse := models.APIError{Status: http.StatusInternalServerError, Title: TitleInternalServerError}

		server := setupMockServer(t, http.StatusInternalServerError, mockErrorResponse, expectedURL, http.MethodPost)
		defer server.Close()

		mockClient := &client.Client{BaseURL: server.URL, HTTPClient: server.Client()}
		api := apis.NewBranchAPI(mockClient)

		branchInfo := models.CreateBranchRequest{
			BranchID:    stubBranchID,
			Description: stubDescription,
		}
		branch, err := api.CreateBranch(context.Background(), stubGroupId, stubArtifactId, &branchInfo)
		assert.Error(t, err)
		assert.Nil(t, branch)

		assertAPIError(t, err, http.StatusInternalServerError, TitleInternalServerError)
	})
}

func TestBranchAPI_GetBranchMetaData(t *testing.T) {
	expectedURL := "/groups/" + stubGroupId + "/artifacts/" + stubArtifactId + "/branches/" + stubBranchID

	t.Run("Success", func(t *testing.T) {
		mockResponse := models.BranchInfo{
			GroupId:       stubGroupId,
			ArtifactId:    stubArtifactId,
			BranchId:      stubBranchID,
			Description:   stubDescription,
			SystemDefined: false,
			CreatedOn:     "2018-02-10T09:30Z",
			ModifiedOn:    "2018-02-10T09:30Z",
			ModifiedBy:    "2018-02-10T09:30Z",
		}

		server := setupMockServer(t, http.StatusOK, mockResponse, expectedURL, http.MethodGet)
		defer server.Close()

		mockClient := &client.Client{BaseURL: server.URL, HTTPClient: server.Client()}
		api := apis.NewBranchAPI(mockClient)

		branch, err := api.GetBranchMetaData(context.Background(), stubGroupId, stubArtifactId, stubBranchID)
		assert.NoError(t, err)
		assert.NotNil(t, branch)
		assert.Equal(t, stubGroupId, branch.GroupId)
		assert.Equal(t, stubArtifactId, branch.ArtifactId)
		assert.Equal(t, stubBranchID, branch.BranchId)
		assert.Equal(t, stubDescription, branch.Description)
	})

	t.Run("Validation Errors", func(t *testing.T) {
		mockClient := &client.Client{BaseURL: "http://mock.server", HTTPClient: http.DefaultClient}
		api := apis.NewBranchAPI(mockClient)

		_, err := api.GetBranchMetaData(context.Background(), "", stubArtifactId, stubBranchID)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "Group ID")

		_, err = api.GetBranchMetaData(context.Background(), stubGroupId, "", stubBranchID)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "Artifact ID")

		_, err = api.GetBranchMetaData(context.Background(), stubGroupId, stubArtifactId, "")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "Branch ID")
	})

	t.Run("Not Found", func(t *testing.T) {
		mockErrorResponse := models.APIError{Status: http.StatusNotFound, Title: TitleNotFound}

		server := setupMockServer(t, http.StatusNotFound, mockErrorResponse, expectedURL, http.MethodGet)
		defer server.Close()

		mockClient := &client.Client{BaseURL: server.URL, HTTPClient: server.Client()}
		api := apis.NewBranchAPI(mockClient)

		branch, err := api.GetBranchMetaData(context.Background(), stubGroupId, stubArtifactId, stubBranchID)
		assert.Error(t, err)
		assert.Nil(t, branch)

		assertAPIError(t, err, http.StatusNotFound, TitleNotFound)
	})

	t.Run("Internal Server Error", func(t *testing.T) {
		mockErrorResponse := models.APIError{Status: http.StatusInternalServerError, Title: TitleInternalServerError}

		server := setupMockServer(t, http.StatusInternalServerError, mockErrorResponse, expectedURL, http.MethodGet)
		defer server.Close()

		mockClient := &client.Client{BaseURL: server.URL, HTTPClient: server.Client()}
		api := apis.NewBranchAPI(mockClient)

		branch, err := api.GetBranchMetaData(context.Background(), stubGroupId, stubArtifactId, stubBranchID)
		assert.Error(t, err)
		assert.Nil(t, branch)

		assertAPIError(t, err, http.StatusInternalServerError, TitleInternalServerError)
	})

}

func TestBranchAPI_UpdateBranchMetaData(t *testing.T) {
	expectedURL := "/groups/" + stubGroupId + "/artifacts/" + stubArtifactId + "/branches/" + stubBranchID

	t.Run("Success", func(t *testing.T) {
		mockResponse := models.BranchInfo{
			GroupId:       stubGroupId,
			ArtifactId:    stubArtifactId,
			BranchId:      stubBranchID,
			Description:   stubDescription,
			SystemDefined: false,
			CreatedOn:     "2018-02-10T09:30Z",
			ModifiedOn:    "2018-02-10T09:30Z",
			ModifiedBy:    "2018-02-10T09:30Z",
		}

		server := setupMockServer(t, http.StatusOK, mockResponse, expectedURL, http.MethodPut)
		defer server.Close()

		mockClient := &client.Client{BaseURL: server.URL, HTTPClient: server.Client()}
		api := apis.NewBranchAPI(mockClient)

		branch, err := api.UpdateBranchMetaData(context.Background(), stubGroupId, stubArtifactId, stubBranchID, stubUpdatedDescription)
		assert.NoError(t, err)
		assert.NotNil(t, branch)
		assert.Equal(t, stubGroupId, branch.GroupId)
		assert.Equal(t, stubArtifactId, branch.ArtifactId)
		assert.Equal(t, stubBranchID, branch.BranchId)
		assert.Equal(t, stubDescription, branch.Description)
	})

	t.Run("Validation Errors", func(t *testing.T) {
		mockClient := &client.Client{BaseURL: "http://mock.server", HTTPClient: http.DefaultClient}
		api := apis.NewBranchAPI(mockClient)

		// Empty GroupID
		_, err := api.UpdateBranchMetaData(context.Background(), "", stubArtifactId, stubBranchID, stubUpdatedDescription)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "Group ID")

		// Empty ArtifactID
		_, err = api.UpdateBranchMetaData(context.Background(), stubGroupId, "", stubBranchID, stubUpdatedDescription)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "Artifact ID")

		// Empty BranchID
		_, err = api.UpdateBranchMetaData(context.Background(), stubGroupId, stubArtifactId, "", stubUpdatedDescription)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "Branch ID")
	})

	t.Run("Not Found", func(t *testing.T) {
		mockErrorResponse := models.APIError{Status: http.StatusNotFound, Title: TitleNotFound}

		server := setupMockServer(t, http.StatusNotFound, mockErrorResponse, expectedURL, http.MethodPut)
		defer server.Close()

		mockClient := &client.Client{BaseURL: server.URL, HTTPClient: server.Client()}
		api := apis.NewBranchAPI(mockClient)

		branch, err := api.UpdateBranchMetaData(context.Background(), stubGroupId, stubArtifactId, stubBranchID, stubUpdatedDescription)
		assert.Error(t, err)
		assert.Nil(t, branch)

		assertAPIError(t, err, http.StatusNotFound, TitleNotFound)
	})

	t.Run("Internal Server Error", func(t *testing.T) {
		mockErrorResponse := models.APIError{Status: http.StatusInternalServerError, Title: TitleInternalServerError}

		server := setupMockServer(t, http.StatusInternalServerError, mockErrorResponse, expectedURL, http.MethodPut)
		defer server.Close()

		mockClient := &client.Client{BaseURL: server.URL, HTTPClient: server.Client()}
		api := apis.NewBranchAPI(mockClient)

		branch, err := api.UpdateBranchMetaData(context.Background(), stubGroupId, stubArtifactId, stubBranchID, stubUpdatedDescription)
		assert.Error(t, err)
		assert.Nil(t, branch)

		assertAPIError(t, err, http.StatusInternalServerError, TitleInternalServerError)
	})
}

func TestBranchAPI_DeleteBranch(t *testing.T) {
	expectedURL := "/groups/" + stubGroupId + "/artifacts/" + stubArtifactId + "/branches/" + stubBranchID

	t.Run("Success", func(t *testing.T) {
		server := setupMockServer(t, http.StatusNoContent, nil, expectedURL, http.MethodDelete)
		defer server.Close()

		mockClient := &client.Client{BaseURL: server.URL, HTTPClient: server.Client()}
		api := apis.NewBranchAPI(mockClient)

		err := api.DeleteBranch(context.Background(), stubGroupId, stubArtifactId, stubBranchID)
		assert.NoError(t, err)
	})

	t.Run("Validation Errors", func(t *testing.T) {
		mockClient := &client.Client{BaseURL: "http://mock.server", HTTPClient: http.DefaultClient}
		api := apis.NewBranchAPI(mockClient)

		// Empty GroupID
		err := api.DeleteBranch(context.Background(), "", stubArtifactId, stubBranchID)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "Group ID")

		// Empty ArtifactID
		err = api.DeleteBranch(context.Background(), stubGroupId, "", stubBranchID)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "Artifact ID")

		// Empty BranchID
		err = api.DeleteBranch(context.Background(), stubGroupId, stubArtifactId, "")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "Branch ID")
	})

	t.Run("Not Found", func(t *testing.T) {
		mockErrorResponse := models.APIError{Status: http.StatusNotFound, Title: TitleNotFound}

		server := setupMockServer(t, http.StatusNotFound, mockErrorResponse, expectedURL, http.MethodDelete)
		defer server.Close()

		mockClient := &client.Client{BaseURL: server.URL, HTTPClient: server.Client()}
		api := apis.NewBranchAPI(mockClient)

		err := api.DeleteBranch(context.Background(), stubGroupId, stubArtifactId, stubBranchID)
		assert.Error(t, err)

		assertAPIError(t, err, http.StatusNotFound, TitleNotFound)
	})

	t.Run("Conflict", func(t *testing.T) {
		mockErrorResponse := models.APIError{Status: http.StatusConflict, Title: TitleConflict}

		server := setupMockServer(t, http.StatusConflict, mockErrorResponse, expectedURL, http.MethodDelete)
		defer server.Close()

		mockClient := &client.Client{BaseURL: server.URL, HTTPClient: server.Client()}
		api := apis.NewBranchAPI(mockClient)

		err := api.DeleteBranch(context.Background(), stubGroupId, stubArtifactId, stubBranchID)
		assert.Error(t, err)

		assertAPIError(t, err, http.StatusConflict, TitleConflict)
	})

	t.Run("Internal Server Error", func(t *testing.T) {
		mockErrorResponse := models.APIError{Status: http.StatusInternalServerError, Title: TitleInternalServerError}

		server := setupMockServer(t, http.StatusInternalServerError, mockErrorResponse, expectedURL, http.MethodDelete)
		defer server.Close()

		mockClient := &client.Client{BaseURL: server.URL, HTTPClient: server.Client()}
		api := apis.NewBranchAPI(mockClient)

		err := api.DeleteBranch(context.Background(), stubGroupId, stubArtifactId, stubBranchID)
		assert.Error(t, err)

		assertAPIError(t, err, http.StatusInternalServerError, TitleInternalServerError)
	})
}

func TestBranchAPI_GetVersionsInBranch(t *testing.T) {
	expectedURL := "/groups/" + stubGroupId + "/artifacts/" + stubArtifactId + "/branches/" + stubBranchID + "/versions"

	t.Run("Success", func(t *testing.T) {
		mockResponse := models.ArtifactVersionListResponse{
			Count: 1,
			Versions: []models.ArtifactVersion{
				{
					CreatedOn:    "2024-12-10T08:56:40Z",
					ArtifactType: models.Json,
					State:        models.StateEnabled,
					GlobalID:     47,
					Version:      stubVersionID2,
					ContentID:    47,
					ArtifactID:   stubArtifactId,
					GroupID:      stubGroupId,
					ModifiedOn:   "2024-12-10T08:56:40Z",
				},
				{
					CreatedOn:    "2024-12-10T08:56:17Z",
					ArtifactType: models.Json,
					State:        models.StateEnabled,
					GlobalID:     46,
					Version:      stubVersionID,
					ContentID:    46,
					ArtifactID:   stubArtifactId,
					GroupID:      stubGroupId,
					ModifiedOn:   "2024-12-10T08:56:17Z",
				},
			},
		}

		server := setupMockServer(t, http.StatusOK, mockResponse, expectedURL, http.MethodGet)
		defer server.Close()

		mockClient := &client.Client{BaseURL: server.URL, HTTPClient: server.Client()}
		api := apis.NewBranchAPI(mockClient)

		versions, err := api.GetVersionsInBranch(context.Background(), stubGroupId, stubArtifactId, stubBranchID, nil)
		assert.NoError(t, err)
		assert.NotNil(t, versions)
		assert.Len(t, versions, 2)
		assert.Equal(t, stubVersionID, versions[1].Version)
		assert.Equal(t, stubVersionID2, versions[0].Version)

	})

	t.Run("Validation Errors", func(t *testing.T) {
		mockClient := &client.Client{BaseURL: "http://mock.server", HTTPClient: http.DefaultClient}
		api := apis.NewBranchAPI(mockClient)

		// Empty GroupID
		_, err := api.GetVersionsInBranch(context.Background(), "", stubArtifactId, stubBranchID, nil)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "Group ID")

		// Empty ArtifactID
		_, err = api.GetVersionsInBranch(context.Background(), stubGroupId, "", stubBranchID, nil)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "Artifact ID")

		// Empty BranchID
		_, err = api.GetVersionsInBranch(context.Background(), stubGroupId, stubArtifactId, "", nil)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "Branch ID")
	})

	t.Run("Not Found", func(t *testing.T) {
		mockErrorResponse := models.APIError{Status: http.StatusNotFound, Title: TitleNotFound}

		server := setupMockServer(t, http.StatusNotFound, mockErrorResponse, expectedURL, http.MethodGet)
		defer server.Close()

		mockClient := &client.Client{BaseURL: server.URL, HTTPClient: server.Client()}
		api := apis.NewBranchAPI(mockClient)

		versions, err := api.GetVersionsInBranch(context.Background(), stubGroupId, stubArtifactId, stubBranchID, nil)
		assert.Error(t, err)
		assert.Nil(t, versions)

		assertAPIError(t, err, http.StatusNotFound, TitleNotFound)
	})

	t.Run("Internal Server Error", func(t *testing.T) {
		mockErrorResponse := models.APIError{Status: http.StatusInternalServerError, Title: TitleInternalServerError}

		server := setupMockServer(t, http.StatusInternalServerError, mockErrorResponse, expectedURL, http.MethodGet)
		defer server.Close()

		mockClient := &client.Client{BaseURL: server.URL, HTTPClient: server.Client()}
		api := apis.NewBranchAPI(mockClient)

		versions, err := api.GetVersionsInBranch(context.Background(), stubGroupId, stubArtifactId, stubBranchID, nil)
		assert.Error(t, err)
		assert.Nil(t, versions)

		assertAPIError(t, err, http.StatusInternalServerError, TitleInternalServerError)
	})
}

func TestBranchAPI_ReplaceVersionsInBranch(t *testing.T) {
	expectedURL := "/groups/" + stubGroupId + "/artifacts/" + stubArtifactId + "/branches/" + stubBranchID + "/versions"

	t.Run("Success", func(t *testing.T) {
		server := setupMockServer(t, http.StatusNoContent, nil, expectedURL, http.MethodPut)
		defer server.Close()

		mockClient := &client.Client{BaseURL: server.URL, HTTPClient: server.Client()}
		api := apis.NewBranchAPI(mockClient)

		err := api.ReplaceVersionsInBranch(context.Background(), stubGroupId, stubArtifactId, stubBranchID, []string{stubVersionID, stubVersionID2})
		assert.NoError(t, err)
	})

	t.Run("Validation Errors", func(t *testing.T) {
		mockClient := &client.Client{BaseURL: "http://mock.server", HTTPClient: http.DefaultClient}
		api := apis.NewBranchAPI(mockClient)

		// Empty GroupID
		err := api.ReplaceVersionsInBranch(context.Background(), "", stubArtifactId, stubBranchID, []string{stubVersionID})
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "Group ID")

		// Empty ArtifactID
		err = api.ReplaceVersionsInBranch(context.Background(), stubGroupId, "", stubBranchID, []string{stubVersionID})
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "Artifact ID")

		// Empty BranchID
		err = api.ReplaceVersionsInBranch(context.Background(), stubGroupId, stubArtifactId, "", []string{stubVersionID})
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "Branch ID")

		// Empty Versions List
		err = api.ReplaceVersionsInBranch(context.Background(), stubGroupId, stubArtifactId, stubBranchID, []string{})
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "versions must not be empty")

		// Invalid Version Format
		err = api.ReplaceVersionsInBranch(context.Background(), stubGroupId, stubArtifactId, stubBranchID, []string{""})
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "Version")
	})

	t.Run("Not Found", func(t *testing.T) {
		mockErrorResponse := models.APIError{Status: http.StatusNotFound, Title: TitleNotFound}

		server := setupMockServer(t, http.StatusNotFound, mockErrorResponse, expectedURL, http.MethodPut)
		defer server.Close()

		mockClient := &client.Client{BaseURL: server.URL, HTTPClient: server.Client()}
		api := apis.NewBranchAPI(mockClient)

		err := api.ReplaceVersionsInBranch(context.Background(), stubGroupId, stubArtifactId, stubBranchID, []string{stubVersionID, stubVersionID2})
		assert.Error(t, err)

		assertAPIError(t, err, http.StatusNotFound, TitleNotFound)
	})

	t.Run("Internal Server Error", func(t *testing.T) {
		mockErrorResponse := models.APIError{Status: http.StatusInternalServerError, Title: TitleInternalServerError}

		server := setupMockServer(t, http.StatusInternalServerError, mockErrorResponse, expectedURL, http.MethodPut)
		defer server.Close()

		mockClient := &client.Client{BaseURL: server.URL, HTTPClient: server.Client()}
		api := apis.NewBranchAPI(mockClient)

		err := api.ReplaceVersionsInBranch(context.Background(), stubGroupId, stubArtifactId, stubBranchID, []string{stubVersionID, stubVersionID2})
		assert.Error(t, err)

		assertAPIError(t, err, http.StatusInternalServerError, TitleInternalServerError)
	})
}

func TestBranchAPI_AddVersionToBranch(t *testing.T) {
	expectedURL := "/groups/" + stubGroupId + "/artifacts/" + stubArtifactId + "/branches/" + stubBranchID + "/versions"

	t.Run("Success", func(t *testing.T) {
		server := setupMockServer(t, http.StatusNoContent, nil, expectedURL, http.MethodPost)
		defer server.Close()

		mockClient := &client.Client{BaseURL: server.URL, HTTPClient: server.Client()}
		api := apis.NewBranchAPI(mockClient)

		err := api.AddVersionToBranch(context.Background(), stubGroupId, stubArtifactId, stubBranchID, stubVersionID)
		assert.NoError(t, err)
	})

	t.Run("Validation Errors", func(t *testing.T) {
		mockClient := &client.Client{BaseURL: "http://mock.server", HTTPClient: http.DefaultClient}
		api := apis.NewBranchAPI(mockClient)

		err := api.AddVersionToBranch(context.Background(), "", stubArtifactId, stubBranchID, stubVersionID)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "Group ID")

		err = api.AddVersionToBranch(context.Background(), stubGroupId, "", stubBranchID, stubVersionID)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "Artifact ID")

		err = api.AddVersionToBranch(context.Background(), stubGroupId, stubArtifactId, "", stubVersionID)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "Branch ID")

		err = api.AddVersionToBranch(context.Background(), stubGroupId, stubArtifactId, stubBranchID, "")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "Version")
	})

	t.Run("Not Found", func(t *testing.T) {
		mockErrorResponse := models.APIError{Status: http.StatusNotFound, Title: TitleNotFound}

		server := setupMockServer(t, http.StatusNotFound, mockErrorResponse, expectedURL, http.MethodPost)
		defer server.Close()

		mockClient := &client.Client{BaseURL: server.URL, HTTPClient: server.Client()}
		api := apis.NewBranchAPI(mockClient)

		err := api.AddVersionToBranch(context.Background(), stubGroupId, stubArtifactId, stubBranchID, stubVersionID)
		assert.Error(t, err)

		assertAPIError(t, err, http.StatusNotFound, TitleNotFound)
	})

	t.Run("Conflict", func(t *testing.T) {
		mockErrorResponse := models.APIError{Status: http.StatusConflict, Title: TitleConflict}

		server := setupMockServer(t, http.StatusConflict, mockErrorResponse, expectedURL, http.MethodPost)
		defer server.Close()

		mockClient := &client.Client{BaseURL: server.URL, HTTPClient: server.Client()}
		api := apis.NewBranchAPI(mockClient)

		err := api.AddVersionToBranch(context.Background(), stubGroupId, stubArtifactId, stubBranchID, stubVersionID)
		assert.Error(t, err)

		assertAPIError(t, err, http.StatusConflict, TitleConflict)
	})

	t.Run("Internal Server Error", func(t *testing.T) {
		mockErrorResponse := models.APIError{Status: http.StatusInternalServerError, Title: TitleInternalServerError}

		server := setupMockServer(t, http.StatusInternalServerError, mockErrorResponse, expectedURL, http.MethodPost)
		defer server.Close()

		mockClient := &client.Client{BaseURL: server.URL, HTTPClient: server.Client()}
		api := apis.NewBranchAPI(mockClient)

		err := api.AddVersionToBranch(context.Background(), stubGroupId, stubArtifactId, stubBranchID, stubVersionID)
		assert.Error(t, err)

		assertAPIError(t, err, http.StatusInternalServerError, TitleInternalServerError)
	})
}

func setupMockServer(t *testing.T, statusCode int, response interface{}, expectedURL string, expectedMethod string) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if expectedURL != "" {
			assert.Equal(t, expectedURL, r.URL.Path, "request URL path should match expected")
		}

		if expectedMethod != "" {
			assert.Equal(t, expectedMethod, r.Method, "request method match expected")
		}

		w.WriteHeader(statusCode)
		if response != nil {
			err := json.NewEncoder(w).Encode(response)
			assert.NoError(t, err, "failed to encode response")
		}
	}))
}

func assertAPIError(t *testing.T, err error, expectedStatus int, expectedTitle string) {
	var apiErr *models.APIError
	ok := errors.As(err, &apiErr)
	assert.True(t, ok, "error should be of type *models.APIError")
	assert.Equal(t, expectedStatus, apiErr.Status)
	assert.Equal(t, expectedTitle, apiErr.Title)
}
