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

type SystemAPI struct {
	Client *client.Client
}

// NewSystemAPI creates a new SystemAPI
func NewSystemAPI(client *client.Client) *SystemAPI {
	return &SystemAPI{
		Client: client,
	}
}

// GetSystemInfo gets the system info
// This operation retrieves information about the running registry system, such as the version of the software and when it was built.
// See https://www.apicur.io/registry/docs/apicurio-registry/3.0.x/assets-attachments/registry-rest-api.htm#tag/System/operation/getSystemInfo
func (api *SystemAPI) GetSystemInfo(ctx context.Context) (*models.SystemInfoResponse, error) {
	url := fmt.Sprintf("%s/system/info", api.Client.BaseURL)
	resp, err := api.executeRequest(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	var systemInfo models.SystemInfoResponse
	if err := handleResponse(resp, http.StatusOK, &systemInfo); err != nil {
		return nil, err
	}

	return &systemInfo, nil
}

// GetResourceLimitInfo gets the resource limit info
// This operation retrieves the list of limitations on used resources, that are applied on the current instance of Registry.
// See https://www.apicur.io/registry/docs/apicurio-registry/3.0.x/assets-attachments/registry-rest-api.htm#tag/System/operation/getResourceLimits
func (api *SystemAPI) GetResourceLimitInfo(ctx context.Context) (*models.SystemResourceLimitInfoResponse, error) {
	url := fmt.Sprintf("%s/system/limits", api.Client.BaseURL)
	resp, err := api.executeRequest(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	var resourceLimitInfo models.SystemResourceLimitInfoResponse
	if err := handleResponse(resp, http.StatusOK, &resourceLimitInfo); err != nil {
		return nil, err
	}

	return &resourceLimitInfo, nil

}

// GetUIConfig gets the UI config
// Returns the UI configuration properties for this server. The registry UI can be connected to a backend using just a URL. The rest of the UI configuration can then be fetched from the backend using this operation. This allows UI and backend to both be configured in the same place.
// See https://www.apicur.io/registry/docs/apicurio-registry/3.0.x/assets-attachments/registry-rest-api.htm#tag/System/operation/getUIConfig
func (api *SystemAPI) GetUIConfig(ctx context.Context) (*models.SystemUIConfigResponse, error) {
	url := fmt.Sprintf("%s/system/uiConfig", api.Client.BaseURL)
	resp, err := api.executeRequest(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	var uiConfig models.SystemUIConfigResponse
	if err := handleResponse(resp, http.StatusOK, &uiConfig); err != nil {
		return nil, err
	}

	return &uiConfig, nil

}

// GetCurrentUser Returns information about the currently authenticated user.
// See https://www.apicur.io/registry/docs/apicurio-registry/3.0.x/assets-attachments/registry-rest-api.htm#tag/Users
func (api *SystemAPI) GetCurrentUser(ctx context.Context) (*models.UserInfo, error) {
	url := fmt.Sprintf("%s/users/me", api.Client.BaseURL)
	resp, err := api.executeRequest(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	var userInfo models.UserInfo
	if err := handleResponse(resp, http.StatusOK, &userInfo); err != nil {
		return nil, err
	}

	return &userInfo, nil
}

// executeRequest handles the creation and execution of an HTTP request.
func (api *SystemAPI) executeRequest(ctx context.Context, method, url string, body interface{}) (*http.Response, error) {
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
