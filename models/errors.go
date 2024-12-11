package models

import "fmt"

var (
	ErrUnknownArtifactType = fmt.Errorf("unknown artifact type")
)

// APIError represents the structure of an error response from the API.
type APIError struct {
	Detail   string `json:"detail"`   // A human-readable explanation specific to the problem
	Type     string `json:"type"`     // A URI reference identifying the problem type
	Title    string `json:"title"`    // A short, human-readable summary of the problem type
	Status   int    `json:"status"`   // The HTTP status code
	Instance string `json:"instance"` // A URI reference identifying the specific occurrence
	Name     string `json:"name"`     // The name of the error (e.g., server exception class name)
}

// Error satisfies the error interface and formats the APIError as a string.
func (e *APIError) Error() string {
	return fmt.Sprintf("[%d] %s: %s (detail: %s, instance: %s, type: %s)",
		e.Status, e.Title, e.Name, e.Detail, e.Instance, e.Type)
}
