// Package client provides the core client implementation for interacting with the Apicurio Registry.
//
// The client package is the entry point for initializing and configuring the Apicurio Registry client.
// It abstracts the complexities of making HTTP requests to the registry and provides a clean interface
// for developers to use in their applications.
//
// Features:
//
// - Centralized configuration for base URL and authentication.
// - Automatic retry mechanisms for transient failures.
// - Thread-safe implementation suitable for concurrent use.
// - Extensible architecture to support additional features and middleware.
//
// Usage:
//
// To use the client, create a new instance using the `NewApicurioClient` function and pass a `Config` structure
// with the required settings such as BaseURL and AuthToken.
//
// Example:
//
//	package main
//
//	import (
//		"fmt"
//		"github.com/mollie/go-apicurio-registry/client"
//	)
//
//	func main() {
//		// Define the client configuration
//		config := client.Config{
//			BaseURL: "https://my-registry.example.com",
//			AuthToken: "my-token",
//		}
//
//		// Initialize the client
//		apiClient := client.NewApicurioClient(config)
//
//		// Check the connection to the registry
//		status, err := apiClient.CheckConnection()
//		if err != nil {
//			fmt.Printf("Error connecting to the registry: %v\n", err)
//			return
//		}
//		fmt.Printf("Registry status: %v\n", status)
//	}
//
// Configuration:
//
// The `Config` struct contains the following fields:
// - BaseURL: The URL of the Apicurio Registry.
// - AuthToken: A token used for authenticating requests.
// - HTTPClient: (Optional) A custom HTTP client for advanced use cases.
//
// Methods:
//
// The client provides several methods to interact with the registry, including:
// - `CheckConnection`: Verifies if the registry is reachable.
// - `DoRequest`: Executes raw HTTP requests for advanced use cases.
//
// Thread Safety:
//
// The client is designed to be thread-safe and can be used in concurrent environments without additional synchronization.
package client
