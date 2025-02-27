package apis_test

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/mollie/go-apicurio-registry/apis"
	"github.com/mollie/go-apicurio-registry/client"
	"github.com/mollie/go-apicurio-registry/models"
	"github.com/stretchr/testify/assert"
)

const (
	version        = "1.0.0"
	newVersion     = "1.1.0"
	commentID      = "test-comment"
	stubContent    = `{"type":"record","name":"TestRecord","fields":[{"name":"field1","type":"string"}]}`
	stubNewContent = `{"type":"record","name":"TestRecord","fields":[{"name":"field1","type":"string"},{"name":"field2","type":"string"}]}`
	stubReference  = `{"groupId":"test-group","artifactId":"ref-artifact","version":"1.0.0"}`
)

func TestVersionsAPI_DeleteArtifactVersion(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		server := setupMockServer(
			t,
			http.StatusNoContent,
			nil,
			"/groups/test-group/artifacts/test-artifact/versions/1.0.0",
			http.MethodDelete,
		)
		defer server.Close()

		mockClient := &client.Client{BaseURL: server.URL, HTTPClient: server.Client()}
		api := apis.NewVersionsAPI(mockClient)

		err := api.DeleteArtifactVersion(
			context.Background(),
			"test-group",
			"test-artifact",
			"1.0.0",
		)
		assert.NoError(t, err)
	})

	t.Run("Not Found", func(t *testing.T) {
		apiError := models.APIError{
			Status: http.StatusNotFound,
			Title:  "Artifact version not found",
		}
		server := setupMockServer(
			t,
			http.StatusNotFound,
			apiError,
			"/groups/test-group/artifacts/test-artifact/versions/1.0.0",
			http.MethodDelete,
		)
		defer server.Close()

		mockClient := &client.Client{BaseURL: server.URL, HTTPClient: server.Client()}
		api := apis.NewVersionsAPI(mockClient)

		err := api.DeleteArtifactVersion(
			context.Background(),
			"test-group",
			"test-artifact",
			"1.0.0",
		)
		assert.Error(t, err)
		assertAPIError(t, err, http.StatusNotFound, "Artifact version not found")
	})

	t.Run("Method Not Allowed", func(t *testing.T) {
		apiError := models.APIError{
			Status: http.StatusMethodNotAllowed,
			Title:  "Method Not Allowed",
		}
		server := setupMockServer(
			t,
			http.StatusMethodNotAllowed,
			apiError,
			"/groups/test-group/artifacts/test-artifact/versions/1.0.0",
			http.MethodDelete,
		)
		defer server.Close()

		mockClient := &client.Client{BaseURL: server.URL, HTTPClient: server.Client()}
		api := apis.NewVersionsAPI(mockClient)

		err := api.DeleteArtifactVersion(
			context.Background(),
			"test-group",
			"test-artifact",
			"1.0.0",
		)
		assert.Error(t, err)
		assertAPIError(t, err, http.StatusMethodNotAllowed, "Method Not Allowed")
	})

	t.Run("Internal Server Error", func(t *testing.T) {
		apiError := models.APIError{
			Status: http.StatusInternalServerError,
			Title:  "Internal Server Error",
		}
		server := setupMockServer(
			t,
			http.StatusInternalServerError,
			apiError,
			"/groups/test-group/artifacts/test-artifact/versions/1.0.0",
			http.MethodDelete,
		)
		defer server.Close()

		mockClient := &client.Client{BaseURL: server.URL, HTTPClient: server.Client()}
		api := apis.NewVersionsAPI(mockClient)

		err := api.DeleteArtifactVersion(
			context.Background(),
			"test-group",
			"test-artifact",
			"1.0.0",
		)
		assert.Error(t, err)
		assertAPIError(t, err, http.StatusInternalServerError, "Internal Server Error")
	})

	t.Run("Validation Error: Empty Group ID", func(t *testing.T) {
		mockClient := &client.Client{}
		api := apis.NewVersionsAPI(mockClient)

		err := api.DeleteArtifactVersion(context.Background(), "", "test-artifact", "1.0.0")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "Group ID")
	})

	t.Run("Validation Error: Empty Artifact ID", func(t *testing.T) {
		mockClient := &client.Client{}
		api := apis.NewVersionsAPI(mockClient)

		err := api.DeleteArtifactVersion(context.Background(), "test-group", "", "1.0.0")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "Artifact ID")
	})

	t.Run("Validation Error: Empty Version Expression", func(t *testing.T) {
		mockClient := &client.Client{}
		api := apis.NewVersionsAPI(mockClient)

		err := api.DeleteArtifactVersion(context.Background(), "test-group", "test-artifact", "")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "Version Expression")
	})
}

func TestVersionsAPI_GetArtifactVersionReferences(t *testing.T) {
	t.Run("Success with Parameters", func(t *testing.T) {
		mockResponse := []models.ArtifactReference{
			{GroupID: "test-group", ArtifactID: "artifact-1", Version: "1", Name: "Reference 1"},
		}

		server := setupMockServer(t, http.StatusOK, mockResponse,
			"/groups/test-group/artifacts/artifact-1/versions/1/references?refType=INBOUND",
			http.MethodGet,
		)
		defer server.Close()

		mockClient := &client.Client{BaseURL: server.URL, HTTPClient: server.Client()}
		api := apis.NewVersionsAPI(mockClient)

		params := &models.ArtifactVersionReferencesParams{RefType: "INBOUND"}
		result, err := api.GetArtifactVersionReferences(
			context.Background(),
			"test-group",
			"artifact-1",
			"1",
			params,
		)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, 1, len(result))
		assert.Equal(t, "Reference 1", (result)[0].Name)
	})

	t.Run("Success without Parameters", func(t *testing.T) {
		mockResponse := []models.ArtifactReference{
			{GroupID: "test-group", ArtifactID: "artifact-1", Version: "1", Name: "Reference 1"},
		}

		server := setupMockServer(t, http.StatusOK, mockResponse,
			"/groups/test-group/artifacts/artifact-1/versions/1/references",
			http.MethodGet,
		)
		defer server.Close()

		mockClient := &client.Client{BaseURL: server.URL, HTTPClient: server.Client()}
		api := apis.NewVersionsAPI(mockClient)

		result, err := api.GetArtifactVersionReferences(
			context.Background(),
			"test-group",
			"artifact-1",
			"1",
			nil,
		)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, 1, len(result))
		assert.Equal(t, "Reference 1", (result)[0].Name)
	})

	t.Run("Bad Request (400)", func(t *testing.T) {
		apiError := models.APIError{
			Status: http.StatusBadRequest,
			Title:  "Invalid refType parameter",
		}
		server := setupMockServer(t, http.StatusBadRequest, apiError,
			"/groups/test-group/artifacts/artifact-1/versions/1/references",
			http.MethodGet,
		)
		defer server.Close()

		mockClient := &client.Client{BaseURL: server.URL, HTTPClient: server.Client()}
		api := apis.NewVersionsAPI(mockClient)

		params := &models.ArtifactVersionReferencesParams{RefType: models.InBound}
		result, err := api.GetArtifactVersionReferences(
			context.Background(),
			"test-group",
			"artifact-1",
			"1",
			params,
		)

		assert.Error(t, err)
		assert.Nil(t, result)
		assertAPIError(t, err, http.StatusBadRequest, "Invalid refType parameter")
	})

	t.Run("Not Found (404)", func(t *testing.T) {
		apiError := models.APIError{Status: http.StatusNotFound, Title: "Artifact not found"}
		server := setupMockServer(t, http.StatusNotFound, apiError,
			"/groups/test-group/artifacts/artifact-1/versions/1/references",
			http.MethodGet,
		)
		defer server.Close()

		mockClient := &client.Client{BaseURL: server.URL, HTTPClient: server.Client()}
		api := apis.NewVersionsAPI(mockClient)

		result, err := api.GetArtifactVersionReferences(
			context.Background(),
			"test-group",
			"artifact-1",
			"1",
			nil,
		)

		assert.Error(t, err)
		assert.Nil(t, result)
		assertAPIError(t, err, http.StatusNotFound, "Artifact not found")
	})

	t.Run("Internal Server Error (500)", func(t *testing.T) {
		apiError := models.APIError{
			Status: http.StatusInternalServerError,
			Title:  "Internal server error",
		}
		server := setupMockServer(t, http.StatusInternalServerError, apiError,
			"/groups/test-group/artifacts/artifact-1/versions/1/references",
			http.MethodGet,
		)
		defer server.Close()

		mockClient := &client.Client{BaseURL: server.URL, HTTPClient: server.Client()}
		api := apis.NewVersionsAPI(mockClient)

		result, err := api.GetArtifactVersionReferences(
			context.Background(),
			"test-group",
			"artifact-1",
			"1",
			nil,
		)

		assert.Error(t, err)
		assert.Nil(t, result)
		assertAPIError(t, err, http.StatusInternalServerError, "Internal server error")
	})

	t.Run("Validation Error: Empty Group ID", func(t *testing.T) {
		mockClient := &client.Client{}
		api := apis.NewVersionsAPI(mockClient)

		_, err := api.GetArtifactVersionReferences(context.Background(), "", "artifact-1", "1", nil)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "Group ID")
	})

	t.Run("Validation Error: Empty Artifact ID", func(t *testing.T) {
		mockClient := &client.Client{}
		api := apis.NewVersionsAPI(mockClient)

		_, err := api.GetArtifactVersionReferences(context.Background(), "test-group", "", "1", nil)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "Artifact ID")
	})

	t.Run("Validation Error: Empty Version Expression", func(t *testing.T) {
		mockClient := &client.Client{}
		api := apis.NewVersionsAPI(mockClient)

		_, err := api.GetArtifactVersionReferences(
			context.Background(),
			"test-group",
			"artifact-1",
			"",
			nil,
		)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "Version Expression")
	})
}

func TestVersionsAPI_GetArtifactVersionComments(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mockResponse := []models.ArtifactComment{
			{
				CommentID: "12345",
				Value:     "This is a comment.",
				Owner:     "user1",
				CreatedOn: "2023-07-01T15:22:01Z",
			},
		}

		server := setupMockServer(t, http.StatusOK, mockResponse,
			"/groups/test-group/artifacts/artifact-1/versions/1/comments",
			http.MethodGet,
		)
		defer server.Close()

		mockClient := &client.Client{BaseURL: server.URL, HTTPClient: server.Client()}
		api := apis.NewVersionsAPI(mockClient)

		result, err := api.GetArtifactVersionComments(
			context.Background(),
			"test-group",
			"artifact-1",
			"1",
		)
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, 1, len(*result))
		assert.Equal(t, "This is a comment.", (*result)[0].Value)
		assert.Equal(t, "user1", (*result)[0].Owner)
	})

	t.Run("Bad Request (400)", func(t *testing.T) {
		apiError := models.APIError{
			Status: http.StatusBadRequest,
			Title:  "Invalid version expression",
		}

		server := setupMockServer(t, http.StatusBadRequest, apiError,
			"/groups/test-group/artifacts/artifact-1/versions/invalid/comments",
			http.MethodGet,
		)
		defer server.Close()

		mockClient := &client.Client{BaseURL: server.URL, HTTPClient: server.Client()}
		api := apis.NewVersionsAPI(mockClient)

		result, err := api.GetArtifactVersionComments(
			context.Background(),
			"test-group",
			"artifact-1",
			"invalid",
		)
		assert.Error(t, err)
		assert.Nil(t, result)
		assertAPIError(t, err, http.StatusBadRequest, "Invalid version expression")
	})

	t.Run("Not Found (404)", func(t *testing.T) {
		apiError := models.APIError{Status: http.StatusNotFound, Title: "Artifact not found"}

		server := setupMockServer(t, http.StatusNotFound, apiError,
			"/groups/test-group/artifacts/non-existent-artifact/versions/1/comments",
			http.MethodGet,
		)
		defer server.Close()

		mockClient := &client.Client{BaseURL: server.URL, HTTPClient: server.Client()}
		api := apis.NewVersionsAPI(mockClient)

		result, err := api.GetArtifactVersionComments(
			context.Background(),
			"test-group",
			"non-existent-artifact",
			"1",
		)
		assert.Error(t, err)
		assert.Nil(t, result)
		assertAPIError(t, err, http.StatusNotFound, "Artifact not found")
	})

	t.Run("Internal Server Error (500)", func(t *testing.T) {
		apiError := models.APIError{
			Status: http.StatusInternalServerError,
			Title:  "Internal server error",
		}

		server := setupMockServer(t, http.StatusInternalServerError, apiError,
			"/groups/test-group/artifacts/artifact-1/versions/1/comments",
			http.MethodGet,
		)
		defer server.Close()

		mockClient := &client.Client{BaseURL: server.URL, HTTPClient: server.Client()}
		api := apis.NewVersionsAPI(mockClient)

		result, err := api.GetArtifactVersionComments(
			context.Background(),
			"test-group",
			"artifact-1",
			"1",
		)
		assert.Error(t, err)
		assert.Nil(t, result)
		assertAPIError(t, err, http.StatusInternalServerError, "Internal server error")
	})

	t.Run("Validation Error: Empty Group ID", func(t *testing.T) {
		mockClient := &client.Client{}
		api := apis.NewVersionsAPI(mockClient)

		_, err := api.GetArtifactVersionComments(context.Background(), "", "artifact-1", "1")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "Group ID")
	})

	t.Run("Validation Error: Empty Artifact ID", func(t *testing.T) {
		mockClient := &client.Client{}
		api := apis.NewVersionsAPI(mockClient)

		_, err := api.GetArtifactVersionComments(context.Background(), "test-group", "", "1")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "Artifact ID")
	})

	t.Run("Validation Error: Empty Version Expression", func(t *testing.T) {
		mockClient := &client.Client{}
		api := apis.NewVersionsAPI(mockClient)

		_, err := api.GetArtifactVersionComments(
			context.Background(),
			"test-group",
			"artifact-1",
			"",
		)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "Version Expression")
	})
}

func TestVersionsAPI_AddArtifactVersionComment(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mockResponse := models.ArtifactComment{
			CommentID: "12345",
			Value:     "This is a new comment on an artifact version.",
			Owner:     "dwayne",
			CreatedOn: "2023-07-01T15:22:01Z",
		}

		server := setupMockServer(t, http.StatusOK, mockResponse,
			"/groups/my-group/artifacts/example-artifact/versions/v1/comments",
			http.MethodPost,
		)
		defer server.Close()

		mockClient := &client.Client{BaseURL: server.URL, HTTPClient: server.Client()}
		api := apis.NewVersionsAPI(mockClient)

		comment := "This is a new comment on an artifact version."
		result, err := api.AddArtifactVersionComment(
			context.Background(),
			"my-group",
			"example-artifact",
			"v1",
			comment,
		)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, mockResponse, *result)
	})

	t.Run("Bad Request (400)", func(t *testing.T) {
		apiError := models.APIError{Status: http.StatusBadRequest, Title: "Invalid input"}

		server := setupMockServer(t, http.StatusBadRequest, apiError,
			"/groups/my-group/artifacts/example-artifact/versions/v1/comments",
			http.MethodPost,
		)
		defer server.Close()

		mockClient := &client.Client{BaseURL: server.URL, HTTPClient: server.Client()}
		api := apis.NewVersionsAPI(mockClient)

		comment := "" // Invalid input
		result, err := api.AddArtifactVersionComment(
			context.Background(),
			"my-group",
			"example-artifact",
			"v1",
			comment,
		)

		assert.Error(t, err)
		assert.Nil(t, result)
		assertAPIError(t, err, http.StatusBadRequest, "Invalid input")
	})

	t.Run("Not Found (404)", func(t *testing.T) {
		apiError := models.APIError{Status: http.StatusNotFound, Title: "Artifact not found"}

		server := setupMockServer(t, http.StatusNotFound, apiError,
			"/groups/invalid-group/artifacts/example-artifact/versions/v1/comments",
			http.MethodPost,
		)
		defer server.Close()

		mockClient := &client.Client{BaseURL: server.URL, HTTPClient: server.Client()}
		api := apis.NewVersionsAPI(mockClient)

		comment := "This is a new comment"
		result, err := api.AddArtifactVersionComment(
			context.Background(),
			"invalid-group",
			"example-artifact",
			"v1",
			comment,
		)

		assert.Error(t, err)
		assert.Nil(t, result)
		assertAPIError(t, err, http.StatusNotFound, "Artifact not found")
	})

	t.Run("Internal Server Error (500)", func(t *testing.T) {
		apiError := models.APIError{
			Status: http.StatusInternalServerError,
			Title:  "Internal server error",
		}

		server := setupMockServer(t, http.StatusInternalServerError, apiError,
			"/groups/my-group/artifacts/example-artifact/versions/v1/comments",
			http.MethodPost,
		)
		defer server.Close()

		mockClient := &client.Client{BaseURL: server.URL, HTTPClient: server.Client()}
		api := apis.NewVersionsAPI(mockClient)

		comment := "This is a new comment"
		result, err := api.AddArtifactVersionComment(
			context.Background(),
			"my-group",
			"example-artifact",
			"v1",
			comment,
		)

		assert.Error(t, err)
		assert.Nil(t, result)
		assertAPIError(t, err, http.StatusInternalServerError, "Internal server error")
	})

	t.Run("Validation Error: Empty Group ID", func(t *testing.T) {
		mockClient := &client.Client{}
		api := apis.NewVersionsAPI(mockClient)

		comment := "Valid comment"
		result, err := api.AddArtifactVersionComment(
			context.Background(),
			"",
			"example-artifact",
			"v1",
			comment,
		)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "Group ID")
	})

	t.Run("Validation Error: Empty Artifact ID", func(t *testing.T) {
		mockClient := &client.Client{}
		api := apis.NewVersionsAPI(mockClient)

		comment := "Valid comment"
		result, err := api.AddArtifactVersionComment(
			context.Background(),
			"my-group",
			"",
			"v1",
			comment,
		)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "Artifact ID")
	})

	t.Run("Validation Error: Empty Version Expression", func(t *testing.T) {
		mockClient := &client.Client{}
		api := apis.NewVersionsAPI(mockClient)

		comment := "Valid comment"
		result, err := api.AddArtifactVersionComment(
			context.Background(),
			"my-group",
			"example-artifact",
			"",
			comment,
		)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "Version Expression")
	})
}

func TestVersionsAPI_UpdateArtifactVersionComment(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		server := setupMockServer(t, http.StatusNoContent, nil,
			"/groups/my-group/artifacts/example-artifact/versions/v1/comments/12345",
			http.MethodPut,
		)
		defer server.Close()

		mockClient := &client.Client{BaseURL: server.URL, HTTPClient: server.Client()}
		api := apis.NewVersionsAPI(mockClient)

		err := api.UpdateArtifactVersionComment(
			context.Background(),
			"my-group",
			"example-artifact",
			"v1",
			"12345",
			"Updated comment value",
		)

		assert.NoError(t, err)
	})

	t.Run("Bad Request (400)", func(t *testing.T) {
		apiError := models.APIError{Status: http.StatusBadRequest, Title: "Invalid input"}

		server := setupMockServer(t, http.StatusBadRequest, apiError,
			"/groups/my-group/artifacts/example-artifact/versions/v1/comments/12345",
			http.MethodPut,
		)
		defer server.Close()

		mockClient := &client.Client{BaseURL: server.URL, HTTPClient: server.Client()}
		api := apis.NewVersionsAPI(mockClient)

		err := api.UpdateArtifactVersionComment(
			context.Background(),
			"my-group",
			"example-artifact",
			"v1",
			"12345",
			"",
		)

		assert.Error(t, err)
		assertAPIError(t, err, http.StatusBadRequest, "Invalid input")
	})

	t.Run("Not Found (404)", func(t *testing.T) {
		apiError := models.APIError{Status: http.StatusNotFound, Title: "Comment not found"}

		server := setupMockServer(t, http.StatusNotFound, apiError,
			"/groups/my-group/artifacts/example-artifact/versions/v1/comments/invalid-id",
			http.MethodPut,
		)
		defer server.Close()

		mockClient := &client.Client{BaseURL: server.URL, HTTPClient: server.Client()}
		api := apis.NewVersionsAPI(mockClient)

		err := api.UpdateArtifactVersionComment(
			context.Background(),
			"my-group",
			"example-artifact",
			"v1",
			"invalid-id",
			"Updated comment",
		)

		assert.Error(t, err)
		assertAPIError(t, err, http.StatusNotFound, "Comment not found")
	})

	t.Run("Internal Server Error (500)", func(t *testing.T) {
		apiError := models.APIError{
			Status: http.StatusInternalServerError,
			Title:  "Internal server error",
		}

		server := setupMockServer(t, http.StatusInternalServerError, apiError,
			"/groups/my-group/artifacts/example-artifact/versions/v1/comments/12345",
			http.MethodPut,
		)
		defer server.Close()

		mockClient := &client.Client{BaseURL: server.URL, HTTPClient: server.Client()}
		api := apis.NewVersionsAPI(mockClient)

		err := api.UpdateArtifactVersionComment(
			context.Background(),
			"my-group",
			"example-artifact",
			"v1",
			"12345",
			"Updated comment",
		)

		assert.Error(t, err)
		assertAPIError(t, err, http.StatusInternalServerError, "Internal server error")
	})

	t.Run("Validation Error: Empty Group ID", func(t *testing.T) {
		mockClient := &client.Client{}
		api := apis.NewVersionsAPI(mockClient)

		err := api.UpdateArtifactVersionComment(
			context.Background(),
			"",
			"example-artifact",
			"v1",
			"12345",
			"Updated comment",
		)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "Group ID")
	})

	t.Run("Validation Error: Empty Artifact ID", func(t *testing.T) {
		mockClient := &client.Client{}
		api := apis.NewVersionsAPI(mockClient)

		err := api.UpdateArtifactVersionComment(
			context.Background(),
			"my-group",
			"",
			"v1",
			"12345",
			"Updated comment",
		)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "Artifact ID")
	})

	t.Run("Validation Error: Empty Version Expression", func(t *testing.T) {
		mockClient := &client.Client{}
		api := apis.NewVersionsAPI(mockClient)

		err := api.UpdateArtifactVersionComment(
			context.Background(),
			"my-group",
			"example-artifact",
			"",
			"12345",
			"Updated comment",
		)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "Version Expression")
	})
}

func TestVersionsAPI_DeleteArtifactVersionComment(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		server := setupMockServer(t, http.StatusNoContent, nil,
			"/groups/my-group/artifacts/example-artifact/versions/v1/comments/12345",
			http.MethodDelete,
		)
		defer server.Close()

		mockClient := &client.Client{BaseURL: server.URL, HTTPClient: server.Client()}
		api := apis.NewVersionsAPI(mockClient)

		err := api.DeleteArtifactVersionComment(
			context.Background(),
			"my-group",
			"example-artifact",
			"v1",
			"12345",
		)
		assert.NoError(t, err)
	})

	t.Run("Bad Request (400)", func(t *testing.T) {
		apiError := models.APIError{Status: http.StatusBadRequest, Title: "Invalid input"}

		server := setupMockServer(t, http.StatusBadRequest, apiError,
			"/groups/my-group/artifacts/example-artifact/versions/v1/comments/12345",
			http.MethodDelete,
		)
		defer server.Close()

		mockClient := &client.Client{BaseURL: server.URL, HTTPClient: server.Client()}
		api := apis.NewVersionsAPI(mockClient)

		err := api.DeleteArtifactVersionComment(
			context.Background(),
			"my-group",
			"example-artifact",
			"v1",
			"12345",
		)
		assert.Error(t, err)
		assertAPIError(t, err, http.StatusBadRequest, "Invalid input")
	})

	t.Run("Not Found (404)", func(t *testing.T) {
		apiError := models.APIError{Status: http.StatusNotFound, Title: "Comment not found"}

		server := setupMockServer(t, http.StatusNotFound, apiError,
			"/groups/my-group/artifacts/example-artifact/versions/v1/comments/invalid-id",
			http.MethodDelete,
		)
		defer server.Close()

		mockClient := &client.Client{BaseURL: server.URL, HTTPClient: server.Client()}
		api := apis.NewVersionsAPI(mockClient)

		err := api.DeleteArtifactVersionComment(
			context.Background(),
			"my-group",
			"example-artifact",
			"v1",
			"invalid-id",
		)
		assert.Error(t, err)
		assertAPIError(t, err, http.StatusNotFound, "Comment not found")
	})

	t.Run("Internal Server Error (500)", func(t *testing.T) {
		apiError := models.APIError{
			Status: http.StatusInternalServerError,
			Title:  "Internal server error",
		}

		server := setupMockServer(t, http.StatusInternalServerError, apiError,
			"/groups/my-group/artifacts/example-artifact/versions/v1/comments/12345",
			http.MethodDelete,
		)
		defer server.Close()

		mockClient := &client.Client{BaseURL: server.URL, HTTPClient: server.Client()}
		api := apis.NewVersionsAPI(mockClient)

		err := api.DeleteArtifactVersionComment(
			context.Background(),
			"my-group",
			"example-artifact",
			"v1",
			"12345",
		)
		assert.Error(t, err)
		assertAPIError(t, err, http.StatusInternalServerError, "Internal server error")
	})

	// Validation Tests
	t.Run("Validation Error: Empty Group ID", func(t *testing.T) {
		mockClient := &client.Client{}
		api := apis.NewVersionsAPI(mockClient)

		err := api.DeleteArtifactVersionComment(
			context.Background(),
			"",
			"example-artifact",
			"v1",
			"12345",
		)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "Group ID")
	})

	t.Run("Validation Error: Empty Artifact ID", func(t *testing.T) {
		mockClient := &client.Client{}
		api := apis.NewVersionsAPI(mockClient)

		err := api.DeleteArtifactVersionComment(context.Background(), "my-group", "", "v1", "12345")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "Artifact ID")
	})

	t.Run("Validation Error: Empty Version Expression", func(t *testing.T) {
		mockClient := &client.Client{}
		api := apis.NewVersionsAPI(mockClient)

		err := api.DeleteArtifactVersionComment(
			context.Background(),
			"my-group",
			"example-artifact",
			"",
			"12345",
		)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "Version Expression")
	})

	t.Run("Validation Error: Empty Comment ID", func(t *testing.T) {
		mockClient := &client.Client{}
		api := apis.NewVersionsAPI(mockClient)

		err := api.DeleteArtifactVersionComment(
			context.Background(),
			"my-group",
			"example-artifact",
			"v1",
			"",
		)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "Comment ID")
	})
}

func TestVersionsAPI_ListArtifactVersions(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mockResponse := models.ArtifactVersionListResponse{
			Count: 2,
			Versions: []models.ArtifactVersion{
				{
					Version:      "2.0.0",
					CreatedOn:    "2024-12-10T08:56:40Z",
					ArtifactType: models.Json,
					GlobalID:     47,
					State:        models.StateEnabled,
					ContentID:    47,
					ArtifactID:   "example-artifact",
					GroupID:      "my-group",
					ModifiedOn:   "2024-12-10T08:56:40Z",
				},
				{
					Version:      "1.0.0",
					CreatedOn:    "2024-12-10T08:56:17Z",
					ArtifactType: models.Json,
					GlobalID:     46,
					State:        models.StateEnabled,
					ContentID:    46,
					ArtifactID:   "example-artifact",
					GroupID:      "my-group",
					ModifiedOn:   "2024-12-10T08:56:17Z",
				},
			},
		}

		server := setupMockServer(t, http.StatusOK, mockResponse,
			"/groups/my-group/artifacts/example-artifact/versions", http.MethodGet)
		defer server.Close()

		mockClient := &client.Client{BaseURL: server.URL, HTTPClient: server.Client()}
		api := apis.NewVersionsAPI(mockClient)

		versions, err := api.ListArtifactVersions(
			context.Background(),
			"my-group",
			"example-artifact",
			nil,
		)
		assert.NoError(t, err)
		assert.NotNil(t, versions)
		assert.Equal(t, 2, len(versions))
		assert.Equal(t, "2.0.0", (versions)[0].Version)
		assert.Equal(t, "1.0.0", (versions)[1].Version)
	})

	t.Run("Invalid Params", func(t *testing.T) {
		params := &models.ListArtifactsVersionsParams{
			Limit:   -1,                              // Invalid: Limit cannot be negative
			OrderBy: models.VersionSortBy("invalid"), // Invalid: Unsupported OrderBy value
		}

		mockClient := &client.Client{BaseURL: "http://example.com", HTTPClient: &http.Client{}}
		api := apis.NewVersionsAPI(mockClient)

		versions, err := api.ListArtifactVersions(
			context.Background(),
			"my-group",
			"example-artifact",
			params,
		)
		assert.Error(t, err)
		assert.Nil(t, versions)

		var validationErrs validator.ValidationErrors
		if errors.As(err, &validationErrs) {
			foundErrors := map[string]string{}
			for _, fieldErr := range validationErrs {
				foundErrors[fieldErr.StructField()] = fieldErr.Tag()
			}

			assert.Equal(t, "gte", foundErrors["Limit"], "Limit must be greater than or equal to 0")
			assert.Equal(t, "oneof", foundErrors["OrderBy"], "OrderBy must match supported values")
		} else {
			t.Errorf("Expected validation errors but got: %v", err)
		}
	})

	t.Run("Not Found (404)", func(t *testing.T) {
		apiError := models.APIError{Status: http.StatusNotFound, Title: "Artifact not found"}

		server := setupMockServer(t, http.StatusNotFound, apiError,
			"/groups/my-group/artifacts/example-artifact/versions", http.MethodGet)
		defer server.Close()

		mockClient := &client.Client{BaseURL: server.URL, HTTPClient: server.Client()}
		api := apis.NewVersionsAPI(mockClient)

		versions, err := api.ListArtifactVersions(
			context.Background(),
			"my-group",
			"example-artifact",
			nil,
		)
		assert.Error(t, err)
		assert.Nil(t, versions)
		assertAPIError(t, err, http.StatusNotFound, "Artifact not found")
	})

	t.Run("Internal Server Error (500)", func(t *testing.T) {
		apiError := models.APIError{
			Status: http.StatusInternalServerError,
			Title:  "Internal server error",
		}

		server := setupMockServer(t, http.StatusInternalServerError, apiError,
			"/groups/my-group/artifacts/example-artifact/versions", http.MethodGet)
		defer server.Close()

		mockClient := &client.Client{BaseURL: server.URL, HTTPClient: server.Client()}
		api := apis.NewVersionsAPI(mockClient)

		versions, err := api.ListArtifactVersions(
			context.Background(),
			"my-group",
			"example-artifact",
			nil,
		)
		assert.Error(t, err)
		assert.Nil(t, versions)
		assertAPIError(t, err, http.StatusInternalServerError, "Internal server error")
	})

	// Validation Scenarios
	t.Run("Validation Error: Empty Group ID", func(t *testing.T) {
		mockClient := &client.Client{}
		api := apis.NewVersionsAPI(mockClient)

		versions, err := api.ListArtifactVersions(context.Background(), "", "example-artifact", nil)
		assert.Error(t, err)
		assert.Nil(t, versions)
		assert.Contains(t, err.Error(), "Group ID")
	})

	t.Run("Validation Error: Empty Artifact ID", func(t *testing.T) {
		mockClient := &client.Client{}
		api := apis.NewVersionsAPI(mockClient)

		versions, err := api.ListArtifactVersions(context.Background(), "my-group", "", nil)
		assert.Error(t, err)
		assert.Nil(t, versions)
		assert.Contains(t, err.Error(), "Artifact ID")
	})
}

func TestVersionsAPI_CreateArtifactVersion(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mockResponse := models.ArtifactVersionDetailed{
			ArtifactVersion: models.ArtifactVersion{
				Version:      "1.0.0",
				CreatedOn:    "2024-12-10T08:56:40Z",
				ArtifactType: models.Json,
				GlobalID:     40,
				State:        models.StateEnabled,
				ContentID:    10,
				ArtifactID:   "example-artifact",
				GroupID:      "my-group",
				ModifiedOn:   "2024-12-10T08:56:40Z",
			},
			Name:        "Artifact Name",
			Description: "Artifact Description",
			Labels: map[string]string{
				"key1": "value1",
				"key2": "value2",
			},
		}

		server := setupMockServer(t, http.StatusOK, mockResponse,
			"/groups/my-group/artifacts/example-artifact/versions", http.MethodPost)
		defer server.Close()

		mockClient := &client.Client{BaseURL: server.URL, HTTPClient: server.Client()}
		api := apis.NewVersionsAPI(mockClient)

		createRequest := &models.CreateVersionRequest{
			Version: "1.0.0",
			Content: models.CreateContentRequest{
				Content:     `{"a": "1"}`,
				ContentType: "application/json",
			},
			Name:        "Artifact Name",
			Description: "Artifact Description",
			Labels:      map[string]string{"key1": "value1", "key2": "value2"},
			IsDraft:     false,
		}

		result, err := api.CreateArtifactVersion(
			context.Background(),
			"my-group",
			"example-artifact",
			createRequest,
			false,
		)
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, "1.0.0", result.Version)
		assert.Equal(t, "Artifact Name", result.Name)
		assert.Equal(t, "Artifact Description", result.Description)
		assert.Equal(t, 2, len(result.Labels))
	})

	t.Run("BadRequest", func(t *testing.T) {
		apiError := models.APIError{Status: http.StatusBadRequest, Title: "Invalid input"}

		server := setupMockServer(t, http.StatusBadRequest, apiError,
			"/groups/my-group/artifacts/example-artifact/versions", http.MethodPost)
		defer server.Close()

		mockClient := &client.Client{BaseURL: server.URL, HTTPClient: server.Client()}
		api := apis.NewVersionsAPI(mockClient)

		createRequest := &models.CreateVersionRequest{}
		result, err := api.CreateArtifactVersion(
			context.Background(),
			"my-group",
			"example-artifact",
			createRequest,
			false,
		)

		assert.Error(t, err)
		assert.Nil(t, result)
		assertAPIError(t, err, http.StatusBadRequest, "Invalid input")
	})

	t.Run("NotFound", func(t *testing.T) {
		apiError := models.APIError{Status: http.StatusNotFound, Title: "Artifact not found"}

		server := setupMockServer(t, http.StatusNotFound, apiError,
			"/groups/my-group/artifacts/example-artifact/versions", http.MethodPost)
		defer server.Close()

		mockClient := &client.Client{BaseURL: server.URL, HTTPClient: server.Client()}
		api := apis.NewVersionsAPI(mockClient)

		createRequest := &models.CreateVersionRequest{
			Version: "1.0.0",
			Content: models.CreateContentRequest{
				Content:     `{"a": "1"}`,
				ContentType: "application/json",
			},
		}

		result, err := api.CreateArtifactVersion(
			context.Background(),
			"my-group",
			"example-artifact",
			createRequest,
			false,
		)
		assert.Error(t, err)
		assert.Nil(t, result)
		assertAPIError(t, err, http.StatusNotFound, "Artifact not found")
	})

	t.Run("Conflict", func(t *testing.T) {
		apiError := models.APIError{Status: http.StatusConflict, Title: "Conflict"}

		server := setupMockServer(t, http.StatusConflict, apiError,
			"/groups/my-group/artifacts/example-artifact/versions", http.MethodPost)
		defer server.Close()

		mockClient := &client.Client{BaseURL: server.URL, HTTPClient: server.Client()}
		api := apis.NewVersionsAPI(mockClient)

		createRequest := &models.CreateVersionRequest{
			Version: "1.0.0",
			Content: models.CreateContentRequest{
				Content:     `{"a": "1"}`,
				ContentType: "application/json",
			},
		}

		result, err := api.CreateArtifactVersion(
			context.Background(),
			"my-group",
			"example-artifact",
			createRequest,
			false,
		)
		assert.Error(t, err)
		assert.Nil(t, result)
		assertAPIError(t, err, http.StatusConflict, "Conflict")
	})

	t.Run("InternalServerError", func(t *testing.T) {
		apiError := models.APIError{
			Status: http.StatusInternalServerError,
			Title:  "Internal server error",
		}

		server := setupMockServer(t, http.StatusInternalServerError, apiError,
			"/groups/my-group/artifacts/example-artifact/versions", http.MethodPost)
		defer server.Close()

		mockClient := &client.Client{BaseURL: server.URL, HTTPClient: server.Client()}
		api := apis.NewVersionsAPI(mockClient)

		createRequest := &models.CreateVersionRequest{
			Version: "1.0.0",
			Content: models.CreateContentRequest{
				Content:     `{"a": "1"}`,
				ContentType: "application/json",
			},
		}

		result, err := api.CreateArtifactVersion(
			context.Background(),
			"my-group",
			"example-artifact",
			createRequest,
			false,
		)
		assert.Error(t, err)
		assert.Nil(t, result)
		assertAPIError(t, err, http.StatusInternalServerError, "Internal server error")
	})

	// Validation Tests
	t.Run("Validation Error: Empty Group ID", func(t *testing.T) {
		mockClient := &client.Client{}
		api := apis.NewVersionsAPI(mockClient)

		result, err := api.CreateArtifactVersion(
			context.Background(),
			"",
			"example-artifact",
			nil,
			false,
		)
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "Group ID")
	})

	t.Run("Validation Error: Empty Artifact ID", func(t *testing.T) {
		mockClient := &client.Client{}
		api := apis.NewVersionsAPI(mockClient)

		result, err := api.CreateArtifactVersion(context.Background(), "my-group", "", nil, false)
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "Artifact ID")
	})
}

func TestVersionsAPI_GetArtifactVersionContent(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mockResponse := `{"a": "1"}`

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(
				t,
				"/groups/my-group/artifacts/example-artifact/versions/1.0.0/content",
				r.URL.Path,
			)
			assert.Equal(t, http.MethodGet, r.Method)
			// Write the response
			w.Header().Set("X-Registry-ArtifactType", string(models.Json))
			w.WriteHeader(http.StatusOK)
			_, err := w.Write([]byte(mockResponse))
			assert.NoError(t, err)

		}))
		defer server.Close()

		mockClient := &client.Client{BaseURL: server.URL, HTTPClient: server.Client()}
		api := apis.NewVersionsAPI(mockClient)

		content, err := api.GetArtifactVersionContent(
			context.Background(),
			"my-group",
			"example-artifact",
			"1.0.0",
			nil,
		)
		assert.NoError(t, err)
		assert.NotEmpty(t, content)
		assert.Equal(t, `{"a": "1"}`, content.Content)
	})

	t.Run("BadRequest", func(t *testing.T) {
		apiError := models.APIError{Status: http.StatusBadRequest, Title: "Invalid request"}
		expectedURL := "/groups/my-group/artifacts/example-artifact/versions/1.0.0/content"

		server := setupMockServer(t, http.StatusBadRequest, apiError, expectedURL, http.MethodGet)
		defer server.Close()

		mockClient := &client.Client{BaseURL: server.URL, HTTPClient: server.Client()}
		api := apis.NewVersionsAPI(mockClient)

		result, err := api.GetArtifactVersionContent(
			context.Background(),
			"my-group",
			"example-artifact",
			"1.0.0",
			nil,
		)

		assert.Error(t, err)
		assert.Nil(t, result)
		assertAPIError(t, err, http.StatusBadRequest, "Invalid request")
	})

	t.Run("NotFound", func(t *testing.T) {
		apiError := models.APIError{Status: http.StatusNotFound, Title: "Artifact not found"}
		expectedURL := "/groups/my-group/artifacts/example-artifact/versions/1.0.0/content"

		server := setupMockServer(t, http.StatusNotFound, apiError, expectedURL, http.MethodGet)
		defer server.Close()

		mockClient := &client.Client{BaseURL: server.URL, HTTPClient: server.Client()}
		api := apis.NewVersionsAPI(mockClient)

		result, err := api.GetArtifactVersionContent(
			context.Background(),
			"my-group",
			"example-artifact",
			"1.0.0",
			nil,
		)

		assert.Error(t, err)
		assert.Nil(t, result)
		assertAPIError(t, err, http.StatusNotFound, "Artifact not found")
	})

	t.Run("InternalServerError", func(t *testing.T) {
		apiError := models.APIError{
			Status: http.StatusInternalServerError,
			Title:  "Internal server error",
		}
		expectedURL := "/groups/my-group/artifacts/example-artifact/versions/1.0.0/content"

		server := setupMockServer(
			t,
			http.StatusInternalServerError,
			apiError,
			expectedURL,
			http.MethodGet,
		)
		defer server.Close()

		mockClient := &client.Client{BaseURL: server.URL, HTTPClient: server.Client()}
		api := apis.NewVersionsAPI(mockClient)

		result, err := api.GetArtifactVersionContent(
			context.Background(),
			"my-group",
			"example-artifact",
			"1.0.0",
			nil,
		)

		assert.Error(t, err)
		assert.Nil(t, result)
		assertAPIError(t, err, http.StatusInternalServerError, "Internal server error")
	})

	// Validation Tests
	t.Run("Validation Error: Empty Group ID", func(t *testing.T) {
		mockClient := &client.Client{}
		api := apis.NewVersionsAPI(mockClient)

		result, err := api.GetArtifactVersionContent(
			context.Background(),
			"",
			"example-artifact",
			"1.0.0",
			nil,
		)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "Group ID")
	})

	t.Run("Validation Error: Empty Artifact ID", func(t *testing.T) {
		mockClient := &client.Client{}
		api := apis.NewVersionsAPI(mockClient)

		result, err := api.GetArtifactVersionContent(
			context.Background(),
			"my-group",
			"",
			"1.0.0",
			nil,
		)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "Artifact ID")
	})

	t.Run("Validation Error: Empty Version Expression", func(t *testing.T) {
		mockClient := &client.Client{}
		api := apis.NewVersionsAPI(mockClient)

		result, err := api.GetArtifactVersionContent(
			context.Background(),
			"my-group",
			"example-artifact",
			"",
			nil,
		)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "Version Expression")
	})
}

func TestVersionsAPI_UpdateArtifactVersionContent(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		expectedURL := "/groups/my-group/artifacts/example-artifact/versions/1.0.0/content"

		server := setupMockServer(t, http.StatusNoContent, nil, expectedURL, http.MethodPut)
		defer server.Close()

		mockClient := &client.Client{BaseURL: server.URL, HTTPClient: server.Client()}
		api := apis.NewVersionsAPI(mockClient)

		content := &models.CreateContentRequest{
			Content:     `{"key": "value"}`,
			ContentType: "application/json",
		}

		err := api.UpdateArtifactVersionContent(
			context.Background(),
			"my-group",
			"example-artifact",
			"1.0.0",
			content,
		)

		assert.NoError(t, err)
	})

	t.Run("BadRequest", func(t *testing.T) {
		apiError := models.APIError{Status: http.StatusBadRequest, Title: "Invalid input"}
		expectedURL := "/groups/my-group/artifacts/example-artifact/versions/1.0.0/content"

		server := setupMockServer(t, http.StatusBadRequest, apiError, expectedURL, http.MethodPut)
		defer server.Close()

		mockClient := &client.Client{BaseURL: server.URL, HTTPClient: server.Client()}
		api := apis.NewVersionsAPI(mockClient)

		content := &models.CreateContentRequest{
			Content:     `{"key": "value"}`,
			ContentType: "application/json",
		}

		err := api.UpdateArtifactVersionContent(
			context.Background(),
			"my-group",
			"example-artifact",
			"1.0.0",
			content,
		)

		assert.Error(t, err)
		assertAPIError(t, err, http.StatusBadRequest, "Invalid input")
	})

	t.Run("NotFound", func(t *testing.T) {
		apiError := models.APIError{Status: http.StatusNotFound, Title: "Artifact not found"}
		expectedURL := "/groups/my-group/artifacts/example-artifact/versions/1.0.0/content"

		server := setupMockServer(t, http.StatusNotFound, apiError, expectedURL, http.MethodPut)
		defer server.Close()

		mockClient := &client.Client{BaseURL: server.URL, HTTPClient: server.Client()}
		api := apis.NewVersionsAPI(mockClient)

		content := &models.CreateContentRequest{
			Content:     `{"key": "value"}`,
			ContentType: "application/json",
		}

		err := api.UpdateArtifactVersionContent(
			context.Background(),
			"my-group",
			"example-artifact",
			"1.0.0",
			content,
		)

		assert.Error(t, err)
		assertAPIError(t, err, http.StatusNotFound, "Artifact not found")
	})

	t.Run("Conflict", func(t *testing.T) {
		apiError := models.APIError{Status: http.StatusConflict, Title: "Conflict"}
		expectedURL := "/groups/my-group/artifacts/example-artifact/versions/1.0.0/content"

		server := setupMockServer(t, http.StatusConflict, apiError, expectedURL, http.MethodPut)
		defer server.Close()

		mockClient := &client.Client{BaseURL: server.URL, HTTPClient: server.Client()}
		api := apis.NewVersionsAPI(mockClient)

		content := &models.CreateContentRequest{
			Content:     `{"key": "value"}`,
			ContentType: "application/json",
		}

		err := api.UpdateArtifactVersionContent(
			context.Background(),
			"my-group",
			"example-artifact",
			"1.0.0",
			content,
		)

		assert.Error(t, err)
		assertAPIError(t, err, http.StatusConflict, "Conflict")
	})

	t.Run("InternalServerError", func(t *testing.T) {
		apiError := models.APIError{
			Status: http.StatusInternalServerError,
			Title:  "Internal server error",
		}
		expectedURL := "/groups/my-group/artifacts/example-artifact/versions/1.0.0/content"

		server := setupMockServer(
			t,
			http.StatusInternalServerError,
			apiError,
			expectedURL,
			http.MethodPut,
		)
		defer server.Close()

		mockClient := &client.Client{BaseURL: server.URL, HTTPClient: server.Client()}
		api := apis.NewVersionsAPI(mockClient)

		content := &models.CreateContentRequest{
			Content:     `{"key": "value"}`,
			ContentType: "application/json",
		}

		err := api.UpdateArtifactVersionContent(
			context.Background(),
			"my-group",
			"example-artifact",
			"1.0.0",
			content,
		)

		assert.Error(t, err)
		assertAPIError(t, err, http.StatusInternalServerError, "Internal server error")
	})

	// Validation Tests
	t.Run("Validation Error: Empty Group ID", func(t *testing.T) {
		mockClient := &client.Client{}
		api := apis.NewVersionsAPI(mockClient)

		content := &models.CreateContentRequest{
			Content:     `{"key": "value"}`,
			ContentType: "application/json",
		}

		err := api.UpdateArtifactVersionContent(
			context.Background(),
			"",
			"example-artifact",
			"1.0.0",
			content,
		)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "Group ID")
	})

	t.Run("Validation Error: Empty Artifact ID", func(t *testing.T) {
		mockClient := &client.Client{}
		api := apis.NewVersionsAPI(mockClient)

		content := &models.CreateContentRequest{
			Content:     `{"key": "value"}`,
			ContentType: "application/json",
		}

		err := api.UpdateArtifactVersionContent(
			context.Background(),
			"my-group",
			"",
			"1.0.0",
			content,
		)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "Artifact ID")
	})

	t.Run("Validation Error: Empty Version Expression", func(t *testing.T) {
		mockClient := &client.Client{}
		api := apis.NewVersionsAPI(mockClient)

		content := &models.CreateContentRequest{
			Content:     `{"key": "value"}`,
			ContentType: "application/json",
		}

		err := api.UpdateArtifactVersionContent(
			context.Background(),
			"my-group",
			"example-artifact",
			"",
			content,
		)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "Version Expression")
	})
}

func TestVersionsAPI_SearchForArtifactVersions(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mockResponse := models.ArtifactVersionListResponse{
			Count: 2,
			Versions: []models.ArtifactVersion{
				{
					CreatedOn:    "2024-12-10T08:56:40Z",
					ArtifactType: models.Json,
					State:        models.StateEnabled,
					GlobalID:     47,
					Version:      "2.0.0",
					ContentID:    47,
					ArtifactID:   "example-artifact",
					GroupID:      "my-group",
					ModifiedOn:   "2024-12-10T08:56:40Z",
				},
				{
					CreatedOn:    "2024-12-10T08:56:17Z",
					ArtifactType: models.Json,
					State:        models.StateEnabled,
					GlobalID:     46,
					Version:      "1.0.0",
					ContentID:    46,
					ArtifactID:   "example-artifact",
					GroupID:      "my-group",
					ModifiedOn:   "2024-12-10T08:56:17Z",
				},
			},
		}

		server := setupMockServer(
			t,
			http.StatusOK,
			mockResponse,
			"/search/versions",
			http.MethodGet,
		)
		defer server.Close()

		mockClient := &client.Client{BaseURL: server.URL, HTTPClient: server.Client()}
		api := apis.NewVersionsAPI(mockClient)

		// Prepare query parameters
		params := &models.SearchVersionParams{
			ArtifactType: models.Json,
			State:        models.StateEnabled,
		}

		// Execute the function
		versions, err := api.SearchForArtifactVersions(context.Background(), params)

		// Assertions
		assert.NoError(t, err)
		assert.NotNil(t, versions)
		assert.Equal(t, 2, len(versions))
		assert.Equal(t, "2.0.0", versions[0].Version)
		assert.Equal(t, "1.0.0", versions[1].Version)
	})

	t.Run("Invalid Params", func(t *testing.T) {
		mockClient := &client.Client{BaseURL: "http://example.com", HTTPClient: &http.Client{}}
		api := apis.NewVersionsAPI(mockClient)

		// Invalid params: negative limit
		params := &models.SearchVersionParams{
			Limit: -10,
		}

		// Execute the function
		versions, err := api.SearchForArtifactVersions(context.Background(), params)

		// Assertions
		assert.Error(t, err)
		assert.Nil(t, versions)

		var validationErrs validator.ValidationErrors
		if errors.As(err, &validationErrs) {
			foundErrors := map[string]string{}
			for _, fieldErr := range validationErrs {
				foundErrors[fieldErr.StructField()] = fieldErr.Tag()
			}
			assert.Equal(
				t,
				"gte",
				foundErrors["Limit"],
				"Expected Limit to fail with 'gte' validation tag",
			)
		} else {
			t.Fatalf("Expected validation error, got: %v", err)
		}
	})

	t.Run("InternalServerError", func(t *testing.T) {
		mockError := models.APIError{
			Status: 500,
			Title:  "Internal server error",
		}

		server := setupMockServer(
			t,
			http.StatusInternalServerError,
			mockError,
			"/search/versions",
			http.MethodGet,
		)
		defer server.Close()

		mockClient := &client.Client{BaseURL: server.URL, HTTPClient: server.Client()}
		api := apis.NewVersionsAPI(mockClient)

		// Prepare query parameters
		params := &models.SearchVersionParams{
			ArtifactType: models.Json,
			State:        models.StateEnabled,
		}

		// Execute the function
		versions, err := api.SearchForArtifactVersions(context.Background(), params)

		// Assertions
		assert.Error(t, err)
		assert.Nil(t, versions)
		assertAPIError(t, err, 500, "Internal server error")
	})
}

func TestVersionsAPI_SearchForArtifactVersionByContent(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mockResponse := models.ArtifactVersionListResponse{
			Count: 2,
			Versions: []models.ArtifactVersion{
				{
					CreatedOn:    "2024-12-10T08:56:40Z",
					ArtifactType: models.Json,
					State:        models.StateEnabled,
					GlobalID:     47,
					Version:      "2.0.0",
					ContentID:    47,
					ArtifactID:   "example-artifact",
					GroupID:      "my-group",
					ModifiedOn:   "2024-12-10T08:56:40Z",
				},
				{
					CreatedOn:    "2024-12-10T08:56:17Z",
					ArtifactType: models.Json,
					State:        models.StateEnabled,
					GlobalID:     46,
					Version:      "1.0.0",
					ContentID:    46,
					ArtifactID:   "example-artifact",
					GroupID:      "my-group",
					ModifiedOn:   "2024-12-10T08:56:17Z",
				},
			},
		}

		server := setupMockServer(
			t,
			http.StatusOK,
			mockResponse,
			"/search/versions",
			http.MethodPost,
		)
		defer server.Close()

		mockClient := &client.Client{BaseURL: server.URL, HTTPClient: server.Client()}
		api := apis.NewVersionsAPI(mockClient)

		params := &models.SearchVersionByContentParams{Limit: 10, Offset: 0}
		content := `{"key": "value"}`

		versions, err := api.SearchForArtifactVersionByContent(
			context.Background(),
			content,
			params,
		)

		assert.NoError(t, err)
		assert.NotNil(t, versions)
		assert.Equal(t, 2, len(versions))
		assert.Equal(t, "2.0.0", versions[0].Version)
		assert.Equal(t, "1.0.0", versions[1].Version)
	})

	t.Run("BadRequest - Empty Content", func(t *testing.T) {
		mockError := models.APIError{
			Status: 400,
			Title:  "Content cannot be empty",
		}

		server := setupMockServer(
			t,
			http.StatusBadRequest,
			mockError,
			"/search/versions",
			http.MethodPost,
		)
		defer server.Close()

		mockClient := &client.Client{BaseURL: server.URL, HTTPClient: server.Client()}
		api := apis.NewVersionsAPI(mockClient)

		versions, err := api.SearchForArtifactVersionByContent(context.Background(), "", nil)

		assert.Error(t, err)
		assert.Nil(t, versions)
		assertAPIError(t, err, 400, "Content cannot be empty")
	})

	t.Run("InternalServerError", func(t *testing.T) {
		mockError := models.APIError{
			Status: 500,
			Title:  "Internal server error",
		}

		server := setupMockServer(
			t,
			http.StatusInternalServerError,
			mockError,
			"/search/versions",
			http.MethodPost,
		)
		defer server.Close()

		mockClient := &client.Client{BaseURL: server.URL, HTTPClient: server.Client()}
		api := apis.NewVersionsAPI(mockClient)

		content := `{"key": "value"}`

		versions, err := api.SearchForArtifactVersionByContent(context.Background(), content, nil)

		assert.Error(t, err)
		assert.Nil(t, versions)
		assertAPIError(t, err, 500, "Internal server error")
	})

	t.Run("ValidationError - Invalid Params", func(t *testing.T) {
		mockClient := &client.Client{BaseURL: "http://example.com", HTTPClient: &http.Client{}}
		api := apis.NewVersionsAPI(mockClient)

		// Invalid params
		params := &models.SearchVersionByContentParams{
			Offset: -1,
			Limit:  -1,
		}
		content := `{"key": "value"}`

		versions, err := api.SearchForArtifactVersionByContent(
			context.Background(),
			content,
			params,
		)

		assert.Error(t, err)
		assert.Nil(t, versions)

		var validationErrs validator.ValidationErrors
		if errors.As(err, &validationErrs) {
			foundErrors := map[string]string{}
			for _, fieldErr := range validationErrs {
				foundErrors[fieldErr.StructField()] = fieldErr.Tag()
			}

			assert.Equal(
				t,
				"gte",
				foundErrors["Offset"],
				"Expected Offset to fail with 'gte' validation tag",
			)
			assert.Equal(
				t,
				"gte",
				foundErrors["Limit"],
				"Expected Limit to fail with 'gte' validation tag",
			)
		} else {
			t.Fatalf("Expected validation error, got: %v", err)
		}
	})
}

func TestVersionsAPI_GetArtifactVersionState(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mockResponse := models.StateResponse{
			State: models.StateEnabled,
		}

		server := setupMockServer(
			t,
			http.StatusOK,
			mockResponse,
			"/groups/my-group/artifacts/example-artifact/versions/1.0/state",
			http.MethodGet,
		)
		defer server.Close()

		mockClient := &client.Client{BaseURL: server.URL, HTTPClient: server.Client()}
		api := apis.NewVersionsAPI(mockClient)

		state, err := api.GetArtifactVersionState(
			context.Background(),
			"my-group",
			"example-artifact",
			"1.0",
		)

		assert.NoError(t, err)
		assert.NotNil(t, state)
		assert.Equal(t, models.StateEnabled, *state)
	})

	t.Run("NotFound", func(t *testing.T) {
		mockError := models.APIError{
			Status: 404,
			Title:  "Artifact version not found",
		}

		server := setupMockServer(
			t,
			http.StatusNotFound,
			mockError,
			"/groups/my-group/artifacts/example-artifact/versions/1.0/state",
			http.MethodGet,
		)
		defer server.Close()

		mockClient := &client.Client{BaseURL: server.URL, HTTPClient: server.Client()}
		api := apis.NewVersionsAPI(mockClient)

		state, err := api.GetArtifactVersionState(
			context.Background(),
			"my-group",
			"example-artifact",
			"1.0",
		)

		assert.Error(t, err)
		assert.Nil(t, state)
		assertAPIError(t, err, 404, "Artifact version not found")
	})

	t.Run("InternalServerError", func(t *testing.T) {
		mockError := models.APIError{
			Status: 500,
			Title:  "Internal server error",
		}

		server := setupMockServer(
			t,
			http.StatusInternalServerError,
			mockError,
			"/groups/my-group/artifacts/example-artifact/versions/1.0/state",
			http.MethodGet,
		)
		defer server.Close()

		mockClient := &client.Client{BaseURL: server.URL, HTTPClient: server.Client()}
		api := apis.NewVersionsAPI(mockClient)

		state, err := api.GetArtifactVersionState(
			context.Background(),
			"my-group",
			"example-artifact",
			"1.0",
		)

		assert.Error(t, err)
		assert.Nil(t, state)
		assertAPIError(t, err, 500, "Internal server error")
	})

	t.Run("Validation Error - Empty Group ID", func(t *testing.T) {
		mockClient := &client.Client{BaseURL: "http://example.com", HTTPClient: &http.Client{}}
		api := apis.NewVersionsAPI(mockClient)

		state, err := api.GetArtifactVersionState(
			context.Background(),
			"",
			"example-artifact",
			"1.0",
		)

		assert.Error(t, err)
		assert.Nil(t, state)
		assert.Contains(t, err.Error(), "Group ID")
	})

	t.Run("Validation Error - Empty Artifact ID", func(t *testing.T) {
		mockClient := &client.Client{BaseURL: "http://example.com", HTTPClient: &http.Client{}}
		api := apis.NewVersionsAPI(mockClient)

		state, err := api.GetArtifactVersionState(context.Background(), "my-group", "", "1.0")

		assert.Error(t, err)
		assert.Nil(t, state)
		assert.Contains(t, err.Error(), "Artifact ID")
	})

	t.Run("Validation Error - Empty Version Expression", func(t *testing.T) {
		mockClient := &client.Client{BaseURL: "http://example.com", HTTPClient: &http.Client{}}
		api := apis.NewVersionsAPI(mockClient)

		state, err := api.GetArtifactVersionState(
			context.Background(),
			"my-group",
			"example-artifact",
			"",
		)

		assert.Error(t, err)
		assert.Nil(t, state)
		assert.Contains(t, err.Error(), "Version Expression")
	})
}

func TestVersionsAPI_UpdateArtifactVersionState(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		server := setupMockServer(t, http.StatusNoContent, nil,
			"/groups/my-group/artifacts/example-artifact/versions/1.0/state", http.MethodPut)
		defer server.Close()

		mockClient := &client.Client{BaseURL: server.URL, HTTPClient: server.Client()}
		api := apis.NewVersionsAPI(mockClient)

		err := api.UpdateArtifactVersionState(
			context.Background(),
			"my-group",
			"example-artifact",
			"1.0",
			models.StateEnabled,
			false,
		)
		assert.NoError(t, err)
	})

	t.Run("BadRequest", func(t *testing.T) {
		mockError := models.APIError{Status: 400, Title: "Invalid state"}
		server := setupMockServer(t, http.StatusBadRequest, mockError,
			"/groups/my-group/artifacts/example-artifact/versions/1.0/state", http.MethodPut)
		defer server.Close()

		mockClient := &client.Client{BaseURL: server.URL, HTTPClient: server.Client()}
		api := apis.NewVersionsAPI(mockClient)

		err := api.UpdateArtifactVersionState(
			context.Background(),
			"my-group",
			"example-artifact",
			"1.0",
			"INVALID_STATE",
			false,
		)

		assert.Error(t, err)
		assertAPIError(t, err, 400, "Invalid state")
	})

	t.Run("Conflict", func(t *testing.T) {
		mockError := models.APIError{Status: 409, Title: "Conflict"}
		server := setupMockServer(t, http.StatusConflict, mockError,
			"/groups/my-group/artifacts/example-artifact/versions/1.0/state", http.MethodPut)
		defer server.Close()

		mockClient := &client.Client{BaseURL: server.URL, HTTPClient: server.Client()}
		api := apis.NewVersionsAPI(mockClient)

		err := api.UpdateArtifactVersionState(
			context.Background(),
			"my-group",
			"example-artifact",
			"1.0",
			models.StateDraft,
			false,
		)

		assert.Error(t, err)
		assertAPIError(t, err, 409, "Conflict")
	})

	t.Run("InternalServerError", func(t *testing.T) {
		mockError := models.APIError{Status: 500, Title: "Internal server error"}
		server := setupMockServer(t, http.StatusInternalServerError, mockError,
			"/groups/my-group/artifacts/example-artifact/versions/1.0/state", http.MethodPut)
		defer server.Close()

		mockClient := &client.Client{BaseURL: server.URL, HTTPClient: server.Client()}
		api := apis.NewVersionsAPI(mockClient)

		err := api.UpdateArtifactVersionState(
			context.Background(),
			"my-group",
			"example-artifact",
			"1.0",
			models.StateEnabled,
			false,
		)

		assert.Error(t, err)
		assertAPIError(t, err, 500, "Internal server error")
	})

	t.Run("Validation Error - Empty Group ID", func(t *testing.T) {
		mockClient := &client.Client{BaseURL: "http://example.com", HTTPClient: &http.Client{}}
		api := apis.NewVersionsAPI(mockClient)

		err := api.UpdateArtifactVersionState(
			context.Background(),
			"",
			"example-artifact",
			"1.0",
			models.StateEnabled,
			false,
		)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "Group ID")
	})

	t.Run("Validation Error - Empty Artifact ID", func(t *testing.T) {
		mockClient := &client.Client{BaseURL: "http://example.com", HTTPClient: &http.Client{}}
		api := apis.NewVersionsAPI(mockClient)

		err := api.UpdateArtifactVersionState(
			context.Background(),
			"my-group",
			"",
			"1.0",
			models.StateEnabled,
			false,
		)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "Artifact ID")
	})

	t.Run("Validation Error - Empty Version Expression", func(t *testing.T) {
		mockClient := &client.Client{BaseURL: "http://example.com", HTTPClient: &http.Client{}}
		api := apis.NewVersionsAPI(mockClient)

		err := api.UpdateArtifactVersionState(
			context.Background(),
			"my-group",
			"example-artifact",
			"",
			models.StateEnabled,
			false,
		)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "Version Expression")
	})
}

func TestVersionsAPI_InputValidation(t *testing.T) {
	t.Run("Empty Group ID", func(t *testing.T) {
		mockClient := &client.Client{}
		api := apis.NewVersionsAPI(mockClient)

		err := api.DeleteArtifactVersion(context.Background(), "", "test-artifact", "1.0.0")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "Group ID")
	})

	t.Run("Empty Artifact ID", func(t *testing.T) {
		mockClient := &client.Client{}
		api := apis.NewVersionsAPI(mockClient)

		err := api.DeleteArtifactVersion(context.Background(), "test-group", "", "1.0.0")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "Artifact ID")
	})

	t.Run("Empty Version Expression", func(t *testing.T) {
		mockClient := &client.Client{}
		api := apis.NewVersionsAPI(mockClient)

		err := api.DeleteArtifactVersion(context.Background(), "test-group", "test-artifact", "")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "Version Expression")
	})

	t.Run("Invalid Group ID", func(t *testing.T) {
		mockClient := &client.Client{}
		api := apis.NewVersionsAPI(mockClient)

		err := api.DeleteArtifactVersion(context.Background(), "", "test-artifact", "1.0.0")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "Group ID")
	})

	t.Run("Invalid Artifact ID", func(t *testing.T) {
		mockClient := &client.Client{}
		api := apis.NewVersionsAPI(mockClient)

		err := api.DeleteArtifactVersion(context.Background(), "test-group", "", "1.0.0")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "Artifact ID")
	})

	t.Run("Invalid Group ID Over 512", func(t *testing.T) {
		invalidGroupID := strings.Repeat("a", 513)
		mockClient := &client.Client{}
		api := apis.NewVersionsAPI(mockClient)

		err := api.DeleteArtifactVersion(
			context.Background(),
			invalidGroupID,
			"test-artifact",
			"1.0.0",
		)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "Group ID")
	})

	t.Run("Invalid Artifact ID Over 512", func(t *testing.T) {
		invalidArtifactId := strings.Repeat("a", 513)
		mockClient := &client.Client{}
		api := apis.NewVersionsAPI(mockClient)

		err := api.DeleteArtifactVersion(
			context.Background(),
			"test-group",
			invalidArtifactId,
			"1.0.0",
		)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "Artifact ID")
	})
}

func TestVersionsAPI_HTTPRequestErrors(t *testing.T) {
	t.Run("Timeout", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			time.Sleep(2 * time.Second) // Simulate a timeout
		}))
		defer server.Close()

		mockClient := &client.Client{
			BaseURL:    server.URL,
			HTTPClient: &http.Client{Timeout: time.Millisecond * 500},
		}
		api := apis.NewVersionsAPI(mockClient)

		err := api.DeleteArtifactVersion(
			context.Background(),
			"test-group",
			"test-artifact",
			"1.0.0",
		)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "Client.Timeout")
	})

	t.Run("Connection Refused", func(t *testing.T) {
		mockClient := &client.Client{
			BaseURL:    "http://localhost:9999",
			HTTPClient: &http.Client{Timeout: time.Second},
		}
		api := apis.NewVersionsAPI(mockClient)

		err := api.DeleteArtifactVersion(
			context.Background(),
			"test-group",
			"test-artifact",
			"1.0.0",
		)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "connection refused")
	})
}

/***********************/
/***** Integration *****/
/***********************/

func setupVersionsAPIClient() *apis.VersionsAPI {
	apiClient := setupHTTPClient()
	return apis.NewVersionsAPI(apiClient)
}

func TestVersionsAPIIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	ctx := context.Background()

	versionsAPI := setupVersionsAPIClient()

	// Prepare test data
	artifactsAPI := apis.NewArtifactsAPI(versionsAPI.Client)

	// Clean up before and after tests
	t.Cleanup(func() { cleanup(t, artifactsAPI) })
	cleanup(t, artifactsAPI)

	// Test CreateArtifactVersion
	t.Run("CreateArtifactVersion", func(t *testing.T) {
		generatedArtifactID, err := generateArtifactForTest(ctx, artifactsAPI)
		if err != nil {
			t.Fatal(err)
		}

		request := &models.CreateVersionRequest{
			Version: newVersion,
			Content: models.CreateContentRequest{
				Content: stubNewContent,
			},
		}

		resp, err := versionsAPI.CreateArtifactVersion(
			ctx,
			stubGroupId,
			generatedArtifactID,
			request,
			false,
		)
		assert.NoError(t, err)
		assert.Equal(t, newVersion, resp.Version)
	})

	// Test ListArtifactVersions
	t.Run("ListArtifactVersions", func(t *testing.T) {
		generatedArtifactID, err := generateArtifactForTest(ctx, artifactsAPI)
		if err != nil {
			t.Fatal(err)
		}

		params := &models.ListArtifactsVersionsParams{}
		resp, err := versionsAPI.ListArtifactVersions(ctx, stubGroupId, generatedArtifactID, params)
		assert.NoError(t, err)
		assert.GreaterOrEqual(t, len(resp), 1)
	})

	// Test GetArtifactVersionReferences
	t.Run("GetArtifactVersionReferences", func(t *testing.T) {
		generatedArtifactID, err := generateArtifactForTest(ctx, artifactsAPI)
		if err != nil {
			t.Fatal(err)
		}

		params := &models.ArtifactVersionReferencesParams{}
		references, err := versionsAPI.GetArtifactVersionReferences(
			ctx,
			stubGroupId,
			generatedArtifactID,
			version,
			params,
		)
		assert.NoError(t, err)
		assert.NotNil(t, references)
	})

	// Test AddArtifactVersionComment
	t.Run("AddArtifactVersionComment", func(t *testing.T) {
		generatedArtifactID, err := generateArtifactForTest(ctx, artifactsAPI)
		if err != nil {
			t.Fatal(err)
		}

		comment, err := versionsAPI.AddArtifactVersionComment(
			ctx,
			stubGroupId,
			generatedArtifactID,
			version,
			"Test comment",
		)
		assert.NoError(t, err)
		assert.Equal(t, "Test comment", comment.Value)
	})

	// Test GetArtifactVersionComments
	t.Run("GetArtifactVersionComments", func(t *testing.T) {
		generatedArtifactID, err := generateArtifactForTest(ctx, artifactsAPI)
		if err != nil {
			t.Fatal(err)
		}

		// Add a comment first
		comment, err := versionsAPI.AddArtifactVersionComment(
			ctx,
			stubGroupId,
			generatedArtifactID,
			version,
			"Test comment",
		)
		assert.NoError(t, err)
		assert.Equal(t, "Test comment", comment.Value)

		// Get comments
		comments, err := versionsAPI.GetArtifactVersionComments(
			ctx,
			stubGroupId,
			generatedArtifactID,
			version,
		)
		assert.NoError(t, err)
		assert.NotNil(t, comments)
	})

	// Test UpdateArtifactVersionComment
	t.Run("UpdateArtifactVersionComment", func(t *testing.T) {
		generatedArtifactID, err := generateArtifactForTest(ctx, artifactsAPI)
		if err != nil {
			t.Fatal(err)
		}

		// Add a comment first
		comment, err := versionsAPI.AddArtifactVersionComment(
			ctx,
			stubGroupId,
			generatedArtifactID,
			version,
			"Initial comment",
		)
		assert.NoError(t, err)

		// Update the comment
		err = versionsAPI.UpdateArtifactVersionComment(
			ctx,
			stubGroupId,
			generatedArtifactID,
			version,
			comment.CommentID,
			"Updated comment",
		)
		assert.NoError(t, err)
	})

	// Test DeleteArtifactVersionComment
	t.Run("DeleteArtifactVersionComment", func(t *testing.T) {
		generatedArtifactID, err := generateArtifactForTest(ctx, artifactsAPI)
		if err != nil {
			t.Fatal(err)
		}

		// Add a comment first
		comment, err := versionsAPI.AddArtifactVersionComment(
			ctx,
			stubGroupId,
			generatedArtifactID,
			version,
			"Temporary comment",
		)
		assert.NoError(t, err)

		// Delete the comment
		err = versionsAPI.DeleteArtifactVersionComment(
			ctx,
			stubGroupId,
			generatedArtifactID,
			version,
			comment.CommentID,
		)
		assert.NoError(t, err)
	})

	// Test DeleteArtifactVersion
	t.Run("DeleteArtifactVersion", func(t *testing.T) {
		generatedArtifactID, err := generateArtifactForTest(ctx, artifactsAPI)
		if err != nil {
			t.Fatal(err)
		}

		err = versionsAPI.DeleteArtifactVersion(ctx, stubGroupId, generatedArtifactID, version)
		assert.NoError(t, err)
	})

	// Test GetArtifactVersionContent
	t.Run("GetArtifactVersionContent", func(t *testing.T) {
		generatedArtifactID, err := generateArtifactForTest(ctx, artifactsAPI)
		if err != nil {
			t.Fatal(err)
		}

		params := &models.ArtifactReferenceParams{}
		content, err := versionsAPI.GetArtifactVersionContent(
			ctx,
			stubGroupId,
			generatedArtifactID,
			version,
			params,
		)
		assert.NoError(t, err)
		assert.NotNil(t, content)
	})

	// Test UpdateArtifactVersionContent
	t.Run("UpdateArtifactVersionContent", func(t *testing.T) {
		generatedArtifactID, err := generateArtifactForTest(ctx, artifactsAPI)
		if err != nil {
			t.Fatal(err)
		}

		content := &models.CreateContentRequest{
			Content:     stubContent,
			ContentType: "application/json",
		}
		err = versionsAPI.UpdateArtifactVersionContent(
			ctx,
			stubGroupId,
			generatedArtifactID,
			version,
			content,
		)
		assert.NoError(t, err)
	})

	// Test SearchForArtifactVersions
	t.Run("SearchForArtifactVersions", func(t *testing.T) {
		_, err := generateArtifactForTest(ctx, artifactsAPI)
		if err != nil {
			t.Fatal(err)
		}

		params := &models.SearchVersionParams{
			Version: version,
		}
		versions, err := versionsAPI.SearchForArtifactVersions(ctx, params)
		assert.NoError(t, err)
		assert.GreaterOrEqual(t, len(versions), 1)
	})

	// Test GetArtifactVersionState
	t.Run("GetArtifactVersionState", func(t *testing.T) {
		generatedArtifactID, err := generateArtifactForTest(ctx, artifactsAPI)
		if err != nil {
			t.Fatal(err)
		}

		state, err := versionsAPI.GetArtifactVersionState(
			ctx,
			stubGroupId,
			generatedArtifactID,
			version,
		)
		assert.NoError(t, err)
		assert.Equal(t, models.StateDraft, *state)
	})

	// Test UpdateArtifactVersionState
	t.Run("UpdateArtifactVersionState", func(t *testing.T) {
		generatedArtifactID, err := generateArtifactForTest(ctx, artifactsAPI)
		if err != nil {
			t.Fatal(err)
		}

		err = versionsAPI.UpdateArtifactVersionState(
			ctx,
			stubGroupId,
			generatedArtifactID,
			version,
			models.StateDeprecated,
			false,
		)
		assert.NoError(t, err)
	})
}
