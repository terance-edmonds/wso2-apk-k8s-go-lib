/*
 *  Copyright (c) 2024, WSO2 LLC. (http://www.wso2.org) All Rights Reserved.
 *
 *  Licensed under the Apache License, Version 2.0 (the "License");
 *  you may not use this file except in compliance with the License.
 *  You may obtain a copy of the License at
 *
 *  http://www.apache.org/licenses/LICENSE-2.0
 *
 *  Unless required by applicable law or agreed to in writing, software
 *  distributed under the License is distributed on an "AS IS" BASIS,
 *  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *  See the License for the specific language governing permissions and
 *  limitations under the License.
 *
 */

package types

// SecretInfo holds the info related to the created secret upon enabling the endpoint security options like basic auth
type SecretInfo struct {
	SecretName     string `yaml:"secretName,omitempty"`
	UsernameKey    string `yaml:"userNameKey,omitempty"`
	PasswordKey    string `yaml:"passwordKey,omitempty"`
	In             string `yaml:"in,omitempty"`
	APIKeyNameKey  string `yaml:"apiKeyNameKey,omitempty"`
	APIKeyValueKey string `yaml:"apiKeyValueKey,omitempty"`
}

// EndpointSecurity comtains the information related to endpoint security configurations enabled by a user for a given API
type EndpointSecurity struct {
	Enabled      bool       `yaml:"enabled,omitempty"`
	SecurityType SecretInfo `yaml:"securityType,omitempty"`
}

// EndpointCertificate struct stores the the alias and the name for a particular endpoint security configuration
type EndpointCertificate struct {
	Name string `yaml:"secretName"`
	Key  string `yaml:"secretKey"`
}

// EndpointDetails represents the details of an endpoint, containing its URL.
type EndpointDetails struct {
	Name         string `json:"name"`
	URL          string `json:"url"`
	Namespace    string `json:"namespace"`
	ServiceEntry bool   `json:"serviceEntry"`
}

// Endpoint struct stores the endpoint configuration for a particular API
type Endpoint interface {
	isEndpoint()
}

// EndpointURL represents a simple string endpoint
type EndpointURL string

func (u EndpointURL) isEndpoint() {}

// K8sService struct stores the name and key for a particular Kubernetes service configuration
type K8sService struct {
	Name      string `yaml:"name"`
	Namespace string `yaml:"namespace"`
	Port      string `yaml:"port"`
	Protocol  string `yaml:"protocol"`
}

func (u K8sService) isEndpoint() {}

// EndpointConfiguration stores the data related to endpoints and their related
type EndpointConfiguration struct {
	Endpoint       Endpoint            `yaml:"endpoint,omitempty"`
	EndCertificate EndpointCertificate `yaml:"certificate,omitempty"`
	EndSecurity    EndpointSecurity    `yaml:"endpointSecurity,omitempty"`
	AIRatelimit    AIRatelimit         `yaml:"aiRatelimit,omitempty"`
}

// AIRatelimit defines the configuration for AI rate limiting,
// including whether rate limiting is enabled and the settings
// for token and request-based limits.
type AIRatelimit struct {
	Enabled bool        `yaml:"enabled"`
	Token   TokenAIRL   `yaml:"token"`
	Request RequestAIRL `yaml:"request"`
}

// TokenAIRL defines the configuration for Token AI rate limit settings.
type TokenAIRL struct {
	PromptLimit     int    `yaml:"promptLimit"`
	CompletionLimit int    `yaml:"completionLimit"`
	TotalLimit      int    `yaml:"totalLimit"`
	Unit            string `yaml:"unit"` // Time unit (Minute, Hour, Day)
}

// RequestAIRL defines the configuration for Request AI rate limit settings.
type RequestAIRL struct {
	RequestLimit int    `yaml:"requestLimit"`
	Unit         string `yaml:"unit"` // Time unit (Minute, Hour, Day)
}

// AdditionalProperty stores the custom properties set by the user for a particular API
type AdditionalProperty struct {
	Name  string `yaml:"name"`
	Value string `yaml:"value"`
}

// Certificate struct stores the the alias and the name for a particular mTLS configuration
type Certificate struct {
	Name string `json:"name"`
	Key  string `json:"key"`
}

// AuthConfiguration represents the security configurations made for the API security
type AuthConfiguration struct {
	Required          string        `yaml:"required,omitempty"`
	AuthType          string        `yaml:"authType,omitempty"`
	HeaderName        string        `yaml:"headerName,omitempty"`
	SendTokenUpStream bool          `yaml:"sendTokenToUpstream,omitempty"`
	Enabled           bool          `yaml:"enabled"`
	QueryParamName    string        `yaml:"queryParamName,omitempty"`
	HeaderEnabled     bool          `yaml:"headerEnable,omitempty"`
	QueryParamEnable  bool          `yaml:"queryParamEnable,omitempty"`
	Certificates      []Certificate `yaml:"certificates,omitempty"`
	Audience          []string      `yaml:"audience,omitempty"`
}

// EndpointConfigurations holds production and sandbox endpoints.
type EndpointConfigurations struct {
	Production *EndpointConfiguration `yaml:"production,omitempty"`
	Sandbox    *EndpointConfiguration `yaml:"sandbox,omitempty"`
}

// OperationPolicy defines policies, including interceptor parameters, for API operations.
type OperationPolicy struct {
	PolicyName    string    `yaml:"policyName,omitempty"`
	PolicyVersion string    `yaml:"policyVersion,omitempty"`
	PolicyID      string    `yaml:"policyId,omitempty"`
	Parameters    Parameter `yaml:"parameters,omitempty"`
}

// Parameter interface is used to define the type of parameters that can be used in an operation policy.
type Parameter interface {
	isParameter()
}

// RedirectPolicy contains the information for redirect request policies
type RedirectPolicy struct {
	URL        string `json:"url,omitempty" yaml:"url,omitempty"`
	StatusCode int    `json:"statusCode,omitempty" yaml:"statusCode,omitempty"`
}

func (u RedirectPolicy) isParameter() {}

// URLList contains the urls for mirror policies
type URLList struct {
	URLs []string `json:"urls,omitempty" yaml:"urls,omitempty"`
}

func (u URLList) isParameter() {}

// Header contains the information for header modification
type Header struct {
	HeaderName  string `yaml:"headerName"`
	HeaderValue string `yaml:"headerValue,omitempty"`
}

func (h Header) isParameter() {}

// InterceptorService holds configuration details for configuring interceptor
// for particular API requests or responses.
type InterceptorService struct {
	BackendURL      string `yaml:"backendUrl,omitempty"`
	HeadersEnabled  bool   `yaml:"headersEnabled,omitempty"`
	BodyEnabled     bool   `yaml:"bodyEnabled,omitempty"`
	TrailersEnabled bool   `yaml:"trailersEnabled,omitempty"`
	ContextEnabled  bool   `yaml:"contextEnabled,omitempty"`
	TLSSecretName   string `yaml:"tlsSecretName,omitempty"`
	TLSSecretKey    string `yaml:"tlsSecretKey,omitempty"`
}

func (s InterceptorService) isParameter() {}

// BackendJWT holds configuration details for configuring JWT for backend
type BackendJWT struct {
	Encoding         string `yaml:"encoding,omitempty"`
	Header           string `yaml:"header,omitempty"`
	SigningAlgorithm string `yaml:"signingAlgorithm,omitempty"`
	TokenTTL         int    `yaml:"tokenTTL,omitempty"`
}

func (j BackendJWT) isParameter() {}

// OperationPolicies organizes request and response policies for an API operation.
type OperationPolicies struct {
	Request  []OperationPolicy `yaml:"request,omitempty"`
	Response []OperationPolicy `yaml:"response,omitempty"`
}

// Operation represents an API operation with target, verb, scopes, security, and associated policies.
type Operation struct {
	Target                 string                  `yaml:"target,omitempty"`
	Verb                   string                  `yaml:"verb,omitempty"`
	Scopes                 []string                `yaml:"scopes,omitempty"`
	Secured                bool                    `yaml:"secured,omitempty"`
	EndpointConfigurations *EndpointConfigurations `yaml:"endpointConfigurations,omitempty"`
	OperationPolicies      *OperationPolicies      `yaml:"operationPolicies,omitempty"`
	RateLimit              *RateLimit              `yaml:"rateLimit,omitempty"`
}

// RateLimit is a placeholder for future rate-limiting configuration.
type RateLimit struct {
	RequestsPerUnit int    `yaml:"requestsPerUnit,omitempty"`
	Unit            string `yaml:"unit,omitempty"`
}

// VHost defines virtual hosts for production and sandbox environments.
type VHost struct {
	Production []string `yaml:"production,omitempty"`
	Sandbox    []string `yaml:"sandbox,omitempty"`
}

// AIProvider represents the AI provider configuration.
type AIProvider struct {
	Name       string `yaml:"name,omitempty"`
	APIVersion string `yaml:"apiVersion,omitempty"`
}

// CORSConfiguration represents the CORS (Cross-Origin Resource Sharing) configuration for an API.
type CORSConfiguration struct {
	CORSConfigurationEnabled      bool     `yaml:"corsConfigurationEnabled"`
	AccessControlAllowOrigins     []string `yaml:"accessControlAllowOrigins"`
	AccessControlAllowCredentials bool     `yaml:"accessControlAllowCredentials"`
	AccessControlAllowHeaders     []string `yaml:"accessControlAllowHeaders"`
	AccessControlAllowMethods     []string `yaml:"accessControlAllowMethods"`
}

// API represents an main API type definition
type APKConf struct {
	Name                   string                  `yaml:"name,omitempty"`
	ID                     string                  `yaml:"id,omitempty"`
	Version                string                  `yaml:"version,omitempty"`
	BasePath               string                  `yaml:"basePath,omitempty"`
	Type                   string                  `yaml:"type,omitempty"`
	Environment            string                  `yaml:"environment,omitempty"`
	DefaultVersion         bool                    `yaml:"defaultVersion,omitempty"`
	DefinitionPath         string                  `yaml:"definitionPath,omitempty"`
	EndpointConfigurations *EndpointConfigurations `yaml:"endpointConfigurations,omitempty"`
	Operations             *[]Operation            `yaml:"operations,omitempty"`
	Authentication         *[]AuthConfiguration    `yaml:"authentication,omitempty"`
	CorsConfig             *CORSConfiguration      `yaml:"corsConfiguration,omitempty"`
	AdditionalProperties   *[]AdditionalProperty   `yaml:"additionalProperties,omitempty"`
	SubscriptionValidation bool                    `yaml:"subscriptionValidation,omitempty"`
	RateLimit              *RateLimit              `yaml:"rateLimit,omitempty"`
	APIPolicies            *OperationPolicies      `yaml:"apiPolicies,omitempty"`
	AIProvider             *AIProvider             `yaml:"aiProvider,omitempty"`
}

// Organization represents an organization configuration.
type Organization struct {
	UUID                     string                 `yaml:"uuid"`
	Name                     string                 `yaml:"name"`
	DisplayName              string                 `yaml:"displayName"`
	OrganizationClaimValue   string                 `yaml:"organizationClaimValue"`
	Enabled                  bool                   `yaml:"enabled"`
	ServiceListingNamespaces []string               `yaml:"serviceListingNamespaces"`
	Properties               []OrganizationProperty `yaml:"properties"`
}

// OrganizationProperty represents a custom property for an organization.
type OrganizationProperty struct {
	Name  string `yaml:"name"`
	Value string `yaml:"value"`
}
