package apis_test

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/mollie/go-apicurio-registry/apis"
	"github.com/mollie/go-apicurio-registry/models"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

const (
	TitleBadRequest          = "Bad request"
	TitleInternalServerError = "Internal server error"
	TitleNotFound            = "Not found"
	TitleConflict            = "Conflict"
	TitleMethodNotAllowed    = "Method Not allowed"

	DefaultBaseURL = "http://localhost:9080/apis/registry/v3"
)

var (
	stubArtifactContent    = `{"type": "record", "name": "Test", "fields": [{"name": "field1", "type": "string"}]}`
	stubArtifactId         = "test-artifact"
	stubGroupId            = "test-group"
	stubBranchID           = "test-branch"
	stubVersionID          = "1.0.0"
	stubVersionID2         = "2.0.0"
	stubDescription        = "description"
	stubUpdatedDescription = "updated-description"
	stubLabels             = map[string]string{"key": "value"}
	stubUpdatedLabels      = map[string]string{"key1": "value1"}
)

func generateArtifactForTest(ctx context.Context, artifactsAPI *apis.ArtifactsAPI) (string, error) {
	// Helper to generate unique artifact IDs
	generateArtifactID := func(prefix string) string {
		return fmt.Sprintf("%s-%d", prefix, time.Now().UnixNano())
	}

	newArtifactID := generateArtifactID("test-artifact")

	artifact := models.CreateArtifactRequest{
		ArtifactID:   newArtifactID,
		ArtifactType: models.Json,
		Name:         newArtifactID,
		FirstVersion: models.CreateVersionRequest{
			Version: "1.0.0",
			Content: models.CreateContentRequest{
				Content:     stubArtifactContent,
				ContentType: "application/json",
			},
			IsDraft: true,
		},
	}
	createParams := &models.CreateArtifactParams{
		IfExists: models.IfExistsFail,
	}
	_, err := artifactsAPI.CreateArtifact(ctx, stubGroupId, artifact, createParams)
	if err != nil {
		return "", err
	}
	return newArtifactID, nil
}

func generateGroupForTest(
	ctx context.Context,
	groupAPI *apis.GroupAPI,
) (*models.GroupInfo, string, error) {
	newGroupID := generateRandomName("test-group")
	resp, err := groupAPI.CreateGroup(ctx, newGroupID, stubDescription, stubLabels)
	return resp, newGroupID, err
}

func deleteGroupAfterTest(ctx context.Context, groupAPI *apis.GroupAPI, groupID string) error {
	return groupAPI.DeleteGroup(ctx, groupID)
}

func generateRandomName(prefix string) string {
	return fmt.Sprintf("%s-%d", prefix, time.Now().UnixNano())
}

func setupMockServer(
	t *testing.T,
	statusCode int,
	response interface{},
	expectedURL string,
	expectedMethod string,
) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if expectedURL != "" {
			assert.Contains(t, expectedURL, r.URL.Path, "request URL path should match expected")
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
