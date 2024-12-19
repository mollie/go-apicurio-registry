// Package apicurio provides a Go client library for interacting with the Apicurio Registry.
//
// The library enables developers to seamlessly integrate with the Apicurio Registry
// for managing, evolving, and validating schemas in a variety of serialization formats.
//
// Features:
//
// - CRUD operations on schemas (artifacts) including creation, retrieval, update, and deletion.
// - Management of schema versions, branches, and metadata.
// - Group-based organization for schemas to support multi-tenancy.
// - Schema validation and compatibility checks for supported formats such as Avro, Protobuf, and JSON Schema.
// - System-level operations such as retrieving registry status and configuration.
//
// Structure:
//
// The library is structured into the following key components:
//
//  1. **Client**: Provides an entry point for interacting with the registry.
//     Use the `client.NewApicurioClient` function to create a new client instance.
//
//  2. **APIs**: Contains modular functions for specific operations such as managing artifacts,
//     branches, versions, groups, and performing administrative tasks.
//
// 3. **Models**: Defines data structures for requests, responses, and errors used across the library.
//
// Example Usage:
//
// The following example demonstrates how to create a new artifact and retrieve its metadata:
//
//	package main
//
//	import (
//		"fmt"
//		"github.com/mollie/go-apicurio-registry/client"
//		"github.com/mollie/go-apicurio-registry/models"
//	)
//
//	func main() {
//		// Initialize the client
//		config := client.Config{
//			BaseURL: "https://my-registry.example.com",
//			AuthToken: "my-token",
//		}
//		apiClient := client.NewApicurioClient(config)
//
//		// Create a new artifact
//		artifact := models.Artifact{
//			GroupID: "example-group",
//			ID: "example-artifact",
//			Content: []byte(`{"type": "record", "name": "Example", "fields": [{"name": "field1", "type": "string"}]}`),
//			ContentType: "application/json",
//		}
//		response, err := apiClient.Artifacts.Create(artifact)
//		if err != nil {
//			fmt.Printf("Error creating artifact: %v\n", err)
//			return
//		}
//		fmt.Printf("Artifact created with ID: %s\n", response.ID)
//
//		// Retrieve artifact metadata
//		metadata, err := apiClient.Artifacts.GetMetadata("example-group", "example-artifact")
//		if err != nil {
//			fmt.Printf("Error retrieving metadata: %v\n", err)
//			return
//		}
//		fmt.Printf("Artifact metadata: %+v\n", metadata)
//	}
package apicurio
