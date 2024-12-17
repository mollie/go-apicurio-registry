package apis_test

import (
	"context"
	"github.com/mollie/go-apicurio-registry/apis"
	"github.com/mollie/go-apicurio-registry/client"
	"github.com/mollie/go-apicurio-registry/models"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func TestAdminAPI_ListGlobalRules(t *testing.T) {
	mockReferences := []models.Rule{models.RuleValidity, models.RuleCompatibility}

	t.Run("Success", func(t *testing.T) {
		server := setupMockServer(t, http.StatusOK, mockReferences, "/admin/rules", http.MethodGet)
		defer server.Close()

		mockClient := &client.Client{BaseURL: server.URL, HTTPClient: server.Client()}
		api := apis.NewAdminAPI(mockClient)

		result, err := api.ListGlobalRules(context.Background())
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Len(t, result, 2)
	})

	t.Run("InternalServerError", func(t *testing.T) {
		errorResponse := models.APIError{Status: http.StatusInternalServerError, Title: TitleInternalServerError}
		server := setupMockServer(t, http.StatusInternalServerError, errorResponse, "/admin/rules", http.MethodGet)
		defer server.Close()

		mockClient := &client.Client{BaseURL: server.URL, HTTPClient: server.Client()}
		api := apis.NewAdminAPI(mockClient)

		res, err := api.ListGlobalRules(context.Background())
		assert.Error(t, err)
		assert.Nil(t, res)
		assertAPIError(t, err, http.StatusInternalServerError, TitleInternalServerError)
	})
}

func TestAdminAPI_CreateGlobalRule(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		server := setupMockServer(t, http.StatusNoContent, nil, "/admin/rules", http.MethodPost)
		defer server.Close()

		mockClient := &client.Client{BaseURL: server.URL, HTTPClient: server.Client()}
		api := apis.NewAdminAPI(mockClient)

		err := api.CreateGlobalRule(context.Background(), models.RuleValidity, models.ValidityLevelFull)
		assert.NoError(t, err)
	})

	t.Run("BadRequest", func(t *testing.T) {
		errorResponse := models.APIError{Status: http.StatusBadRequest, Title: TitleBadRequest}
		server := setupMockServer(t, http.StatusBadRequest, errorResponse, "/admin/rules", http.MethodPost)
		defer server.Close()

		mockClient := &client.Client{BaseURL: server.URL, HTTPClient: server.Client()}
		api := apis.NewAdminAPI(mockClient)

		err := api.CreateGlobalRule(context.Background(), models.RuleValidity, models.ValidityLevelFull)
		assertAPIError(t, err, http.StatusBadRequest, TitleBadRequest)
	})

	t.Run("Conflict", func(t *testing.T) {
		errorResponse := models.APIError{Status: http.StatusConflict, Title: TitleConflict}
		server := setupMockServer(t, http.StatusConflict, errorResponse, "/admin/rules", http.MethodPost)
		defer server.Close()

		mockClient := &client.Client{BaseURL: server.URL, HTTPClient: server.Client()}
		api := apis.NewAdminAPI(mockClient)

		err := api.CreateGlobalRule(context.Background(), models.RuleValidity, models.ValidityLevelFull)
		assertAPIError(t, err, http.StatusConflict, TitleConflict)
	})

	t.Run("InternalServerError", func(t *testing.T) {
		errorResponse := models.APIError{Status: http.StatusInternalServerError, Title: TitleInternalServerError}
		server := setupMockServer(t, http.StatusInternalServerError, errorResponse, "/admin/rules", http.MethodPost)
		defer server.Close()

		mockClient := &client.Client{BaseURL: server.URL, HTTPClient: server.Client()}
		api := apis.NewAdminAPI(mockClient)

		err := api.CreateGlobalRule(context.Background(), models.RuleValidity, models.ValidityLevelFull)
		assertAPIError(t, err, http.StatusInternalServerError, TitleInternalServerError)
	})
}

func TestAdminAPI_DeleteAllGlobalRule(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		server := setupMockServer(t, http.StatusNoContent, nil, "/admin/rules", http.MethodDelete)
		defer server.Close()

		mockClient := &client.Client{BaseURL: server.URL, HTTPClient: server.Client()}
		api := apis.NewAdminAPI(mockClient)

		err := api.DeleteAllGlobalRule(context.Background())
		assert.NoError(t, err)
	})

	t.Run("InternalServerError", func(t *testing.T) {
		errorResponse := models.APIError{Status: http.StatusInternalServerError, Title: TitleInternalServerError}
		server := setupMockServer(t, http.StatusInternalServerError, errorResponse, "/admin/rules", http.MethodDelete)
		defer server.Close()

		mockClient := &client.Client{BaseURL: server.URL, HTTPClient: server.Client()}
		api := apis.NewAdminAPI(mockClient)

		err := api.DeleteAllGlobalRule(context.Background())
		assertAPIError(t, err, http.StatusInternalServerError, TitleInternalServerError)
	})
}

func TestAdminAPI_GetGlobalRule(t *testing.T) {
	mockResponse := models.RuleResponse{RuleType: models.RuleValidity, Config: models.ValidityLevelFull}

	t.Run("Success", func(t *testing.T) {
		server := setupMockServer(t, http.StatusOK, mockResponse, "/admin/rules/VALIDITY", http.MethodGet)
		defer server.Close()

		mockClient := &client.Client{BaseURL: server.URL, HTTPClient: server.Client()}
		api := apis.NewAdminAPI(mockClient)

		result, err := api.GetGlobalRule(context.Background(), models.RuleValidity)
		assert.NoError(t, err)
		assert.Equal(t, models.ValidityLevelFull, result)
	})

	t.Run("NotFound", func(t *testing.T) {
		errorResponse := models.APIError{Status: http.StatusNotFound, Title: TitleNotFound}
		server := setupMockServer(t, http.StatusNotFound, errorResponse, "/admin/rules/VALIDITY", http.MethodGet)
		defer server.Close()

		mockClient := &client.Client{BaseURL: server.URL, HTTPClient: server.Client()}
		api := apis.NewAdminAPI(mockClient)

		result, err := api.GetGlobalRule(context.Background(), models.RuleValidity)
		assert.Error(t, err)
		assert.Empty(t, result)
		assertAPIError(t, err, http.StatusNotFound, TitleNotFound)
	})

	t.Run("InternalServerError", func(t *testing.T) {
		errorResponse := models.APIError{Status: http.StatusInternalServerError, Title: TitleInternalServerError}
		server := setupMockServer(t, http.StatusInternalServerError, errorResponse, "/admin/rules/VALIDITY", http.MethodGet)
		defer server.Close()

		mockClient := &client.Client{BaseURL: server.URL, HTTPClient: server.Client()}
		api := apis.NewAdminAPI(mockClient)

		result, err := api.GetGlobalRule(context.Background(), models.RuleValidity)
		assert.Error(t, err)
		assert.Empty(t, result)
		assertAPIError(t, err, http.StatusInternalServerError, TitleInternalServerError)
	})
}

func TestAdminAPI_UpdateGlobalRule(t *testing.T) {
	mockResponse := models.RuleResponse{
		RuleType: models.RuleValidity,
		Config:   models.ValidityLevelFull,
	}

	t.Run("Success", func(t *testing.T) {
		server := setupMockServer(t, http.StatusOK, mockResponse, "/admin/rules/VALIDITY", http.MethodPut)
		defer server.Close()

		mockClient := &client.Client{BaseURL: server.URL, HTTPClient: server.Client()}
		api := apis.NewAdminAPI(mockClient)

		err := api.UpdateGlobalRule(context.Background(), models.RuleValidity, models.ValidityLevelFull)
		assert.NoError(t, err)
	})

	t.Run("NotFound", func(t *testing.T) {
		errorResponse := models.APIError{Status: http.StatusNotFound, Title: TitleNotFound}
		server := setupMockServer(t, http.StatusNotFound, errorResponse, "/admin/rules/VALIDITY", http.MethodPut)
		defer server.Close()

		mockClient := &client.Client{BaseURL: server.URL, HTTPClient: server.Client()}
		api := apis.NewAdminAPI(mockClient)

		err := api.UpdateGlobalRule(context.Background(), models.RuleValidity, models.ValidityLevelFull)
		assertAPIError(t, err, http.StatusNotFound, TitleNotFound)
	})

	t.Run("BadRequest", func(t *testing.T) {
		errorResponse := models.APIError{Status: http.StatusBadRequest, Title: TitleBadRequest}
		server := setupMockServer(t, http.StatusBadRequest, errorResponse, "/admin/rules/VALIDITY", http.MethodPut)
		defer server.Close()

		mockClient := &client.Client{BaseURL: server.URL, HTTPClient: server.Client()}
		api := apis.NewAdminAPI(mockClient)

		err := api.UpdateGlobalRule(context.Background(), models.RuleValidity, models.ValidityLevelFull)
		assertAPIError(t, err, http.StatusBadRequest, TitleBadRequest)
	})

	t.Run("InternalServerError", func(t *testing.T) {
		errorResponse := models.APIError{Status: http.StatusInternalServerError, Title: TitleInternalServerError}
		server := setupMockServer(t, http.StatusInternalServerError, errorResponse, "/admin/rules/VALIDITY", http.MethodPut)
		defer server.Close()

		mockClient := &client.Client{BaseURL: server.URL, HTTPClient: server.Client()}
		api := apis.NewAdminAPI(mockClient)

		err := api.UpdateGlobalRule(context.Background(), models.RuleValidity, models.ValidityLevelFull)
		assertAPIError(t, err, http.StatusInternalServerError, TitleInternalServerError)
	})
}

func TestAdminAPI_DeleteGlobalRule(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		server := setupMockServer(t, http.StatusNoContent, nil, "/admin/rules/VALIDITY", http.MethodDelete)
		defer server.Close()

		mockClient := &client.Client{BaseURL: server.URL, HTTPClient: server.Client()}
		api := apis.NewAdminAPI(mockClient)

		err := api.DeleteGlobalRule(context.Background(), models.RuleValidity)
		assert.NoError(t, err)
	})

	t.Run("NotFound", func(t *testing.T) {
		errorResponse := models.APIError{Status: http.StatusNotFound, Title: TitleNotFound}
		server := setupMockServer(t, http.StatusNotFound, errorResponse, "/admin/rules/VALIDITY", http.MethodDelete)
		defer server.Close()

		mockClient := &client.Client{BaseURL: server.URL, HTTPClient: server.Client()}
		api := apis.NewAdminAPI(mockClient)

		err := api.DeleteGlobalRule(context.Background(), models.RuleValidity)
		assertAPIError(t, err, http.StatusNotFound, TitleNotFound)
	})

	t.Run("InternalServerError", func(t *testing.T) {
		errorResponse := models.APIError{Status: http.StatusInternalServerError, Title: TitleInternalServerError}
		server := setupMockServer(t, http.StatusInternalServerError, errorResponse, "/admin/rules/VALIDITY", http.MethodDelete)
		defer server.Close()

		mockClient := &client.Client{BaseURL: server.URL, HTTPClient: server.Client()}
		api := apis.NewAdminAPI(mockClient)

		err := api.DeleteGlobalRule(context.Background(), models.RuleValidity)
		assertAPIError(t, err, http.StatusInternalServerError, TitleInternalServerError)
	})
}

func TestAdminAPI_ListArtifactTypes(t *testing.T) {
	mockResponse := []models.ArtifactTypeResponse{
		{Name: models.Avro},
		{Name: models.Protobuf},
	}

	t.Run("Success", func(t *testing.T) {
		server := setupMockServer(t, http.StatusOK, mockResponse, "/admin/config/artifactTypes", http.MethodGet)
		defer server.Close()

		mockClient := &client.Client{BaseURL: server.URL, HTTPClient: server.Client()}
		api := apis.NewAdminAPI(mockClient)

		result, err := api.ListArtifactTypes(context.Background())
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Len(t, result, 2)
	})

	t.Run("InternalServerError", func(t *testing.T) {
		errorResponse := models.APIError{Status: http.StatusInternalServerError, Title: TitleInternalServerError}
		server := setupMockServer(t, http.StatusInternalServerError, errorResponse, "/admin/config/artifactTypes", http.MethodGet)
		defer server.Close()

		mockClient := &client.Client{BaseURL: server.URL, HTTPClient: server.Client()}
		api := apis.NewAdminAPI(mockClient)

		result, err := api.ListArtifactTypes(context.Background())
		assert.Error(t, err)
		assert.Nil(t, result)
		assertAPIError(t, err, http.StatusInternalServerError, TitleInternalServerError)
	})
}

/***********************/
/***** Integration *****/
/***********************/

func TestAdminAPI_Rules_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	ctx := context.Background()
	adminAPI := setupAdminAPIClient()

	// Test ListGlobalRules
	t.Run("ListGlobalRulesToGetEmptyResults", func(t *testing.T) {
		// Delete all rules
		err := adminAPI.DeleteAllGlobalRule(ctx)
		assert.NoError(t, err)

		rules, err := adminAPI.ListGlobalRules(ctx)
		assert.NoError(t, err)
		assert.NotNil(t, rules)
	})

	// Test CreateGlobalRule
	t.Run("CreateGlobalValidityRuleAndCheckIfItApplied", func(t *testing.T) {
		// Delete all rules
		err := adminAPI.DeleteAllGlobalRule(ctx)
		assert.NoError(t, err)

		// Create a new rule
		err = adminAPI.CreateGlobalRule(ctx, models.RuleValidity, models.ValidityLevelFull)
		assert.NoError(t, err)

		// Verify rule creation
		rules, err := adminAPI.ListGlobalRules(ctx)
		assert.NoError(t, err)
		assert.Contains(t, rules, models.RuleValidity)
	})

	// Test GetGlobalRule
	t.Run("GetGlobalRule", func(t *testing.T) {
		// Delete all rules
		err := adminAPI.DeleteAllGlobalRule(ctx)
		assert.NoError(t, err)

		// Create a new rule
		err = adminAPI.CreateGlobalRule(ctx, models.RuleValidity, models.ValidityLevelFull)
		assert.NoError(t, err)

		level, err := adminAPI.GetGlobalRule(ctx, models.RuleValidity)
		assert.NoError(t, err)
		assert.Equal(t, models.ValidityLevelFull, level)
	})

	// Test UpdateGlobalRule
	t.Run("UpdateGlobalRule", func(t *testing.T) {
		// Delete all rules
		err := adminAPI.DeleteAllGlobalRule(ctx)
		assert.NoError(t, err)

		// Create a new rule
		err = adminAPI.CreateGlobalRule(ctx, models.RuleValidity, models.ValidityLevelFull)
		assert.NoError(t, err)

		// Update the rule
		err = adminAPI.UpdateGlobalRule(ctx, models.RuleValidity, models.ValidityLevelSyntaxOnly)
		assert.NoError(t, err)

		// Verify rule update
		level, err := adminAPI.GetGlobalRule(ctx, models.RuleValidity)
		assert.NoError(t, err)
		assert.Equal(t, models.ValidityLevelSyntaxOnly, level)
	})

	// Test DeleteGlobalRule
	t.Run("DeleteGlobalRuleAndDeleteAllGlobalRules", func(t *testing.T) {
		// Delete all rules
		err := adminAPI.DeleteAllGlobalRule(ctx)
		assert.NoError(t, err)

		// Create a new rules
		err = adminAPI.CreateGlobalRule(ctx, models.RuleValidity, models.ValidityLevelFull)
		assert.NoError(t, err)

		err = adminAPI.CreateGlobalRule(ctx, models.RuleCompatibility, models.CompatibilityLevelFull)
		assert.NoError(t, err)

		err = adminAPI.CreateGlobalRule(ctx, models.RuleIntegrity, models.IntegrityLevelFull)
		assert.NoError(t, err)

		// Delete the validity rule
		err = adminAPI.DeleteGlobalRule(ctx, models.RuleValidity)
		assert.NoError(t, err)

		// Verify rule deletion
		rules, err := adminAPI.ListGlobalRules(ctx)
		assert.NoError(t, err)
		assert.Equal(t, 2, len(rules))
		assert.NotContains(t, rules, models.RuleValidity)

		// Delete all rules again
		// Delete all rules
		err = adminAPI.DeleteAllGlobalRule(ctx)
		assert.NoError(t, err)

		// Verify all rules are deleted
		rules, err = adminAPI.ListGlobalRules(ctx)
		assert.NoError(t, err)
		assert.Empty(t, rules)

	})

	// List Artifact Types
	t.Run("ListArtifactTypes", func(t *testing.T) {
		list, err := adminAPI.ListArtifactTypes(ctx)
		assert.NoError(t, err)
		assert.NotNil(t, list)
		assert.NotEmpty(t, list)
		assert.Contains(t, list, models.Protobuf)
		assert.Contains(t, list, models.OpenAPI)
		assert.Contains(t, list, models.AsyncAPI)
		assert.Contains(t, list, models.Json)
		assert.Contains(t, list, models.Avro)
		assert.Contains(t, list, models.GraphQL)
		assert.Contains(t, list, models.KConnect)
		assert.Contains(t, list, models.WSDL)
		assert.Contains(t, list, models.XSD)
		assert.Contains(t, list, models.XML)

	})
}

func setupAdminAPIClient() *apis.AdminAPI {
	apiClient := setupHTTPClient()
	return apis.NewAdminAPI(apiClient)
}
