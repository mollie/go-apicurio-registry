package apis_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/mollie/go-apicurio-registry/apis"
	"github.com/mollie/go-apicurio-registry/client"
	"github.com/mollie/go-apicurio-registry/models"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

func setupSystemAPIClient() *apis.SystemAPI {
	apiClient := setupHTTPClient()
	return apis.NewSystemAPI(apiClient)
}

func TestSystemAPI_GetSystemInfo(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mockResponse := models.SystemInfoResponse{
			Name:        "Apicurio Registry (SQL)",
			Description: "The Apicurio Registry application.",
			Version:     "3.0.5",
			BuiltOn:     "2021-03-19T12:55:00Z",
		}

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Contains(t, r.URL.Path, "/system/info")
			assert.Equal(t, http.MethodGet, r.Method)

			w.WriteHeader(http.StatusOK)
			err := json.NewEncoder(w).Encode(mockResponse)
			assert.NoError(t, err)
		}))
		defer server.Close()

		mockClient := &client.Client{BaseURL: server.URL, HTTPClient: server.Client()}
		api := apis.NewSystemAPI(mockClient)

		result, err := api.GetSystemInfo(context.Background())
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, mockResponse.Version, result.Version)
		assert.Equal(t, mockResponse.BuiltOn, result.BuiltOn)
	})

	t.Run("InternalServerError", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Contains(t, r.URL.Path, "/system/info")
			assert.Equal(t, http.MethodGet, r.Method)

			w.WriteHeader(http.StatusInternalServerError)
			err := json.NewEncoder(w).Encode(models.APIError{Status: http.StatusInternalServerError, Title: "Internal Server Error"})
			assert.NoError(t, err)
		}))
		defer server.Close()

		mockClient := &client.Client{BaseURL: server.URL, HTTPClient: server.Client()}
		api := apis.NewSystemAPI(mockClient)

		result, err := api.GetSystemInfo(context.Background())
		assert.Error(t, err)
		assert.Nil(t, result)

		var apiErr *models.APIError
		ok := errors.As(err, &apiErr)
		assert.True(t, ok)
		assert.Equal(t, http.StatusInternalServerError, apiErr.Status)
		assert.Equal(t, "Internal Server Error", apiErr.Title)
	})
}

func TestSystemAPI_GetResourceLimitInfo(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mockResponse := models.SystemResourceLimitInfoResponse{
			MaxTotalSchemasCount:              -1,
			MaxSchemaSizeBytes:                -1,
			MaxArtifactsCount:                 -1,
			MaxVersionsPerArtifactCount:       -1,
			MaxArtifactPropertiesCount:        -1,
			MaxPropertyKeySizeBytes:           -1,
			MaxPropertyValueSizeBytes:         -1,
			MaxArtifactLabelsCount:            -1,
			MaxLabelSizeBytes:                 -1,
			MaxArtifactNameLengthChars:        -1,
			MaxArtifactDescriptionLengthChars: -1,
			MaxRequestsPerSecondCount:         -1,
		}

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Contains(t, r.URL.Path, "/system/limits")
			assert.Equal(t, http.MethodGet, r.Method)

			w.WriteHeader(http.StatusOK)
			err := json.NewEncoder(w).Encode(mockResponse)
			assert.NoError(t, err)
		}))
		defer server.Close()

		mockClient := &client.Client{BaseURL: server.URL, HTTPClient: server.Client()}
		api := apis.NewSystemAPI(mockClient)

		result, err := api.GetResourceLimitInfo(context.Background())
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, mockResponse, *result)
	})

	t.Run("InternalServerError", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Contains(t, r.URL.Path, "/system/limits")
			assert.Equal(t, http.MethodGet, r.Method)

			w.WriteHeader(http.StatusInternalServerError)
			err := json.NewEncoder(w).Encode(models.APIError{Status: http.StatusInternalServerError, Title: "Internal Server Error"})
			assert.NoError(t, err)
		}))
		defer server.Close()

		mockClient := &client.Client{BaseURL: server.URL, HTTPClient: server.Client()}
		api := apis.NewSystemAPI(mockClient)

		result, err := api.GetResourceLimitInfo(context.Background())
		assert.Error(t, err)
		assert.Nil(t, result)

		var apiErr *models.APIError
		ok := errors.As(err, &apiErr)
		assert.True(t, ok)
		assert.Equal(t, http.StatusInternalServerError, apiErr.Status)
		assert.Equal(t, "Internal Server Error", apiErr.Title)
	})
}

func TestSystemAPI_GetUIConfig(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mockResponse := models.SystemUIConfigResponse{
			Ui: models.UIConfig{
				ContextPath:   "/",
				NavPrefixPath: "/",
				OaiDocsUrl:    "https://registry.apicur.io/docs",
			},
			Auth: models.AuthConfig{
				Type:        "oidc",
				RbacEnabled: true,
				ObacEnabled: false,
				Options: models.AuthOptions{
					Url:         "https://auth.apicur.io/realms/apicurio",
					RedirectUri: "http://registry.apicur.io",
					ClientId:    "apicurio-registry-ui",
				},
			},
			Features: models.FeatureFlags{
				ReadOnly:       false,
				Breadcrumbs:    true,
				RoleManagement: false,
				Settings:       true,
			},
		}

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Contains(t, r.URL.Path, "/system/uiConfig")
			assert.Equal(t, http.MethodGet, r.Method)

			w.WriteHeader(http.StatusOK)
			err := json.NewEncoder(w).Encode(mockResponse)
			assert.NoError(t, err)
		}))
		defer server.Close()

		mockClient := &client.Client{BaseURL: server.URL, HTTPClient: server.Client()}
		api := apis.NewSystemAPI(mockClient)

		result, err := api.GetUIConfig(context.Background())
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, mockResponse, *result)
	})

	t.Run("InternalServerError", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Contains(t, r.URL.Path, "/system/uiConfig")
			assert.Equal(t, http.MethodGet, r.Method)

			w.WriteHeader(http.StatusInternalServerError)
			err := json.NewEncoder(w).Encode(models.APIError{Status: http.StatusInternalServerError, Title: "Internal Server Error"})
			assert.NoError(t, err)
		}))
		defer server.Close()

		mockClient := &client.Client{BaseURL: server.URL, HTTPClient: server.Client()}
		api := apis.NewSystemAPI(mockClient)

		result, err := api.GetUIConfig(context.Background())
		assert.Error(t, err)
		assert.Nil(t, result)

		var apiErr *models.APIError
		ok := errors.As(err, &apiErr)
		assert.True(t, ok)
		assert.Equal(t, http.StatusInternalServerError, apiErr.Status)
		assert.Equal(t, "Internal Server Error", apiErr.Title)
	})
}

func TestSystemAPI_GetCurrentUser(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mockResponse := models.UserInfo{
			Username:    "test-user",
			DisplayName: "my-test-user",
			Admin:       true,
			Developer:   true,
			Viewer:      true,
		}

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Contains(t, r.URL.Path, "/users/me")
			assert.Equal(t, http.MethodGet, r.Method)

			w.WriteHeader(http.StatusOK)
			err := json.NewEncoder(w).Encode(mockResponse)
			assert.NoError(t, err)
		}))
		defer server.Close()

		mockClient := &client.Client{BaseURL: server.URL, HTTPClient: server.Client()}
		api := apis.NewSystemAPI(mockClient)

		result, err := api.GetCurrentUser(context.Background())
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, mockResponse, *result)
	})

	t.Run("Unauthorized", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Contains(t, r.URL.Path, "/users/me")
			assert.Equal(t, http.MethodGet, r.Method)

			w.WriteHeader(http.StatusUnauthorized)
			err := json.NewEncoder(w).Encode(models.APIError{Status: http.StatusUnauthorized, Title: "Unauthorized"})
			assert.NoError(t, err)
		}))
		defer server.Close()

		mockClient := &client.Client{BaseURL: server.URL, HTTPClient: server.Client()}
		api := apis.NewSystemAPI(mockClient)

		result, err := api.GetCurrentUser(context.Background())
		assert.Error(t, err)
		assert.Nil(t, result)

		var apiErr *models.APIError
		ok := errors.As(err, &apiErr)
		assert.True(t, ok)
		assert.Equal(t, http.StatusUnauthorized, apiErr.Status)
		assert.Equal(t, "Unauthorized", apiErr.Title)
	})

	t.Run("InvalidJSON", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Contains(t, r.URL.Path, "/users/me")
			assert.Equal(t, http.MethodGet, r.Method)

			w.WriteHeader(http.StatusOK)
			_, err := w.Write([]byte(`{ invalid json }`))
			assert.NoError(t, err)
		}))
		defer server.Close()

		mockClient := &client.Client{BaseURL: server.URL, HTTPClient: server.Client()}
		api := apis.NewSystemAPI(mockClient)

		result, err := api.GetCurrentUser(context.Background())
		assert.Error(t, err)
		assert.Nil(t, result)
	})
}

/***********************/
/***** Integration *****/
/***********************/

func TestSystemAPI_All_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	ctx := context.Background()
	api := setupSystemAPIClient()

	t.Run("GetSystemInfo", func(t *testing.T) {
		expected := &models.SystemInfoResponse{
			Name:        "Apicurio Registry (In Memory)",
			Description: "High performance, runtime registry for schemas and API designs.",
			Version:     "3.0.5",
			BuiltOn:     "2024-12-03T12:31:33Z",
		}

		result, err := api.GetSystemInfo(ctx)
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, expected, result)
	})

	t.Run("GetResourceLimitInfo", func(t *testing.T) {
		expected := &models.SystemResourceLimitInfoResponse{
			MaxTotalSchemasCount:              -1,
			MaxSchemaSizeBytes:                -1,
			MaxArtifactsCount:                 -1,
			MaxVersionsPerArtifactCount:       -1,
			MaxArtifactPropertiesCount:        -1,
			MaxPropertyKeySizeBytes:           -1,
			MaxPropertyValueSizeBytes:         -1,
			MaxArtifactLabelsCount:            -1,
			MaxLabelSizeBytes:                 -1,
			MaxArtifactNameLengthChars:        -1,
			MaxArtifactDescriptionLengthChars: -1,
			MaxRequestsPerSecondCount:         -1,
		}

		result, err := api.GetResourceLimitInfo(ctx)
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.EqualValues(t, expected, result)
	})

	t.Run("GetUIConfig", func(t *testing.T) {
		expected := &models.SystemUIConfigResponse{
			Ui: models.UIConfig{
				ContextPath:   "/",
				NavPrefixPath: "/",
				OaiDocsUrl:    "/docs/",
			},
			Auth: models.AuthConfig{
				Type:        "none",
				RbacEnabled: false,
				ObacEnabled: false,
			},
			Features: models.FeatureFlags{
				ReadOnly:        false,
				Breadcrumbs:     true,
				RoleManagement:  false,
				Settings:        true,
				DeleteGroup:     true,
				DeleteArtifact:  true,
				DeleteVersion:   true,
				DraftMutability: true,
			},
		}

		result, err := api.GetUIConfig(ctx)
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.EqualValues(t, expected, result)
	})
}
