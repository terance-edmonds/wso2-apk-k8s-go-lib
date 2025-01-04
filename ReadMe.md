# Kubernetes Resource Generator Library

This Go library provides a framework for generating Kubernetes Gateway API-specific resources, including HTTPRoute and gRPC resources. The library offers default implementations for generating various components of these resources while allowing developers to override specific methods to suit their use cases.

## Features

- Generate Kubernetes resources for HTTP and gRPC configurations.
- Support for generating HTTPRoute and gRPC-specific Custom Resources (CRs).
- Flexible method overriding for custom implementations.
- Default implementations for common resource generation tasks.

## Installation

To use the library, include it in your Go project by adding the following import:

```go
import "path/to/gw_artifacts"
```

Ensure that the library and its dependencies are properly vendored in your project.

## Usage

### Initializing the Generator

#### HTTPRoute Generator

Create an instance of the HTTPRoute generator:

```go
import http_generator "gw_artifacts/pkg/generators/http"

gen := http_generator.Generator()
```

#### gRPC Generator

Create an instance of the gRPC generator:

```go
import grpc_generator "gw_artifacts/pkg/generators/grpc"

gen := grpc_generator.Generator()
```

These initialize the respective generators with default implementations for all functions.

### Generating HTTPRoute Resources

Use the HTTPRoute generator to create an HTTPRoute by calling the desired methods:

```go
httpRoute, err := gen.GenerateHTTPRoute(*apkConf, organization, gatewayConfig, *apkConf.Operations, &endpoint, constants.PRODUCTION_TYPE, "unique-route-id", 1)
if err != nil {
    log.Fatalf("Failed to generate HTTP route: %v", err)
}
```

### Generating gRPC Resources

Similarly, use the gRPC generator to create gRPC-specific resources:

```go
grpcResource, err := gen.GenerateGRPCRoute(*apkConf, organization, gatewayConfig, *apkConf.Operations, &endpoint, constants.PRODUCTION_TYPE, "unique-grpc-id")
if err != nil {
    log.Fatalf("Failed to generate gRPC resource: %v", err)
}
```

### Overriding Default Implementations

To customize the behavior of the generator, you can override specific methods:

#### HTTPRoute Example

```go
gen := http_generator.Generator()
gen.GenerateHTTPRouteRules = myCustomHttpRouteRuleImplementation
```

#### gRPC Example

```go
gen := grpc_generator.Generator()
gen.GenerateGRPCRouteRules = myCustomGrpcRouteRuleImplementation
```

This allows you to replace the default implementations with your own.

### Example

Examples of using the library are available in the following files:

- `examples/http/main.go`: Demonstrates HTTPRoute generation.
- `examples/grpc/main.go`: Demonstrates gRPC resource generation.

## API Reference

### HTTPRoute Generator Functions

```go
GenerateHTTPRouteRules(apkConf types.APKConf, operations []types.Operation, endpoint *types.EndpointDetails, endpointType string) ([]gwapiv1.HTTPRouteRule, error)
GenerateHTTPRouteRule(apkConf types.APKConf, operation types.Operation, endpoint *types.EndpointDetails, endpointType string) (*gwapiv1.HTTPRouteRule, error)
GenerateAndRetrieveParentRefs(gatewayConfig types.GatewayConfigurations, uniqueId string) []gwapiv1.ParentReference
GenerateHTTPRouteFilters(apkConf types.APKConf, endpointToUse types.EndpointDetails, operation types.Operation, endpointType string) ([]gwapiv1.HTTPRouteFilter, bool)
ExtractHTTPRouteFilter(apkConf *types.APKConf, endpoint types.EndpointDetails, operation types.Operation, operationPolicies []types.OperationPolicy, isRequest bool) ([]gwapiv1.HTTPRouteFilter, bool)
GetHostNames(apkConf types.APKConf, endpointType string, organization types.Organization) []gwapiv1.Hostname
RetrieveHTTPMatches(apkConf types.APKConf, operation types.Operation) ([]gwapiv1.HTTPRouteMatch, error)
RetrieveHTTPMatch(apkConf types.APKConf, operation types.Operation) (gwapiv1.HTTPRouteMatch, error)
GenerateHTTPBackEndRef(endpoint types.EndpointDetails, operation types.Operation, endpointType string) []gwapiv1.HTTPBackendRef
```

### gRPC Generator Functions

```go
GenerateGRPCRouteRules(apkConf types.APKConf, operations []types.Operation, endpoint *types.EndpointDetails, endpointType string) ([]gwapiv1.GRPCRouteRule, error)
GenerateGRPCRouteRule(apkConf types.APKConf, operation types.Operation, endpoint *types.EndpointDetails, endpointType string) (*gwapiv1.GRPCRouteRule, error)
GenerateAndRetrieveParentRefs(gatewayConfig types.GatewayConfigurations, uniqueId string) []gwapiv1.ParentReference
GetHostNames(apkConf types.APKConf, endpointType string, organization types.Organization) []gwapiv1.Hostname
RetrieveGRPCMatches(operation types.Operation) []gwapiv1.GRPCRouteMatch
RetrieveGRPCMatch(operation types.Operation) gwapiv1.GRPCRouteMatch
GenerateGRPCBackEndRef(endpoint types.EndpointDetails, operation types.Operation) []gwapiv1.GRPCBackendRef
```

### Function: `Generator`

Creates and initializes a new generator instance with default implementations for HTTP or gRPC resources.

## Directory Structure

- `pkg/generators/http`: Contains HTTPRoute-specific generator logic.
- `pkg/generators/grpc`: Contains gRPC-specific generator logic.
