package apis_test

import (
	"context"
	"fmt"
	"github.com/mollie/go-apicurio-registry/apis"
	"github.com/mollie/go-apicurio-registry/client"
	"github.com/mollie/go-apicurio-registry/models"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func TestGroupAPI_ListGroups(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mockGroups := []models.GroupInfo{{GroupId: "group1"}, {GroupId: "group2"}}
		mockResponse := models.GroupInfoResponse{Groups: mockGroups}

		server := setupMockServer(t, http.StatusOK, mockResponse, "/groups", http.MethodGet)
		defer server.Close()

		mockClient := &client.Client{BaseURL: server.URL, HTTPClient: server.Client()}
		groupAPI := apis.NewGroupAPI(mockClient)

		result, err := groupAPI.ListGroups(context.Background(), nil)
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Len(t, result, 2)
	})

	t.Run("Internal Server Error", func(t *testing.T) {
		errorResponse := models.APIError{Status: http.StatusInternalServerError, Title: TitleInternalServerError}

		server := setupMockServer(t, http.StatusInternalServerError, errorResponse, "/groups", http.MethodGet)
		defer server.Close()

		mockClient := &client.Client{BaseURL: server.URL, HTTPClient: server.Client()}
		groupAPI := apis.NewGroupAPI(mockClient)

		result, err := groupAPI.ListGroups(context.Background(), nil)
		assert.Error(t, err)
		assert.Nil(t, result)
		assertAPIError(t, err, http.StatusInternalServerError, TitleInternalServerError)
	})
}

func TestGroupAPI_CreateGroup(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mockGroup := models.GroupInfo{GroupId: "group1"}

		server := setupMockServer(t, http.StatusOK, mockGroup, "/groups", http.MethodPost)
		defer server.Close()

		mockClient := &client.Client{BaseURL: server.URL, HTTPClient: server.Client()}
		groupAPI := apis.NewGroupAPI(mockClient)

		result, err := groupAPI.CreateGroup(context.Background(), "group1", "description", nil)
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, "group1", result.GroupId)
	})

	t.Run("Validation: Empty Group ID", func(t *testing.T) {
		mockClient := &client.Client{BaseURL: "http://example.com", HTTPClient: http.DefaultClient}
		groupAPI := apis.NewGroupAPI(mockClient)

		result, err := groupAPI.CreateGroup(context.Background(), "", "description", nil)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "Group ID=''")
		assert.Nil(t, result)
	})

	t.Run("Conflict", func(t *testing.T) {
		errorResponse := models.APIError{Status: http.StatusConflict, Title: TitleConflict}

		server := setupMockServer(t, http.StatusConflict, errorResponse, "/groups", http.MethodPost)
		defer server.Close()

		mockClient := &client.Client{BaseURL: server.URL, HTTPClient: server.Client()}
		groupAPI := apis.NewGroupAPI(mockClient)

		result, err := groupAPI.CreateGroup(context.Background(), "group1", "description", nil)
		assert.Error(t, err)
		assert.Nil(t, result)
		assertAPIError(t, err, http.StatusConflict, TitleConflict)
	})

	t.Run("Internal Server Error", func(t *testing.T) {
		errorResponse := models.APIError{Status: http.StatusInternalServerError, Title: TitleInternalServerError}

		server := setupMockServer(t, http.StatusInternalServerError, errorResponse, "/groups", http.MethodPost)
		defer server.Close()

		mockClient := &client.Client{BaseURL: server.URL, HTTPClient: server.Client()}
		groupAPI := apis.NewGroupAPI(mockClient)

		result, err := groupAPI.CreateGroup(context.Background(), "group1", "description", nil)
		assert.Error(t, err)
		assert.Nil(t, result)
		assertAPIError(t, err, http.StatusInternalServerError, TitleInternalServerError)
	})
}

func TestGroupAPI_GetGroupById(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mockGroup := models.GroupInfo{GroupId: "group1"}

		server := setupMockServer(t, http.StatusOK, mockGroup, "/groups/group1", http.MethodGet)
		defer server.Close()

		mockClient := &client.Client{BaseURL: server.URL, HTTPClient: server.Client()}
		groupAPI := apis.NewGroupAPI(mockClient)

		result, err := groupAPI.GetGroupById(context.Background(), "group1")
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, "group1", result.GroupId)
	})

	t.Run("Validation: Empty Group ID", func(t *testing.T) {
		mockClient := &client.Client{BaseURL: "http://example.com", HTTPClient: http.DefaultClient}
		groupAPI := apis.NewGroupAPI(mockClient)

		result, err := groupAPI.GetGroupById(context.Background(), "")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "Group ID=''")
		assert.Nil(t, result)
	})

	t.Run("NotFound", func(t *testing.T) {
		errorResponse := models.APIError{Status: http.StatusNotFound, Title: TitleNotFound}

		server := setupMockServer(t, http.StatusNotFound, errorResponse, "/groups/group1", http.MethodGet)
		defer server.Close()

		mockClient := &client.Client{BaseURL: server.URL, HTTPClient: server.Client()}
		groupAPI := apis.NewGroupAPI(mockClient)

		result, err := groupAPI.GetGroupById(context.Background(), "group1")
		assert.Error(t, err)
		assert.Nil(t, result)
		assertAPIError(t, err, http.StatusNotFound, TitleNotFound)
	})

	t.Run("Internal Server Error", func(t *testing.T) {
		errorResponse := models.APIError{Status: http.StatusInternalServerError, Title: TitleInternalServerError}

		server := setupMockServer(t, http.StatusInternalServerError, errorResponse, "/groups/group1", http.MethodGet)
		defer server.Close()

		mockClient := &client.Client{BaseURL: server.URL, HTTPClient: server.Client()}
		groupAPI := apis.NewGroupAPI(mockClient)

		result, err := groupAPI.GetGroupById(context.Background(), "group1")
		assert.Error(t, err)
		assert.Nil(t, result)
		assertAPIError(t, err, http.StatusInternalServerError, TitleInternalServerError)
	})
}

func TestGroupAPI_UpdateGroupMetadata(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		server := setupMockServer(t, http.StatusNoContent, nil,
			"/groups/group1", http.MethodPut)
		defer server.Close()

		mockClient := &client.Client{BaseURL: server.URL, HTTPClient: server.Client()}
		groupAPI := apis.NewGroupAPI(mockClient)

		err := groupAPI.UpdateGroupMetadata(context.Background(), "group1", "description", nil)
		assert.NoError(t, err)
	})

	t.Run("Validation: Empty Group ID", func(t *testing.T) {
		mockClient := &client.Client{BaseURL: "http://example.com", HTTPClient: http.DefaultClient}
		groupAPI := apis.NewGroupAPI(mockClient)

		err := groupAPI.UpdateGroupMetadata(context.Background(), "", "description", nil)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "Group ID=''")
	})

	t.Run("NotFound", func(t *testing.T) {
		errorResponse := models.APIError{Status: http.StatusNotFound, Title: TitleNotFound}

		server := setupMockServer(t, http.StatusNotFound, errorResponse,
			"/groups/group1", http.MethodPut)
		defer server.Close()

		mockClient := &client.Client{BaseURL: server.URL, HTTPClient: server.Client()}
		groupAPI := apis.NewGroupAPI(mockClient)

		err := groupAPI.UpdateGroupMetadata(context.Background(), "group1", "description", nil)
		assert.Error(t, err)
		assertAPIError(t, err, http.StatusNotFound, TitleNotFound)
	})

	t.Run("Internal Server Error", func(t *testing.T) {
		errorResponse := models.APIError{Status: http.StatusInternalServerError, Title: TitleInternalServerError}

		server := setupMockServer(t, http.StatusInternalServerError, errorResponse,
			"/groups/group1", http.MethodPut)
		defer server.Close()

		mockClient := &client.Client{BaseURL: server.URL, HTTPClient: server.Client()}
		groupAPI := apis.NewGroupAPI(mockClient)

		err := groupAPI.UpdateGroupMetadata(context.Background(), "group1", "description", nil)
		assert.Error(t, err)
		assertAPIError(t, err, http.StatusInternalServerError, TitleInternalServerError)
	})
}

func TestGroupAPI_DeleteGroup(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		server := setupMockServer(t, http.StatusNoContent, nil,
			"/groups/group1", http.MethodDelete)
		defer server.Close()

		mockClient := &client.Client{BaseURL: server.URL, HTTPClient: server.Client()}
		groupAPI := apis.NewGroupAPI(mockClient)

		err := groupAPI.DeleteGroup(context.Background(), "group1")
		assert.NoError(t, err)
	})

	t.Run("Validation: Empty Group ID", func(t *testing.T) {
		mockClient := &client.Client{BaseURL: "http://example.com", HTTPClient: http.DefaultClient}
		groupAPI := apis.NewGroupAPI(mockClient)

		err := groupAPI.DeleteGroup(context.Background(), "")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "Group ID=''")
	})

	t.Run("NotFound", func(t *testing.T) {
		errorResponse := models.APIError{Status: http.StatusNotFound, Title: TitleNotFound}

		server := setupMockServer(t, http.StatusNotFound, errorResponse,
			"/groups/group1", http.MethodDelete)
		defer server.Close()

		mockClient := &client.Client{BaseURL: server.URL, HTTPClient: server.Client()}
		groupAPI := apis.NewGroupAPI(mockClient)

		err := groupAPI.DeleteGroup(context.Background(), "group1")
		assert.Error(t, err)
		assertAPIError(t, err, http.StatusNotFound, TitleNotFound)
	})

	t.Run("Not Allowed", func(t *testing.T) {
		errorResponse := models.APIError{Status: http.StatusMethodNotAllowed, Title: TitleMethodNotAllowed}

		server := setupMockServer(t, http.StatusMethodNotAllowed, errorResponse,
			"/groups/group1", http.MethodDelete)
		defer server.Close()

		mockClient := &client.Client{BaseURL: server.URL, HTTPClient: server.Client()}
		groupAPI := apis.NewGroupAPI(mockClient)

		err := groupAPI.DeleteGroup(context.Background(), "group1")
		assert.Error(t, err)
		assertAPIError(t, err, http.StatusMethodNotAllowed, TitleMethodNotAllowed)
	})

	t.Run("Internal Server Error", func(t *testing.T) {
		errorResponse := models.APIError{Status: http.StatusInternalServerError, Title: TitleInternalServerError}

		server := setupMockServer(t, http.StatusInternalServerError, errorResponse,
			"/groups/group1", http.MethodDelete)
		defer server.Close()

		mockClient := &client.Client{BaseURL: server.URL, HTTPClient: server.Client()}
		groupAPI := apis.NewGroupAPI(mockClient)

		err := groupAPI.DeleteGroup(context.Background(), "group1")
		assert.Error(t, err)
		assertAPIError(t, err, http.StatusInternalServerError, TitleInternalServerError)
	})
}

func TestGroupAPI_SearchGroups(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mockGroups := []models.GroupInfo{{GroupId: "group1"}, {GroupId: "group2"}}
		mockResponse := models.GroupInfoResponse{Groups: mockGroups}

		server := setupMockServer(t, http.StatusOK, mockResponse, "/search/groups", http.MethodGet)
		defer server.Close()

		mockClient := &client.Client{BaseURL: server.URL, HTTPClient: server.Client()}
		groupAPI := apis.NewGroupAPI(mockClient)

		result, err := groupAPI.SearchGroups(context.Background(), nil)
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Len(t, result, 2)
	})

	t.Run("Internal Server Error", func(t *testing.T) {
		errorResponse := models.APIError{Status: http.StatusInternalServerError, Title: TitleInternalServerError}

		server := setupMockServer(t, http.StatusInternalServerError, errorResponse, "/search/groups", http.MethodGet)
		defer server.Close()

		mockClient := &client.Client{BaseURL: server.URL, HTTPClient: server.Client()}
		groupAPI := apis.NewGroupAPI(mockClient)

		result, err := groupAPI.SearchGroups(context.Background(), nil)
		assert.Error(t, err)
		assert.Nil(t, result)
		assertAPIError(t, err, http.StatusInternalServerError, TitleInternalServerError)
	})
}

func TestGroupsAPI_ListGroupRules(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mockReferences := []models.Rule{models.RuleValidity, models.RuleCompatibility}

		server := setupMockServer(t, http.StatusOK, mockReferences,
			fmt.Sprintf("/groups/%s/rules", stubGroupId), http.MethodGet)
		defer server.Close()

		mockClient := &client.Client{BaseURL: server.URL, HTTPClient: server.Client()}
		api := apis.NewGroupAPI(mockClient)

		result, err := api.ListGroupRules(context.Background(), stubGroupId)
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Len(t, result, 2)
	})

	t.Run("Validation: Empty Group ID", func(t *testing.T) {
		mockClient := &client.Client{BaseURL: "http://example.com", HTTPClient: http.DefaultClient}
		api := apis.NewGroupAPI(mockClient)

		result, err := api.ListGroupRules(context.Background(), "")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "Group ID=''")
		assert.Nil(t, result)
	})

	t.Run("Not Found", func(t *testing.T) {
		errorResponse := models.APIError{Status: http.StatusNotFound, Title: TitleNotFound}

		server := setupMockServer(t, http.StatusNotFound, errorResponse,
			fmt.Sprintf("/groups/%s/rules", stubGroupId), http.MethodGet)
		defer server.Close()

		mockClient := &client.Client{BaseURL: server.URL, HTTPClient: server.Client()}
		api := apis.NewGroupAPI(mockClient)

		result, err := api.ListGroupRules(context.Background(), stubGroupId)
		assert.Error(t, err)
		assert.Nil(t, result)
		assertAPIError(t, err, http.StatusNotFound, TitleNotFound)
	})

	t.Run("Internal Server Error", func(t *testing.T) {
		errorResponse := models.APIError{Status: http.StatusInternalServerError, Title: TitleInternalServerError}

		server := setupMockServer(t, http.StatusInternalServerError, errorResponse,
			fmt.Sprintf("/groups/%s/rules", stubGroupId), http.MethodGet)
		defer server.Close()

		mockClient := &client.Client{BaseURL: server.URL, HTTPClient: server.Client()}
		api := apis.NewGroupAPI(mockClient)

		result, err := api.ListGroupRules(context.Background(), stubGroupId)
		assert.Error(t, err)
		assert.Nil(t, result)
		assertAPIError(t, err, http.StatusInternalServerError, TitleInternalServerError)
	})
}

func TestGroupsAPI_CreateGroupRule(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		server := setupMockServer(t, http.StatusNoContent, nil,
			fmt.Sprintf("/groups/%s/rules", stubGroupId), http.MethodPost)
		defer server.Close()

		mockClient := &client.Client{BaseURL: server.URL, HTTPClient: server.Client()}
		api := apis.NewGroupAPI(mockClient)

		err := api.CreateGroupRule(context.Background(), stubGroupId, models.RuleValidity, models.ValidityLevelFull)
		assert.NoError(t, err)
	})

	t.Run("Validation: Empty Group ID", func(t *testing.T) {
		mockClient := &client.Client{BaseURL: "http://example.com", HTTPClient: http.DefaultClient}
		api := apis.NewGroupAPI(mockClient)

		err := api.CreateGroupRule(context.Background(), "", models.RuleValidity, models.ValidityLevelFull)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "Group ID=''")
	})

	t.Run("BadRequest", func(t *testing.T) {
		errorResponse := models.APIError{Status: http.StatusBadRequest, Title: TitleBadRequest}

		server := setupMockServer(t, http.StatusBadRequest, errorResponse,
			fmt.Sprintf("/groups/%s/rules", stubGroupId), http.MethodPost)
		defer server.Close()

		mockClient := &client.Client{BaseURL: server.URL, HTTPClient: server.Client()}
		api := apis.NewGroupAPI(mockClient)

		err := api.CreateGroupRule(context.Background(), stubGroupId, models.RuleValidity, models.ValidityLevelFull)
		assert.Error(t, err)
		assertAPIError(t, err, http.StatusBadRequest, TitleBadRequest)
	})

	t.Run("Conflict", func(t *testing.T) {
		errorResponse := models.APIError{Status: http.StatusConflict, Title: TitleConflict}

		server := setupMockServer(t, http.StatusConflict, errorResponse,
			fmt.Sprintf("/groups/%s/rules", stubGroupId), http.MethodPost)
		defer server.Close()

		mockClient := &client.Client{BaseURL: server.URL, HTTPClient: server.Client()}
		api := apis.NewGroupAPI(mockClient)

		err := api.CreateGroupRule(context.Background(), stubGroupId, models.RuleValidity, models.ValidityLevelFull)
		assert.Error(t, err)
		assertAPIError(t, err, http.StatusConflict, TitleConflict)
	})

	t.Run("NotFound", func(t *testing.T) {
		errorResponse := models.APIError{Status: http.StatusNotFound, Title: TitleNotFound}

		server := setupMockServer(t, http.StatusNotFound, errorResponse,
			fmt.Sprintf("/groups/%s/rules", stubGroupId), http.MethodPost)
		defer server.Close()

		mockClient := &client.Client{BaseURL: server.URL, HTTPClient: server.Client()}
		api := apis.NewGroupAPI(mockClient)

		err := api.CreateGroupRule(context.Background(), stubGroupId, models.RuleValidity, models.ValidityLevelFull)
		assert.Error(t, err)
		assertAPIError(t, err, http.StatusNotFound, TitleNotFound)
	})

	t.Run("InternalServerError", func(t *testing.T) {
		errorResponse := models.APIError{Status: http.StatusInternalServerError, Title: TitleInternalServerError}

		server := setupMockServer(t, http.StatusInternalServerError, errorResponse,
			fmt.Sprintf("/groups/%s/rules", stubGroupId), http.MethodPost)
		defer server.Close()

		mockClient := &client.Client{BaseURL: server.URL, HTTPClient: server.Client()}
		api := apis.NewGroupAPI(mockClient)

		err := api.CreateGroupRule(context.Background(), stubGroupId, models.RuleValidity, models.ValidityLevelFull)
		assert.Error(t, err)
		assertAPIError(t, err, http.StatusInternalServerError, TitleInternalServerError)
	})
}

func TestGroupsAPI_DeleteAllGroupRule(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		server := setupMockServer(t, http.StatusNoContent, nil,
			fmt.Sprintf("/groups/%s/rules", stubGroupId), http.MethodDelete)
		defer server.Close()

		mockClient := &client.Client{BaseURL: server.URL, HTTPClient: server.Client()}
		api := apis.NewGroupAPI(mockClient)

		err := api.DeleteAllGroupRule(context.Background(), stubGroupId)
		assert.NoError(t, err)
	})

	t.Run("Validation: Empty Group ID", func(t *testing.T) {
		mockClient := &client.Client{BaseURL: "http://example.com", HTTPClient: http.DefaultClient}
		api := apis.NewGroupAPI(mockClient)

		err := api.DeleteAllGroupRule(context.Background(), "")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "Group ID=''")
	})

	t.Run("NotFound", func(t *testing.T) {
		errorResponse := models.APIError{Status: http.StatusNotFound, Title: TitleNotFound}

		server := setupMockServer(t, http.StatusNotFound, errorResponse,
			fmt.Sprintf("/groups/%s/rules", stubGroupId), http.MethodDelete)
		defer server.Close()

		mockClient := &client.Client{BaseURL: server.URL, HTTPClient: server.Client()}
		api := apis.NewGroupAPI(mockClient)

		err := api.DeleteAllGroupRule(context.Background(), stubGroupId)
		assert.Error(t, err)
		assertAPIError(t, err, http.StatusNotFound, TitleNotFound)
	})

	t.Run("InternalServerError", func(t *testing.T) {
		errorResponse := models.APIError{Status: http.StatusInternalServerError, Title: TitleInternalServerError}

		server := setupMockServer(t, http.StatusInternalServerError, errorResponse,
			fmt.Sprintf("/groups/%s/rules", stubGroupId), http.MethodDelete)
		defer server.Close()

		mockClient := &client.Client{BaseURL: server.URL, HTTPClient: server.Client()}
		api := apis.NewGroupAPI(mockClient)

		err := api.DeleteAllGroupRule(context.Background(), stubGroupId)
		assert.Error(t, err)
		assertAPIError(t, err, http.StatusInternalServerError, TitleInternalServerError)
	})
}

func TestGroupsAPI_GetGroupRule(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mockResponse := models.RuleResponse{
			RuleType: models.RuleValidity,
			Config:   models.ValidityLevelFull,
		}
		mockRule := models.RuleValidity

		server := setupMockServer(t, http.StatusOK, mockResponse,
			fmt.Sprintf("/groups/%s/rules/%s", stubGroupId, mockRule), http.MethodGet)
		defer server.Close()

		mockClient := &client.Client{BaseURL: server.URL, HTTPClient: server.Client()}
		api := apis.NewGroupAPI(mockClient)

		result, err := api.GetGroupRule(context.Background(), stubGroupId, mockRule)
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, models.ValidityLevelFull, result)
	})

	t.Run("Validation: Empty Group ID", func(t *testing.T) {
		mockClient := &client.Client{BaseURL: "http://example.com", HTTPClient: http.DefaultClient}
		api := apis.NewGroupAPI(mockClient)

		result, err := api.GetGroupRule(context.Background(), "", models.RuleValidity)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "Group ID=''")
		assert.Empty(t, result)
	})

	t.Run("NotFound", func(t *testing.T) {
		mockRule := models.RuleValidity
		errorResponse := models.APIError{Status: http.StatusNotFound, Title: TitleNotFound}

		server := setupMockServer(t, http.StatusNotFound, errorResponse,
			fmt.Sprintf("/groups/%s/rules/%s", stubGroupId, mockRule), http.MethodGet)
		defer server.Close()

		mockClient := &client.Client{BaseURL: server.URL, HTTPClient: server.Client()}
		api := apis.NewGroupAPI(mockClient)

		result, err := api.GetGroupRule(context.Background(), stubGroupId, mockRule)
		assert.Error(t, err)
		assert.Empty(t, result)
		assertAPIError(t, err, http.StatusNotFound, TitleNotFound)
	})

	t.Run("InternalServerError", func(t *testing.T) {
		mockRule := models.RuleValidity
		errorResponse := models.APIError{Status: http.StatusInternalServerError, Title: TitleInternalServerError}

		server := setupMockServer(t, http.StatusInternalServerError, errorResponse,
			fmt.Sprintf("/groups/%s/rules/%s", stubGroupId, mockRule), http.MethodGet)
		defer server.Close()

		mockClient := &client.Client{BaseURL: server.URL, HTTPClient: server.Client()}
		api := apis.NewGroupAPI(mockClient)

		result, err := api.GetGroupRule(context.Background(), stubGroupId, mockRule)
		assert.Error(t, err)
		assert.Empty(t, result)
		assertAPIError(t, err, http.StatusInternalServerError, TitleInternalServerError)
	})
}

func TestGroupsAPI_UpdateGroupRule(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mockRule := models.RuleValidity
		mockResponse := models.RuleResponse{
			RuleType: mockRule,
			Config:   models.ValidityLevelFull,
		}

		server := setupMockServer(t, http.StatusOK, mockResponse,
			fmt.Sprintf("/groups/%s/rules/%s", stubGroupId, mockRule), http.MethodPut)
		defer server.Close()

		mockClient := &client.Client{BaseURL: server.URL, HTTPClient: server.Client()}
		api := apis.NewGroupAPI(mockClient)

		err := api.UpdateGroupRule(context.Background(), stubGroupId, mockRule, models.ValidityLevelFull)
		assert.NoError(t, err)
	})

	t.Run("Validation: Empty Group ID", func(t *testing.T) {
		mockClient := &client.Client{BaseURL: "http://example.com", HTTPClient: http.DefaultClient}
		api := apis.NewGroupAPI(mockClient)

		err := api.UpdateGroupRule(context.Background(), "", models.RuleValidity, models.ValidityLevelFull)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "Group ID=''")
	})

	t.Run("NotFound", func(t *testing.T) {
		mockRule := models.RuleValidity
		errorResponse := models.APIError{Status: http.StatusNotFound, Title: TitleNotFound}

		server := setupMockServer(t, http.StatusNotFound, errorResponse,
			fmt.Sprintf("/groups/%s/rules/%s", stubGroupId, mockRule), http.MethodPut)
		defer server.Close()

		mockClient := &client.Client{BaseURL: server.URL, HTTPClient: server.Client()}
		api := apis.NewGroupAPI(mockClient)

		err := api.UpdateGroupRule(context.Background(), stubGroupId, mockRule, models.ValidityLevelFull)
		assert.Error(t, err)
		assertAPIError(t, err, http.StatusNotFound, TitleNotFound)
	})

	t.Run("InternalServerError", func(t *testing.T) {
		mockRule := models.RuleValidity
		errorResponse := models.APIError{Status: http.StatusInternalServerError, Title: TitleInternalServerError}

		server := setupMockServer(t, http.StatusInternalServerError, errorResponse,
			fmt.Sprintf("/groups/%s/rules/%s", stubGroupId, mockRule), http.MethodPut)
		defer server.Close()

		mockClient := &client.Client{BaseURL: server.URL, HTTPClient: server.Client()}
		api := apis.NewGroupAPI(mockClient)

		err := api.UpdateGroupRule(context.Background(), stubGroupId, mockRule, models.ValidityLevelFull)
		assert.Error(t, err)
		assertAPIError(t, err, http.StatusInternalServerError, TitleInternalServerError)
	})
}

func TestGroupsAPI_DeleteGroupRule(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mockRule := models.RuleValidity
		server := setupMockServer(t, http.StatusNoContent, nil,
			fmt.Sprintf("/groups/%s/rules/%s", stubGroupId, mockRule), http.MethodDelete)
		defer server.Close()

		mockClient := &client.Client{BaseURL: server.URL, HTTPClient: server.Client()}
		api := apis.NewGroupAPI(mockClient)

		err := api.DeleteGroupRule(context.Background(), stubGroupId, mockRule)
		assert.NoError(t, err)
	})

	t.Run("Validation: Empty Group ID", func(t *testing.T) {
		mockClient := &client.Client{BaseURL: "http://example.com", HTTPClient: http.DefaultClient}
		api := apis.NewGroupAPI(mockClient)

		err := api.DeleteGroupRule(context.Background(), "", models.RuleValidity)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "Group ID=''")
	})

	t.Run("NotFound", func(t *testing.T) {
		mockRule := models.RuleValidity
		errorResponse := models.APIError{Status: http.StatusNotFound, Title: TitleNotFound}

		server := setupMockServer(t, http.StatusNotFound, errorResponse,
			fmt.Sprintf("/groups/%s/rules/%s", stubGroupId, mockRule), http.MethodDelete)
		defer server.Close()

		mockClient := &client.Client{BaseURL: server.URL, HTTPClient: server.Client()}
		api := apis.NewGroupAPI(mockClient)

		err := api.DeleteGroupRule(context.Background(), stubGroupId, mockRule)
		assert.Error(t, err)
		assertAPIError(t, err, http.StatusNotFound, TitleNotFound)
	})

	t.Run("InternalServerError", func(t *testing.T) {
		mockRule := models.RuleValidity
		errorResponse := models.APIError{Status: http.StatusInternalServerError, Title: TitleInternalServerError}

		server := setupMockServer(t, http.StatusInternalServerError, errorResponse,
			fmt.Sprintf("/groups/%s/rules/%s", stubGroupId, mockRule), http.MethodDelete)
		defer server.Close()

		mockClient := &client.Client{BaseURL: server.URL, HTTPClient: server.Client()}
		api := apis.NewGroupAPI(mockClient)

		err := api.DeleteGroupRule(context.Background(), stubGroupId, mockRule)
		assert.Error(t, err)
		assertAPIError(t, err, http.StatusInternalServerError, TitleInternalServerError)
	})
}

/***********************/
/***** Integration *****/
/***********************/

func TestGroupsAPIIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	ctx := context.Background()

	groupAPI := setupGroupsAPIClient()

	// Prepare test data
	artifactsAPI := apis.NewArtifactsAPI(groupAPI.Client)

	// Clean up before and after tests
	t.Cleanup(func() { cleanup(t, artifactsAPI) })
	cleanup(t, artifactsAPI)

	// Test CreateGroup
	t.Run("CreateGroup", func(t *testing.T) {
		groupInfo, randomGroupID, err := generateGroupForTest(ctx, groupAPI)
		assert.NoError(t, err)
		assert.Equal(t, randomGroupID, groupInfo.GroupId)
		assert.Equal(t, stubDescription, groupInfo.Description)
		// Clean up
		err = deleteGroupAfterTest(ctx, groupAPI, randomGroupID)
		assert.NoError(t, err)
	})

	// Test ListGroups
	t.Run("ListGroups", func(t *testing.T) {
		_, randomGroupID, err := generateGroupForTest(ctx, groupAPI)
		assert.NoError(t, err)

		resp, err := groupAPI.ListGroups(ctx, nil)
		assert.NoError(t, err)
		assert.GreaterOrEqual(t, len(resp), 1)

		// Clean up
		err = deleteGroupAfterTest(ctx, groupAPI, randomGroupID)
		assert.NoError(t, err)
	})

	// Test GetGroupByID
	t.Run("GetGroupByID", func(t *testing.T) {
		groupInfo, randomGroupID, err := generateGroupForTest(ctx, groupAPI)
		assert.NoError(t, err)

		resp, err := groupAPI.GetGroupById(ctx, randomGroupID)
		assert.NoError(t, err)
		assert.Equal(t, randomGroupID, resp.GroupId)
		assert.Equal(t, groupInfo.Description, resp.Description)
		assert.Equal(t, groupInfo.Labels, resp.Labels)

		// Clean up
		err = deleteGroupAfterTest(ctx, groupAPI, randomGroupID)
		assert.NoError(t, err)
	})

	// Test UpdateGroupMetadata
	t.Run("UpdateGroupMetadata", func(t *testing.T) {
		// Create a group
		_, randomGroupID, err := generateGroupForTest(ctx, groupAPI)
		assert.NoError(t, err)

		// Update the group
		err = groupAPI.UpdateGroupMetadata(ctx, randomGroupID, stubUpdatedDescription, stubUpdatedLabels)
		assert.NoError(t, err)

		// Get the group and verify the update
		resp, err := groupAPI.GetGroupById(ctx, randomGroupID)
		assert.Equal(t, randomGroupID, resp.GroupId)
		assert.Equal(t, stubUpdatedDescription, resp.Description)
		assert.Equal(t, stubUpdatedLabels, resp.Labels)

		// Clean up
		err = deleteGroupAfterTest(ctx, groupAPI, randomGroupID)
		assert.NoError(t, err)
	})

	// Test SearchGroups
	t.Run("SearchGroups", func(t *testing.T) {
		// Create a group
		_, randomGroupID, err := generateGroupForTest(ctx, groupAPI)
		assert.NoError(t, err)

		// Update the group
		res, err := groupAPI.SearchGroups(ctx, &models.SearchGroupsParams{
			GroupID: randomGroupID,
		})
		assert.NoError(t, err)
		assert.GreaterOrEqual(t, len(res), 1)

		// Clean up
		err = deleteGroupAfterTest(ctx, groupAPI, randomGroupID)
		assert.NoError(t, err)
	})

	// Test ListGroupRules
	t.Run("AllGroupRules", func(t *testing.T) {
		// Create a group
		_, randomGroupID, err := generateGroupForTest(ctx, groupAPI)
		assert.NoError(t, err)

		// List group rules
		rules, err := groupAPI.ListGroupRules(ctx, randomGroupID)
		assert.NoError(t, err)
		assert.Len(t, rules, 0)

		// Create a rule
		err = groupAPI.CreateGroupRule(ctx, randomGroupID, models.RuleValidity, models.ValidityLevelFull)
		assert.NoError(t, err)

		// List group rules
		rules, err = groupAPI.ListGroupRules(ctx, randomGroupID)
		assert.NoError(t, err)
		assert.Len(t, rules, 1)
		assert.Equal(t, models.RuleValidity, rules[0])

		// Get the rule
		rule, err := groupAPI.GetGroupRule(ctx, randomGroupID, models.RuleValidity)
		assert.NoError(t, err)
		assert.Equal(t, models.ValidityLevelFull, rule)

		// Update the rule
		err = groupAPI.UpdateGroupRule(ctx, randomGroupID, models.RuleValidity, models.ValidityLevelSyntaxOnly)
		assert.NoError(t, err)

		// Get the rule
		rule, err = groupAPI.GetGroupRule(ctx, randomGroupID, models.RuleValidity)
		assert.NoError(t, err)
		assert.Equal(t, models.ValidityLevelSyntaxOnly, rule)

		// Delete the rule
		err = groupAPI.DeleteGroupRule(ctx, randomGroupID, models.RuleValidity)
		assert.NoError(t, err)

		// List group rules
		rules, err = groupAPI.ListGroupRules(ctx, randomGroupID)
		assert.NoError(t, err)
		assert.Len(t, rules, 0)

		// Create three rules
		err = groupAPI.CreateGroupRule(ctx, randomGroupID, models.RuleValidity, models.ValidityLevelFull)
		assert.NoError(t, err)
		err = groupAPI.CreateGroupRule(ctx, randomGroupID, models.RuleCompatibility, models.CompatibilityLevelFull)
		assert.NoError(t, err)
		err = groupAPI.CreateGroupRule(ctx, randomGroupID, models.RuleIntegrity, models.IntegrityLevelFull)
		assert.NoError(t, err)

		// List group rules
		rules, err = groupAPI.ListGroupRules(ctx, randomGroupID)
		assert.NoError(t, err)
		assert.Len(t, rules, 3)

		// Delete all rules
		err = groupAPI.DeleteAllGroupRule(ctx, randomGroupID)
		assert.NoError(t, err)

		// List group rules
		rules, err = groupAPI.ListGroupRules(ctx, randomGroupID)
		assert.NoError(t, err)
		assert.Len(t, rules, 0)

		// Clean up
		err = deleteGroupAfterTest(ctx, groupAPI, randomGroupID)
		assert.NoError(t, err)
	})

}

func setupGroupsAPIClient() *apis.GroupAPI {
	apiClient := setupHTTPClient()
	return apis.NewGroupAPI(apiClient)
}
