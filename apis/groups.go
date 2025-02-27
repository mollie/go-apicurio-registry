package apis

import (
	"context"
	"fmt"
	"net/http"
	"net/url"

	"github.com/mollie/go-apicurio-registry/client"
	"github.com/mollie/go-apicurio-registry/models"
	"github.com/pkg/errors"
)

type GroupAPI struct {
	Client *client.Client
}

func NewGroupAPI(client *client.Client) *GroupAPI {
	return &GroupAPI{
		Client: client,
	}
}

// ListGroups Returns a list of all groups. This list is paged.
// See https://www.apicur.io/registry/docs/apicurio-registry/3.0.x/assets-attachments/registry-rest-api.htm#tag/Groups/operation/listGroups
func (api *GroupAPI) ListGroups(
	ctx context.Context,
	params *models.ListGroupsParams,
) ([]models.GroupInfo, error) {
	query := ""
	if params != nil {
		if err := params.Validate(); err != nil {
			return nil, errors.Wrap(err, "invalid parameters provided")
		}
		query = "?" + params.ToQuery().Encode()
	}

	urlPath := fmt.Sprintf("%s/groups%s", api.Client.BaseURL, query)
	resp, err := api.executeRequest(ctx, http.MethodGet, urlPath, nil)
	if err != nil {
		return nil, err
	}

	var result models.GroupInfoResponse
	if err := handleResponse(resp, http.StatusOK, &result); err != nil {
		return nil, err
	}

	return result.Groups, nil
}

// CreateGroup Creates a new group.
// See https://www.apicur.io/registry/docs/apicurio-registry/3.0.x/assets-attachments/registry-rest-api.htm#tag/Groups/operation/createGroup
func (api *GroupAPI) CreateGroup(
	ctx context.Context,
	groupId, description string,
	labels map[string]string,
) (*models.GroupInfo, error) {
	if err := validateInput(groupId, regexGroupIDArtifactID, "Group ID"); err != nil {
		return nil, err
	}

	urlPath := fmt.Sprintf("%s/groups", api.Client.BaseURL)
	body := models.CreateGroupRequest{
		GroupID:     groupId,
		Description: description,
		Labels:      labels,
	}

	resp, err := api.executeRequest(ctx, http.MethodPost, urlPath, body)
	if err != nil {
		return nil, err
	}

	var groupInfo models.GroupInfo
	err = handleResponse(resp, http.StatusOK, &groupInfo)
	if err != nil {
		return nil, err
	}

	return &groupInfo, nil

}

// GetGroupById Returns the group with the specified ID.
// See https://www.apicur.io/registry/docs/apicurio-registry/3.0.x/assets-attachments/registry-rest-api.htm#tag/Groups/operation/getGroupById
func (api *GroupAPI) GetGroupById(ctx context.Context, groupId string) (*models.GroupInfo, error) {
	if err := validateInput(groupId, regexGroupIDArtifactID, "Group ID"); err != nil {
		return nil, err
	}
	urlPath := fmt.Sprintf("%s/groups/%s", api.Client.BaseURL, url.PathEscape(groupId))

	resp, err := api.executeRequest(ctx, http.MethodGet, urlPath, nil)
	if err != nil {
		return nil, err
	}

	var groupInfo models.GroupInfo

	err = handleResponse(resp, http.StatusOK, &groupInfo)
	if err != nil {
		return nil, err
	}

	return &groupInfo, nil
}

// UpdateGroupMetadata Updates the metadata of the group with the specified ID.
// See https://www.apicur.io/registry/docs/apicurio-registry/3.0.x/assets-attachments/registry-rest-api.htm#tag/Groups/operation/updateGroupById
func (api *GroupAPI) UpdateGroupMetadata(
	ctx context.Context,
	groupId string,
	description string,
	labels map[string]string,
) error {
	if err := validateInput(groupId, regexGroupIDArtifactID, "Group ID"); err != nil {
		return err
	}

	urlPath := fmt.Sprintf("%s/groups/%s", api.Client.BaseURL, url.PathEscape(groupId))
	body := models.UpdateGroupRequest{
		Description: description,
		Labels:      labels,
	}

	resp, err := api.executeRequest(ctx, http.MethodPut, urlPath, body)
	if err != nil {
		return err
	}

	return handleResponse(resp, http.StatusNoContent, nil)
}

// DeleteGroup Deletes the group with the specified ID.
// See https://www.apicur.io/registry/docs/apicurio-registry/3.0.x/assets-attachments/registry-rest-api.htm#tag/Groups/operation/deleteGroupById
func (api *GroupAPI) DeleteGroup(ctx context.Context, groupId string) error {
	if err := validateInput(groupId, regexGroupIDArtifactID, "Group ID"); err != nil {
		return err
	}

	urlPath := fmt.Sprintf("%s/groups/%s", api.Client.BaseURL, url.PathEscape(groupId))

	resp, err := api.executeRequest(ctx, http.MethodDelete, urlPath, nil)
	if err != nil {
		return err
	}

	return handleResponse(resp, http.StatusNoContent, nil)

}

// SearchGroups Returns a list of groups that match the specified criteria. This list is paged.
// See https://www.apicur.io/registry/docs/apicurio-registry/3.0.x/assets-attachments/registry-rest-api.htm#tag/Groups/operation/searchGroups
func (api *GroupAPI) SearchGroups(
	ctx context.Context,
	params *models.SearchGroupsParams,
) ([]models.GroupInfo, error) {
	query := ""
	if params != nil {
		if err := params.Validate(); err != nil {
			return nil, errors.Wrap(err, "invalid parameters provided")
		}
		query = "?" + params.ToQuery().Encode()
	}

	urlPath := fmt.Sprintf("%s/search/groups%s", api.Client.BaseURL, query)
	resp, err := api.executeRequest(ctx, http.MethodGet, urlPath, nil)
	if err != nil {
		return nil, err
	}

	var result models.GroupInfoResponse
	if err := handleResponse(resp, http.StatusOK, &result); err != nil {
		return nil, err
	}

	return result.Groups, nil

}

// ListGroupRules Returns a list of all rules configured for the group.
// The set of rules determines how the content of an artifact in the group can evolve over time.
// If no rules are configured for a group, the set of globally configured rules are used.
// See https://www.apicur.io/registry/docs/apicurio-registry/3.0.x/assets-attachments/registry-rest-api.htm#tag/Group-rules/operation/listGroupRules
func (api *GroupAPI) ListGroupRules(ctx context.Context, groupID string) ([]models.Rule, error) {
	if err := validateInput(groupID, regexGroupIDArtifactID, "Group ID"); err != nil {
		return nil, err
	}

	urlPath := fmt.Sprintf("%s/groups/%s/rules", api.Client.BaseURL, url.PathEscape(groupID))
	resp, err := api.executeRequest(ctx, http.MethodGet, urlPath, nil)
	if err != nil {
		return nil, err
	}

	var rules []models.Rule
	if err := handleResponse(resp, http.StatusOK, &rules); err != nil {
		return nil, err
	}

	return rules, nil
}

// CreateGroupRule Adds a rule to the list of rules that get applied to an artifact in the group when adding new versions.
// All configured rules must pass to successfully add a new artifact version.
// See https://www.apicur.io/registry/docs/apicurio-registry/3.0.x/assets-attachments/registry-rest-api.htm#tag/Group-rules/operation/createGroupRule
func (api *GroupAPI) CreateGroupRule(
	ctx context.Context,
	groupID string,
	rule models.Rule,
	level models.RuleLevel,
) error {
	if err := validateInput(groupID, regexGroupIDArtifactID, "Group ID"); err != nil {
		return err
	}

	urlPath := fmt.Sprintf("%s/groups/%s/rules", api.Client.BaseURL, url.PathEscape(groupID))

	// Prepare the request body
	body := models.CreateUpdateRuleRequest{
		RuleType: rule,
		Config:   level,
	}
	resp, err := api.executeRequest(ctx, http.MethodPost, urlPath, body)
	if err != nil {
		return err
	}

	return handleResponse(resp, http.StatusNoContent, nil)
}

// DeleteAllGroupRule Deletes all the rules configured for the group.
// After this is done, the global rules apply to artifacts in the group again.
// See https://www.apicur.io/registry/docs/apicurio-registry/3.0.x/assets-attachments/registry-rest-api.htm#tag/Group-rules/operation/deleteGroupRules
func (api *GroupAPI) DeleteAllGroupRule(ctx context.Context, groupID string) error {
	if err := validateInput(groupID, regexGroupIDArtifactID, "Group ID"); err != nil {
		return err
	}

	urlPath := fmt.Sprintf("%s/groups/%s/rules", api.Client.BaseURL, url.PathEscape(groupID))
	resp, err := api.executeRequest(ctx, http.MethodDelete, urlPath, nil)
	if err != nil {
		return err
	}

	return handleResponse(resp, http.StatusNoContent, nil)
}

// GetGroupRule returns the configuration of a single rule for the group.
// Returns information about a single rule configured for a group.
// This is useful when you want to know what the current configuration settings are for a specific rule.
// See https://www.apicur.io/registry/docs/apicurio-registry/3.0.x/assets-attachments/registry-rest-api.htm#tag/Group-rules/operation/getGroupRuleConfig
func (api *GroupAPI) GetGroupRule(
	ctx context.Context,
	groupID string,
	rule models.Rule,
) (models.RuleLevel, error) {
	if err := validateInput(groupID, regexGroupIDArtifactID, "Group ID"); err != nil {
		return "", err
	}

	urlPath := fmt.Sprintf(
		"%s/groups/%s/rules/%s",
		api.Client.BaseURL,
		url.PathEscape(groupID),
		rule,
	)
	resp, err := api.executeRequest(ctx, http.MethodGet, urlPath, nil)
	if err != nil {
		return "", err
	}

	var globalRule models.RuleResponse
	if err := handleResponse(resp, http.StatusOK, &globalRule); err != nil {
		return "", err
	}

	return globalRule.Config, nil
}

// UpdateGroupRule Updates the configuration of a single rule for the group.
// The configuration data is specific to each rule type, so the configuration of the COMPATIBILITY rule is in a different format from the configuration of the VALIDITY rule.
// See https://www.apicur.io/registry/docs/apicurio-registry/3.0.x/assets-attachments/registry-rest-api.htm#tag/Group-rules/operation/updateGroupRuleConfig
func (api *GroupAPI) UpdateGroupRule(
	ctx context.Context,
	groupID string,
	rule models.Rule,
	level models.RuleLevel,
) error {
	if err := validateInput(groupID, regexGroupIDArtifactID, "Group ID"); err != nil {
		return err
	}

	urlPath := fmt.Sprintf(
		"%s/groups/%s/rules/%s",
		api.Client.BaseURL,
		url.PathEscape(groupID),
		rule,
	)

	// Prepare the request body
	body := models.CreateUpdateRuleRequest{
		RuleType: rule,
		Config:   level,
	}
	resp, err := api.executeRequest(ctx, http.MethodPut, urlPath, body)
	if err != nil {
		return err
	}

	var globalRule models.RuleResponse
	if err := handleResponse(resp, http.StatusOK, &globalRule); err != nil {
		return err
	}

	return nil
}

// DeleteGroupRule deletes the rule for a given group.
// See https://www.apicur.io/registry/docs/apicurio-registry/3.0.x/assets-attachments/registry-rest-api.htm#tag/Group-rules/operation/deleteGroupRule
func (api *GroupAPI) DeleteGroupRule(ctx context.Context, groupID string, rule models.Rule) error {
	if err := validateInput(groupID, regexGroupIDArtifactID, "Group ID"); err != nil {
		return err
	}

	urlPath := fmt.Sprintf(
		"%s/groups/%s/rules/%s",
		api.Client.BaseURL,
		url.PathEscape(groupID),
		rule,
	)
	resp, err := api.executeRequest(ctx, http.MethodDelete, urlPath, nil)
	if err != nil {
		return err
	}

	return handleResponse(resp, http.StatusNoContent, nil)
}

// executeRequest handles the creation and execution of an HTTP request.
func (api *GroupAPI) executeRequest(
	ctx context.Context,
	method, url string,
	body interface{},
) (*http.Response, error) {
	return executeRequest(ctx, api.Client, method, url, body)
}
