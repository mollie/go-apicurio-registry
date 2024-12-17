package apis_test

import (
	"context"
	"encoding/json"
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

const (
	DefaultBaseURL = "http://localhost:9080/apis/registry/v3"
	groupID        = "test-group"
	artifactID     = "test-artifact"
)

var (
	stubArtifactContent = `{"type": "record", "name": "Test", "fields": [{"name": "field1", "type": "string"}]}`
	stubArtifactId      = "test-artifact"
	stubGroupId         = "test-group"
	stubBranchID        = "test-branch"
	stubVersionID       = "1.0.0"
	stubVersionID2      = "2.0.0"
)

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
	err := artifactsAPI.DeleteArtifactsInGroup(ctx, groupID)
	if err != nil {
		var APIError *models.APIError
		if errors.As(err, &APIError) && APIError.Status == 404 {
			return
		}
		t.Fatalf("Failed to clean up artifacts: %v", err)
	}
}

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
		params := models.GetArtifactByGlobalIDParams{
			ReturnArtifactType: true,
		}
		result, err := api.GetArtifactByGlobalID(context.Background(), 1, &params)
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, stubArtifactContent, result.Content)
		assert.Equal(t, models.Json, result.ArtifactType)
	})

	t.Run("Not Found", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusNotFound)
			err := json.NewEncoder(w).Encode(models.APIError{Status: http.StatusNotFound, Title: TitleNotFound})
			assert.NoError(t, err)
		}))
		defer server.Close()

		mockClient := &client.Client{BaseURL: server.URL, HTTPClient: server.Client()}
		api := apis.NewArtifactsAPI(mockClient)

		result, err := api.GetArtifactByGlobalID(context.Background(), 1, nil)
		assert.Error(t, err)
		assert.Nil(t, result)

		var apiErr *models.APIError
		ok := errors.As(err, &apiErr)
		assert.True(t, ok)
		assert.Equal(t, http.StatusNotFound, apiErr.Status)
		assert.Equal(t, TitleNotFound, apiErr.Title)
	})

	t.Run("Internal Server Error", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
			err := json.NewEncoder(w).Encode(models.APIError{Status: http.StatusInternalServerError, Title: TitleInternalServerError})
			assert.NoError(t, err)
		}))
		defer server.Close()

		mockClient := &client.Client{BaseURL: server.URL, HTTPClient: server.Client()}
		api := apis.NewArtifactsAPI(mockClient)

		result, err := api.GetArtifactByGlobalID(context.Background(), 1, nil)
		assert.Error(t, err)
		assert.Nil(t, result)

		var apiErr *models.APIError
		ok := errors.As(err, &apiErr)
		assert.True(t, ok)
		assert.Equal(t, http.StatusInternalServerError, apiErr.Status)
		assert.Equal(t, TitleInternalServerError, apiErr.Title)
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
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "/search/artifacts", r.URL.Path)
			assert.Equal(t, http.MethodGet, r.Method)

			w.WriteHeader(http.StatusOK)
			err := json.NewEncoder(w).Encode(mockResponse)
			assert.NoError(t, err)
		}))
		defer server.Close()

		mockClient := &client.Client{BaseURL: server.URL, HTTPClient: server.Client()}
		api := apis.NewArtifactsAPI(mockClient)

		params := &models.SearchArtifactsParams{Name: "Test Artifact"}
		result, err := api.SearchArtifacts(context.Background(), params)
		assert.NoError(t, err)
		assert.NotNil(t, result)
	})

	t.Run("Internal Server Error", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
			err := json.NewEncoder(w).Encode(models.APIError{Status: http.StatusInternalServerError, Title: TitleInternalServerError})
			assert.NoError(t, err)
		}))
		defer server.Close()

		mockClient := &client.Client{BaseURL: server.URL, HTTPClient: server.Client()}
		api := apis.NewArtifactsAPI(mockClient)

		params := &models.SearchArtifactsParams{}
		result, err := api.SearchArtifacts(context.Background(), params)
		assert.Error(t, err)
		assert.Nil(t, result)

		var apiErr *models.APIError
		ok := errors.As(err, &apiErr)
		assert.True(t, ok)
		assert.Equal(t, http.StatusInternalServerError, apiErr.Status)
		assert.Equal(t, TitleInternalServerError, apiErr.Title)
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
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "/search/artifacts", r.URL.Path)
			assert.Equal(t, http.MethodPost, r.Method)

			w.WriteHeader(http.StatusOK)
			err := json.NewEncoder(w).Encode(mockResponse)
			assert.NoError(t, err)
		}))
		defer server.Close()

		mockClient := &client.Client{BaseURL: server.URL, HTTPClient: server.Client()}
		api := apis.NewArtifactsAPI(mockClient)

		params := &models.SearchArtifactsByContentParams{Canonical: true}
		result, err := api.SearchArtifactsByContent(context.Background(), []byte("{\"key\":\"value\"}"), params)
		assert.NoError(t, err)
		assert.NotNil(t, result)
	})

	t.Run("Invalid Content", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusBadRequest)
			err := json.NewEncoder(w).Encode(models.APIError{Status: http.StatusBadRequest, Title: TitleBadRequest})
			assert.NoError(t, err)
		}))
		defer server.Close()

		mockClient := &client.Client{BaseURL: server.URL, HTTPClient: server.Client()}
		api := apis.NewArtifactsAPI(mockClient)

		params := &models.SearchArtifactsByContentParams{}
		result, err := api.SearchArtifactsByContent(context.Background(), []byte(""), params)
		assert.Error(t, err)
		assert.Nil(t, result)

		var apiErr *models.APIError
		ok := errors.As(err, &apiErr)
		assert.True(t, ok)
		assert.Equal(t, http.StatusBadRequest, apiErr.Status)
		assert.Equal(t, TitleBadRequest, apiErr.Title)
	})
}

func TestArtifactsAPI_ListArtifactReferences(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mockReferences := []models.ArtifactReference{
			{GroupID: "group-1", ArtifactID: "artifact-1", Version: "v1"},
		}
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Contains(t, r.URL.Path, "/references")
			assert.Equal(t, http.MethodGet, r.Method)

			w.WriteHeader(http.StatusOK)
			err := json.NewEncoder(w).Encode(mockReferences)
			assert.NoError(t, err)
		}))
		defer server.Close()

		mockClient := &client.Client{BaseURL: server.URL, HTTPClient: server.Client()}
		api := apis.NewArtifactsAPI(mockClient)

		result, err := api.ListArtifactReferences(context.Background(), 123)
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Len(t, *result, 1)
	})

	t.Run("Not Found", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusNotFound)
			err := json.NewEncoder(w).Encode(models.APIError{Status: http.StatusNotFound, Title: TitleNotFound})
			assert.NoError(t, err)
		}))
		defer server.Close()

		mockClient := &client.Client{BaseURL: server.URL, HTTPClient: server.Client()}
		api := apis.NewArtifactsAPI(mockClient)

		result, err := api.ListArtifactReferences(context.Background(), 123)
		assert.Error(t, err)
		assert.Nil(t, result)

		var apiErr *models.APIError
		ok := errors.As(err, &apiErr)
		assert.True(t, ok)
		assert.Equal(t, http.StatusNotFound, apiErr.Status)
		assert.Equal(t, TitleNotFound, apiErr.Title)
	})
}

func TestArtifactsAPI_ListArtifactReferencesByGlobalID(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mockReferences := []models.ArtifactReference{
			{GroupID: "group-1", ArtifactID: "artifact-1", Version: "v1"},
		}
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Contains(t, r.URL.Path, "ids/globalIds")
			assert.Equal(t, http.MethodGet, r.Method)

			w.WriteHeader(http.StatusOK)
			err := json.NewEncoder(w).Encode(mockReferences)
			assert.NoError(t, err)
		}))
		defer server.Close()

		mockClient := &client.Client{BaseURL: server.URL, HTTPClient: server.Client()}
		api := apis.NewArtifactsAPI(mockClient)

		params := &models.ListArtifactReferencesByGlobalIDParams{RefType: models.OutBound}
		result, err := api.ListArtifactReferencesByGlobalID(context.Background(), 123, params)
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Len(t, *result, 1)
	})

	t.Run("Internal Server Error", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
			err := json.NewEncoder(w).Encode(models.APIError{Status: http.StatusInternalServerError, Title: TitleInternalServerError})
			assert.NoError(t, err)
		}))
		defer server.Close()

		mockClient := &client.Client{BaseURL: server.URL, HTTPClient: server.Client()}
		api := apis.NewArtifactsAPI(mockClient)

		params := &models.ListArtifactReferencesByGlobalIDParams{}
		result, err := api.ListArtifactReferencesByGlobalID(context.Background(), 123, params)
		assert.Error(t, err)
		assert.Nil(t, result)

		var apiErr *models.APIError
		ok := errors.As(err, &apiErr)
		assert.True(t, ok)
		assert.Equal(t, http.StatusInternalServerError, apiErr.Status)
		assert.Equal(t, TitleInternalServerError, apiErr.Title)
	})
}

func TestArtifactsAPI_ListArtifactReferencesByHash(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mockReferences := []models.ArtifactReference{
			{GroupID: "group-1", ArtifactID: "artifact-1", Version: "v1"},
		}
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Contains(t, r.URL.Path, "/contentHashes")
			assert.Equal(t, http.MethodGet, r.Method)

			w.WriteHeader(http.StatusOK)
			err := json.NewEncoder(w).Encode(mockReferences)
			assert.NoError(t, err)
		}))
		defer server.Close()

		mockClient := &client.Client{BaseURL: server.URL, HTTPClient: server.Client()}
		api := apis.NewArtifactsAPI(mockClient)

		result, err := api.ListArtifactReferencesByHash(context.Background(), "hash-123")
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Len(t, *result, 1)
	})

	t.Run("Internal Server Error", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
			err := json.NewEncoder(w).Encode(models.APIError{Status: http.StatusInternalServerError, Title: TitleInternalServerError})
			assert.NoError(t, err)
		}))
		defer server.Close()

		mockClient := &client.Client{BaseURL: server.URL, HTTPClient: server.Client()}
		api := apis.NewArtifactsAPI(mockClient)

		result, err := api.ListArtifactReferencesByHash(context.Background(), "hash-123")
		assert.Error(t, err)
		assert.Nil(t, result)

		var apiErr *models.APIError
		ok := errors.As(err, &apiErr)
		assert.True(t, ok)
		assert.Equal(t, http.StatusInternalServerError, apiErr.Status)
		assert.Equal(t, TitleInternalServerError, apiErr.Title)
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
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Contains(t, r.URL.Path, "/groups/group-1/artifacts")
			assert.Equal(t, http.MethodGet, r.Method)

			w.WriteHeader(http.StatusOK)
			err := json.NewEncoder(w).Encode(mockResponse)
			assert.NoError(t, err)
		}))
		defer server.Close()

		mockClient := &client.Client{BaseURL: server.URL, HTTPClient: server.Client()}
		api := apis.NewArtifactsAPI(mockClient)

		params := &models.ListArtifactsInGroupParams{Limit: 10, Offset: 0, Order: "asc"}
		result, err := api.ListArtifactsInGroup(context.Background(), "group-1", params)
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Len(t, result.Artifacts, 1)
	})

	t.Run("Internal Server Error", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
			err := json.NewEncoder(w).Encode(models.APIError{Status: http.StatusInternalServerError, Title: TitleInternalServerError})
			assert.NoError(t, err)
		}))
		defer server.Close()

		mockClient := &client.Client{BaseURL: server.URL, HTTPClient: server.Client()}
		api := apis.NewArtifactsAPI(mockClient)

		params := &models.ListArtifactsInGroupParams{}
		result, err := api.ListArtifactsInGroup(context.Background(), "group-1", params)
		assert.Error(t, err)
		assert.Nil(t, result)

		var apiErr *models.APIError
		ok := errors.As(err, &apiErr)
		assert.True(t, ok)
		assert.Equal(t, http.StatusInternalServerError, apiErr.Status)
		assert.Equal(t, TitleInternalServerError, apiErr.Title)
	})
}

func TestArtifactsAPI_GetArtifactContentByHash(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mockContent := models.ArtifactContent{
			Content:      "{\"key\":\"value\"}",
			ArtifactType: models.Json,
		}
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Contains(t, r.URL.Path, "/contentHashes/hash-123")
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
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusNotFound)
			err := json.NewEncoder(w).Encode(models.APIError{Status: http.StatusNotFound, Title: TitleNotFound})
			assert.NoError(t, err)
		}))
		defer server.Close()

		mockClient := &client.Client{BaseURL: server.URL, HTTPClient: server.Client()}
		api := apis.NewArtifactsAPI(mockClient)

		result, err := api.GetArtifactContentByHash(context.Background(), "hash-123")
		assert.Error(t, err)
		assert.Nil(t, result)

		var apiErr *models.APIError
		ok := errors.As(err, &apiErr)
		assert.True(t, ok)
		assert.Equal(t, http.StatusNotFound, apiErr.Status)
		assert.Equal(t, TitleNotFound, apiErr.Title)
	})

	t.Run("Internal Server Error", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
			err := json.NewEncoder(w).Encode(models.APIError{Status: http.StatusInternalServerError, Title: TitleInternalServerError})
			assert.NoError(t, err)
		}))
		defer server.Close()

		mockClient := &client.Client{BaseURL: server.URL, HTTPClient: server.Client()}
		api := apis.NewArtifactsAPI(mockClient)

		result, err := api.GetArtifactContentByHash(context.Background(), "hash-123")
		assert.Error(t, err)
		assert.Nil(t, result)

		var apiErr *models.APIError
		ok := errors.As(err, &apiErr)
		assert.True(t, ok)
		assert.Equal(t, http.StatusInternalServerError, apiErr.Status)
		assert.Equal(t, TitleInternalServerError, apiErr.Title)
	})
}

func TestArtifactsAPI_GetArtifactContentByID(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mockContent := models.ArtifactContent{
			Content:      "{\"key\":\"value\"}",
			ArtifactType: models.Json,
		}
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Contains(t, r.URL.Path, "/contentIds/123")
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
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusNotFound)
			err := json.NewEncoder(w).Encode(models.APIError{Status: http.StatusNotFound, Title: TitleNotFound})
			assert.NoError(t, err)
		}))
		defer server.Close()

		mockClient := &client.Client{BaseURL: server.URL, HTTPClient: server.Client()}
		api := apis.NewArtifactsAPI(mockClient)

		result, err := api.GetArtifactContentByID(context.Background(), 123)
		assert.Error(t, err)
		assert.Nil(t, result)

		var apiErr *models.APIError
		ok := errors.As(err, &apiErr)
		assert.True(t, ok)
		assert.Equal(t, http.StatusNotFound, apiErr.Status)
		assert.Equal(t, TitleNotFound, apiErr.Title)
	})

	t.Run("Internal Server Error", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
			err := json.NewEncoder(w).Encode(models.APIError{Status: http.StatusInternalServerError, Title: TitleInternalServerError})
			assert.NoError(t, err)
		}))
		defer server.Close()

		mockClient := &client.Client{BaseURL: server.URL, HTTPClient: server.Client()}
		api := apis.NewArtifactsAPI(mockClient)

		result, err := api.GetArtifactContentByID(context.Background(), 123)
		assert.Error(t, err)
		assert.Nil(t, result)

		var apiErr *models.APIError
		ok := errors.As(err, &apiErr)
		assert.True(t, ok)
		assert.Equal(t, http.StatusInternalServerError, apiErr.Status)
		assert.Equal(t, TitleInternalServerError, apiErr.Title)
	})
}

func TestArtifactsAPI_DeleteArtifactsInGroup(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusNoContent)
		}))
		defer server.Close()

		mockClient := &client.Client{BaseURL: server.URL, HTTPClient: server.Client()}
		api := apis.NewArtifactsAPI(mockClient)

		err := api.DeleteArtifactsInGroup(context.Background(), "test-group")
		assert.NoError(t, err)
	})

	t.Run("Internal Server Error", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
			err := json.NewEncoder(w).Encode(models.APIError{Status: http.StatusInternalServerError, Title: TitleInternalServerError})
			assert.NoError(t, err)
		}))
		defer server.Close()

		mockClient := &client.Client{BaseURL: server.URL, HTTPClient: server.Client()}
		api := apis.NewArtifactsAPI(mockClient)

		err := api.DeleteArtifactsInGroup(context.Background(), "test-group")
		assert.Error(t, err)

		var apiErr *models.APIError
		ok := errors.As(err, &apiErr)
		assert.True(t, ok)
		assert.Equal(t, http.StatusInternalServerError, apiErr.Status)
		assert.Equal(t, TitleInternalServerError, apiErr.Title)
	})
}

func TestArtifactsAPI_DeleteArtifact(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusNoContent)
		}))
		defer server.Close()

		mockClient := &client.Client{BaseURL: server.URL, HTTPClient: server.Client()}
		api := apis.NewArtifactsAPI(mockClient)

		err := api.DeleteArtifact(context.Background(), "test-group", "artifact-1")
		assert.NoError(t, err)
	})

	t.Run("Not Found", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusNotFound)
			err := json.NewEncoder(w).Encode(models.APIError{Status: http.StatusNotFound, Title: TitleNotFound})
			assert.NoError(t, err)
		}))
		defer server.Close()

		mockClient := &client.Client{BaseURL: server.URL, HTTPClient: server.Client()}
		api := apis.NewArtifactsAPI(mockClient)

		err := api.DeleteArtifact(context.Background(), "test-group", "artifact-1")
		assert.Error(t, err)

		var apiErr *models.APIError
		ok := errors.As(err, &apiErr)
		assert.True(t, ok)
		assert.Equal(t, http.StatusNotFound, apiErr.Status)
		assert.Equal(t, TitleNotFound, apiErr.Title)
	})

	t.Run("Not Allowed", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusMethodNotAllowed)
			err := json.NewEncoder(w).Encode(models.APIError{Status: http.StatusMethodNotAllowed, Title: TitleMethodNotAllowed})
			assert.NoError(t, err)
		}))
		defer server.Close()

		mockClient := &client.Client{BaseURL: server.URL, HTTPClient: server.Client()}
		api := apis.NewArtifactsAPI(mockClient)

		err := api.DeleteArtifact(context.Background(), "test-group", "artifact-1")
		assert.Error(t, err)

		var apiErr *models.APIError
		ok := errors.As(err, &apiErr)
		assert.True(t, ok)
		assert.Equal(t, http.StatusMethodNotAllowed, apiErr.Status)
		assert.Equal(t, TitleMethodNotAllowed, apiErr.Title)
	})

	t.Run("Internal Server Error", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
			err := json.NewEncoder(w).Encode(models.APIError{Status: http.StatusInternalServerError, Title: TitleInternalServerError})
			assert.NoError(t, err)
		}))
		defer server.Close()

		mockClient := &client.Client{BaseURL: server.URL, HTTPClient: server.Client()}
		api := apis.NewArtifactsAPI(mockClient)

		err := api.DeleteArtifact(context.Background(), "test-group", "artifact-1")
		assert.Error(t, err)

		var apiErr *models.APIError
		ok := errors.As(err, &apiErr)
		assert.True(t, ok)
		assert.Equal(t, http.StatusInternalServerError, apiErr.Status)
		assert.Equal(t, TitleInternalServerError, apiErr.Title)
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
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Contains(t, r.URL.Path, "/artifacts")
			assert.Equal(t, http.MethodPost, r.Method)

			w.WriteHeader(http.StatusOK)
			err := json.NewEncoder(w).Encode(mockResponse)
			assert.NoError(t, err)
		}))
		defer server.Close()

		mockClient := &client.Client{BaseURL: server.URL, HTTPClient: server.Client()}
		api := apis.NewArtifactsAPI(mockClient)

		artifact := models.CreateArtifactRequest{
			ArtifactType: models.Json,
			FirstVersion: models.CreateVersionRequest{
				Version: "1.0.0",
				Content: models.CreateContentRequest{
					Content: "{\"key\":\"value\"}",
				},
			},
		}
		params := &models.CreateArtifactParams{
			IfExists: models.IfExistsCreate,
		}
		result, err := api.CreateArtifact(context.Background(), "test-group", artifact, params)
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, "artifact-1", result.ArtifactID)
	})

	t.Run("Bad Request", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusBadRequest)
			err := json.NewEncoder(w).Encode(models.APIError{Status: http.StatusBadRequest, Title: TitleBadRequest})
			assert.NoError(t, err)
		}))
		defer server.Close()

		mockClient := &client.Client{BaseURL: server.URL, HTTPClient: server.Client()}
		api := apis.NewArtifactsAPI(mockClient)

		artifact := models.CreateArtifactRequest{
			ArtifactType: models.Json,
			FirstVersion: models.CreateVersionRequest{
				Version: "1.0.0",
				Content: models.CreateContentRequest{
					Content: "{\"key\":\"value\"}",
				},
			},
		}

		params := &models.CreateArtifactParams{
			IfExists:  models.IfExistsCreate,
			Canonical: false,
			DryRun:    false,
		}
		result, err := api.CreateArtifact(context.Background(), "test-group", artifact, params)
		assert.Error(t, err)
		assert.Nil(t, result)

		fmt.Println(err)

		var apiErr *models.APIError
		ok := errors.As(err, &apiErr)
		assert.True(t, ok)
		assert.Equal(t, http.StatusBadRequest, apiErr.Status)
		assert.Equal(t, TitleBadRequest, apiErr.Title)
	})

	t.Run("Conflict", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusConflict)
			err := json.NewEncoder(w).Encode(models.APIError{Status: http.StatusConflict, Title: TitleConflict})
			assert.NoError(t, err)
		}))
		defer server.Close()

		mockClient := &client.Client{BaseURL: server.URL, HTTPClient: server.Client()}
		api := apis.NewArtifactsAPI(mockClient)

		artifact := models.CreateArtifactRequest{
			ArtifactType: models.Json,
			FirstVersion: models.CreateVersionRequest{
				Version: "1.0.0",
				Content: models.CreateContentRequest{
					Content: "{\"key\":\"value\"}",
				},
			},
		}

		params := &models.CreateArtifactParams{
			IfExists:  models.IfExistsCreate,
			Canonical: false,
			DryRun:    false,
		}
		result, err := api.CreateArtifact(context.Background(), "test-group", artifact, params)
		assert.Error(t, err)
		assert.Nil(t, result)

		var apiErr *models.APIError
		ok := errors.As(err, &apiErr)
		assert.True(t, ok)
		assert.Equal(t, http.StatusConflict, apiErr.Status)
		assert.Equal(t, TitleConflict, apiErr.Title)
	})

	t.Run("Internal Server Error", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
			err := json.NewEncoder(w).Encode(models.APIError{Status: http.StatusInternalServerError, Title: TitleInternalServerError})
			assert.NoError(t, err)
		}))
		defer server.Close()

		mockClient := &client.Client{BaseURL: server.URL, HTTPClient: server.Client()}
		api := apis.NewArtifactsAPI(mockClient)

		artifact := models.CreateArtifactRequest{
			ArtifactType: models.Json,
			FirstVersion: models.CreateVersionRequest{
				Version: "1.0.0",
				Content: models.CreateContentRequest{
					Content: "{\"key\":\"value\"}",
				},
			},
		}

		params := &models.CreateArtifactParams{
			IfExists:  models.IfExistsCreate,
			Canonical: false,
			DryRun:    false,
		}
		result, err := api.CreateArtifact(context.Background(), "test-group", artifact, params)
		assert.Error(t, err)
		assert.Nil(t, result)

		var apiErr *models.APIError
		ok := errors.As(err, &apiErr)
		assert.True(t, ok)
		assert.Equal(t, http.StatusInternalServerError, apiErr.Status)
		assert.Equal(t, TitleInternalServerError, apiErr.Title)
	})
}

func TestArtifactsAPI_ListArtifactRules(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mockReferences := []models.Rule{models.RuleValidity, models.RuleCompatibility}
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Contains(t, r.URL.Path, fmt.Sprintf("/groups/%s/artifacts/%s/rules", stubGroupId, stubArtifactId))
			assert.Equal(t, http.MethodGet, r.Method)

			w.WriteHeader(http.StatusOK)
			err := json.NewEncoder(w).Encode(mockReferences)
			assert.NoError(t, err)
		}))
		defer server.Close()

		mockClient := &client.Client{BaseURL: server.URL, HTTPClient: server.Client()}
		api := apis.NewArtifactsAPI(mockClient)

		result, err := api.ListArtifactRules(context.Background(), stubGroupId, stubArtifactId)
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Len(t, result, 2)
	})

	t.Run("Not Found", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusNotFound)
			err := json.NewEncoder(w).Encode(models.APIError{Status: http.StatusNotFound, Title: TitleNotFound})
			assert.NoError(t, err)
		}))
		defer server.Close()

		mockClient := &client.Client{BaseURL: server.URL, HTTPClient: server.Client()}
		api := apis.NewArtifactsAPI(mockClient)

		result, err := api.ListArtifactRules(context.Background(), stubGroupId, stubArtifactId)
		assert.Error(t, err)
		assert.Nil(t, result)

		var apiErr *models.APIError
		ok := errors.As(err, &apiErr)
		assert.True(t, ok)
		assert.Equal(t, http.StatusNotFound, apiErr.Status)
		assert.Equal(t, TitleNotFound, apiErr.Title)
	})

	t.Run("Internal Server Error", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
			err := json.NewEncoder(w).Encode(models.APIError{Status: http.StatusInternalServerError, Title: TitleInternalServerError})
			assert.NoError(t, err)
		}))
		defer server.Close()

		mockClient := &client.Client{BaseURL: server.URL, HTTPClient: server.Client()}
		api := apis.NewArtifactsAPI(mockClient)

		result, err := api.ListArtifactRules(context.Background(), stubGroupId, stubArtifactId)
		assert.Error(t, err)
		assert.Nil(t, result)

		var apiErr *models.APIError
		ok := errors.As(err, &apiErr)
		assert.True(t, ok)
		assert.Equal(t, http.StatusInternalServerError, apiErr.Status)
		assert.Equal(t, TitleInternalServerError, apiErr.Title)
	})
}

func TestArtifactsAPI_CreateArtifactRule(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Contains(t, r.URL.Path, fmt.Sprintf("/groups/%s/artifacts/%s/rules", stubGroupId, stubArtifactId))
			assert.Equal(t, http.MethodPost, r.Method)

			w.WriteHeader(http.StatusNoContent)
		}))
		defer server.Close()

		mockClient := &client.Client{BaseURL: server.URL, HTTPClient: server.Client()}
		api := apis.NewArtifactsAPI(mockClient)

		err := api.CreateArtifactRule(context.Background(), stubGroupId, stubArtifactId, models.RuleValidity, models.ValidityLevelFull)
		assert.NoError(t, err)
	})

	t.Run("BadRequest", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Contains(t, r.URL.Path, fmt.Sprintf("/groups/%s/artifacts/%s/rules", stubGroupId, stubArtifactId))
			assert.Equal(t, http.MethodPost, r.Method)

			w.WriteHeader(http.StatusBadRequest)
			err := json.NewEncoder(w).Encode(models.APIError{Status: http.StatusBadRequest, Title: TitleBadRequest})
			assert.NoError(t, err)
		}))
		defer server.Close()

		mockClient := &client.Client{BaseURL: server.URL, HTTPClient: server.Client()}
		api := apis.NewArtifactsAPI(mockClient)
		err := api.CreateArtifactRule(context.Background(), stubGroupId, stubArtifactId, models.RuleValidity, models.ValidityLevelFull)

		assert.Error(t, err)

		var apiErr *models.APIError
		ok := errors.As(err, &apiErr)
		assert.True(t, ok)
		assert.Equal(t, http.StatusBadRequest, apiErr.Status)
		assert.Equal(t, TitleBadRequest, apiErr.Title)
	})

	t.Run("Conflict", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Contains(t, r.URL.Path, fmt.Sprintf("/groups/%s/artifacts/%s/rules", stubGroupId, stubArtifactId))
			assert.Equal(t, http.MethodPost, r.Method)

			w.WriteHeader(http.StatusInternalServerError)
			err := json.NewEncoder(w).Encode(models.APIError{Status: http.StatusConflict, Title: TitleConflict})
			assert.NoError(t, err)
		}))
		defer server.Close()

		mockClient := &client.Client{BaseURL: server.URL, HTTPClient: server.Client()}
		api := apis.NewArtifactsAPI(mockClient)
		err := api.CreateArtifactRule(context.Background(), stubGroupId, stubArtifactId, models.RuleValidity, models.ValidityLevelFull)

		assert.Error(t, err)

		var apiErr *models.APIError
		ok := errors.As(err, &apiErr)
		assert.True(t, ok)
		assert.Equal(t, http.StatusConflict, apiErr.Status)
		assert.Equal(t, TitleConflict, apiErr.Title)
	})

	t.Run("NotFound", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Contains(t, r.URL.Path, fmt.Sprintf("/groups/%s/artifacts/%s/rules", stubGroupId, stubArtifactId))
			assert.Equal(t, http.MethodPost, r.Method)

			w.WriteHeader(http.StatusNotFound)
			err := json.NewEncoder(w).Encode(models.APIError{Status: http.StatusNotFound, Title: TitleNotFound})
			assert.NoError(t, err)
		}))
		defer server.Close()

		mockClient := &client.Client{BaseURL: server.URL, HTTPClient: server.Client()}
		api := apis.NewArtifactsAPI(mockClient)
		err := api.CreateArtifactRule(context.Background(), stubGroupId, stubArtifactId, models.RuleValidity, models.ValidityLevelFull)

		assert.Error(t, err)

		var apiErr *models.APIError
		ok := errors.As(err, &apiErr)
		assert.True(t, ok)
		assert.Equal(t, http.StatusNotFound, apiErr.Status)
		assert.Equal(t, TitleNotFound, apiErr.Title)
	})

	t.Run("InternalServerError", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Contains(t, r.URL.Path, fmt.Sprintf("/groups/%s/artifacts/%s/rules", stubGroupId, stubArtifactId))
			assert.Equal(t, http.MethodPost, r.Method)

			w.WriteHeader(http.StatusInternalServerError)
			err := json.NewEncoder(w).Encode(models.APIError{Status: http.StatusInternalServerError, Title: "Internal server error"})
			assert.NoError(t, err)
		}))
		defer server.Close()

		mockClient := &client.Client{BaseURL: server.URL, HTTPClient: server.Client()}
		api := apis.NewArtifactsAPI(mockClient)
		err := api.CreateArtifactRule(context.Background(), stubGroupId, stubArtifactId, models.RuleValidity, models.ValidityLevelFull)

		assert.Error(t, err)

		var apiErr *models.APIError
		ok := errors.As(err, &apiErr)
		assert.True(t, ok)
		assert.Equal(t, http.StatusInternalServerError, apiErr.Status)
		assert.Equal(t, TitleInternalServerError, apiErr.Title)
	})
}

func TestArtifactsAPI_DeleteAllArtifactRule(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Contains(t, r.URL.Path, fmt.Sprintf("/groups/%s/artifacts/%s/rules", stubGroupId, stubArtifactId))
			assert.Equal(t, http.MethodDelete, r.Method)

			w.WriteHeader(http.StatusNoContent)
		}))
		defer server.Close()

		mockClient := &client.Client{BaseURL: server.URL, HTTPClient: server.Client()}
		api := apis.NewArtifactsAPI(mockClient)
		err := api.DeleteAllArtifactRule(context.Background(), stubGroupId, stubArtifactId)
		assert.NoError(t, err)
	})

	t.Run("NotFound", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Contains(t, r.URL.Path, fmt.Sprintf("/groups/%s/artifacts/%s/rules", stubGroupId, stubArtifactId))
			assert.Equal(t, http.MethodDelete, r.Method)

			w.WriteHeader(http.StatusNotFound)
			err := json.NewEncoder(w).Encode(models.APIError{Status: http.StatusNotFound, Title: TitleNotFound})
			assert.NoError(t, err)
		}))
		defer server.Close()

		mockClient := &client.Client{BaseURL: server.URL, HTTPClient: server.Client()}
		api := apis.NewArtifactsAPI(mockClient)
		err := api.DeleteAllArtifactRule(context.Background(), stubGroupId, stubArtifactId)
		assert.Error(t, err)

		var apiErr *models.APIError
		ok := errors.As(err, &apiErr)
		assert.True(t, ok)
		assert.Equal(t, http.StatusNotFound, apiErr.Status)
		assert.Equal(t, TitleNotFound, apiErr.Title)
	})

	t.Run("InternalServerError", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Contains(t, r.URL.Path, fmt.Sprintf("/groups/%s/artifacts/%s/rules", stubGroupId, stubArtifactId))
			assert.Equal(t, http.MethodDelete, r.Method)

			w.WriteHeader(http.StatusInternalServerError)
			err := json.NewEncoder(w).Encode(models.APIError{Status: http.StatusInternalServerError, Title: "Internal server error"})
			assert.NoError(t, err)
		}))
		defer server.Close()

		mockClient := &client.Client{BaseURL: server.URL, HTTPClient: server.Client()}
		api := apis.NewArtifactsAPI(mockClient)
		err := api.DeleteAllArtifactRule(context.Background(), stubGroupId, stubArtifactId)
		assert.Error(t, err)

		var apiErr *models.APIError
		ok := errors.As(err, &apiErr)
		assert.True(t, ok)
		assert.Equal(t, http.StatusInternalServerError, apiErr.Status)
		assert.Equal(t, TitleInternalServerError, apiErr.Title)
	})
}

func TestArtifactsAPI_GetArtifactRule(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mockResponse := models.RuleResponse{
			RuleType: models.RuleValidity,
			Config:   models.ValidityLevelFull,
		}
		mockRule := models.RuleValidity
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Contains(t, r.URL.Path, fmt.Sprintf("/groups/%s/artifacts/%s/rules/%s", stubGroupId, stubArtifactId, mockRule))
			assert.Equal(t, http.MethodGet, r.Method)

			w.WriteHeader(http.StatusOK)
			err := json.NewEncoder(w).Encode(mockResponse)
			assert.NoError(t, err)
		}))
		defer server.Close()

		mockClient := &client.Client{BaseURL: server.URL, HTTPClient: server.Client()}
		api := apis.NewArtifactsAPI(mockClient)
		result, err := api.GetArtifactRule(context.Background(), stubGroupId, stubArtifactId, mockRule)
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, models.ValidityLevelFull, result)
	})

	t.Run("NotFound", func(t *testing.T) {
		mockRule := models.RuleValidity
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Contains(t, r.URL.Path, fmt.Sprintf("/groups/%s/artifacts/%s/rules/%s", stubGroupId, stubArtifactId, mockRule))
			assert.Equal(t, http.MethodGet, r.Method)

			w.WriteHeader(http.StatusNotFound)
			err := json.NewEncoder(w).Encode(models.APIError{Status: http.StatusNotFound, Title: TitleNotFound})
			assert.NoError(t, err)
		}))
		defer server.Close()

		mockClient := &client.Client{BaseURL: server.URL, HTTPClient: server.Client()}
		api := apis.NewArtifactsAPI(mockClient)
		result, err := api.GetArtifactRule(context.Background(), stubGroupId, stubArtifactId, mockRule)
		assert.Error(t, err)
		assert.Empty(t, result)

		var apiErr *models.APIError
		ok := errors.As(err, &apiErr)
		assert.True(t, ok)
		assert.Equal(t, http.StatusNotFound, apiErr.Status)
		assert.Equal(t, TitleNotFound, apiErr.Title)
	})

	t.Run("InternalServerError", func(t *testing.T) {
		mockRule := models.RuleValidity
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Contains(t, r.URL.Path, fmt.Sprintf("/groups/%s/artifacts/%s/rules/%s", stubGroupId, stubArtifactId, mockRule))
			assert.Equal(t, http.MethodGet, r.Method)

			w.WriteHeader(http.StatusInternalServerError)
			err := json.NewEncoder(w).Encode(models.APIError{Status: http.StatusInternalServerError, Title: "Internal server error"})
			assert.NoError(t, err)
		}))
		defer server.Close()

		mockClient := &client.Client{BaseURL: server.URL, HTTPClient: server.Client()}
		api := apis.NewArtifactsAPI(mockClient)
		result, err := api.GetArtifactRule(context.Background(), stubGroupId, stubArtifactId, mockRule)
		assert.Error(t, err)
		assert.Empty(t, result)

		var apiErr *models.APIError
		ok := errors.As(err, &apiErr)
		assert.True(t, ok)
		assert.Equal(t, http.StatusInternalServerError, apiErr.Status)
		assert.Equal(t, TitleInternalServerError, apiErr.Title)
	})
}

func TestArtifactsAPI_UpdateArtifactRule(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mockRule := models.RuleValidity
		mockResponse := models.RuleResponse{
			RuleType: mockRule,
			Config:   models.ValidityLevelFull,
		}
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Contains(t, r.URL.Path, fmt.Sprintf("/groups/%s/artifacts/%s/rules/%s", stubGroupId, stubArtifactId, mockRule))
			assert.Equal(t, http.MethodPut, r.Method)

			w.WriteHeader(http.StatusOK)
			err := json.NewEncoder(w).Encode(mockResponse)
			assert.NoError(t, err)
		}))
		defer server.Close()

		mockClient := &client.Client{BaseURL: server.URL, HTTPClient: server.Client()}
		api := apis.NewArtifactsAPI(mockClient)
		err := api.UpdateArtifactRule(context.Background(), stubGroupId, stubArtifactId, mockRule, models.ValidityLevelFull)
		assert.NoError(t, err)
	})

	t.Run("NotFound", func(t *testing.T) {
		mockRule := models.RuleValidity
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Contains(t, r.URL.Path, fmt.Sprintf("/groups/%s/artifacts/%s/rules/%s", stubGroupId, stubArtifactId, mockRule))
			assert.Equal(t, http.MethodPut, r.Method)

			w.WriteHeader(http.StatusNotFound)
			err := json.NewEncoder(w).Encode(models.APIError{Status: http.StatusNotFound, Title: TitleNotFound})
			assert.NoError(t, err)
		}))
		defer server.Close()

		mockClient := &client.Client{BaseURL: server.URL, HTTPClient: server.Client()}
		api := apis.NewArtifactsAPI(mockClient)
		err := api.UpdateArtifactRule(context.Background(), stubGroupId, stubArtifactId, mockRule, models.ValidityLevelFull)
		assert.Error(t, err)

		var apiErr *models.APIError
		ok := errors.As(err, &apiErr)
		assert.True(t, ok)
		assert.Equal(t, http.StatusNotFound, apiErr.Status)
		assert.Equal(t, TitleNotFound, apiErr.Title)
	})

	t.Run("InternalServerError", func(t *testing.T) {
		mockRule := models.RuleValidity
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Contains(t, r.URL.Path, fmt.Sprintf("/groups/%s/artifacts/%s/rules/%s", stubGroupId, stubArtifactId, mockRule))
			assert.Equal(t, http.MethodPut, r.Method)

			w.WriteHeader(http.StatusInternalServerError)
			err := json.NewEncoder(w).Encode(models.APIError{Status: http.StatusInternalServerError, Title: "Internal server error"})
			assert.NoError(t, err)
		}))
		defer server.Close()

		mockClient := &client.Client{BaseURL: server.URL, HTTPClient: server.Client()}
		api := apis.NewArtifactsAPI(mockClient)
		err := api.UpdateArtifactRule(context.Background(), stubGroupId, stubArtifactId, mockRule, models.ValidityLevelFull)
		assert.Error(t, err)

		var apiErr *models.APIError
		ok := errors.As(err, &apiErr)
		assert.True(t, ok)
		assert.Equal(t, http.StatusInternalServerError, apiErr.Status)
		assert.Equal(t, TitleInternalServerError, apiErr.Title)
	})
}

func TestArtifactsAPI_DeleteArtifactRule(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mockRule := models.RuleValidity
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Contains(t, r.URL.Path, fmt.Sprintf("/groups/%s/artifacts/%s/rules/%s", stubGroupId, stubArtifactId, mockRule))
			assert.Equal(t, http.MethodDelete, r.Method)

			w.WriteHeader(http.StatusNoContent)
		}))
		defer server.Close()

		mockClient := &client.Client{BaseURL: server.URL, HTTPClient: server.Client()}
		api := apis.NewArtifactsAPI(mockClient)
		err := api.DeleteArtifactRule(context.Background(), stubGroupId, stubArtifactId, mockRule)
		assert.NoError(t, err)
	})

	t.Run("NotFound", func(t *testing.T) {
		mockRule := models.RuleValidity
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Contains(t, r.URL.Path, fmt.Sprintf("/groups/%s/artifacts/%s/rules/%s", stubGroupId, stubArtifactId, mockRule))
			assert.Equal(t, http.MethodDelete, r.Method)

			w.WriteHeader(http.StatusNotFound)
			err := json.NewEncoder(w).Encode(models.APIError{Status: http.StatusNotFound, Title: TitleNotFound})
			assert.NoError(t, err)
		}))
		defer server.Close()

		mockClient := &client.Client{BaseURL: server.URL, HTTPClient: server.Client()}
		api := apis.NewArtifactsAPI(mockClient)
		err := api.DeleteArtifactRule(context.Background(), stubGroupId, stubArtifactId, mockRule)
		assert.Error(t, err)

		var apiErr *models.APIError
		ok := errors.As(err, &apiErr)
		assert.True(t, ok)
		assert.Equal(t, http.StatusNotFound, apiErr.Status)
		assert.Equal(t, TitleNotFound, apiErr.Title)
	})

	t.Run("InternalServerError", func(t *testing.T) {
		mockRule := models.RuleValidity
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Contains(t, r.URL.Path, fmt.Sprintf("/groups/%s/artifacts/%s/rules/%s", stubGroupId, stubArtifactId, mockRule))
			assert.Equal(t, http.MethodDelete, r.Method)

			w.WriteHeader(http.StatusInternalServerError)
			err := json.NewEncoder(w).Encode(models.APIError{Status: http.StatusInternalServerError, Title: "Internal server error"})
			assert.NoError(t, err)
		}))
		defer server.Close()

		mockClient := &client.Client{BaseURL: server.URL, HTTPClient: server.Client()}
		api := apis.NewArtifactsAPI(mockClient)
		err := api.DeleteArtifactRule(context.Background(), stubGroupId, stubArtifactId, mockRule)
		assert.Error(t, err)

		var apiErr *models.APIError
		ok := errors.As(err, &apiErr)
		assert.True(t, ok)
		assert.Equal(t, http.StatusInternalServerError, apiErr.Status)
		assert.Equal(t, TitleInternalServerError, apiErr.Title)
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
			ArtifactID:   artifactID,
			Name:         artifactID,
			FirstVersion: models.CreateVersionRequest{
				Version: "1.0.0",
				Content: models.CreateContentRequest{
					Content: stubArtifactContent,
				},
			},
		}

		params := &models.CreateArtifactParams{
			IfExists: models.IfExistsFail,
		}

		resp, err := artifactsAPI.CreateArtifact(ctx, groupID, artifact, params)
		assert.NoError(t, err)
		assert.Equal(t, groupID, resp.GroupID)
		assert.Equal(t, artifactID, resp.Name)
	})

	// Test SearchArtifacts
	t.Run("SearchArtifacts", func(t *testing.T) {
		params := &models.SearchArtifactsParams{
			Name: artifactID,
		}
		resp, err := artifactsAPI.SearchArtifacts(ctx, params)
		assert.NoError(t, err)
		assert.GreaterOrEqual(t, len(*resp), 1)
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
		resp, err := artifactsAPI.ListArtifactsInGroup(ctx, groupID, params)
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
		err := artifactsAPI.DeleteArtifactsInGroup(ctx, groupID)
		assert.NoError(t, err)
	})

	// Test DeleteArtifact
	t.Run("DeleteArtifact", func(t *testing.T) {

		// Re-create the artifact
		artifact := models.CreateArtifactRequest{
			ArtifactType: models.Json,
			ArtifactID:   artifactID,
			Name:         artifactID,
			FirstVersion: models.CreateVersionRequest{
				Version: "1.0.0",
				Content: models.CreateContentRequest{
					Content: stubArtifactContent,
				},
			},
		}
		params := &models.CreateArtifactParams{
			IfExists: models.IfExistsFail,
		}

		resp, err := artifactsAPI.CreateArtifact(ctx, groupID, artifact, params)
		assert.NoError(t, err)
		assert.Equal(t, groupID, resp.GroupID)
		assert.Equal(t, artifactID, resp.Name)

		// Delete the artifact
		err = artifactsAPI.DeleteArtifact(ctx, groupID, artifactID)
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
