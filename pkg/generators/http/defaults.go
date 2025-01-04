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

package http_generator

import (
	"errors"
	"fmt"
	"gw_artifacts/config/types"
	"gw_artifacts/pkg/utils"

	gwapiv1 "sigs.k8s.io/gateway-api/apis/v1"
)

// generateHTTPRouteRules generates a list of HTTPRouteRules based on the provided configurations.
func (g *httpRouteGenerator) generateHTTPRouteRules(apkConf types.APKConf, operations []types.Operation, endpoint *types.EndpointDetails, endpointType string) ([]gwapiv1.HTTPRouteRule, error) {
	var httpRouteRules []gwapiv1.HTTPRouteRule
	for _, operation := range operations {
		httpRouteRule, err := g.GenerateHTTPRouteRule(apkConf, operation, endpoint, endpointType)
		if err != nil {
			return nil, err
		} else {
			httpRouteRules = append(httpRouteRules, *httpRouteRule)
		}
	}
	return httpRouteRules, nil
}

// generateRouteRule generates a route rule based on the operation and endpoint details.
func (g *httpRouteGenerator) generateHTTPRouteRule(apkConf types.APKConf, operation types.Operation, endpoint *types.EndpointDetails, endpointType string) (*gwapiv1.HTTPRouteRule, error) {
	var endpointToUse *types.EndpointDetails = utils.GetEndpointToUse(operation.EndpointConfigurations, endpointType)
	if endpointToUse == nil && endpoint != nil {
		endpointToUse = endpoint
	}
	if endpointToUse != nil {
		var filters []gwapiv1.HTTPRouteFilter
		var hasRedirectPolicy bool
		filters, hasRedirectPolicy = g.GenerateHTTPRouteFilters(apkConf, *endpointToUse, operation, endpointType)
		matches, err := g.RetrieveHTTPMatches(apkConf, operation)
		if err != nil {
			return nil, err
		}

		httpRouteRule := gwapiv1.HTTPRouteRule{
			Matches: matches,
			Filters: filters,
		}
		if !hasRedirectPolicy {
			httpRouteRule.BackendRefs = g.GenerateHTTPBackEndRef(*endpointToUse, operation, endpointType)
		}
		return &httpRouteRule, nil
	} else {
		return nil, errors.New("invalid endpoint specified")
	}
}

// generateAndRetrieveParentRefs generates and retrieves the parent references for the HTTPRoute.
func (g *httpRouteGenerator) generateAndRetrieveParentRefs(gatewayConfig types.GatewayConfigurations, uniqueId string) []gwapiv1.ParentReference {
	var parentRefs = make([]gwapiv1.ParentReference, 0)
	gatewayName := gatewayConfig.Name
	listenerName := gwapiv1.SectionName(gatewayConfig.ListenerName)
	parentGroup := gwapiv1.Group("gateway.networking.k8s.io")
	parentKind := gwapiv1.Kind("Gateway")

	parentRef := gwapiv1.ParentReference{
		Group:       &parentGroup,
		Kind:        &parentKind,
		Name:        gwapiv1.ObjectName(gatewayName),
		SectionName: &listenerName,
	}
	parentRefs = append(parentRefs, parentRef)
	return parentRefs
}

// generateHTTPBackEndRef generates a list of HTTPBackendRefs based on the provided configurations.
func (g *httpRouteGenerator) generateHTTPBackEndRef(endpoint types.EndpointDetails, operation types.Operation, endpointType string) []gwapiv1.HTTPBackendRef {
	kind := gwapiv1.Kind("Service")
	httpBackEndRef := gwapiv1.HTTPBackendRef{
		BackendRef: gwapiv1.BackendRef{
			BackendObjectReference: gwapiv1.BackendObjectReference{
				Kind: &kind,
				Name: gwapiv1.ObjectName(endpoint.Name),
			},
		},
	}
	return []gwapiv1.HTTPBackendRef{httpBackEndRef}
}

// generateHTTPRouteFilters generates a list of HTTPRouteFilters based on the provided configurations.
func (g *httpRouteGenerator) generateHTTPRouteFilters(apkConf types.APKConf, endpointToUse types.EndpointDetails, operation types.Operation, endpointType string) ([]gwapiv1.HTTPRouteFilter, bool) {
	routeFilters := make([]gwapiv1.HTTPRouteFilter, 0)
	var operationPoliciesToUse *types.OperationPolicies
	hasRedirectPolicy := false
	operationPolicies := apkConf.APIPolicies
	if operationPolicies != nil {
		if operationPolicies.Request != nil || operationPolicies.Response != nil {
			operationPoliciesToUse = operationPolicies
		}
	} else {
		operationPoliciesToUse = operation.OperationPolicies
	}

	if operationPoliciesToUse != nil {
		requestPolicies := operationPoliciesToUse.Request
		responsePolicies := operationPoliciesToUse.Response

		if len(requestPolicies) > 0 {
			var requestHttpRouteFilters []gwapiv1.HTTPRouteFilter
			requestHttpRouteFilters, hasRedirectPolicy = g.ExtractHTTPRouteFilter(&apkConf, endpointToUse, operation, requestPolicies, true)
			routeFilters = append(routeFilters, requestHttpRouteFilters...)
		}
		if len(responsePolicies) > 0 {
			var responseHttpRouteFilters []gwapiv1.HTTPRouteFilter
			responseHttpRouteFilters, _ = g.ExtractHTTPRouteFilter(&apkConf, endpointToUse, operation, responsePolicies, false)
			routeFilters = append(routeFilters, responseHttpRouteFilters...)
		}
	}
	if !hasRedirectPolicy {
		generatedPath := utils.GeneratePrefixMatch(endpointToUse, operation)
		replacePathFilter := gwapiv1.HTTPRouteFilter{
			Type: "URLRewrite",
			URLRewrite: &gwapiv1.HTTPURLRewriteFilter{
				Path: &gwapiv1.HTTPPathModifier{
					Type:            gwapiv1.FullPathHTTPPathModifier,
					ReplaceFullPath: &generatedPath,
				},
			},
		}
		routeFilters = append(routeFilters, replacePathFilter)
	}
	return routeFilters, hasRedirectPolicy
}

// extractHTTPRouteFilter extracts the HTTPRouteFilters based on the provided configurations.
func (g *httpRouteGenerator) extractHTTPRouteFilter(apkConf *types.APKConf, endpoint types.EndpointDetails, operation types.Operation, operationPolicies []types.OperationPolicy, isRequest bool) ([]gwapiv1.HTTPRouteFilter, bool) {
	var httpRouteFilters = make([]gwapiv1.HTTPRouteFilter, 0)
	var addHeaders = make([]gwapiv1.HTTPHeader, 0)
	var setHeaders = make([]gwapiv1.HTTPHeader, 0)
	var removeHeaders = make([]string, 0)
	var hasRedirectPolicy bool = false

	for _, policy := range operationPolicies {
		if policyParameters, ok := policy.Parameters.(types.Header); ok {
			switch policy.PolicyName {
			case "AddHeader":
				addHeader := gwapiv1.HTTPHeader{
					Name:  gwapiv1.HTTPHeaderName(policyParameters.HeaderName),
					Value: policyParameters.HeaderValue}
				addHeaders = append(addHeaders, addHeader)
			case "SetHeader":
				setHeader := gwapiv1.HTTPHeader{
					Name:  gwapiv1.HTTPHeaderName(policyParameters.HeaderName),
					Value: policyParameters.HeaderValue}
				setHeaders = append(setHeaders, setHeader)
			case "RemoveHeader":
				removeHeaders = append(removeHeaders, policyParameters.HeaderName)
			}
		} else if policyParameters, ok := policy.Parameters.(types.URLList); ok {
			urls := policyParameters.URLs
			for _, url := range urls {
				mirrorFilter := gwapiv1.HTTPRouteFilter{Type: "RequestMirror"}
				if !isRequest {
					fmt.Println("Mirror filter cannot be appended as a response policy.")
				}
				port := utils.GetPort(url)
				if port > 0 {
					backendRef := g.GenerateHTTPBackEndRef(endpoint, operation, "")[0]
					mirrorFilter.RequestMirror = &gwapiv1.HTTPRequestMirrorFilter{
						BackendRef: gwapiv1.BackendObjectReference{
							Name:      backendRef.Name,
							Namespace: backendRef.Namespace,
							Group:     backendRef.Group,
							Kind:      backendRef.Kind,
							Port:      backendRef.Port,
						},
					}
					httpRouteFilters = append(httpRouteFilters, mirrorFilter)
				}
			}
		} else if policyParameters, ok := policy.Parameters.(types.RedirectPolicy); ok {
			hasRedirectPolicy = true
			if !isRequest {
				fmt.Println("Redirect filter cannot be appended as a response policy.")
			}
			url := policyParameters.URL
			redirectFilter := gwapiv1.HTTPRouteFilter{Type: "RequestRedirect"}
			port := utils.GetPort(url)
			if port > 0 {
				host := gwapiv1.PreciseHostname(utils.GetHost(types.EndpointURL(url)))
				schema := utils.GetProtocol(url)
				replaceFullPath := utils.GetPath(url)
				redirectFilter.RequestRedirect = &gwapiv1.HTTPRequestRedirectFilter{
					Hostname: &host,
					Scheme:   &schema,
					Path: &gwapiv1.HTTPPathModifier{
						Type:            "ReplaceFullPath",
						ReplaceFullPath: &replaceFullPath,
					},
				}
				if policyParameters.StatusCode > 0 {
					redirectFilter.RequestRedirect.StatusCode = &policyParameters.StatusCode
				}
			}
			httpRouteFilters = append(httpRouteFilters, redirectFilter)
		}
	}
	var headerModifier gwapiv1.HTTPHeaderFilter
	if len(addHeaders) != 0 {
		headerModifier.Add = addHeaders
	}
	if len(setHeaders) != 0 {
		headerModifier.Set = setHeaders
	}
	if len(removeHeaders) != 0 {
		headerModifier.Remove = removeHeaders
	}
	var headerModifierFilter gwapiv1.HTTPRouteFilter
	if isRequest {
		headerModifierFilter = gwapiv1.HTTPRouteFilter{
			Type:                  "RequestHeaderModifier",
			RequestHeaderModifier: &headerModifier,
		}
	} else {
		headerModifierFilter = gwapiv1.HTTPRouteFilter{
			Type:                   "ResponseHeaderModifier",
			ResponseHeaderModifier: &headerModifier,
		}
	}
	if len(addHeaders) > 0 || len(setHeaders) > 0 || len(removeHeaders) > 0 {
		httpRouteFilters = append(httpRouteFilters, headerModifierFilter)
	}
	return httpRouteFilters, hasRedirectPolicy
}

// retrieveHTTPMatches retrieves the HTTPRouteMatches based on the provided configurations.
func (g *httpRouteGenerator) retrieveHTTPMatches(apkConf types.APKConf, operation types.Operation) ([]gwapiv1.HTTPRouteMatch, error) {
	var httpRouteMatches []gwapiv1.HTTPRouteMatch
	httpRouteMatch, err := g.RetrieveHTTPMatch(apkConf, operation)
	if err != nil {
		return nil, err
	}
	httpRouteMatches = append(httpRouteMatches, httpRouteMatch)
	return httpRouteMatches, nil
}

// retrieveHTTPMatch retrieves the HTTPRouteMatch based on the provided configurations.
func (g *httpRouteGenerator) retrieveHTTPMatch(apkConf types.APKConf, operation types.Operation) (gwapiv1.HTTPRouteMatch, error) {
	method := gwapiv1.HTTPMethod(operation.Verb)
	pathType := gwapiv1.PathMatchRegularExpression
	operationTarget := "/*"
	if operation.Target != "" {
		operationTarget = operation.Target
	}
	pathValue := utils.RetrievePathPrefix(operationTarget, apkConf.BasePath)
	httpRouteMatch := gwapiv1.HTTPRouteMatch{
		Method: &method,
		Path: &gwapiv1.HTTPPathMatch{
			Type:  &pathType,
			Value: &pathValue,
		},
	}
	return httpRouteMatch, nil
}
