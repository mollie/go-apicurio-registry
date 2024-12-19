// Package apis provides multiple APIs to interact with the Apicurio Registry.
//
// The apis package exposes different high-level APIs to manage resources in the Apicurio Registry,
// such as artifacts, metadata, and administrative operations. Each API abstracts specific
// functionality of the Apicurio Registry, providing a clean and easy-to-use interface.
//
// Available APIs:
//   - AdminAPI: Provides administrative operations, such as clearing the cache and managing system-level settings.
//   - ArtifactsAPI: Manages artifact lifecycles, including creation, retrieval, versioning, and deletion of schemas.
//   - MetadataAPI: Handles schema metadata operations, such as tagging, labeling, and descriptions.
//   - RulesAPI: Allows configuration and management of rules for compatibility and validation.
//   - BranchesAPI: Allows configuration and management of branches for artifact versioning.
//   - GroupsAPI: Manages groups related resource in apicurio registry.
//   - VersionsAPI: Manages versions related resource in apicurio registry.
//   - SystemsAPI: Manages systems related resource in apicurio registry.
//   - ...
//
// Features:
//   - Register, retrieve, and delete artifacts.
//   - Perform schema versioning and compatibility checks.
//   - Manage metadata, tags, and labels for artifacts.
//   - Execute administrative actions on the Apicurio Registry.
//
// Example usage:
//
//	import "github.com/mollie/go-apicurio-registry/apis"
//
//	func main() {
//	    // Initialize the Artifacts API
//	    artifactsAPI := apis.NewArtifactsAPI("https://registry.example.com")
//
//	    // Register a new artifact
//	    err := artifactsAPI.RegisterArtifact("example-subject", `{"type":"record","name":"example"}`)
//	    if err != nil {
//	        log.Fatal("Error registering artifact:", err)
//	    }
//	    fmt.Println("Artifact registered successfully")
//
//	    // Retrieve an Artifact
//	    artifact, err := artifactsAPI.GetArtifact("example-subject", "latest")
//	    if err != nil {
//	        log.Fatal("Error fetching artifact:", err)
//	    }
//	    fmt.Println("Retrieved artifact:", artifact)
//	}
package apis
