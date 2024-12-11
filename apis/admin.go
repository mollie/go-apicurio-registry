package apis

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/mollie/go-apicurio-registry/client"
	"github.com/mollie/go-apicurio-registry/models"
	"github.com/pkg/errors"
	"net/http"
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
func (api *AdminAPI) CreateGlobalRule(ctx context.Context, rule models.Rule, level models.RuleLevel) error {
	url := fmt.Sprintf("%s/admin/rules", api.Client.BaseURL)

	// Prepare the request body
	body := models.CreateUpdateGlobalRuleRequest{
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
func (api *AdminAPI) GetGlobalRule(ctx context.Context, rule models.Rule) (models.RuleLevel, error) {
	url := fmt.Sprintf("%s/admin/rules/%s", api.Client.BaseURL, rule)
	resp, err := api.executeRequest(ctx, http.MethodGet, url, nil)
	if err != nil {
		return "", err
	}

	var globalRule models.GlobalRuleResponse
	if err := handleResponse(resp, http.StatusOK, &globalRule); err != nil {
		return "", err
	}

	return globalRule.Config, nil
}

// UpdateGlobalRule Updates the configuration of the named globally configured rule.
// PUT /admin/rules/{rule}
// See https://www.apicur.io/registry/docs/apicurio-registry/3.0.x/assets-attachments/registry-rest-api.htm#tag/Global-rules/operation/updateGlobalRuleConfig
func (api *AdminAPI) UpdateGlobalRule(ctx context.Context, rule models.Rule, level models.RuleLevel) error {
	url := fmt.Sprintf("%s/admin/rules/%s", api.Client.BaseURL, rule)

	// Prepare the request body
	body := models.CreateUpdateGlobalRuleRequest{
		RuleType: rule,
		Config:   level,
	}
	resp, err := api.executeRequest(ctx, http.MethodPut, url, body)
	if err != nil {
		return err
	}

	var globalRule models.GlobalRuleResponse
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

// executeRequest handles the creation and execution of an HTTP request.
func (api *AdminAPI) executeRequest(ctx context.Context, method, url string, body interface{}) (*http.Response, error) {
	var reqBody []byte
	var err error
	contentType := "*/*"

	switch v := body.(type) {
	case string:
		reqBody = []byte(v)
		contentType = "*/*"
	case []byte:
		reqBody = v
		contentType = "*/*"
	default:
		contentType = "application/json"
		reqBody, err = json.Marshal(body)
		if err != nil {
			return nil, errors.Wrap(err, "failed to marshal request body as JSON")
		}
	}

	// Create the HTTP request
	req, err := http.NewRequestWithContext(ctx, method, url, bytes.NewReader(reqBody))
	if err != nil {
		return nil, errors.Wrap(err, "failed to create HTTP request")
	}

	// Set appropriate Content-Type header
	if body != nil {
		req.Header.Set("Content-Type", contentType)
	}

	// Execute the request
	resp, err := api.Client.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "failed to execute HTTP request")
	}

	return resp, nil
}
