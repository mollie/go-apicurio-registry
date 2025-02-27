package apis_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/mollie/go-apicurio-registry/apis"
	"github.com/mollie/go-apicurio-registry/client"
	"github.com/mollie/go-apicurio-registry/models"
	"github.com/stretchr/testify/assert"
)

func TestSystemAPI_GetSystemInfo(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mockResponse := models.SystemInfoResponse{
			Name:        "Apicurio Registry (SQL)",
			Description: "The Apicurio Registry application.",
			Version:     "3.0.5",
			BuiltOn:     "2021-03-19T12:55:00Z",
		}

		server := setupMockServer(t, http.StatusOK, mockResponse, "/system/info", http.MethodGet)
		defer server.Close()

		mockClient := &client.Client{BaseURL: server.URL, HTTPClient: server.Client()}
		api := apis.NewSystemAPI(mockClient)

		result, err := api.GetSystemInfo(context.Background())

		assert.NoError(t, err)
		assert.Equal(t, mockResponse.Version, result.Version)
	})

	t.Run("InternalServerError", func(t *testing.T) {
		server := setupMockServer(t, http.StatusInternalServerError, models.APIError{
			Status: http.StatusInternalServerError, Title: "Internal Server Error",
		}, "/system/info", http.MethodGet)
		defer server.Close()

		mockClient := &client.Client{BaseURL: server.URL, HTTPClient: server.Client()}
		api := apis.NewSystemAPI(mockClient)

		result, err := api.GetSystemInfo(context.Background())

		assertAPIError(t, err, http.StatusInternalServerError, "Internal Server Error")
		assert.Nil(t, result)
	})
}

func TestSystemAPI_GetResourceLimitInfo(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mockResponse := models.SystemResourceLimitInfoResponse{MaxArtifactsCount: -1}
		server := setupMockServer(t, http.StatusOK, mockResponse, "/system/limits", http.MethodGet)
		defer server.Close()

		mockClient := &client.Client{BaseURL: server.URL, HTTPClient: server.Client()}
		api := apis.NewSystemAPI(mockClient)

		result, err := api.GetResourceLimitInfo(context.Background())

		assert.NoError(t, err)
		assert.Equal(t, mockResponse, *result)
	})

	t.Run("InternalServerError", func(t *testing.T) {
		server := setupMockServer(t, http.StatusInternalServerError, models.APIError{
			Status: http.StatusInternalServerError, Title: "Internal Server Error",
		}, "/system/limits", http.MethodGet)
		defer server.Close()

		mockClient := &client.Client{BaseURL: server.URL, HTTPClient: server.Client()}
		api := apis.NewSystemAPI(mockClient)

		result, err := api.GetResourceLimitInfo(context.Background())

		assertAPIError(t, err, http.StatusInternalServerError, "Internal Server Error")
		assert.Nil(t, result)
	})
}

func TestSystemAPI_GetUIConfig(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mockResponse := models.SystemUIConfigResponse{Ui: models.UIConfig{ContextPath: "/"}}
		server := setupMockServer(
			t,
			http.StatusOK,
			mockResponse,
			"/system/uiConfig",
			http.MethodGet,
		)
		defer server.Close()

		mockClient := &client.Client{BaseURL: server.URL, HTTPClient: server.Client()}
		api := apis.NewSystemAPI(mockClient)

		result, err := api.GetUIConfig(context.Background())

		assert.NoError(t, err)
		assert.Equal(t, mockResponse, *result)
	})

	t.Run("InternalServerError", func(t *testing.T) {
		server := setupMockServer(t, http.StatusInternalServerError, models.APIError{
			Status: http.StatusInternalServerError, Title: "Internal Server Error",
		}, "/system/uiConfig", http.MethodGet)
		defer server.Close()

		mockClient := &client.Client{BaseURL: server.URL, HTTPClient: server.Client()}
		api := apis.NewSystemAPI(mockClient)

		result, err := api.GetUIConfig(context.Background())

		assertAPIError(t, err, http.StatusInternalServerError, "Internal Server Error")
		assert.Nil(t, result)
	})
}

func TestSystemAPI_GetCurrentUser(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mockResponse := models.UserInfo{Username: "test-user"}
		server := setupMockServer(t, http.StatusOK, mockResponse, "/users/me", http.MethodGet)
		defer server.Close()

		mockClient := &client.Client{BaseURL: server.URL, HTTPClient: server.Client()}
		api := apis.NewSystemAPI(mockClient)

		result, err := api.GetCurrentUser(context.Background())

		assert.NoError(t, err)
		assert.Equal(t, mockResponse, *result)
	})

	t.Run("Unauthorized", func(t *testing.T) {
		server := setupMockServer(t, http.StatusUnauthorized, models.APIError{
			Status: http.StatusUnauthorized, Title: "Unauthorized",
		}, "/users/me", http.MethodGet)
		defer server.Close()

		mockClient := &client.Client{BaseURL: server.URL, HTTPClient: server.Client()}
		api := apis.NewSystemAPI(mockClient)

		result, err := api.GetCurrentUser(context.Background())

		assertAPIError(t, err, http.StatusUnauthorized, "Unauthorized")
		assert.Nil(t, result)
	})

	t.Run("InvalidJSON", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			_, err := w.Write([]byte(`{ invalid json }`)) // Invalid JSON Response
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

func setupSystemAPIClient() *apis.SystemAPI {
	apiClient := setupHTTPClient()
	return apis.NewSystemAPI(apiClient)
}
