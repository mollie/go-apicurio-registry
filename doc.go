// Package apicurio provides a client library for interacting with the Apicurio Registry.
// It allows developers to manage schemas, validate payloads, and perform schema evolution.
//
// Features:
// - CRUD operations on schemas.
// - Schema validation and compatibility checks.
// - Support for multiple serialization formats like Avro, Protobuf, and JSON Schema.
//
// Example usage:
//
//	import "github.com/mollie/go-apicurio-registry/client"
//
//	func main() {
//	    client := client.NewApicurioClient("https://registry.example.com")
//	    schema, err := client.GetSchema("example-subject", "latest")
//	    if err != nil {
//	        log.Fatal(err)
//	    }
//	    fmt.Println("Schema:", schema)
//	}
package apicurio
