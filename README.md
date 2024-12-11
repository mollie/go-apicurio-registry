# Go-APICURIO-SDK

[![Go Reference](https://pkg.go.dev/badge/github.com/mollie/go-apicurio-registry.svg)](https://pkg.go.dev/github.com/mollie/go-apicurio-registry)
[![Integration Tests](https://github.com/mollie/go-apicurio-registry/actions/workflows/tests.yml/badge.svg)](https://github.com/mollie/go-apicurio-registry/actions)

**Go-APICURIO-SDK** is an open-source Go SDK for interacting with [Apicurio Schema Registry v3.0 APIs](https://www.apicur.io/). This library provides an idiomatic Go interface to manage schemas, validate data, and handle schema evolution in Go applications. It is designed to make it easy for developers to integrate the Apicurio Schema Registry into their Go-based microservices or event-driven systems.

## 🚧 Work In Progress

This SDK is currently under active development. Key features and improvements are being added incrementally. While contributions and feedback are welcome, please note that some APIs and features may change as the project evolves.

## Features

- **Schema Management**: Create, update, delete, and retrieve schemas.
- **Validation**: Validate payloads against registered schemas.
- **Schema Evolution**: Tools for managing schema compatibility and evolution.
- **Integration**: Works seamlessly with the Apicurio Schema Registry.

## Getting Started

### Prerequisites

- Go 1.18+ installed on your system.
- Access to an Apicurio Schema Registry instance.

### Installation

To use the SDK, add it to your project using `go get`:

```bash
go get github.com/mollie/go-apicurio-registry
```

### Development
Running Locally with Docker
This project includes a docker-compose.yml file for setting up a local Apicurio Schema Registry instance. Run the following command to start the registry:

```bash
docker-compose up
```

### Running Tests

The repository includes unit and integration tests. Use the following command to run all tests:

```bash
go test ./...
```

## Contribution Guidelines
Contributions are welcome! Please see the CONTRIBUTING.md (to be added) for details on the process.


## License
This project is licensed under the Apache License 2.0.
