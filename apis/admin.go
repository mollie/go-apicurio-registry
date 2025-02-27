package apis

import (
	"context"
	"fmt"
	"net/http"

	"github.com/mollie/go-apicurio-registry/client"
	"github.com/mollie/go-apicurio-registry/models"
)

type AdminAPI struct {
	Client *client.Client
}

func NewAdminAPI(client *client.Client) *AdminAPI {
	return &AdminAPI{
		Client: client,
	}
}

// ListGlobalRules Gets a list of all the currently configured global rules (if any).
// GET /admin/rules
// See: https://www.apicur.io/registry/docs/apicurio-registry/3.0.x/assets-attachments/registry-rest-api.htm#tag/Global-rules/operation/listGlobalRules
func (api *AdminAPI) ListGlobalRules(ctx context.Context) ([]models.Rule, error) {
	url := fmt.Sprintf("%s/admin/rules", api.Client.BaseURL)
	resp, err := api.executeRequest(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	var rules []models.Rule
	if err := handleResponse(resp, http.StatusOK, &rules); err != nil {
		return nil, err
	}

	return rules, nil
}

// CreateGlobalRule Creates a new global rule.
// POST /admin/rules
// See: https://www.apicur.io/registry/docs/apicurio-registry/3.0.x/assets-attachments/registry-rest-api.htm#tag/Global-rules/operation/createGlobalRule
func (api *AdminAPI) CreateGlobalRule(
	ctx context.Context,
	rule models.Rule,
	level models.RuleLevel,
) error {
	url := fmt.Sprintf("%s/admin/rules", api.Client.BaseURL)

	// Prepare the request body
	body := models.CreateUpdateRuleRequest{
		RuleType: rule,
		Config:   level,
	}
	resp, err := api.executeRequest(ctx, http.MethodPost, url, body)
	if err != nil {
		return err
	}

	return handleResponse(resp, http.StatusNoContent, nil)
}

// DeleteAllGlobalRule Adds a rule to the list of globally configured rules.
// DELETE /admin/rules
// See: https://www.apicur.io/registry/docs/apicurio-registry/3.0.x/assets-attachments/registry-rest-api.htm#tag/Global-rules/operation/deleteAllGlobalRules
func (api *AdminAPI) DeleteAllGlobalRule(ctx context.Context) error {
	url := fmt.Sprintf("%s/admin/rules", api.Client.BaseURL)
	resp, err := api.executeRequest(ctx, http.MethodDelete, url, nil)
	if err != nil {
		return err
	}

	return handleResponse(resp, http.StatusNoContent, nil)
}

// GetGlobalRule Returns information about the named globally configured rule.
// GET /admin/rules/{rule}
// See: https://www.apicur.io/registry/docs/apicurio-registry/3.0.x/assets-attachments/registry-rest-api.htm#tag/Global-rules/operation/getGlobalRuleConfig
func (api *AdminAPI) GetGlobalRule(
	ctx context.Context,
	rule models.Rule,
) (models.RuleLevel, error) {
	url := fmt.Sprintf("%s/admin/rules/%s", api.Client.BaseURL, rule)
	resp, err := api.executeRequest(ctx, http.MethodGet, url, nil)
	if err != nil {
		return "", err
	}

	var globalRule models.RuleResponse
	if err := handleResponse(resp, http.StatusOK, &globalRule); err != nil {
		return "", err
	}

	return globalRule.Config, nil
}

// UpdateGlobalRule Updates the configuration of the named globally configured rule.
// PUT /admin/rules/{rule}
// See https://www.apicur.io/registry/docs/apicurio-registry/3.0.x/assets-attachments/registry-rest-api.htm#tag/Global-rules/operation/updateGlobalRuleConfig
func (api *AdminAPI) UpdateGlobalRule(
	ctx context.Context,
	rule models.Rule,
	level models.RuleLevel,
) error {
	url := fmt.Sprintf("%s/admin/rules/%s", api.Client.BaseURL, rule)

	// Prepare the request body
	body := models.CreateUpdateRuleRequest{
		RuleType: rule,
		Config:   level,
	}
	resp, err := api.executeRequest(ctx, http.MethodPut, url, body)
	if err != nil {
		return err
	}

	var globalRule models.RuleResponse
	if err := handleResponse(resp, http.StatusOK, &globalRule); err != nil {
		return err
	}

	return nil
}

// DeleteGlobalRule Deletes the named globally configured rule.
// DELETE /admin/rules/{rule}
// See https://www.apicur.io/registry/docs/apicurio-registry/3.0.x/assets-attachments/registry-rest-api.htm#tag/Global-rules/operation/deleteGlobalRule
func (api *AdminAPI) DeleteGlobalRule(ctx context.Context, rule models.Rule) error {
	url := fmt.Sprintf("%s/admin/rules/%s", api.Client.BaseURL, rule)
	resp, err := api.executeRequest(ctx, http.MethodDelete, url, nil)
	if err != nil {
		return err
	}

	return handleResponse(resp, http.StatusNoContent, nil)
}

// ListArtifactTypes Gets a list of all the currently configured artifact types (if any).
// GET admin/config/artifactTypes
// See https://www.apicur.io/registry/docs/apicurio-registry/3.0.x/assets-attachments/registry-rest-api.htm#tag/Artifact-Type/operation/listArtifactTypes
func (api *AdminAPI) ListArtifactTypes(ctx context.Context) ([]models.ArtifactType, error) {
	url := fmt.Sprintf("%s/admin/config/artifactTypes", api.Client.BaseURL)
	resp, err := api.executeRequest(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	var artifactTypesResponse []models.ArtifactTypeResponse
	if err := handleResponse(resp, http.StatusOK, &artifactTypesResponse); err != nil {
		return nil, err
	}

	var artifactTypes []models.ArtifactType
	for _, item := range artifactTypesResponse {
		artifactTypes = append(artifactTypes, item.Name)
	}

	return artifactTypes, nil

}

// executeRequest handles the creation and execution of an HTTP request.
func (api *AdminAPI) executeRequest(
	ctx context.Context,
	method, url string,
	body interface{},
) (*http.Response, error) {
	return executeRequest(ctx, api.Client, method, url, body)
}
