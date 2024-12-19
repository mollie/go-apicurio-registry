// Package apis provides multiple APIs to interact with the Apicurio Registry.
//
// The apis package exposes high-level APIs to manage resources in the Apicurio Registry,
// such as artifacts, versions, metadata, and administrative operations. Each API
// encapsulates specific functionalities of the Apicurio Registry, ensuring a modular and
// user-friendly interface.
//
// Available APIs:
//
//  1. **ArtifactsAPI**: Manages artifacts (schemas), including creation, retrieval,
//     updating, and deletion.
//
//  2. **VersionsAPI**: Handles operations related to artifact versions, such as retrieving
//     and managing version history.
//
//  3. **GroupsAPI**: Supports grouping of artifacts to enable multi-tenancy and
//     organizational categorization.
//
//  4. **MetadataAPI**: Provides methods to manage artifact metadata and perform
//     compatibility checks.
//
//  5. **AdminAPI**: Enables administrative operations like clearing caches, monitoring
//     system health, and managing global settings.
//
// 6. **SystemAPI**: Offers methods for querying registry status and configuration.
//
// Example Usage:
//
// The following example demonstrates how to use the ArtifactsAPI to create and retrieve an artifact:
//
//	package main
//
//	import (
//		"fmt"
//		"github.com/mollie/go-apicurio-registry/apis"
//		"github.com/mollie/go-apicurio-registry/client"
//		"github.com/mollie/go-apicurio-registry/models"
//	)
//
//	func main() {
//		// Initialize the API client
//		config := client.Config{
//			BaseURL: "https://my-registry.example.com",
//			AuthToken: "my-token",
//		}
//		apiClient := client.NewApicurioClient(config)
//
//		// Access the ArtifactsAPI
//		artifactsAPI := apis.NewArtifactsAPI(apiClient)
//
//		// Create a new artifact
//		artifact := models.Artifact{
//			GroupID: "example-group",
//			ID: "example-artifact",
//			Content: []byte(`{"type": "record", "name": "Example", "fields": [{"name": "field1", "type": "string"}]}`),
//			ContentType: "application/json",
//		}
//		response, err := artifactsAPI.Create(artifact)
//		if err != nil {
//			fmt.Printf("Error creating artifact: %v\n", err)
//			return
//		}
//		fmt.Printf("Artifact created with ID: %s\n", response.ID)
//
//		// Retrieve artifact metadata
//		metadata, err := artifactsAPI.GetMetadata("example-group", "example-artifact")
//		if err != nil {
//			fmt.Printf("Error retrieving metadata: %v\n", err)
//			return
//		}
//		fmt.Printf("Artifact metadata: %+v\n", metadata)
//	}
//
// Note:
//
// Each API can be accessed via its respective constructor function (e.g., `NewArtifactsAPI`,
// `NewAdminAPI`). These APIs are designed to integrate seamlessly with the client package.
package apis
