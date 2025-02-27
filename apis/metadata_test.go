package apis_test

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/mollie/go-apicurio-registry/apis"
	"github.com/mollie/go-apicurio-registry/client"
	"github.com/mollie/go-apicurio-registry/models"
	"github.com/stretchr/testify/assert"
)

const (
	versionExpression = "1.0.0"
)

func TestGetArtifactVersionMetadata(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mockMetadata := models.ArtifactVersionMetadata{
			BaseMetadata: models.BaseMetadata{
				GroupID:      "test-group",
				ArtifactID:   "artifact-1",
				Name:         "Test Artifact",
				Description:  "Test Description",
				ArtifactType: "JSON",
			},
			Version:   "1.0",
			GlobalID:  12345,
			ContentID: 67890,
		}

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Contains(t, r.URL.Path, "/groups/test-group/artifacts/artifact-1/versions/1.0")
			assert.Equal(t, http.MethodGet, r.Method)

			w.WriteHeader(http.StatusOK)
			err := json.NewEncoder(w).Encode(mockMetadata)
			assert.NoError(t, err)
		}))
		defer server.Close()

		mockClient := &client.Client{BaseURL: server.URL, HTTPClient: server.Client()}
		api := apis.NewMetadataAPI(mockClient)

		result, err := api.GetArtifactVersionMetadata(
			context.Background(),
			"test-group",
			"artifact-1",
			"1.0",
		)
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, "Test Artifact", result.Name)
		assert.Equal(t, "1.0", result.Version)
	})

	t.Run("InvalidInputs", func(t *testing.T) {
		ctx := context.Background()

		mockClient := &client.Client{BaseURL: "http://example.com", HTTPClient: http.DefaultClient}
		api := apis.NewMetadataAPI(mockClient)

		tests := []struct {
			groupID       string
			artifactID    string
			version       string
			expectedError string
		}{
			{"", "artifact-1", "1.0", "Group ID"},
			{"test-group", "", "1.0", "Artifact ID"},
			{"test-group", "artifact-1", "", "Version"},
		}

		for _, test := range tests {
			_, err := api.GetArtifactVersionMetadata(
				ctx,
				test.groupID,
				test.artifactID,
				test.version,
			)
			assert.Error(t, err)
			assert.Contains(t, err.Error(), test.expectedError)
		}
	})

	t.Run("Not Found", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusNotFound)
		}))
		defer server.Close()

		mockClient := &client.Client{BaseURL: server.URL, HTTPClient: server.Client()}
		api := apis.NewMetadataAPI(mockClient)

		result, err := api.GetArtifactVersionMetadata(
			context.Background(),
			"test-group",
			"artifact-1",
			"1.0",
		)
		assert.Error(t, err)
		assert.Nil(t, result)
	})
}

func TestUpdateArtifactVersionMetadata(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		server := setupMockServer(t, http.StatusNoContent, nil,
			"/groups/test-group/artifacts/artifact-1/versions/1.0.0", http.MethodPut)
		defer server.Close()

		mockClient := &client.Client{BaseURL: server.URL, HTTPClient: server.Client()}
		api := apis.NewMetadataAPI(mockClient)

		metadata := models.UpdateArtifactMetadataRequest{
			Name:        "Updated Artifact",
			Description: "Updated Description",
		}

		err := api.UpdateArtifactVersionMetadata(
			context.Background(),
			"test-group",
			"artifact-1",
			"1.0.0",
			metadata,
		)
		assert.NoError(t, err)
	})

	t.Run("Invalid Inputs", func(t *testing.T) {
		ctx := context.Background()

		mockClient := &client.Client{BaseURL: "http://example.com", HTTPClient: http.DefaultClient}
		api := apis.NewMetadataAPI(mockClient)

		tests := []struct {
			groupID       string
			artifactID    string
			version       string
			metadata      models.UpdateArtifactMetadataRequest
			expectedError string
		}{
			{
				"",
				"artifact-1",
				"1.0",
				models.UpdateArtifactMetadataRequest{Name: "Updated"},
				"Group ID",
			},
			{
				"test-group",
				"",
				"1.0",
				models.UpdateArtifactMetadataRequest{Name: "Updated"},
				"Artifact ID",
			},
			{
				"test-group",
				"artifact-1",
				"",
				models.UpdateArtifactMetadataRequest{Name: "Updated"},
				"Version",
			},
		}

		for _, test := range tests {
			err := api.UpdateArtifactVersionMetadata(
				ctx,
				test.groupID,
				test.artifactID,
				test.version,
				test.metadata,
			)
			assert.Error(t, err)
			assert.Contains(t, err.Error(), test.expectedError)
		}
	})

	t.Run("Invalid Input", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusBadRequest)
		}))
		defer server.Close()

		mockClient := &client.Client{BaseURL: server.URL, HTTPClient: server.Client()}
		api := apis.NewMetadataAPI(mockClient)

		metadata := models.UpdateArtifactMetadataRequest{
			Name: "",
		}

		err := api.UpdateArtifactVersionMetadata(
			context.Background(),
			"test-group",
			"artifact-1",
			"1.0",
			metadata,
		)
		assert.Error(t, err)
	})
}

func TestGetArtifactMetadata(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mockMetadata := models.ArtifactMetadata{
			BaseMetadata: models.BaseMetadata{
				GroupID:      "test-group",
				ArtifactID:   "artifact-1",
				Name:         "Test Artifact",
				Description:  "Test Description",
				ArtifactType: "JSON",
			},
			ModifiedBy: "user-1",
			ModifiedOn: "2024-12-09",
		}

		server := setupMockServer(t, http.StatusOK, mockMetadata,
			"/groups/test-group/artifacts/artifact-1", http.MethodGet)
		defer server.Close()

		mockClient := &client.Client{BaseURL: server.URL, HTTPClient: server.Client()}
		api := apis.NewMetadataAPI(mockClient)

		result, err := api.GetArtifactMetadata(context.Background(), "test-group", "artifact-1")
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, "Test Artifact", result.Name)
		assert.Equal(t, "user-1", result.ModifiedBy)
	})

	t.Run("Validation: Invalid Inputs", func(t *testing.T) {
		ctx := context.Background()
		mockClient := &client.Client{BaseURL: "http://example.com", HTTPClient: http.DefaultClient}
		api := apis.NewMetadataAPI(mockClient)

		tests := []struct {
			groupID       string
			artifactID    string
			expectedError string
		}{
			{"", "artifact-1", "Group ID"},
			{"test-group", "", "Artifact ID"},
		}

		for _, test := range tests {
			_, err := api.GetArtifactMetadata(ctx, test.groupID, test.artifactID)
			assert.Error(t, err)
			assert.Contains(t, err.Error(), test.expectedError)
		}
	})

	t.Run("Artifact Not Found", func(t *testing.T) {
		server := setupMockServer(t, http.StatusNotFound,
			models.APIError{Status: http.StatusNotFound, Title: TitleNotFound},
			"/groups/test-group/artifacts/artifact-1", http.MethodGet)
		defer server.Close()

		mockClient := &client.Client{BaseURL: server.URL, HTTPClient: server.Client()}
		api := apis.NewMetadataAPI(mockClient)

		result, err := api.GetArtifactMetadata(context.Background(), "test-group", "artifact-1")
		assert.Error(t, err)
		assert.Nil(t, result)
		assertAPIError(t, err, http.StatusNotFound, TitleNotFound)
	})
}

func TestUpdateArtifactMetadata(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		server := setupMockServer(t, http.StatusNoContent, nil,
			"/groups/test-group/artifacts/artifact-1", http.MethodPut)
		defer server.Close()

		mockClient := &client.Client{BaseURL: server.URL, HTTPClient: server.Client()}
		api := apis.NewMetadataAPI(mockClient)

		metadata := models.UpdateArtifactMetadataRequest{
			Name:        "Updated Artifact",
			Description: "Updated Description",
			Labels:      map[string]string{"env": "prod"},
		}

		err := api.UpdateArtifactMetadata(
			context.Background(),
			"test-group",
			"artifact-1",
			metadata,
		)
		assert.NoError(t, err)
	})

	t.Run("Invalid Inputs", func(t *testing.T) {
		ctx := context.Background()

		mockClient := &client.Client{BaseURL: "http://example.com", HTTPClient: http.DefaultClient}
		api := apis.NewMetadataAPI(mockClient)

		tests := []struct {
			groupID       string
			artifactID    string
			metadata      models.UpdateArtifactMetadataRequest
			expectedError string
		}{
			{"", "artifact-1", models.UpdateArtifactMetadataRequest{Name: "Updated"}, "Group ID"},
			{
				"test-group",
				"",
				models.UpdateArtifactMetadataRequest{Name: "Updated"},
				"Artifact ID",
			},
		}

		for _, test := range tests {
			err := api.UpdateArtifactMetadata(ctx, test.groupID, test.artifactID, test.metadata)
			assert.Error(t, err)
			assert.Contains(t, err.Error(), test.expectedError)
		}
	})

	t.Run("Invalid Input", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusBadRequest)
		}))
		defer server.Close()

		mockClient := &client.Client{BaseURL: server.URL, HTTPClient: server.Client()}
		api := apis.NewMetadataAPI(mockClient)

		metadata := models.UpdateArtifactMetadataRequest{}

		err := api.UpdateArtifactMetadata(
			context.Background(),
			"test-group",
			"artifact-1",
			metadata,
		)
		assert.Error(t, err)
	})
}

/***********************/
/***** Integration *****/
/***********************/

func setupMetadataAPIClient() *apis.MetadataAPI {
	apiClient := setupHTTPClient()
	return apis.NewMetadataAPI(apiClient)
}

func TestMetadataAPIIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	metadataAPI := setupMetadataAPIClient()

	ctx := context.Background()

	// Prepare test data
	artifactsAPI := apis.NewArtifactsAPI(metadataAPI.Client)

	// Clean up before and after tests
	t.Cleanup(func() { cleanup(t, artifactsAPI) })
	cleanup(t, artifactsAPI)

	artifact := models.CreateArtifactRequest{
		ArtifactType: models.Json,
		ArtifactID:   stubArtifactId,
		Name:         stubArtifactId,
		FirstVersion: models.CreateVersionRequest{
			Version: versionExpression,
			Content: models.CreateContentRequest{
				Content:     stubArtifactContent,
				ContentType: "application/json",
			},
		},
	}
	createParams := &models.CreateArtifactParams{
		IfExists: models.IfExistsFail,
	}
	_, err := artifactsAPI.CreateArtifact(ctx, stubGroupId, artifact, createParams)
	if err != nil {
		t.Fatalf("Failed to create artifact: %v", err)
	}

	// Test GetArtifactVersionMetadata
	t.Run("GetArtifactVersionMetadata", func(t *testing.T) {
		result, err := metadataAPI.GetArtifactVersionMetadata(
			ctx,
			stubGroupId,
			stubArtifactId,
			versionExpression,
		)
		assert.NoError(t, err)
		assert.NotNil(t, result)
		fmt.Println(result)
		assert.Equal(t, stubArtifactId, result.ArtifactID)
		assert.Equal(t, versionExpression, result.Version)
	})

	// Test UpdateArtifactVersionMetadata
	t.Run("UpdateArtifactVersionMetadata", func(t *testing.T) {
		updateRequest := models.UpdateArtifactMetadataRequest{
			Name:        "Updated Artifact Version Name",
			Description: "Updated Artifact Version Description",
		}

		err := metadataAPI.UpdateArtifactVersionMetadata(
			ctx,
			stubGroupId,
			stubArtifactId,
			versionExpression,
			updateRequest,
		)
		assert.NoError(t, err)

		// Verify the update
		updatedMetadata, err := metadataAPI.GetArtifactVersionMetadata(
			ctx,
			stubGroupId,
			stubArtifactId,
			versionExpression,
		)
		assert.NoError(t, err)
		assert.Equal(t, "Updated Artifact Version Name", updatedMetadata.Name)
		assert.Equal(t, "Updated Artifact Version Description", updatedMetadata.Description)
	})

	// Test GetArtifactMetadata
	t.Run("GetArtifactMetadata", func(t *testing.T) {
		result, err := metadataAPI.GetArtifactMetadata(ctx, stubGroupId, stubArtifactId)
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, stubArtifactId, result.ArtifactID)
	})

	// Test UpdateArtifactMetadata
	t.Run("UpdateArtifactMetadata", func(t *testing.T) {
		updateRequest := models.UpdateArtifactMetadataRequest{
			Name:        "Updated Artifact Name",
			Description: "Updated Artifact Description",
			Labels: map[string]string{
				"env": "production",
			},
		}

		err := metadataAPI.UpdateArtifactMetadata(ctx, stubGroupId, stubArtifactId, updateRequest)
		assert.NoError(t, err)

		// Verify the update
		updatedMetadata, err := metadataAPI.GetArtifactMetadata(ctx, stubGroupId, stubArtifactId)
		assert.NoError(t, err)
		assert.Equal(t, "Updated Artifact Name", updatedMetadata.Name)
		assert.Equal(t, "Updated Artifact Description", updatedMetadata.Description)
		assert.Equal(t, "production", updatedMetadata.Labels["env"])
	})
}
