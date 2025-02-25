package apis

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/mollie/go-apicurio-registry/client"
	"github.com/mollie/go-apicurio-registry/models"
	"github.com/pkg/errors"
	"io"
	"net/http"
	"regexp"
)

const (
	ContentTypeJSON = "application/json"
	ContentTypeAll  = "*/*"
)

var (
	regexGroupIDArtifactID = regexp.MustCompile(`^.{1,512}$`)
	regexVersion           = regexp.MustCompile(`[a-zA-Z0-9._\-+]{1,256}`)
	regexBranchID          = regexp.MustCompile(`[a-zA-Z0-9._\-+]{1,256}`)

	ErrInvalidInput = errors.New("input did not pass validation with regex")
)

// ErrInvalidInput is returned when an input validation fails.
func validateInput(input string, regex *regexp.Regexp, name string) error {
	if match := regex.MatchString(input); !match {
		return errors.Wrapf(ErrInvalidInput, "%s='%s', regex=%s", name, input, regex.String())
	}
	return nil
}

// parseAPIError parses an API error response and returns an APIError struct.
func parseAPIError(resp *http.Response) (*models.APIError, error) {
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read error response body: %w", err)
	}

	var apiError models.APIError
	if err := json.Unmarshal(body, &apiError); err != nil {
		return nil, fmt.Errorf("failed to parse error response: %w", err)
	}

	return &apiError, nil
}

func parseArtifactTypeHeader(resp *http.Response) (models.ArtifactType, error) {
	artifactTypeHeader := resp.Header.Get("X-Registry-ArtifactType")
	artifactType, err := models.ParseArtifactType(artifactTypeHeader)
	if err != nil {
		return "", errors.Wrapf(err, "invalid artifact type in response header: %s", artifactTypeHeader)
	}
	return artifactType, nil
}

// handleResponse reads the response body and checks the status code.
func handleResponse(resp *http.Response, expectedStatus int, result interface{}) error {
	defer resp.Body.Close()

	if resp.StatusCode != expectedStatus {
		apiError, parseErr := parseAPIError(resp)
		if parseErr != nil {
			return errors.Wrapf(parseErr, "unexpected server error: %d", resp.StatusCode)
		}
		return apiError
	}

	if result != nil && resp.StatusCode == expectedStatus {
		if err := json.NewDecoder(resp.Body).Decode(result); err != nil {
			return errors.Wrap(err, "failed to parse response body")
		}
	}

	return nil
}

// handleRawResponse reads the response body and checks the status code.
func handleRawResponse(resp *http.Response, expectedStatus int) (string, error) {
	defer resp.Body.Close()
	if resp.StatusCode != expectedStatus {
		apiError, parseErr := parseAPIError(resp)
		if parseErr != nil {
			return "", errors.Wrap(parseErr, "unexpected server error")
		}
		return "", apiError
	}

	content, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", errors.Wrap(err, "failed to read response body")
	}

	return string(content), nil
}

// executeRequest handles the creation and execution of an HTTP request.
func executeRequest(ctx context.Context, client *client.Client, method, url string, body interface{}) (*http.Response, error) {
	var reqBody io.Reader
	contentType := ""

	if body != nil {
		switch v := body.(type) {
		case string:
			reqBody = bytes.NewReader([]byte(v))
			contentType = "*/*"
		case []byte:
			reqBody = bytes.NewReader(v)
			contentType = "*/*"
		default:
			contentType = "application/json"
			jsonData, err := json.Marshal(body)
			if err != nil {
				return nil, errors.Wrap(err, "failed to marshal request body as JSON")
			}
			reqBody = bytes.NewReader(jsonData)
		}
	} else {
		reqBody = nil // Send request without body
	}

	// Create the HTTP request
	req, err := http.NewRequestWithContext(ctx, method, url, reqBody)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create HTTP request")
	}

	// Set Content-Type header only if there is a body
	if contentType != "" {
		req.Header.Set("Content-Type", contentType)
	}

	// Execute the request
	resp, err := client.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "failed to execute HTTP request")
	}

	return resp, nil
}
