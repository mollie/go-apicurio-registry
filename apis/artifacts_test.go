package apis_test

import (
	"context"
	"fmt"
	"github.com/mollie/go-apicurio-registry/apis"
	"github.com/mollie/go-apicurio-registry/client"
	"github.com/mollie/go-apicurio-registry/models"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"
)

func TestArtifactsAPI_GetArtifactByGlobalID(t *testing.T) {
	t.Run("Success-Without-ArtifactType", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "/ids/globalIds/1", r.URL.Path)
			assert.Equal(t, http.MethodGet, r.Method)

			w.WriteHeader(http.StatusOK)
			_, err := w.Write([]byte(stubArtifactContent))
			assert.NoError(t, err)
		}))
		defer server.Close()

		mockClient := &client.Client{BaseURL: server.URL, HTTPClient: server.Client()}
		api := apis.NewArtifactsAPI(mockClient)

		result, err := api.GetArtifactByGlobalID(context.Background(), 1, nil)
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, stubArtifactContent, result.Content)
	})

	t.Run("Success-With-ArtifactType", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "/ids/globalIds/1", r.URL.Path)
			assert.Equal(t, http.MethodGet, r.Method)

			w.Header().Set("X-Registry-ArtifactType", "JSON")
			w.WriteHeader(http.StatusOK)
			_, err := w.Write([]byte(stubArtifactContent))
			assert.NoError(t, err)
		}))
		defer server.Close()

		mockClient := &client.Client{BaseURL: server.URL, HTTPClient: server.Client()}
		api := apis.NewArtifactsAPI(mockClient)

		params := models.GetArtifactByGlobalIDParams{ReturnArtifactType: true}
		result, err := api.GetArtifactByGlobalID(context.Background(), 1, &params)
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, stubArtifactContent, result.Content)
		assert.Equal(t, models.Json, result.ArtifactType)
	})

	t.Run("Not Found", func(t *testing.T) {
		errorResponse := models.APIError{Status: http.StatusNotFound, Title: TitleNotFound}
		server := setupMockServer(t, http.StatusNotFound, errorResponse, "/ids/globalIds/1", http.MethodGet)
		defer server.Close()

		mockClient := &client.Client{BaseURL: server.URL, HTTPClient: server.Client()}
		api := apis.NewArtifactsAPI(mockClient)

		result, err := api.GetArtifactByGlobalID(context.Background(), 1, nil)
		assert.Error(t, err)
		assert.Nil(t, result)
		assertAPIError(t, err, http.StatusNotFound, TitleNotFound)
	})

	t.Run("Internal Server Error", func(t *testing.T) {
		errorResponse := models.APIError{Status: http.StatusInternalServerError, Title: TitleInternalServerError}
		server := setupMockServer(t, http.StatusInternalServerError, errorResponse, "/ids/globalIds/1", http.MethodGet)
		defer server.Close()

		mockClient := &client.Client{BaseURL: server.URL, HTTPClient: server.Client()}
		api := apis.NewArtifactsAPI(mockClient)

		result, err := api.GetArtifactByGlobalID(context.Background(), 1, nil)
		assert.Error(t, err)
		assert.Nil(t, result)
		assertAPIError(t, err, http.StatusInternalServerError, TitleInternalServerError)
	})
}

func TestArtifactsAPI_SearchArtifacts(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mockResponse := models.SearchArtifactsAPIResponse{
			Artifacts: []models.SearchedArtifact{
				{
					GroupId:      "test-group",
					ArtifactId:   "artifact-1",
					Name:         "Test Artifact",
					ArtifactType: models.Avro,
				},
			},
			Count: 1,
		}

		server := setupMockServer(t, http.StatusOK, mockResponse, "/search/artifacts", http.MethodGet)
		defer server.Close()

		mockClient := &client.Client{BaseURL: server.URL, HTTPClient: server.Client()}
		api := apis.NewArtifactsAPI(mockClient)

		params := &models.SearchArtifactsParams{Name: "Test Artifact"}
		result, err := api.SearchArtifacts(context.Background(), params)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, 1, len(result))
		assert.Equal(t, "test-group", result[0].GroupId)
		assert.Equal(t, "artifact-1", result[0].ArtifactId)
		assert.Equal(t, models.Avro, result[0].ArtifactType)
	})

	t.Run("Internal Server Error", func(t *testing.T) {
		errorResponse := models.APIError{Status: http.StatusInternalServerError, Title: TitleInternalServerError}
		server := setupMockServer(t, http.StatusInternalServerError, errorResponse, "/search/artifacts", http.MethodGet)
		defer server.Close()

		mockClient := &client.Client{BaseURL: server.URL, HTTPClient: server.Client()}
		api := apis.NewArtifactsAPI(mockClient)

		params := &models.SearchArtifactsParams{}
		result, err := api.SearchArtifacts(context.Background(), params)

		assert.Error(t, err)
		assert.Nil(t, result)
		assertAPIError(t, err, http.StatusInternalServerError, TitleInternalServerError)
	})
}

func TestArtifactsAPI_SearchArtifactsByContent(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mockResponse := models.SearchArtifactsAPIResponse{
			Artifacts: []models.SearchedArtifact{
				{
					GroupId:      "test-group",
					ArtifactId:   "artifact-1",
					Name:         "Test Artifact",
					ArtifactType: models.Avro,
				},
			},
			Count: 1,
		}

		server := setupMockServer(t, http.StatusOK, mockResponse, "/search/artifacts", http.MethodPost)
		defer server.Close()

		mockClient := &client.Client{BaseURL: server.URL, HTTPClient: server.Client()}
		api := apis.NewArtifactsAPI(mockClient)

		params := &models.SearchArtifactsByContentParams{Canonical: true}
		result, err := api.SearchArtifactsByContent(context.Background(), []byte("{\"key\":\"value\"}"), params)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, 1, len(result))
		assert.Equal(t, "test-group", result[0].GroupId)
		assert.Equal(t, "artifact-1", result[0].ArtifactId)
		assert.Equal(t, models.Avro, result[0].ArtifactType)
	})

	t.Run("Invalid Content", func(t *testing.T) {
		errorResponse := models.APIError{Status: http.StatusBadRequest, Title: TitleBadRequest}
		server := setupMockServer(t, http.StatusBadRequest, errorResponse, "/search/artifacts", http.MethodPost)
		defer server.Close()

		mockClient := &client.Client{BaseURL: server.URL, HTTPClient: server.Client()}
		api := apis.NewArtifactsAPI(mockClient)

		params := &models.SearchArtifactsByContentParams{}
		result, err := api.SearchArtifactsByContent(context.Background(), []byte(""), params)

		assert.Error(t, err)
		assert.Nil(t, result)
		assertAPIError(t, err, http.StatusBadRequest, TitleBadRequest)
	})
}

func TestArtifactsAPI_ListArtifactReferences(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mockReferences := []models.ArtifactReference{
			{GroupID: "group-1", ArtifactID: "artifact-1", Version: "v1"},
		}

		server := setupMockServer(t, http.StatusOK, mockReferences, "/ids/contentId/123/references", http.MethodGet)
		defer server.Close()

		mockClient := &client.Client{BaseURL: server.URL, HTTPClient: server.Client()}
		api := apis.NewArtifactsAPI(mockClient)

		result, err := api.ListArtifactReferences(context.Background(), 123)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Len(t, *result, 1)
		assert.Equal(t, "group-1", (*result)[0].GroupID)
		assert.Equal(t, "artifact-1", (*result)[0].ArtifactID)
		assert.Equal(t, "v1", (*result)[0].Version)
	})

	t.Run("Not Found", func(t *testing.T) {
		errorResponse := models.APIError{Status: http.StatusNotFound, Title: TitleNotFound}
		server := setupMockServer(t, http.StatusNotFound, errorResponse, "/ids/contentId/123/references", http.MethodGet)
		defer server.Close()

		mockClient := &client.Client{BaseURL: server.URL, HTTPClient: server.Client()}
		api := apis.NewArtifactsAPI(mockClient)

		result, err := api.ListArtifactReferences(context.Background(), 123)

		assert.Error(t, err)
		assert.Nil(t, result)
		assertAPIError(t, err, http.StatusNotFound, TitleNotFound)
	})
}

func TestArtifactsAPI_ListArtifactReferencesByGlobalID(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mockReferences := []models.ArtifactReference{
			{GroupID: "group-1", ArtifactID: "artifact-1", Version: "v1"},
		}

		server := setupMockServer(t, http.StatusOK, mockReferences, "/ids/globalIds/123/references", http.MethodGet)
		defer server.Close()

		mockClient := &client.Client{BaseURL: server.URL, HTTPClient: server.Client()}
		api := apis.NewArtifactsAPI(mockClient)

		params := &models.ListArtifactReferencesByGlobalIDParams{RefType: models.OutBound}
		result, err := api.ListArtifactReferencesByGlobalID(context.Background(), 123, params)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Len(t, *result, 1)
		assert.Equal(t, "group-1", (*result)[0].GroupID)
		assert.Equal(t, "artifact-1", (*result)[0].ArtifactID)
		assert.Equal(t, "v1", (*result)[0].Version)
	})

	t.Run("Internal Server Error", func(t *testing.T) {
		errorResponse := models.APIError{Status: http.StatusInternalServerError, Title: TitleInternalServerError}
		server := setupMockServer(t, http.StatusInternalServerError, errorResponse, "/ids/globalIds/123/references", http.MethodGet)
		defer server.Close()

		mockClient := &client.Client{BaseURL: server.URL, HTTPClient: server.Client()}
		api := apis.NewArtifactsAPI(mockClient)

		params := &models.ListArtifactReferencesByGlobalIDParams{}
		result, err := api.ListArtifactReferencesByGlobalID(context.Background(), 123, params)

		assert.Error(t, err)
		assert.Nil(t, result)
		assertAPIError(t, err, http.StatusInternalServerError, TitleInternalServerError)
	})
}

func TestArtifactsAPI_ListArtifactReferencesByHash(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mockReferences := []models.ArtifactReference{
			{GroupID: "group-1", ArtifactID: "artifact-1", Version: "v1"},
		}

		server := setupMockServer(t, http.StatusOK, mockReferences, "/ids/contentHashes/hash-123/references", http.MethodGet)
		defer server.Close()

		mockClient := &client.Client{BaseURL: server.URL, HTTPClient: server.Client()}
		api := apis.NewArtifactsAPI(mockClient)

		result, err := api.ListArtifactReferencesByHash(context.Background(), "hash-123")

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Len(t, result, 1)
		assert.Equal(t, "group-1", result[0].GroupID)
		assert.Equal(t, "artifact-1", result[0].ArtifactID)
		assert.Equal(t, "v1", result[0].Version)
	})

	t.Run("Internal Server Error", func(t *testing.T) {
		errorResponse := models.APIError{Status: http.StatusInternalServerError, Title: TitleInternalServerError}
		server := setupMockServer(t, http.StatusInternalServerError, errorResponse, "/ids/contentHashes/hash-123/references", http.MethodGet)
		defer server.Close()

		mockClient := &client.Client{BaseURL: server.URL, HTTPClient: server.Client()}
		api := apis.NewArtifactsAPI(mockClient)

		result, err := api.ListArtifactReferencesByHash(context.Background(), "hash-123")

		assert.Error(t, err)
		assert.Nil(t, result)
		assertAPIError(t, err, http.StatusInternalServerError, TitleInternalServerError)
	})
}

func TestArtifactsAPI_ListArtifactsInGroup(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mockResponse := models.ListArtifactsResponse{
			Artifacts: []models.SearchedArtifact{
				{
					GroupId:      "group-1",
					ArtifactId:   "artifact-1",
					Name:         "Test Artifact",
					ArtifactType: models.XML,
				},
			},
			Count: 1,
		}

		server := setupMockServer(t, http.StatusOK, mockResponse, "/groups/group-1/artifacts", http.MethodGet)
		defer server.Close()

		mockClient := &client.Client{BaseURL: server.URL, HTTPClient: server.Client()}
		api := apis.NewArtifactsAPI(mockClient)

		params := &models.ListArtifactsInGroupParams{Limit: 10, Offset: 0, Order: "asc"}
		result, err := api.ListArtifactsInGroup(context.Background(), "group-1", params)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Len(t, result.Artifacts, 1)
		assert.Equal(t, "group-1", result.Artifacts[0].GroupId)
		assert.Equal(t, "artifact-1", result.Artifacts[0].ArtifactId)
		assert.Equal(t, models.XML, result.Artifacts[0].ArtifactType)
	})

	t.Run("Internal Server Error", func(t *testing.T) {
		errorResponse := models.APIError{Status: http.StatusInternalServerError, Title: TitleInternalServerError}
		server := setupMockServer(t, http.StatusInternalServerError, errorResponse, "/groups/group-1/artifacts", http.MethodGet)
		defer server.Close()

		mockClient := &client.Client{BaseURL: server.URL, HTTPClient: server.Client()}
		api := apis.NewArtifactsAPI(mockClient)

		params := &models.ListArtifactsInGroupParams{}
		result, err := api.ListArtifactsInGroup(context.Background(), "group-1", params)

		assert.Error(t, err)
		assert.Nil(t, result)
		assertAPIError(t, err, http.StatusInternalServerError, TitleInternalServerError)
	})
}

func TestArtifactsAPI_GetArtifactContentByHash(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mockContent := models.ArtifactContent{
			Content:      "{\"key\":\"value\"}",
			ArtifactType: models.Json,
		}
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Contains(t, r.URL.Path, "/ids/contentHashes/hash-123")
			assert.Equal(t, http.MethodGet, r.Method)

			w.Header().Set("X-Registry-ArtifactType", "JSON")
			w.WriteHeader(http.StatusOK)
			_, err := w.Write([]byte(mockContent.Content))
			assert.NoError(t, err)
		}))
		defer server.Close()

		mockClient := &client.Client{BaseURL: server.URL, HTTPClient: server.Client()}
		api := apis.NewArtifactsAPI(mockClient)

		result, err := api.GetArtifactContentByHash(context.Background(), "hash-123")
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, "{\"key\":\"value\"}", result.Content)
		assert.Equal(t, models.Json, result.ArtifactType)
	})

	t.Run("Not Found", func(t *testing.T) {
		errorResponse := models.APIError{Status: http.StatusNotFound, Title: TitleNotFound}
		server := setupMockServer(t, http.StatusNotFound, errorResponse, "/ids/contentHashes/hash-123", http.MethodGet)
		defer server.Close()

		mockClient := &client.Client{BaseURL: server.URL, HTTPClient: server.Client()}
		api := apis.NewArtifactsAPI(mockClient)

		result, err := api.GetArtifactContentByHash(context.Background(), "hash-123")
		assert.Error(t, err)
		assert.Nil(t, result)
		assertAPIError(t, err, http.StatusNotFound, TitleNotFound)
	})

	t.Run("Internal Server Error", func(t *testing.T) {
		errorResponse := models.APIError{Status: http.StatusInternalServerError, Title: TitleInternalServerError}
		server := setupMockServer(t, http.StatusInternalServerError, errorResponse, "/ids/contentHashes/hash-123", http.MethodGet)
		defer server.Close()

		mockClient := &client.Client{BaseURL: server.URL, HTTPClient: server.Client()}
		api := apis.NewArtifactsAPI(mockClient)

		result, err := api.GetArtifactContentByHash(context.Background(), "hash-123")
		assert.Error(t, err)
		assert.Nil(t, result)
		assertAPIError(t, err, http.StatusInternalServerError, TitleInternalServerError)
	})
}

func TestArtifactsAPI_GetArtifactContentByID(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mockContent := models.ArtifactContent{
			Content:      "{\"key\":\"value\"}",
			ArtifactType: models.Json,
		}
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Contains(t, r.URL.Path, "/ids/contentIds/123")
			assert.Equal(t, http.MethodGet, r.Method)

			w.Header().Set("X-Registry-ArtifactType", "JSON")
			w.WriteHeader(http.StatusOK)
			_, err := w.Write([]byte(mockContent.Content))
			assert.NoError(t, err)
		}))
		defer server.Close()

		mockClient := &client.Client{BaseURL: server.URL, HTTPClient: server.Client()}
		api := apis.NewArtifactsAPI(mockClient)

		result, err := api.GetArtifactContentByID(context.Background(), 123)
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, "{\"key\":\"value\"}", result.Content)
		assert.Equal(t, models.Json, result.ArtifactType)
	})

	t.Run("Not Found", func(t *testing.T) {
		errorResponse := models.APIError{Status: http.StatusNotFound, Title: TitleNotFound}
		server := setupMockServer(t, http.StatusNotFound, errorResponse, "/ids/contentIds/123", http.MethodGet)
		defer server.Close()

		mockClient := &client.Client{BaseURL: server.URL, HTTPClient: server.Client()}
		api := apis.NewArtifactsAPI(mockClient)

		result, err := api.GetArtifactContentByID(context.Background(), 123)
		assert.Error(t, err)
		assert.Nil(t, result)
		assertAPIError(t, err, http.StatusNotFound, TitleNotFound)
	})

	t.Run("Internal Server Error", func(t *testing.T) {
		errorResponse := models.APIError{Status: http.StatusInternalServerError, Title: TitleInternalServerError}
		server := setupMockServer(t, http.StatusInternalServerError, errorResponse, "/ids/contentIds/123", http.MethodGet)
		defer server.Close()

		mockClient := &client.Client{BaseURL: server.URL, HTTPClient: server.Client()}
		api := apis.NewArtifactsAPI(mockClient)

		result, err := api.GetArtifactContentByID(context.Background(), 123)
		assert.Error(t, err)
		assert.Nil(t, result)
		assertAPIError(t, err, http.StatusInternalServerError, TitleInternalServerError)
	})
}

func TestArtifactsAPI_DeleteArtifactsInGroup(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		server := setupMockServer(t, http.StatusNoContent, nil, "/groups/group-1/artifacts", http.MethodDelete)
		defer server.Close()

		mockClient := &client.Client{BaseURL: server.URL, HTTPClient: server.Client()}
		api := apis.NewArtifactsAPI(mockClient)

		err := api.DeleteArtifactsInGroup(context.Background(), "group-1")
		assert.NoError(t, err)
	})

	t.Run("Not Found", func(t *testing.T) {
		errorResponse := models.APIError{Status: http.StatusNotFound, Title: TitleNotFound}
		server := setupMockServer(t, http.StatusNotFound, errorResponse, "/groups/group-1/artifacts", http.MethodDelete)
		defer server.Close()

		mockClient := &client.Client{BaseURL: server.URL, HTTPClient: server.Client()}
		api := apis.NewArtifactsAPI(mockClient)

		err := api.DeleteArtifactsInGroup(context.Background(), "group-1")
		assert.Error(t, err)
		assertAPIError(t, err, http.StatusNotFound, TitleNotFound)
	})

	t.Run("Internal Server Error", func(t *testing.T) {
		errorResponse := models.APIError{Status: http.StatusInternalServerError, Title: TitleInternalServerError}
		server := setupMockServer(t, http.StatusInternalServerError, errorResponse, "/groups/group-1/artifacts", http.MethodDelete)
		defer server.Close()

		mockClient := &client.Client{BaseURL: server.URL, HTTPClient: server.Client()}
		api := apis.NewArtifactsAPI(mockClient)

		err := api.DeleteArtifactsInGroup(context.Background(), "group-1")
		assert.Error(t, err)
		assertAPIError(t, err, http.StatusInternalServerError, TitleInternalServerError)
	})
}

func TestArtifactsAPI_DeleteArtifact(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		server := setupMockServer(t, http.StatusNoContent, nil, "/groups/test-group/artifacts/artifact-1", http.MethodDelete)
		defer server.Close()

		mockClient := &client.Client{BaseURL: server.URL, HTTPClient: server.Client()}
		api := apis.NewArtifactsAPI(mockClient)

		err := api.DeleteArtifact(context.Background(), "test-group", "artifact-1")
		assert.NoError(t, err)
	})

	t.Run("Not Found", func(t *testing.T) {
		errorResponse := models.APIError{Status: http.StatusNotFound, Title: TitleNotFound}
		server := setupMockServer(t, http.StatusNotFound, errorResponse, "/groups/test-group/artifacts/artifact-1", http.MethodDelete)
		defer server.Close()

		mockClient := &client.Client{BaseURL: server.URL, HTTPClient: server.Client()}
		api := apis.NewArtifactsAPI(mockClient)

		err := api.DeleteArtifact(context.Background(), "test-group", "artifact-1")
		assert.Error(t, err)
		assertAPIError(t, err, http.StatusNotFound, TitleNotFound)
	})

	t.Run("Not Allowed", func(t *testing.T) {
		errorResponse := models.APIError{Status: http.StatusMethodNotAllowed, Title: TitleMethodNotAllowed}
		server := setupMockServer(t, http.StatusMethodNotAllowed, errorResponse, "/groups/test-group/artifacts/artifact-1", http.MethodDelete)
		defer server.Close()

		mockClient := &client.Client{BaseURL: server.URL, HTTPClient: server.Client()}
		api := apis.NewArtifactsAPI(mockClient)

		err := api.DeleteArtifact(context.Background(), "test-group", "artifact-1")
		assert.Error(t, err)
		assertAPIError(t, err, http.StatusMethodNotAllowed, TitleMethodNotAllowed)
	})

	t.Run("Internal Server Error", func(t *testing.T) {
		errorResponse := models.APIError{Status: http.StatusInternalServerError, Title: TitleInternalServerError}
		server := setupMockServer(t, http.StatusInternalServerError, errorResponse, "/groups/test-group/artifacts/artifact-1", http.MethodDelete)
		defer server.Close()

		mockClient := &client.Client{BaseURL: server.URL, HTTPClient: server.Client()}
		api := apis.NewArtifactsAPI(mockClient)

		err := api.DeleteArtifact(context.Background(), "test-group", "artifact-1")
		assert.Error(t, err)
		assertAPIError(t, err, http.StatusInternalServerError, TitleInternalServerError)
	})
}

func TestArtifactsAPI_CreateArtifact(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mockResponse := models.CreateArtifactResponse{
			Artifact: models.ArtifactDetail{
				GroupID:     "test-group",
				ArtifactID:  "artifact-1",
				Name:        "New Artifact",
				Description: "Test Description",
			},
		}

		server := setupMockServer(t, http.StatusOK, mockResponse, "/groups/test-group/artifacts", http.MethodPost)
		defer server.Close()

		mockClient := &client.Client{BaseURL: server.URL, HTTPClient: server.Client()}
		api := apis.NewArtifactsAPI(mockClient)

		artifact := models.CreateArtifactRequest{
			ArtifactID:   stubArtifactId,
			ArtifactType: models.Json,
			FirstVersion: models.CreateVersionRequest{
				Version: "1.0.0",
				Content: models.CreateContentRequest{Content: "{\"key\":\"value\"}", ContentType: "application/json"},
			},
		}
		params := &models.CreateArtifactParams{IfExists: models.IfExistsCreate}

		result, err := api.CreateArtifact(context.Background(), "test-group", artifact, params)
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, "artifact-1", result.ArtifactID)
		assert.Equal(t, "New Artifact", result.Name)
	})

	t.Run("Invalid Artifact", func(t *testing.T) {
		mockResponse := models.CreateArtifactResponse{
			Artifact: models.ArtifactDetail{
				GroupID:     "test-group",
				ArtifactID:  "artifact-1",
				Name:        "New Artifact",
				Description: "Test Description",
			},
		}

		server := setupMockServer(t, http.StatusOK, mockResponse, "/groups/test-group/artifacts", http.MethodPost)
		defer server.Close()

		mockClient := &client.Client{BaseURL: server.URL, HTTPClient: server.Client()}
		api := apis.NewArtifactsAPI(mockClient)

		artifact := models.CreateArtifactRequest{
			ArtifactType: models.Json,
			FirstVersion: models.CreateVersionRequest{
				Version: "1.0.0",
				Content: models.CreateContentRequest{Content: "{\"key\":\"value\"}"},
			},
		}
		params := &models.CreateArtifactParams{IfExists: models.IfExistsCreate}

		result, err := api.CreateArtifact(context.Background(), "test-group", artifact, params)
		assert.Error(t, err)
		assert.Nil(t, result)
	})

	t.Run("Bad Request", func(t *testing.T) {
		errorResponse := models.APIError{Status: http.StatusBadRequest, Title: TitleBadRequest}
		server := setupMockServer(t, http.StatusBadRequest, errorResponse, "/groups/test-group/artifacts", http.MethodPost)
		defer server.Close()

		mockClient := &client.Client{BaseURL: server.URL, HTTPClient: server.Client()}
		api := apis.NewArtifactsAPI(mockClient)

		artifact := models.CreateArtifactRequest{
			ArtifactID:   stubArtifactId,
			ArtifactType: models.Json,
			FirstVersion: models.CreateVersionRequest{
				Version: "1.0.0",
				Content: models.CreateContentRequest{Content: "{\"key\":\"value\"}", ContentType: "application/json"},
			},
		}
		params := &models.CreateArtifactParams{IfExists: models.IfExistsCreate}

		result, err := api.CreateArtifact(context.Background(), "test-group", artifact, params)
		assert.Error(t, err)
		assert.Nil(t, result)
		assertAPIError(t, err, http.StatusBadRequest, TitleBadRequest)
	})

	t.Run("Conflict", func(t *testing.T) {
		errorResponse := models.APIError{Status: http.StatusConflict, Title: TitleConflict}
		server := setupMockServer(t, http.StatusConflict, errorResponse, "/groups/test-group/artifacts", http.MethodPost)
		defer server.Close()

		mockClient := &client.Client{BaseURL: server.URL, HTTPClient: server.Client()}
		api := apis.NewArtifactsAPI(mockClient)

		artifact := models.CreateArtifactRequest{
			ArtifactID:   stubArtifactId,
			ArtifactType: models.Json,
			FirstVersion: models.CreateVersionRequest{
				Version: "1.0.0",
				Content: models.CreateContentRequest{Content: "{\"key\":\"value\"}", ContentType: "application/json"},
			},
		}
		params := &models.CreateArtifactParams{IfExists: models.IfExistsCreate}

		result, err := api.CreateArtifact(context.Background(), "test-group", artifact, params)
		assert.Error(t, err)
		assert.Nil(t, result)
		assertAPIError(t, err, http.StatusConflict, TitleConflict)
	})

	t.Run("Internal Server Error", func(t *testing.T) {
		errorResponse := models.APIError{Status: http.StatusInternalServerError, Title: TitleInternalServerError}
		server := setupMockServer(t, http.StatusInternalServerError, errorResponse, "/groups/test-group/artifacts", http.MethodPost)
		defer server.Close()

		mockClient := &client.Client{BaseURL: server.URL, HTTPClient: server.Client()}
		api := apis.NewArtifactsAPI(mockClient)

		artifact := models.CreateArtifactRequest{
			ArtifactID:   stubArtifactId,
			ArtifactType: models.Json,
			FirstVersion: models.CreateVersionRequest{
				Version: "1.0.0",
				Content: models.CreateContentRequest{Content: "{\"key\":\"value\"}", ContentType: "application/json"},
			},
		}
		params := &models.CreateArtifactParams{IfExists: models.IfExistsCreate}

		result, err := api.CreateArtifact(context.Background(), "test-group", artifact, params)
		assert.Error(t, err)
		assert.Nil(t, result)
		assertAPIError(t, err, http.StatusInternalServerError, TitleInternalServerError)
	})
}

func TestArtifactsAPI_ListArtifactRules(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mockReferences := []models.Rule{models.RuleValidity, models.RuleCompatibility}

		server := setupMockServer(t, http.StatusOK, mockReferences,
			fmt.Sprintf("/groups/%s/artifacts/%s/rules", stubGroupId, stubArtifactId), http.MethodGet)
		defer server.Close()

		mockClient := &client.Client{BaseURL: server.URL, HTTPClient: server.Client()}
		api := apis.NewArtifactsAPI(mockClient)

		result, err := api.ListArtifactRules(context.Background(), stubGroupId, stubArtifactId)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Len(t, result, 2)
		assert.Contains(t, result, models.RuleValidity)
		assert.Contains(t, result, models.RuleCompatibility)
	})

	t.Run("Not Found", func(t *testing.T) {
		errorResponse := models.APIError{Status: http.StatusNotFound, Title: TitleNotFound}

		server := setupMockServer(t, http.StatusNotFound, errorResponse,
			fmt.Sprintf("/groups/%s/artifacts/%s/rules", stubGroupId, stubArtifactId), http.MethodGet)
		defer server.Close()

		mockClient := &client.Client{BaseURL: server.URL, HTTPClient: server.Client()}
		api := apis.NewArtifactsAPI(mockClient)

		result, err := api.ListArtifactRules(context.Background(), stubGroupId, stubArtifactId)

		assert.Error(t, err)
		assert.Nil(t, result)
		assertAPIError(t, err, http.StatusNotFound, TitleNotFound)
	})

	t.Run("Internal Server Error", func(t *testing.T) {
		errorResponse := models.APIError{Status: http.StatusInternalServerError, Title: TitleInternalServerError}

		server := setupMockServer(t, http.StatusInternalServerError, errorResponse,
			fmt.Sprintf("/groups/%s/artifacts/%s/rules", stubGroupId, stubArtifactId), http.MethodGet)
		defer server.Close()

		mockClient := &client.Client{BaseURL: server.URL, HTTPClient: server.Client()}
		api := apis.NewArtifactsAPI(mockClient)

		result, err := api.ListArtifactRules(context.Background(), stubGroupId, stubArtifactId)

		assert.Error(t, err)
		assert.Nil(t, result)
		assertAPIError(t, err, http.StatusInternalServerError, TitleInternalServerError)
	})
}

func TestArtifactsAPI_CreateArtifactRule(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		server := setupMockServer(t, http.StatusNoContent, nil,
			fmt.Sprintf("/groups/%s/artifacts/%s/rules", stubGroupId, stubArtifactId), http.MethodPost)
		defer server.Close()

		mockClient := &client.Client{BaseURL: server.URL, HTTPClient: server.Client()}
		api := apis.NewArtifactsAPI(mockClient)

		err := api.CreateArtifactRule(context.Background(), stubGroupId, stubArtifactId, models.RuleValidity, models.ValidityLevelFull)
		assert.NoError(t, err)
	})

	t.Run("BadRequest", func(t *testing.T) {
		errorResponse := models.APIError{Status: http.StatusBadRequest, Title: TitleBadRequest}
		server := setupMockServer(t, http.StatusBadRequest, errorResponse,
			fmt.Sprintf("/groups/%s/artifacts/%s/rules", stubGroupId, stubArtifactId), http.MethodPost)
		defer server.Close()

		mockClient := &client.Client{BaseURL: server.URL, HTTPClient: server.Client()}
		api := apis.NewArtifactsAPI(mockClient)

		err := api.CreateArtifactRule(context.Background(), stubGroupId, stubArtifactId, models.RuleValidity, models.ValidityLevelFull)
		assert.Error(t, err)
		assertAPIError(t, err, http.StatusBadRequest, TitleBadRequest)
	})

	t.Run("Conflict", func(t *testing.T) {
		errorResponse := models.APIError{Status: http.StatusConflict, Title: TitleConflict}
		server := setupMockServer(t, http.StatusConflict, errorResponse,
			fmt.Sprintf("/groups/%s/artifacts/%s/rules", stubGroupId, stubArtifactId), http.MethodPost)
		defer server.Close()

		mockClient := &client.Client{BaseURL: server.URL, HTTPClient: server.Client()}
		api := apis.NewArtifactsAPI(mockClient)

		err := api.CreateArtifactRule(context.Background(), stubGroupId, stubArtifactId, models.RuleValidity, models.ValidityLevelFull)
		assert.Error(t, err)
		assertAPIError(t, err, http.StatusConflict, TitleConflict)
	})

	t.Run("NotFound", func(t *testing.T) {
		errorResponse := models.APIError{Status: http.StatusNotFound, Title: TitleNotFound}
		server := setupMockServer(t, http.StatusNotFound, errorResponse,
			fmt.Sprintf("/groups/%s/artifacts/%s/rules", stubGroupId, stubArtifactId), http.MethodPost)
		defer server.Close()

		mockClient := &client.Client{BaseURL: server.URL, HTTPClient: server.Client()}
		api := apis.NewArtifactsAPI(mockClient)

		err := api.CreateArtifactRule(context.Background(), stubGroupId, stubArtifactId, models.RuleValidity, models.ValidityLevelFull)
		assert.Error(t, err)
		assertAPIError(t, err, http.StatusNotFound, TitleNotFound)
	})

	t.Run("InternalServerError", func(t *testing.T) {
		errorResponse := models.APIError{Status: http.StatusInternalServerError, Title: TitleInternalServerError}
		server := setupMockServer(t, http.StatusInternalServerError, errorResponse,
			fmt.Sprintf("/groups/%s/artifacts/%s/rules", stubGroupId, stubArtifactId), http.MethodPost)
		defer server.Close()

		mockClient := &client.Client{BaseURL: server.URL, HTTPClient: server.Client()}
		api := apis.NewArtifactsAPI(mockClient)

		err := api.CreateArtifactRule(context.Background(), stubGroupId, stubArtifactId, models.RuleValidity, models.ValidityLevelFull)
		assert.Error(t, err)
		assertAPIError(t, err, http.StatusInternalServerError, TitleInternalServerError)
	})
}

func TestArtifactsAPI_DeleteAllArtifactRule(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		server := setupMockServer(t, http.StatusNoContent, nil,
			fmt.Sprintf("/groups/%s/artifacts/%s/rules", stubGroupId, stubArtifactId), http.MethodDelete)
		defer server.Close()

		mockClient := &client.Client{BaseURL: server.URL, HTTPClient: server.Client()}
		api := apis.NewArtifactsAPI(mockClient)

		err := api.DeleteAllArtifactRule(context.Background(), stubGroupId, stubArtifactId)
		assert.NoError(t, err)
	})

	t.Run("NotFound", func(t *testing.T) {
		errorResponse := models.APIError{Status: http.StatusNotFound, Title: TitleNotFound}
		server := setupMockServer(t, http.StatusNotFound, errorResponse,
			fmt.Sprintf("/groups/%s/artifacts/%s/rules", stubGroupId, stubArtifactId), http.MethodDelete)
		defer server.Close()

		mockClient := &client.Client{BaseURL: server.URL, HTTPClient: server.Client()}
		api := apis.NewArtifactsAPI(mockClient)

		err := api.DeleteAllArtifactRule(context.Background(), stubGroupId, stubArtifactId)
		assert.Error(t, err)
		assertAPIError(t, err, http.StatusNotFound, TitleNotFound)
	})

	t.Run("InternalServerError", func(t *testing.T) {
		errorResponse := models.APIError{Status: http.StatusInternalServerError, Title: TitleInternalServerError}
		server := setupMockServer(t, http.StatusInternalServerError, errorResponse,
			fmt.Sprintf("/groups/%s/artifacts/%s/rules", stubGroupId, stubArtifactId), http.MethodDelete)
		defer server.Close()

		mockClient := &client.Client{BaseURL: server.URL, HTTPClient: server.Client()}
		api := apis.NewArtifactsAPI(mockClient)

		err := api.DeleteAllArtifactRule(context.Background(), stubGroupId, stubArtifactId)
		assert.Error(t, err)
		assertAPIError(t, err, http.StatusInternalServerError, TitleInternalServerError)
	})
}

func TestArtifactsAPI_GetArtifactRule(t *testing.T) {
	mockRule := models.RuleValidity
	successResponse := models.RuleResponse{
		RuleType: models.RuleValidity,
		Config:   models.ValidityLevelFull,
	}

	t.Run("Success", func(t *testing.T) {
		server := setupMockServer(t, http.StatusOK, successResponse,
			fmt.Sprintf("/groups/%s/artifacts/%s/rules/%s", stubGroupId, stubArtifactId, mockRule), http.MethodGet)
		defer server.Close()

		mockClient := &client.Client{BaseURL: server.URL, HTTPClient: server.Client()}
		api := apis.NewArtifactsAPI(mockClient)

		result, err := api.GetArtifactRule(context.Background(), stubGroupId, stubArtifactId, mockRule)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, models.ValidityLevelFull, result)
	})

	t.Run("NotFound", func(t *testing.T) {
		errorResponse := models.APIError{Status: http.StatusNotFound, Title: TitleNotFound}
		server := setupMockServer(t, http.StatusNotFound, errorResponse,
			fmt.Sprintf("/groups/%s/artifacts/%s/rules/%s", stubGroupId, stubArtifactId, mockRule), http.MethodGet)
		defer server.Close()

		mockClient := &client.Client{BaseURL: server.URL, HTTPClient: server.Client()}
		api := apis.NewArtifactsAPI(mockClient)

		result, err := api.GetArtifactRule(context.Background(), stubGroupId, stubArtifactId, mockRule)

		assert.Error(t, err)
		assert.Empty(t, result)
		assertAPIError(t, err, http.StatusNotFound, TitleNotFound)
	})

	t.Run("InternalServerError", func(t *testing.T) {
		errorResponse := models.APIError{Status: http.StatusInternalServerError, Title: TitleInternalServerError}
		server := setupMockServer(t, http.StatusInternalServerError, errorResponse,
			fmt.Sprintf("/groups/%s/artifacts/%s/rules/%s", stubGroupId, stubArtifactId, mockRule), http.MethodGet)
		defer server.Close()

		mockClient := &client.Client{BaseURL: server.URL, HTTPClient: server.Client()}
		api := apis.NewArtifactsAPI(mockClient)

		result, err := api.GetArtifactRule(context.Background(), stubGroupId, stubArtifactId, mockRule)

		assert.Error(t, err)
		assert.Empty(t, result)
		assertAPIError(t, err, http.StatusInternalServerError, TitleInternalServerError)
	})
}

func TestArtifactsAPI_UpdateArtifactRule(t *testing.T) {
	mockRule := models.RuleValidity
	successResponse := models.RuleResponse{
		RuleType: mockRule,
		Config:   models.ValidityLevelFull,
	}

	t.Run("Success", func(t *testing.T) {
		server := setupMockServer(t, http.StatusOK, successResponse,
			fmt.Sprintf("/groups/%s/artifacts/%s/rules/%s", stubGroupId, stubArtifactId, mockRule), http.MethodPut)
		defer server.Close()

		mockClient := &client.Client{BaseURL: server.URL, HTTPClient: server.Client()}
		api := apis.NewArtifactsAPI(mockClient)

		err := api.UpdateArtifactRule(context.Background(), stubGroupId, stubArtifactId, mockRule, models.ValidityLevelFull)
		assert.NoError(t, err)
	})

	t.Run("NotFound", func(t *testing.T) {
		errorResponse := models.APIError{Status: http.StatusNotFound, Title: TitleNotFound}
		server := setupMockServer(t, http.StatusNotFound, errorResponse,
			fmt.Sprintf("/groups/%s/artifacts/%s/rules/%s", stubGroupId, stubArtifactId, mockRule), http.MethodPut)
		defer server.Close()

		mockClient := &client.Client{BaseURL: server.URL, HTTPClient: server.Client()}
		api := apis.NewArtifactsAPI(mockClient)

		err := api.UpdateArtifactRule(context.Background(), stubGroupId, stubArtifactId, mockRule, models.ValidityLevelFull)
		assert.Error(t, err)
		assertAPIError(t, err, http.StatusNotFound, TitleNotFound)
	})

	t.Run("InternalServerError", func(t *testing.T) {
		errorResponse := models.APIError{Status: http.StatusInternalServerError, Title: TitleInternalServerError}
		server := setupMockServer(t, http.StatusInternalServerError, errorResponse,
			fmt.Sprintf("/groups/%s/artifacts/%s/rules/%s", stubGroupId, stubArtifactId, mockRule), http.MethodPut)
		defer server.Close()

		mockClient := &client.Client{BaseURL: server.URL, HTTPClient: server.Client()}
		api := apis.NewArtifactsAPI(mockClient)

		err := api.UpdateArtifactRule(context.Background(), stubGroupId, stubArtifactId, mockRule, models.ValidityLevelFull)
		assert.Error(t, err)
		assertAPIError(t, err, http.StatusInternalServerError, TitleInternalServerError)
	})
}

func TestArtifactsAPI_DeleteArtifactRule(t *testing.T) {
	mockRule := models.RuleValidity

	t.Run("Success", func(t *testing.T) {
		server := setupMockServer(t, http.StatusNoContent, nil,
			fmt.Sprintf("/groups/%s/artifacts/%s/rules/%s", stubGroupId, stubArtifactId, mockRule), http.MethodDelete)
		defer server.Close()

		mockClient := &client.Client{BaseURL: server.URL, HTTPClient: server.Client()}
		api := apis.NewArtifactsAPI(mockClient)

		err := api.DeleteArtifactRule(context.Background(), stubGroupId, stubArtifactId, mockRule)
		assert.NoError(t, err)
	})

	t.Run("NotFound", func(t *testing.T) {
		errorResponse := models.APIError{Status: http.StatusNotFound, Title: TitleNotFound}

		server := setupMockServer(t, http.StatusNotFound, errorResponse,
			fmt.Sprintf("/groups/%s/artifacts/%s/rules/%s", stubGroupId, stubArtifactId, mockRule), http.MethodDelete)
		defer server.Close()

		mockClient := &client.Client{BaseURL: server.URL, HTTPClient: server.Client()}
		api := apis.NewArtifactsAPI(mockClient)

		err := api.DeleteArtifactRule(context.Background(), stubGroupId, stubArtifactId, mockRule)
		assert.Error(t, err)
		assertAPIError(t, err, http.StatusNotFound, TitleNotFound)
	})

	t.Run("InternalServerError", func(t *testing.T) {
		errorResponse := models.APIError{Status: http.StatusInternalServerError, Title: TitleInternalServerError}

		server := setupMockServer(t, http.StatusInternalServerError, errorResponse,
			fmt.Sprintf("/groups/%s/artifacts/%s/rules/%s", stubGroupId, stubArtifactId, mockRule), http.MethodDelete)
		defer server.Close()

		mockClient := &client.Client{BaseURL: server.URL, HTTPClient: server.Client()}
		api := apis.NewArtifactsAPI(mockClient)

		err := api.DeleteArtifactRule(context.Background(), stubGroupId, stubArtifactId, mockRule)
		assert.Error(t, err)
		assertAPIError(t, err, http.StatusInternalServerError, TitleInternalServerError)
	})
}

/***********************/
/***** Integration *****/
/***********************/
func TestArtifactsAPIIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	artifactsAPI := setupArtifactAPIClient()

	// Clean up before and after tests
	t.Cleanup(func() { cleanup(t, artifactsAPI) })
	cleanup(t, artifactsAPI)

	ctx := context.Background()

	// Test CreateArtifact
	t.Run("CreateArtifact", func(t *testing.T) {
		artifact := models.CreateArtifactRequest{
			ArtifactType: models.Json,
			ArtifactID:   stubArtifactId,
			Name:         stubArtifactId,
			FirstVersion: models.CreateVersionRequest{
				Version: "1.0.0",
				Content: models.CreateContentRequest{
					Content:     stubArtifactContent,
					ContentType: "application/json",
				},
			},
		}

		params := &models.CreateArtifactParams{
			IfExists: models.IfExistsFail,
		}

		resp, err := artifactsAPI.CreateArtifact(ctx, stubGroupId, artifact, params)
		assert.NoError(t, err)
		assert.Equal(t, stubGroupId, resp.GroupID)
		assert.Equal(t, stubArtifactId, resp.Name)
	})

	// Test SearchArtifacts
	t.Run("SearchArtifacts", func(t *testing.T) {
		params := &models.SearchArtifactsParams{
			Name: stubArtifactId,
		}
		resp, err := artifactsAPI.SearchArtifacts(ctx, params)
		assert.NoError(t, err)
		assert.GreaterOrEqual(t, len(resp), 1)
	})

	// Test ListArtifactReferences
	t.Run("ListArtifactReferences", func(t *testing.T) {
		contentID := int64(12345) // Replace with a valid content ID for your tests
		_, err := artifactsAPI.ListArtifactReferences(ctx, contentID)
		assert.Error(t, err) // Expect an error since no content ID exists
	})

	// Test ListArtifactReferencesByGlobalID
	t.Run("ListArtifactReferencesByGlobalID", func(t *testing.T) {
		globalID := int64(12345) // Replace with a valid global ID for your tests
		params := &models.ListArtifactReferencesByGlobalIDParams{}
		_, err := artifactsAPI.ListArtifactReferencesByGlobalID(ctx, globalID, params)
		assert.Error(t, err) // Expect an error since no global ID exists
	})

	// Test ListArtifactReferencesByHash
	t.Run("ListArtifactReferencesByHash", func(t *testing.T) {
		contentHash := "invalidhash" // Replace with a valid content hash for your tests
		_, err := artifactsAPI.ListArtifactReferencesByHash(ctx, contentHash)
		assert.Error(t, err) // Expect an error since no hash exists
	})

	// Test ListArtifactsInGroup
	t.Run("ListArtifactsInGroup", func(t *testing.T) {
		params := &models.ListArtifactsInGroupParams{}
		resp, err := artifactsAPI.ListArtifactsInGroup(ctx, stubGroupId, params)
		assert.NoError(t, err)
		assert.GreaterOrEqual(t, len(resp.Artifacts), 1)
	})

	// Test GetArtifactContentByHash
	t.Run("GetArtifactContentByHash", func(t *testing.T) {
		contentHash := "invalidhash" // Replace with a valid content hash for your tests
		_, err := artifactsAPI.GetArtifactContentByHash(ctx, contentHash)
		assert.Error(t, err) // Expect an error since no hash exists
	})

	// Test GetArtifactContentByID
	t.Run("GetArtifactContentByID", func(t *testing.T) {
		contentID := int64(12345) // Replace with a valid content ID for your tests
		_, err := artifactsAPI.GetArtifactContentByID(ctx, contentID)
		assert.Error(t, err) // Expect an error since no content ID exists
	})

	// Test DeleteArtifactsInGroup
	t.Run("DeleteArtifactsInGroup", func(t *testing.T) {
		err := artifactsAPI.DeleteArtifactsInGroup(ctx, stubGroupId)
		assert.NoError(t, err)
	})

	// Test DeleteArtifact
	t.Run("DeleteArtifact", func(t *testing.T) {

		// Re-create the artifact
		artifact := models.CreateArtifactRequest{
			ArtifactType: models.Json,
			ArtifactID:   stubArtifactId,
			Name:         stubArtifactId,
			FirstVersion: models.CreateVersionRequest{
				Version: "1.0.0",
				Content: models.CreateContentRequest{
					Content:     stubArtifactContent,
					ContentType: "application/json",
				},
			},
		}
		params := &models.CreateArtifactParams{
			IfExists: models.IfExistsFail,
		}

		resp, err := artifactsAPI.CreateArtifact(ctx, stubGroupId, artifact, params)
		assert.NoError(t, err)
		assert.Equal(t, stubGroupId, resp.GroupID)
		assert.Equal(t, stubArtifactId, resp.Name)

		// Delete the artifact
		err = artifactsAPI.DeleteArtifact(ctx, stubGroupId, stubArtifactId)
		assert.NoError(t, err)
	})

	// Test Artifact rules
	t.Run("AllArtifactRules", func(t *testing.T) {
		randomArtifactID, err := generateArtifactForTest(ctx, artifactsAPI)
		assert.NoError(t, err)

		// List Artifact rules
		rules, err := artifactsAPI.ListArtifactRules(ctx, stubGroupId, randomArtifactID)
		assert.NoError(t, err)
		assert.Len(t, rules, 0)

		// Create a rule
		err = artifactsAPI.CreateArtifactRule(ctx, stubGroupId, randomArtifactID, models.RuleValidity, models.ValidityLevelFull)
		assert.NoError(t, err)

		// List artifact rules
		rules, err = artifactsAPI.ListArtifactRules(ctx, stubGroupId, randomArtifactID)
		assert.NoError(t, err)
		assert.Len(t, rules, 1)
		assert.Equal(t, models.RuleValidity, rules[0])

		// Get the rule
		rule, err := artifactsAPI.GetArtifactRule(ctx, stubGroupId, randomArtifactID, models.RuleValidity)
		assert.NoError(t, err)
		assert.Equal(t, models.ValidityLevelFull, rule)

		// Update the rule
		err = artifactsAPI.UpdateArtifactRule(ctx, stubGroupId, randomArtifactID, models.RuleValidity, models.ValidityLevelSyntaxOnly)
		assert.NoError(t, err)

		// Get the rule
		rule, err = artifactsAPI.GetArtifactRule(ctx, stubGroupId, randomArtifactID, models.RuleValidity)
		assert.NoError(t, err)
		assert.Equal(t, models.ValidityLevelSyntaxOnly, rule)

		// Delete the rule
		err = artifactsAPI.DeleteArtifactRule(ctx, stubGroupId, randomArtifactID, models.RuleValidity)
		assert.NoError(t, err)

		// List artifact rules
		rules, err = artifactsAPI.ListArtifactRules(ctx, stubGroupId, randomArtifactID)
		assert.NoError(t, err)
		assert.Len(t, rules, 0)

		// Create three rules
		err = artifactsAPI.CreateArtifactRule(ctx, stubGroupId, randomArtifactID, models.RuleValidity, models.ValidityLevelFull)
		assert.NoError(t, err)
		err = artifactsAPI.CreateArtifactRule(ctx, stubGroupId, randomArtifactID, models.RuleCompatibility, models.CompatibilityLevelFull)
		assert.NoError(t, err)
		err = artifactsAPI.CreateArtifactRule(ctx, stubGroupId, randomArtifactID, models.RuleIntegrity, models.IntegrityLevelFull)
		assert.NoError(t, err)

		// List artifact rules
		rules, err = artifactsAPI.ListArtifactRules(ctx, stubGroupId, randomArtifactID)
		assert.NoError(t, err)
		assert.Len(t, rules, 3)

		// Delete all rules
		err = artifactsAPI.DeleteAllArtifactRule(ctx, stubGroupId, randomArtifactID)
		assert.NoError(t, err)

		// List artifact rules
		rules, err = artifactsAPI.ListArtifactRules(ctx, stubGroupId, randomArtifactID)
		assert.NoError(t, err)
		assert.Len(t, rules, 0)
	})
}

func setupHTTPClient() *client.Client {
	baseURL := os.Getenv("APICURIO_BASE_URL")
	if baseURL == "" {
		baseURL = DefaultBaseURL
	}
	httpClient := &http.Client{Timeout: 10 * time.Second}
	apiClient := client.NewClient(baseURL, client.WithHTTPClient(httpClient))
	return apiClient
}

func setupArtifactAPIClient() *apis.ArtifactsAPI {
	apiClient := setupHTTPClient()
	return apis.NewArtifactsAPI(apiClient)
}

func cleanup(t *testing.T, artifactsAPI *apis.ArtifactsAPI) {
	ctx := context.Background()
	err := artifactsAPI.DeleteArtifactsInGroup(ctx, stubGroupId)
	if err != nil {
		var APIError *models.APIError
		if errors.As(err, &APIError) && APIError.Status == 404 {
			return
		}
		t.Fatalf("Failed to clean up artifacts: %v", err)
	}
}
